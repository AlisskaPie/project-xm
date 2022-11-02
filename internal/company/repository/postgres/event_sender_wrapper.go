package postgres

import (
	"context"
	"fmt"

	"github.com/AlisskaPie/project-xm/pkg/domain"

	"github.com/google/uuid"
)

type eventSenderWrapper struct {
	repo        domain.CompanyRepository
	eventSender domain.CompanyEventSender
}

// Create implements domain.CompanyRepository
func (r *eventSenderWrapper) Create(ctx context.Context, c domain.CreateCompany) error {
	if err := r.repo.Create(ctx, c); err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	if err := r.eventSender.Send(ctx, domain.CompanyEvent{
		Action: domain.InsertEventActionType,
		ID:     c.ID,
		State:  domain.Company(c),
	}); err != nil {
		return fmt.Errorf("failed to send patch event: %w", err)
	}

	return nil
}

// Delete implements domain.CompanyRepository
func (r *eventSenderWrapper) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("repo.Delete: %w", err)
	}

	if err := r.eventSender.Send(ctx, domain.CompanyEvent{
		Action: domain.DeleteEventActionType,
		ID:     id,
	}); err != nil {
		return fmt.Errorf("failed to send patch event: %w", err)
	}

	return nil
}

// GetByID implements domain.CompanyRepository
func (r *eventSenderWrapper) GetByID(ctx context.Context, id uuid.UUID) (domain.Company, error) {
	return r.repo.GetByID(ctx, id)
}

// Patch implements domain.CompanyRepository
func (r *eventSenderWrapper) Patch(ctx context.Context, id uuid.UUID, c domain.PatchCompany) (domain.Company, error) {
	company, err := r.repo.Patch(ctx, id, c)
	if err != nil {
		return company, fmt.Errorf("repo.Patch: %w", err)
	}

	if err := r.eventSender.Send(ctx, domain.CompanyEvent{
		Action: domain.UpdateEventActionType,
		ID:     id,
		State:  company,
	}); err != nil {
		return domain.Company{}, fmt.Errorf("failed to send patch event: %w", err)
	}

	return company, nil
}

func NewEventSenderWrapper(repo domain.CompanyRepository, eventSender domain.CompanyEventSender) domain.CompanyRepository {
	return &eventSenderWrapper{
		eventSender: eventSender,
		repo:        repo,
	}
}
