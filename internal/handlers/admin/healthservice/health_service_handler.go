package healthservice

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHealthServiceHandler struct {
	Service services.HealthServiceService
}

func NewAdminHealthServiceHandler(
	svc services.HealthServiceService) *AdminHealthServiceHandler {
	return &AdminHealthServiceHandler{
		Service: svc,
	}
}

func (h *AdminHealthServiceHandler) GetAllHealthServices(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	healthServices, err := h.Service.FindAll(c.Request.Context())
	if err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: healthServices,
	})
}

func (h *AdminHealthServiceHandler) GetHealthServiceByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("healthServiceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	healthService, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: healthService})
}

func (h *AdminHealthServiceHandler) GetHealthServicesByName(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	name := c.Query("name")

	var (
		healthServices []models.HealthServiceAdminTableResponse
		err            error
	)

	if name == "" {
		healthServices, err = h.Service.FindAll(c.Request.Context())
	} else {
		healthServices, err = h.Service.FindByName(c.Request.Context(), name)
	}

	if err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: healthServices})
}

func (h *AdminHealthServiceHandler) CreateHealthService(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	var newHealthService models.HealthServiceCreateInput

	if errMsg, valid := validations.Validate(c, localizer,
		&newHealthService); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	if !newHealthService.Type.IsValid() {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: responses.GetResponse(
					localizer, responses.HealthServiceInvalidType),
			})
		return
	}

	healthService, err := h.Service.Create(
		c.Request.Context(), newHealthService)
	if err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			},
		)
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data: healthService,
		Message: responses.GetResponse(
			localizer,
			responses.HealthServiceCreationSuccess),
	})
}

func (h *AdminHealthServiceHandler) UpdateHealthService(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("healthServiceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	var healthServiceUpdateInput models.HealthServiceUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &healthServiceUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			})
		return
	}

	healthServiceUpdated, err := h.Service.Update(
		c.Request.Context(), id, healthServiceUpdateInput)
	if err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: healthServiceUpdated,
	})
}

func (h *AdminHealthServiceHandler) DeleteHealthService(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("healthServiceId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleHealthServiceError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			},
		)
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(
			localizer, responses.HealthServiceDeleted),
	})
}
