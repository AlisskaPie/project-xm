package domain

import (
	"context"

	"github.com/google/uuid"
)

// CompanyUsecase represent the company's usecases
type CompanyUsecase interface {
	Create(ctx context.Context, c CreateCompany) error
	GetByID(ctx context.Context, id uuid.UUID) (Company, error)
	Patch(ctx context.Context, id uuid.UUID, c PatchCompany) (Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PatchCompany struct {
	Name              *string
	Description       *string
	AmountOfEmployees *uint32
	Registered        *bool
	CompanyType       *CompanyType
}

type CreateCompany struct {
	ID                uuid.UUID
	Name              string
	Description       string
	AmountOfEmployees uint32
	Registered        bool
	CompanyType       CompanyType
}
