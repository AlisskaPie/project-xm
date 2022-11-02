package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/AlisskaPie/project-xm/pkg/domain"
	"github.com/AlisskaPie/project-xm/pkg/domain/mocks"
)

func TestCreateSuccess(t *testing.T) {
	var mockCompanyPostRequest CompanyPostRequest
	err := gofakeit.Struct(&mockCompanyPostRequest)
	mockCompanyPostRequest.CompanyType = domain.CorporationsType
	assert.NoError(t, err)
	js, err := json.Marshal(mockCompanyPostRequest)
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("Create", mock.Anything, mockCompanyPostRequest.ToCreateCompany()).
		Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/companies", bytes.NewReader(js))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Create(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUseCase.AssertExpectations(t)
}

func TestCreateFailed_Validation(t *testing.T) {
	mockCompanyPatchRequest := CompanyPostRequest{}

	js, err := json.Marshal(mockCompanyPatchRequest)
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/companies", bytes.NewReader(js))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Create(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with invalid request parameters"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestCreateFailed_InternalError(t *testing.T) {
	var mockCompanyPostRequest CompanyPostRequest
	err := gofakeit.Struct(&mockCompanyPostRequest)
	mockCompanyPostRequest.CompanyType = domain.CorporationsType
	assert.NoError(t, err)
	js, err := json.Marshal(mockCompanyPostRequest)
	assert.NoError(t, err)

	expError := errors.New("some error")
	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("Create", mock.Anything, mock.Anything).Return(expError)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/companies", bytes.NewReader(js))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Create(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with internal error"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestGetByIDSuccess(t *testing.T) {
	var mockCompany domain.Company
	err := gofakeit.Struct(&mockCompany)
	assert.NoError(t, err)

	var mockCompanyIDRequest IDPathRequest
	err = gofakeit.Struct(&mockCompanyIDRequest)
	assert.NoError(t, err)
	mockCompany.ID = mockCompanyIDRequest.ID

	js, err := json.Marshal(GetCompanyResponseFromDomain(mockCompany))
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("GetByID", mock.Anything, mockCompanyIDRequest.ID).Return(mockCompany, nil)

	e := echo.New()
	req, err := http.NewRequest(
		echo.GET,
		fmt.Sprintf("/companies/%s", mockCompanyIDRequest.ID.String()),
		bytes.NewReader(js),
	)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/companies/:id")
	c.SetParamNames("id")
	c.SetParamValues(mockCompanyIDRequest.ID.String())
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t,
		strings.Trim(string(js), " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestGetByIDFailed_Validation(t *testing.T) {
	var mockCompany domain.Company
	err := gofakeit.Struct(&mockCompany)
	assert.NoError(t, err)

	var mockCompanyIDRequest IDPathRequest
	err = gofakeit.Struct(&mockCompanyIDRequest)
	assert.NoError(t, err)
	mockCompanyIDRequest.ID = uuid.Nil
	mockCompany.ID = mockCompanyIDRequest.ID

	js, err := json.Marshal(mockCompany)
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}

	e := echo.New()
	req, err := http.NewRequest(
		echo.GET,
		"/companies/",
		bytes.NewReader(js),
	)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/companies/:id")
	c.SetParamNames("id")
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with invalid request parameters"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestGetByIDFailed_InternalError(t *testing.T) {
	var mockCompany domain.Company
	err := gofakeit.Struct(&mockCompany)
	assert.NoError(t, err)

	var mockCompanyIDRequest IDPathRequest
	err = gofakeit.Struct(&mockCompanyIDRequest)
	assert.NoError(t, err)
	mockCompany.ID = mockCompanyIDRequest.ID
	js, err := json.Marshal(mockCompany)
	assert.NoError(t, err)

	expError := errors.New("some error")
	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("GetByID", mock.Anything, mock.Anything).Return(domain.Company{}, expError)

	e := echo.New()
	req, err := http.NewRequest(
		echo.GET,
		"/companies/",
		bytes.NewReader(js),
	)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/companies/:id")
	c.SetParamNames("id")
	c.SetParamValues(mockCompanyIDRequest.ID.String())
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.GetByID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with internal error"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestPatchSuccess(t *testing.T) {
	var mockCompanyPatchRequest CompanyPatchRequest
	err := gofakeit.Struct(&mockCompanyPatchRequest)
	assert.NoError(t, err)
	js1, err := json.Marshal(mockCompanyPatchRequest)
	assert.NoError(t, err)

	var mockCompany domain.Company
	err = gofakeit.Struct(&mockCompany)
	assert.NoError(t, err)
	js2, err := json.Marshal(GetCompanyResponseFromDomain(mockCompany))
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("Patch", mock.Anything, mockCompanyPatchRequest.ID, mockCompanyPatchRequest.ToPatchCompany()).
		Return(mockCompany, nil)

	req, err := http.NewRequest(
		echo.PATCH,
		fmt.Sprintf("/companies/%s", mockCompanyPatchRequest.ID.String()),
		bytes.NewReader(js1),
	)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	c.SetPath("/companies/:id")
	c.SetParamNames("id")
	c.SetParamValues(mockCompanyPatchRequest.ID.String())

	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Patch(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t,
		strings.Trim(string(js2), " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestPatchFailed_Validation(t *testing.T) {
	mockCompanyPatchRequest := CompanyPatchRequest{}
	js, err := json.Marshal(mockCompanyPatchRequest)
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}

	e := echo.New()
	req, err := http.NewRequest(echo.PATCH, "/companies/123", bytes.NewReader(js))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Patch(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with invalid request parameters"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestPatchFailed_InternalError(t *testing.T) {
	var mockCompanyPatchRequest CompanyPatchRequest
	err := gofakeit.Struct(&mockCompanyPatchRequest)
	assert.NoError(t, err)
	js, err := json.Marshal(mockCompanyPatchRequest)
	assert.NoError(t, err)

	expError := errors.New("some error")
	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("Patch", mock.Anything, mock.Anything, mock.Anything).Return(domain.Company{}, expError)

	e := echo.New()
	req, err := http.NewRequest(echo.PATCH, "/companies/123", bytes.NewReader(js))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Patch(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with internal error"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestDeleteSuccess(t *testing.T) {
	mockIDRequest := &IDPathRequest{}
	err := gofakeit.Struct(&mockIDRequest)
	assert.NoError(t, err)

	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("Delete", mock.Anything, mockIDRequest.ID).Return(nil)

	req, err := http.NewRequest(
		echo.DELETE,
		fmt.Sprintf("/companies/%s", mockIDRequest.ID.String()),
		nil,
	)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	e := echo.New()

	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	c := e.NewContext(req, rec)
	c.SetPath("/companies/:id")
	c.SetParamNames("id")
	c.SetParamValues(mockIDRequest.ID.String())
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUseCase.AssertExpectations(t)
}

func TestDeleteFailed_Validation(t *testing.T) {
	mockUseCase := &mocks.CompanyUsecase{}

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/companies/123", nil)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with invalid request parameters"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}

func TestDeleteFailed_InternalError(t *testing.T) {
	var mockIDPathRequest IDPathRequest
	err := gofakeit.Struct(&mockIDPathRequest)
	assert.NoError(t, err)

	expError := errors.New("some error")
	mockUseCase := &mocks.CompanyUsecase{}
	mockUseCase.On("Delete", mock.Anything, mock.Anything).Return(expError)

	e := echo.New()
	req, err := http.NewRequest(
		echo.DELETE,
		fmt.Sprintf("/companies/%s", mockIDPathRequest.ID.String()),
		nil,
	)
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/companies/:id")
	c.SetParamNames("id")
	c.SetParamValues(mockIDPathRequest.ID.String())
	handler := NewCompanyHandler(e, mockUseCase, nil, zerolog.New(io.Discard))
	err = handler.Delete(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t,
		strings.Trim(`{"message":"failed with internal error"}`, " \n"),
		strings.Trim(rec.Body.String(), " \n"),
	)
	mockUseCase.AssertExpectations(t)
}
