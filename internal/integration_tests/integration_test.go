package integrationtests

import (
	"encoding/json"
	"net/http"
	"testing"

	delivery "github.com/AlisskaPie/project-xm/internal/company/delivery/http"
	"github.com/AlisskaPie/project-xm/pkg/domain"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// generated by jwt.io with secret key
	jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.zbgd5BNF1cqQ_prCEqIvBTjSxMS8bDLnJAE_wE-0Cxg"
)

func TestIntegration_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := ClientSetup()

	id := uuid.New()
	companyParams := delivery.CompanyPostRequest{
		ID:                id,
		Name:              gofakeit.LetterN(uint(14)),
		Description:       gofakeit.Phone(),
		AmountOfEmployees: uint32(gofakeit.IntRange(1, 20)),
		Registered:        gofakeit.Bool(),
		CompanyType:       domain.CorporationsType,
	}

	// generated by jwt.io with secret key
	t.Run("Create passed", func(t *testing.T) {
		resp, err := client.Create(companyParams, jwt)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Create duplicate error", func(t *testing.T) {
		resp, err := client.Create(companyParams, jwt)
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Create failed: invalid token", func(t *testing.T) {
		invalidJWT := "dsjhd"
		resp, err := client.Create(companyParams, invalidJWT)
		require.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Create failed: no authorization", func(t *testing.T) {
		emptyJWT := ""
		resp, err := client.Create(companyParams, emptyJWT)
		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Create failed: required fields not passed", func(t *testing.T) {
		newCompanyParams := delivery.CompanyPostRequest{
			Description: "",
		}
		resp, err := client.Create(newCompanyParams, jwt)
		require.Error(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestIntegration_GetByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := ClientSetup()

	id := uuid.New()
	companyParams := delivery.CompanyPostRequest{
		ID:                id,
		Name:              gofakeit.LetterN(uint(14)),
		Description:       gofakeit.Phone(),
		AmountOfEmployees: uint32(gofakeit.IntRange(1, 20)),
		Registered:        gofakeit.Bool(),
		CompanyType:       domain.CooperativeType,
	}

	_, err := client.Create(companyParams, jwt)
	require.NoError(t, err)

	t.Run("GetByID passed", func(t *testing.T) {
		idReq := delivery.IDPathRequest{
			ID: id,
		}
		resp, err := client.GetByID(idReq)
		require.NoError(t, err)

		company, err := toCompany(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, id, company.ID)
	})

	t.Run("GetByID failed: id not exist, internal error", func(t *testing.T) {
		idReq := delivery.IDPathRequest{
			ID: uuid.New(),
		}
		resp, err := client.GetByID(idReq)
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("GetByID failed: validation error", func(t *testing.T) {
		idReq := delivery.IDPathRequest{}
		resp, err := client.GetByID(idReq)
		require.Error(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

func TestIntegration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := ClientSetup()

	id := uuid.New()
	companyParams := delivery.CompanyPostRequest{
		ID:                id,
		Name:              gofakeit.LetterN(uint(14)),
		Description:       gofakeit.Phone(),
		AmountOfEmployees: uint32(gofakeit.IntRange(1, 20)),
		Registered:        gofakeit.Bool(),
		CompanyType:       domain.CooperativeType,
	}

	_, err := client.Create(companyParams, jwt)
	require.NoError(t, err)

	t.Run("Delete failed: no authorization", func(t *testing.T) {
		emptyJWT := ""
		idReq := delivery.IDPathRequest{
			ID: id,
		}
		resp, err := client.Delete(idReq, emptyJWT)
		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Delete failed: invalid token", func(t *testing.T) {
		invalidJWT := "bvc"
		idReq := delivery.IDPathRequest{
			ID: id,
		}
		resp, err := client.Delete(idReq, invalidJWT)
		require.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Delete failed: validation error", func(t *testing.T) {
		idReq := delivery.IDPathRequest{}
		resp, err := client.Delete(idReq, jwt)
		require.Error(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("Delete passed", func(t *testing.T) {
		idReq := delivery.IDPathRequest{
			ID: id,
		}
		resp, err := client.Delete(idReq, jwt)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func TestIntegration_Patch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := ClientSetup()

	id := uuid.New()
	companyParams := delivery.CompanyPostRequest{
		ID:                id,
		Name:              "some-name",
		Description:       gofakeit.Phone(),
		AmountOfEmployees: 9,
		Registered:        true,
		CompanyType:       domain.NonProfitType,
	}

	_, err := client.Create(companyParams, jwt)
	require.NoError(t, err)

	t.Run("Patch failed: no authorization", func(t *testing.T) {
		emptyJWT := ""
		patchReq := makeStandardPatchTemplate(id)
		resp, err := client.Patch(patchReq, emptyJWT)
		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Patch failed: invalid token", func(t *testing.T) {
		invalidJWT := "bvc"
		patchReq := makeStandardPatchTemplate(id)
		resp, err := client.Patch(patchReq, invalidJWT)
		require.Error(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Patch passed", func(t *testing.T) {
		patchReq := makeStandardPatchTemplate(id)

		resp, err := client.Patch(patchReq, jwt)
		require.NoError(t, err)

		company, err := toCompany(resp.Body)
		require.NoError(t, err)

		assert.Equal(t, patchReq.ID, company.ID)
		assert.Equal(t, *patchReq.Name, company.Name)
		assert.Equal(t, *patchReq.Description, company.Description)
		assert.Equal(t, *patchReq.AmountOfEmployees, company.AmountOfEmployees)
		assert.Equal(t, *patchReq.Registered, company.Registered)
		assert.Equal(t, *patchReq.CompanyType, company.CompanyType)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Patch failed: no updates", func(t *testing.T) {
		patchReq := delivery.CompanyPatchRequest{
			ID: id,
		}

		resp, err := client.Patch(patchReq, jwt)
		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func toCompany(b []byte) (delivery.CompanyResponse, error) {
	newCompany := new(delivery.CompanyResponse)
	err := json.Unmarshal(b, &newCompany)
	if err != nil {
		return delivery.CompanyResponse{}, err
	}
	return *newCompany, nil
}

func makeStandardPatchTemplate(id uuid.UUID) delivery.CompanyPatchRequest {
	name := gofakeit.LetterN(uint(14))
	description := gofakeit.Phone()
	amount := uint32(gofakeit.IntRange(1, 10))
	registered := true
	cType := domain.SoleProprietorshipType

	return delivery.CompanyPatchRequest{
		ID:                id,
		Name:              &name,
		Description:       &description,
		AmountOfEmployees: &amount,
		Registered:        &registered,
		CompanyType:       &cType,
	}
}
