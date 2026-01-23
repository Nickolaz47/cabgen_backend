package country

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

type AdminCountryHandler struct {
	Service services.CountryService
}

func NewAdminCountryHandler(svc services.CountryService) *AdminCountryHandler {
	return &AdminCountryHandler{
		Service: svc,
	}
}

func (h *AdminCountryHandler) GetCountries(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)

	countries, err := h.Service.FindAll(c.Request.Context(), language)
	if err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: countries})
}

func (h *AdminCountryHandler) GetCountriesByName(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	language := translation.GetLanguageFromContext(c)
	name := c.Query("name")

	var (
		countries []models.CountryFormResponse
		err       error
	)

	if name == "" {
		countries, err = h.Service.FindAll(c.Request.Context(), language)
	} else {
		countries, err = h.Service.FindByName(c.Request.Context(), name, language)
	}

	if err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: countries})
}

func (h *AdminCountryHandler) GetCountryByCode(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	code := c.Param("code")

	country, err := h.Service.FindByCode(c.Request.Context(), code)
	if err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: country})
}

func (h *AdminCountryHandler) CreateCountry(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	var newCountry models.CountryCreateInput

	if errMsg, valid := validations.Validate(c, localizer, &newCountry); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{Error: errMsg})
		return
	}

	errMsg, ok := validations.ValidateTranslationMap(c, "country", newCountry.Names)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	country, err := h.Service.Create(c.Request.Context(), newCountry)
	if err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Data:    country,
		Message: responses.GetResponse(localizer, responses.CountryCreateSuccess),
	})
}

func (h *AdminCountryHandler) UpdateCountry(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	code := c.Param("code")

	var countryUpdateInput models.CountryUpdateInput
	errMsg, ok := validations.Validate(c, localizer, &countryUpdateInput)
	if !ok {
		c.JSON(http.StatusBadRequest,
			responses.APIResponse{
				Error: errMsg,
			},
		)
		return
	}

	errMsg, ok = validations.ValidateTranslationMap(c, "country", countryUpdateInput.Names)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	countryUpdated, err := h.Service.Update(c.Request.Context(), code, countryUpdateInput)
	if err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Data: countryUpdated,
	})
}

func (h *AdminCountryHandler) DeleteCountry(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	code := c.Param("code")

	if err := h.Service.Delete(c.Request.Context(), code); err != nil {
		code, errMsg := handlererrors.HandleCountryError(err)
		c.JSON(
			code,
			responses.APIResponse{
				Error: responses.GetResponse(localizer, errMsg),
			})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.CountryDeleted),
	})
}
