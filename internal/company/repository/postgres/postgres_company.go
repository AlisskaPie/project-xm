package postgres

import (
	"context"
	"fmt"

	"github.com/AlisskaPie/project-xm/pkg/domain"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type companyRepository struct {
	ctx context.Context
	db  *sqlx.DB
}

// Create implements domain.CompanyRepository
func (r *companyRepository) Create(ctx context.Context, c domain.CreateCompany) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}

	query, _, err := goqu.Insert("company").Rows(Company(c)).ToSQL()
	if err != nil {
		return fmt.Errorf("cannot build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("ExecContext: %w", err)
	}

	return nil
}

// Delete implements domain.CompanyRepository
func (r *companyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q, _, err := goqu.Delete("company").Where(goqu.Ex{"id": id.String()}).ToSQL()
	if err != nil {
		return fmt.Errorf("cannot build query: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("ExecContext: %w", err)
	}

	return nil
}

// GetByID implements domain.CompanyRepository
func (r *companyRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Company, error) {
	q, _, err := goqu.From("company").Where(goqu.Ex{"id": id.String()}).ToSQL()
	if err != nil {
		return domain.Company{}, fmt.Errorf("cannot build query: %w", err)
	}

	var res Company
	err = r.db.QueryRowxContext(ctx, q).StructScan(&res)
	if err != nil {
		return domain.Company{}, fmt.Errorf("QueryRowxContext: %w", err)
	}

	return domain.Company(res), nil
}

// Patch implements domain.CompanyRepository
func (r *companyRepository) Patch(ctx context.Context, id uuid.UUID, c domain.PatchCompany) (domain.Company, error) {
	updates := map[string]any{}

	if c.Description != nil {
		updates["description"] = *c.Description
	}
	if c.Name != nil {
		updates["name"] = *c.Name
	}

	if c.AmountOfEmployees != nil {
		updates["amount_of_employees"] = *c.AmountOfEmployees
	}

	if c.Registered != nil {
		updates["registered"] = *c.Registered
	}

	if c.CompanyType != nil {
		updates["type"] = string(*c.CompanyType)
	}

	q, _, err := goqu.Update("company").
		Set(updates).
		Where(goqu.Ex{"id": id.String()}).
		Returning(goqu.T("company").All()).
		ToSQL()
	if err != nil {
		return domain.Company{}, fmt.Errorf("cannot build query: %w", err)
	}

	var res Company
	err = r.db.QueryRowxContext(ctx, q).StructScan(&res)
	if err != nil {
		return domain.Company{}, fmt.Errorf("QueryRowxContext: %w", err)
	}

	return domain.Company(res), nil
}

// NewCompanyRepository creates an object that represent the company.Repository interface
func NewCompanyRepository(ctx context.Context, db *sqlx.DB) domain.CompanyRepository {
	return &companyRepository{
		ctx: ctx,
		db:  db,
	}
}
