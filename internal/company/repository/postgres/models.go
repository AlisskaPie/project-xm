package postgres

import (
	"github.com/AlisskaPie/project-xm/pkg/domain"

	"github.com/google/uuid"
)

type Company struct {
	ID                uuid.UUID          `db:"id"`
	Name              string             `db:"name"`
	Description       string             `db:"description"`
	AmountOfEmployees uint32             `db:"amount_of_employees"`
	Registered        bool               `db:"registered"`
	CompanyType       domain.CompanyType `db:"type"`
}

type PatchCompany struct {
	Name              *string             `db:"name"`
	Description       *string             `db:"description"`
	AmountOfEmployees *uint32             `db:"amount_of_employees"`
	Registered        *bool               `db:"registered"`
	CompanyType       *domain.CompanyType `db:"type"`
}
