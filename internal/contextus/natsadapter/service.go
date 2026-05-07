package natsadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/JamesPagetButler/contextus/internal/contextus/synthesis"
	"github.com/JamesPagetButler/contextus/pkg/types"
)

// Service binds a Subscriber to a synthesis.Agent.
//
// Service is constructed by the cmd binary; tests construct it directly with
// a fake Subscriber. Service does not own the Agent or the Subscriber — both
// are injected — so the binary controls lifecycle.
type Service struct {
	Sub    Subscriber
	Agent  *synthesis.Agent
	Logger *slog.Logger

	subs []Unsubscriber
}

// New returns a Service ready to Start. The logger may be nil; in that case
// the package falls back to slog.Default.
func New(sub Subscriber, agent *synthesis.Agent, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		Sub:    sub,
		Agent:  agent,
		Logger: logger,
	}
}

// Start subscribes to the three Spec v1.3 §5.3 subjects and registers
// per-subject handlers that JSON-decode payloads into the typed §11.3
// structs and dispatch into the synthesis.Agent. Start does not block.
//
// Start is not idempotent: calling it twice creates duplicate subscriptions.
// Stop must be called before Start may be called again.
func (s *Service) Start(ctx context.Context) error {
	type binding struct {
		subject string
		handler func(payload []byte)
	}
	bindings := []binding{
		{
			subject: synthesis.SubjectEdgeBoundary,
			handler: s.handleEdgeBoundary,
		},
		{
			subject: synthesis.SubjectCorpusDiversity,
			handler: s.handleCorpusDiversity,
		},
		{
			subject: synthesis.SubjectBridgeIntervention,
			handler: s.handleBridgeIntervention,
		},
	}
	for _, b := range bindings {
		u, err := s.Sub.Subscribe(ctx, b.subject, b.handler)
		if err != nil {
			// Roll back partial subscriptions.
			s.unsubscribeAll()
			return fmt.Errorf("natsadapter: start subscribe %q: %w", b.subject, err)
		}
		s.subs = append(s.subs, u)
		s.Logger.Info("natsadapter: subscribed", slog.String("subject", b.subject))
	}
	return nil
}

// Stop unsubscribes from all bound subjects. It does NOT drain the underlying
// Subscriber — the caller owns that. Stop is idempotent.
func (s *Service) Stop() {
	s.unsubscribeAll()
}

func (s *Service) unsubscribeAll() {
	for _, u := range s.subs {
		if err := u.Unsubscribe(); err != nil {
			s.Logger.Warn("natsadapter: unsubscribe failed", slog.String("err", err.Error()))
		}
	}
	s.subs = nil
}

// handleEdgeBoundary decodes payload as EdgeScoutFlag and dispatches to the
// Agent. Decode errors and Agent errors are logged at the top of the call
// stack per go-coding-guide.md "One log per error, at the top of the call
// stack."
func (s *Service) handleEdgeBoundary(payload []byte) {
	var flag types.EdgeScoutFlag
	if err := json.Unmarshal(payload, &flag); err != nil {
		s.Logger.Error("natsadapter: edge-boundary decode failed",
			slog.String("err", err.Error()),
			slog.Int("bytes", len(payload)),
		)
		return
	}
	id, ok, err := s.Agent.HandleEdgeBoundary(context.Background(), flag)
	if err != nil {
		s.Logger.Error("natsadapter: edge-boundary handle failed",
			slog.String("err", err.Error()),
			slog.String("session_id", flag.SessionID),
		)
		return
	}
	s.logDispatch("edge-boundary", flag.SessionID, ok, id)
}

func (s *Service) handleCorpusDiversity(payload []byte) {
	var report types.CorpusDiversityReport
	if err := json.Unmarshal(payload, &report); err != nil {
		s.Logger.Error("natsadapter: corpus-diversity decode failed",
			slog.String("err", err.Error()),
			slog.Int("bytes", len(payload)),
		)
		return
	}
	id, ok, err := s.Agent.HandleCorpusDiversity(context.Background(), report)
	if err != nil {
		s.Logger.Error("natsadapter: corpus-diversity handle failed",
			slog.String("err", err.Error()),
			slog.String("search_id", report.SearchID),
		)
		return
	}
	s.logDispatch("corpus-diversity", report.SearchID, ok, id)
}

func (s *Service) handleBridgeIntervention(payload []byte) {
	var iv types.BridgeIntervention
	if err := json.Unmarshal(payload, &iv); err != nil {
		s.Logger.Error("natsadapter: bridge-intervention decode failed",
			slog.String("err", err.Error()),
			slog.Int("bytes", len(payload)),
		)
		return
	}
	id, ok, err := s.Agent.HandleBridgeIntervention(context.Background(), iv)
	if err != nil {
		s.Logger.Error("natsadapter: bridge-intervention handle failed",
			slog.String("err", err.Error()),
			slog.String("session_id", iv.SessionID),
		)
		return
	}
	s.logDispatch("bridge-intervention", iv.SessionID, ok, id)
}

func (s *Service) logDispatch(source, sessionOrSearchID string, promoted bool, id types.Addr) {
	if !promoted {
		s.Logger.Debug("natsadapter: below threshold; ephemeral",
			slog.String("source", source),
			slog.String("id", sessionOrSearchID),
		)
		return
	}
	s.Logger.Info("natsadapter: minted insight signal",
		slog.String("source", source),
		slog.String("session_or_search_id", sessionOrSearchID),
		slog.String("signal_id", fmt.Sprintf("%x", id)),
	)
}
