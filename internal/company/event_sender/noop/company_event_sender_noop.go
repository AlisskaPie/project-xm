package noop

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/AlisskaPie/project-xm/pkg/domain"
)

// No operation (just logging) implementation of company event sender
type companyEventSenderNoop struct {
	log zerolog.Logger
}

func NewCompanyEventSenderNoop(log zerolog.Logger) domain.CompanyEventSender {
	return &companyEventSenderNoop{
		log: log,
	}
}

func (n *companyEventSenderNoop) Send(_ context.Context, event domain.CompanyEvent) error {
	n.log.Info().Interface("event", event).Msg("noop event has been sent")

	return nil
}
