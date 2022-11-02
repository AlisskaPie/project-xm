package http

import (
	"net/http"

	"github.com/AlisskaPie/project-xm/pkg/domain"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// CompanyHandler represent the httphandler for company
type CompanyHandler struct {
	Usecase domain.CompanyUsecase
	log     zerolog.Logger
}

// NewCompanyHandler will initialize the companies resources endpoint
func NewCompanyHandler(e *echo.Echo, us domain.CompanyUsecase, auth echo.MiddlewareFunc, log zerolog.Logger) *CompanyHandler {
	handler := &CompanyHandler{
		Usecase: us,
		log:     log,
	}
	e.PATCH("/companies/:id", handler.Patch, auth)
	e.POST("/companies", handler.Create, auth)
	e.GET("/companies/:id", handler.GetByID)
	e.DELETE("/companies/:id", handler.Delete, auth)

	return handler
}

// Create creates the company by given request body
func (h *CompanyHandler) Create(c echo.Context) error {
	req := CompanyPostRequest{}

	if err := req.BindValidate(c); err != nil {
		h.log.Err(err).Msg("error while create binding")
		return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(domain.ErrBadRequest))
	}

	if err := h.Usecase.Create(c.Request().Context(), req.ToCreateCompany()); err != nil {
		h.log.Err(err).Msg("failed to create company by use case")
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(domain.ErrInternalError))
	}

	return c.NoContent(http.StatusCreated)
}

// GetByID gets company by given id
func (h *CompanyHandler) GetByID(c echo.Context) error {
	idReq := &IDPathRequest{}
	if err := idReq.BindValidate(c); err != nil {
		h.log.Err(err).Msg("failed to bind IDPathRequest")
		return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(domain.ErrBadRequest))
	}

	company, err := h.Usecase.GetByID(c.Request().Context(), idReq.ID)
	if err != nil {
		h.log.Err(err).Msg("GetByID error")
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(domain.ErrInternalError))
	}

	return c.JSON(http.StatusOK, GetCompanyResponseFromDomain(company))
}

// Patch patches the company by given request body
func (h *CompanyHandler) Patch(c echo.Context) (err error) {
	req := &CompanyPatchRequest{}
	if err := req.BindValidate(c); err != nil {
		h.log.Err(err).Msg("failed to bind CompanyPatchRequest")
		return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(domain.ErrBadRequest))
	}

	company, err := h.Usecase.Patch(c.Request().Context(), req.ID, req.ToPatchCompany())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(domain.ErrInternalError))
	}

	return c.JSON(http.StatusOK, GetCompanyResponseFromDomain(company))
}

// Delete deletes company by given param
func (h *CompanyHandler) Delete(c echo.Context) error {
	idReq := &IDPathRequest{}
	if err := idReq.BindValidate(c); err != nil {
		h.log.Err(err).Msg("failed to bind IDPathRequest")
		return c.JSON(http.StatusUnprocessableEntity, NewErrorResponse(domain.ErrBadRequest))
	}

	if err := h.Usecase.Delete(c.Request().Context(), idReq.ID); err != nil {
		h.log.Err(err).Msg("failed to delete by usecase")
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(domain.ErrInternalError))
	}

	return c.NoContent(http.StatusNoContent)
}
