package http

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/AlisskaPie/project-xm/pkg/domain"
)

type CompanyPostRequest struct {
	ID                uuid.UUID          `json:"id"`
	Name              string             `json:"name" validate:"required"`
	Description       string             `json:"description,omitempty"`
	AmountOfEmployees uint32             `json:"amount_of_employees" validate:"required"`
	Registered        bool               `json:"registered"`
	CompanyType       domain.CompanyType `json:"type" validate:"required"`
}

func (c *CompanyPostRequest) BindValidate(ctx echo.Context) error {
	if err := ctx.Bind(c); err != nil {
		return fmt.Errorf("failed to bind CompanyPostRequest: %w", err)
	}

	return c.Validate()
}

func (c *CompanyPostRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (c *CompanyPostRequest) ToCreateCompany() domain.CreateCompany {
	return domain.CreateCompany{
		ID:                c.ID,
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		CompanyType:       c.CompanyType,
	}
}

type CompanyPatchRequest struct {
	ID                uuid.UUID           `param:"id" validate:"required"`
	Name              *string             `json:"name"`
	Description       *string             `json:"description,omitempty"`
	AmountOfEmployees *uint32             `json:"amount_of_employees"`
	Registered        *bool               `json:"registered"`
	CompanyType       *domain.CompanyType `json:"type"`
}

func (c *CompanyPatchRequest) BindValidate(ctx echo.Context) error {
	if err := ctx.Bind(c); err != nil {
		return fmt.Errorf("failed to bind CompanyPatchRequest: %w", err)
	}

	return c.Validate()
}

func (c *CompanyPatchRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (c *CompanyPatchRequest) ToPatchCompany() domain.PatchCompany {
	return domain.PatchCompany{
		Name:              c.Name,
		Description:       c.Description,
		AmountOfEmployees: c.AmountOfEmployees,
		Registered:        c.Registered,
		CompanyType:       c.CompanyType,
	}
}

type IDPathRequest struct {
	ID uuid.UUID `param:"id" validate:"required"`
}

func (i *IDPathRequest) BindValidate(ctx echo.Context) error {
	if err := ctx.Bind(i); err != nil {
		return fmt.Errorf("failed to bind IDPathRequest: %w", err)
	}

	return i.Validate()
}

func (c *IDPathRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

type CompanyResponse struct {
	ID                uuid.UUID          `json:"id" validate:"required"`
	Name              string             `json:"name" validate:"required"`
	Description       string             `json:"description,omitempty"`
	AmountOfEmployees uint32             `json:"amount_of_employees" validate:"required"`
	Registered        bool               `json:"registered" validate:"required"`
	CompanyType       domain.CompanyType `json:"type" validate:"required"`
}

func GetCompanyResponseFromDomain(d domain.Company) CompanyResponse {
	return CompanyResponse{
		ID:                d.ID,
		Name:              d.Name,
		Description:       d.Description,
		AmountOfEmployees: d.AmountOfEmployees,
		Registered:        d.Registered,
		CompanyType:       d.CompanyType,
	}
}
