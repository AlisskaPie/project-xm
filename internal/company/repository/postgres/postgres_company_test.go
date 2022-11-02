package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/AlisskaPie/project-xm/pkg/domain"
)

func TestPostgresCompanyCreate(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name          string
		rf            registerFunc
		createCompany domain.CreateCompany
		wantErr       error
	}{
		{
			name: "Success",
			createCompany: domain.CreateCompany{
				ID:                uuid.MustParse("10000000-0000-0000-0000-000000000000"),
				Name:              "1",
				Description:       "2",
				AmountOfEmployees: 3,
				Registered:        true,
				CompanyType:       domain.NonProfitType,
			},
			rf: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`^INSERT INTO "company" (.*) VALUES \(3, '2', '10000000-0000-0000-0000-000000000000', '1', TRUE, 'NonProfit'\)$`).
					WillReturnResult(driver.RowsAffected(1))
			},
			wantErr: nil,
		},
		{
			name:          "Failed",
			createCompany: domain.CreateCompany{},
			rf: func(s sqlmock.Sqlmock) {
				s.ExpectExec("^(.+)").
					WillReturnError(testErr)
			},
			wantErr: fmt.Errorf("ExecContext: %w", testErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			assert.NoError(t, err)
			tt.rf(dbMock)
			dbMock.ExpectClose()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := NewCompanyRepository(context.TODO(), sqlxDB)
			err = r.Create(context.TODO(), tt.createCompany)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPostgresCompanyDelete(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name    string
		rf      registerFunc
		uuid    uuid.UUID
		wantErr error
	}{
		{
			name: "Success",
			uuid: uuid.MustParse("10000000-0000-0000-0000-000000000000"),
			rf: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`^DELETE FROM "company" WHERE \("id" = '10000000-0000-0000-0000-000000000000'\)$`).
					WillReturnResult(driver.RowsAffected(1))
			},
			wantErr: nil,
		},
		{
			name: "Failed",
			uuid: uuid.Nil,
			rf: func(s sqlmock.Sqlmock) {
				s.ExpectExec("^DELETE (.+)").
					WillReturnError(testErr)
			},
			wantErr: fmt.Errorf("ExecContext: %w", testErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			assert.NoError(t, err)
			tt.rf(dbMock)
			dbMock.ExpectClose()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := NewCompanyRepository(context.TODO(), sqlxDB)
			err = r.Delete(context.TODO(), tt.uuid)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestPostgresCompanyPatch(t *testing.T) {
	tests := []struct {
		name         string
		rf           registerFunc
		patchCompany domain.PatchCompany
		company      domain.Company
		uuid         uuid.UUID
		wantErr      string
	}{
		{
			name: "Success",
			uuid: testUUID,
			patchCompany: domain.PatchCompany{
				Name:              getPointer("1"),
				Description:       getPointer("2"),
				AmountOfEmployees: getPointer(uint32(3)),
				Registered:        getPointer(true),
				CompanyType:       getPointer(domain.SoleProprietorshipType),
			},
			company: domain.Company{
				ID:                testUUID,
				Name:              "1",
				Description:       "2",
				AmountOfEmployees: 3,
				Registered:        true,
				CompanyType:       domain.SoleProprietorshipType,
			},
			rf: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "description", "amount_of_employees", "registered", "type",
				})
				rows.AddRow(
					testUUID.String(),
					"1", "2", 3, true, domain.SoleProprietorshipType,
				)
				s.ExpectQuery(`^UPDATE "company" SET "amount_of_employees"=3,"description"='2',"name"='1',"registered"=TRUE,"type"='Sole Proprietorship'`).
					WillReturnRows(rows)
			},
			wantErr: "",
		},
		{
			name: "Failed",
			uuid: uuid.Nil,
			rf: func(s sqlmock.Sqlmock) {
				s.ExpectExec("^(.+)").
					WillReturnError(errors.New("test error"))
			},
			wantErr: "cannot build query: goqu: no update values provided",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			assert.NoError(t, err)
			tt.rf(dbMock)
			dbMock.ExpectClose()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := NewCompanyRepository(context.TODO(), sqlxDB)
			c, err := r.Patch(context.TODO(), tt.uuid, tt.patchCompany)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
			assert.Equal(t, tt.company, c)
		})
	}
}

func TestPostgresCompanyGetByID(t *testing.T) {
	testErr := errors.New("test error")
	tests := []struct {
		name    string
		rf      registerFunc
		company domain.Company
		uuid    uuid.UUID
		wantErr error
	}{
		{
			name: "Success",
			uuid: testUUID,
			company: domain.Company{
				ID:                testUUID,
				Name:              "1",
				Description:       "2",
				AmountOfEmployees: 3,
				Registered:        true,
				CompanyType:       domain.CooperativeType,
			},
			rf: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "description", "amount_of_employees", "registered", "type",
				})
				rows.AddRow(
					testUUID.String(),
					"1", "2", 3, true, domain.CooperativeType,
				)
				s.ExpectQuery(`^SELECT \* FROM "company" WHERE \("id" = '10000000-0000-0000-0000-000000000000'\)`).
					WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name: "Failed",
			uuid: uuid.Nil,
			rf: func(s sqlmock.Sqlmock) {
				s.ExpectQuery("^(.+)").
					WillReturnError(testErr)
			},
			wantErr: fmt.Errorf("QueryRowxContext: %w", testErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			tt.rf(dbMock)
			dbMock.ExpectClose()

			sqlxDB := sqlx.NewDb(db, "sqlmock")
			r := NewCompanyRepository(context.TODO(), sqlxDB)
			company, err := r.GetByID(context.TODO(), tt.uuid)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.company, company)
		})
	}
}

type registerFunc func(sqlmock.Sqlmock)

// Generic to get pointer
func getPointer[T any](value T) *T {
	return &value
}
