package usecase

import (
	"context"
	"fmt"

	"github.com/AlisskaPie/project-xm/pkg/domain"

	"github.com/google/uuid"
)

type companyUsecase struct {
	companyRepo domain.CompanyRepository
}

// Create implements domain.CompanyUsecase
func (u *companyUsecase) Create(ctx context.Context, c domain.CreateCompany) error {
	if err := u.companyRepo.Create(ctx, c); err != nil {
		return fmt.Errorf("companyRepo.Create: %w", err)
	}
	return nil
}

// Delete implements domain.CompanyUsecase
func (u *companyUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.companyRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("companyRepo.Delete: %w", err)
	}
	return nil
}

// GetByID implements domain.CompanyUsecase
func (u *companyUsecase) GetByID(ctx context.Context, id uuid.UUID) (domain.Company, error) {
	res, err := u.companyRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Company{}, fmt.Errorf("companyRepo.GetByID: %w", err)
	}
	return res, nil
}

// Patch implements domain.CompanyUsecase
func (u *companyUsecase) Patch(ctx context.Context, id uuid.UUID, c domain.PatchCompany) (domain.Company, error) {
	company, err := u.companyRepo.Patch(ctx, id, c)
	if err != nil {
		return domain.Company{}, fmt.Errorf("companyRepo.Patch: %w", err)
	}
	return company, nil
}

// NewCompanyUsecase creates new usecase object representation of domain.CompanyUsecase interface
func NewCompanyUsecase(r domain.CompanyRepository) domain.CompanyUsecase {
	return &companyUsecase{
		companyRepo: r,
	}
}
