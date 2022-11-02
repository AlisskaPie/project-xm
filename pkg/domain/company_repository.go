package domain

import (
	"context"

	"github.com/google/uuid"
)

// CompanyRepository represent the company's repository contract
type CompanyRepository interface {
	Create(ctx context.Context, c CreateCompany) error
	GetByID(ctx context.Context, id uuid.UUID) (Company, error)
	Patch(ctx context.Context, id uuid.UUID, c PatchCompany) (Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
