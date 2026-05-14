package natsadapter

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/JamesPagetButler/contextus/internal/contextus/synthesis"
	"github.com/JamesPagetButler/contextus/pkg/types"
)

// fakeSubscriber records subscriptions and lets tests deliver synthetic
// payloads to the registered handlers. Goroutine-safe for the parallel
// test paths.
type fakeSubscriber struct {
	mu       sync.Mutex
	handlers map[string]func([]byte)
	closed   bool
	drained  bool
}

func newFakeSubscriber() *fakeSubscriber {
	return &fakeSubscriber{handlers: map[string]func([]byte){}}
}

func (f *fakeSubscriber) Subscribe(_ context.Context, subject string, handler func(payload []byte)) (Unsubscriber, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return nil, ErrSubscriberClosed
	}
	f.handlers[subject] = handler
	return &fakeUnsub{f: f, subject: subject}, nil
}

func (f *fakeSubscriber) Drain(_ context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.drained = true
	f.closed = true
	f.handlers = nil
	return nil
}

func (f *fakeSubscriber) deliver(subject string, payload []byte) {
	f.mu.Lock()
	h, ok := f.handlers[subject]
	f.mu.Unlock()
	if !ok {
		return
	}
	h(payload)
}

type fakeUnsub struct {
	f       *fakeSubscriber
	subject string
}

func (u *fakeUnsub) Unsubscribe() error {
	u.f.mu.Lock()
	defer u.f.mu.Unlock()
	delete(u.f.handlers, u.subject)
	return nil
}

func newTestService(t *testing.T) (*Service, *fakeSubscriber, *synthesis.NoopPersister) {
	t.Helper()
	fake := newFakeSubscriber()
	persister := &synthesis.NoopPersister{}
	agent := synthesis.New(persister)
	agent.Now = func() time.Time { return time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC) }
	logger := slog.New(slog.NewTextHandler(io.Discard, nil)) // silent during tests
	svc := New(fake, agent, logger)
	if err := svc.Start(context.Background()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	return svc, fake, persister
}

func TestService_StartSubscribesToThreeSubjects(t *testing.T) {
	_, fake, _ := newTestService(t)
	fake.mu.Lock()
	defer fake.mu.Unlock()
	want := []string{
		synthesis.SubjectEdgeBoundary,
		synthesis.SubjectCorpusDiversity,
		synthesis.SubjectBridgeIntervention,
	}
	for _, s := range want {
		if _, ok := fake.handlers[s]; !ok {
			t.Errorf("expected handler registered for %q; missing", s)
		}
	}
}

func TestService_EdgeBoundaryAboveThreshold_MintsStructural(t *testing.T) {
	_, fake, persister := newTestService(t)
	flag := types.EdgeScoutFlag{
		SessionID:        "s1",
		BoundaryNodeID:   "n-stream-soda-butte",
		UnexploredDomain: "vegetation-phenology",
		Significance:     0.92,
	}
	payload, err := json.Marshal(flag)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	fake.deliver(synthesis.SubjectEdgeBoundary, payload)
	if persister.Calls != 1 {
		t.Fatalf("Persister.MintSignal: want 1 call, got %d", persister.Calls)
	}
	if got := persister.LastSignal.AnomalyType; got != types.AnomalyStructural {
		t.Errorf("anomaly_type: want structural, got %s", got)
	}
}

func TestService_EdgeBoundaryBelowThreshold_NoMint(t *testing.T) {
	_, fake, persister := newTestService(t)
	flag := types.EdgeScoutFlag{Significance: 0.30}
	payload, _ := json.Marshal(flag)
	fake.deliver(synthesis.SubjectEdgeBoundary, payload)
	if persister.Calls != 0 {
		t.Errorf("Persister.MintSignal: want 0 calls, got %d", persister.Calls)
	}
}

func TestService_BridgeInterventionAboveThreshold_MintsNarrative(t *testing.T) {
	_, fake, persister := newTestService(t)
	iv := types.BridgeIntervention{
		SessionID:          "s1",
		ScoutFlagID:        "flag-snowpack-streamflow-2000",
		KnowledgeDomainGap: "vegetation-as-passive-vs-active",
		Confidence:         0.88,
	}
	payload, _ := json.Marshal(iv)
	fake.deliver(synthesis.SubjectBridgeIntervention, payload)
	if persister.Calls != 1 {
		t.Fatalf("Persister.MintSignal: want 1 call, got %d", persister.Calls)
	}
	if got := persister.LastSignal.AnomalyType; got != types.AnomalyNarrative {
		t.Errorf("anomaly_type: want narrative, got %s", got)
	}
}

func TestService_CorpusDiversityConverged_MintsStructural(t *testing.T) {
	_, fake, persister := newTestService(t)
	report := types.CorpusDiversityReport{
		SearchID:            "search-1",
		ConvergenceDetected: true,
		VocabConcentration:  0.85,
	}
	payload, _ := json.Marshal(report)
	fake.deliver(synthesis.SubjectCorpusDiversity, payload)
	if persister.Calls != 1 {
		t.Fatalf("Persister.MintSignal: want 1 call, got %d", persister.Calls)
	}
	if got := persister.LastSignal.AnomalyType; got != types.AnomalyStructural {
		t.Errorf("anomaly_type: want structural, got %s", got)
	}
}

func TestService_MalformedPayload_LoggedNotPanicked(t *testing.T) {
	// Decode failure must not crash the handler goroutine. Persister.Calls
	// must remain 0.
	_, fake, persister := newTestService(t)
	fake.deliver(synthesis.SubjectEdgeBoundary, []byte("{not valid json"))
	if persister.Calls != 0 {
		t.Errorf("Persister.MintSignal: want 0 calls (decode error), got %d", persister.Calls)
	}
}

func TestService_StopUnsubscribesAll(t *testing.T) {
	svc, fake, _ := newTestService(t)
	svc.Stop()
	fake.mu.Lock()
	defer fake.mu.Unlock()
	if len(fake.handlers) != 0 {
		t.Errorf("expected all handlers removed after Stop; %d remaining", len(fake.handlers))
	}
}

func TestService_StartTwiceWithoutStopFailsCleanly(t *testing.T) {
	// Start is documented as not idempotent. We don't promise cleanup if
	// Start is called twice; we promise at minimum that the Service does
	// not panic and returns a comprehensible error or duplicates handlers.
	// (Current implementation duplicates handlers; flag this as a known
	// limitation rather than enforce a strict contract.)
	svc, fake, _ := newTestService(t)
	if err := svc.Start(context.Background()); err != nil {
		t.Logf("Start (second call) returned error (acceptable): %v", err)
	}
	// Either way: subjects still resolve.
	fake.mu.Lock()
	defer fake.mu.Unlock()
	for _, s := range []string{
		synthesis.SubjectEdgeBoundary,
		synthesis.SubjectCorpusDiversity,
		synthesis.SubjectBridgeIntervention,
	} {
		if _, ok := fake.handlers[s]; !ok {
			t.Errorf("subject %q dropped after second Start; not acceptable", s)
		}
	}
}
