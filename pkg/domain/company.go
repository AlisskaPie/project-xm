package domain

import (
	"github.com/google/uuid"
)

// Company implements domain
type Company struct {
	ID                uuid.UUID
	Name              string
	Description       string
	AmountOfEmployees uint32
	Registered        bool
	CompanyType       CompanyType
}

// CompanyType implements enum for type
type CompanyType string

// Scope of CompanyType values
const (
	CorporationsType       CompanyType = "Corporations"
	NonProfitType          CompanyType = "NonProfit"
	CooperativeType        CompanyType = "Cooperative"
	SoleProprietorshipType CompanyType = "Sole Proprietorship"
)
