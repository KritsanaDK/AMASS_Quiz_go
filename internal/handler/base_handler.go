package handler

import (
	"amass/internal/infra/logger"
	"amass/internal/models"
	"amass/internal/static"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type handlerFunc func(req interface{}, tracking *models.Tracking) (events.APIGatewayProxyResponse, error)

type BaseHandler struct {
	Debug     bool
	SkipCache bool
}

type Message struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
	Now     string `json:"now"`
}

func (h *BaseHandler) Health(c echo.Context) error {

	return c.JSON(http.StatusOK, Message{
		Name:    os.Getenv("NAME"),
		Version: os.Getenv("VERSION"),
		Env:     os.Getenv("ENV"),
		Now:     time.Now().Format("2006-01-02 15:04:05"),
	})
}

func (h *BaseHandler) ValidatorParam(req interface{}) error {
	v := validator.New()
	err := v.Struct(req)
	if err != nil {
		msg := ""
		for _, e := range err.(validator.ValidationErrors) {
			msg += fmt.Sprintf("%v ", e)
		}

		return fmt.Errorf("%s", msg)
	}

	return nil
}

func (h *BaseHandler) GenSuccessResponse(response interface{}, c echo.Context, tracking *models.Tracking) error {
	templateResponse := models.TemplateResponse{
		Code:    static.CODE_SUCCESS,
		Message: static.SUCCESS,
		Data:    response,
	}
	tracking.Response = templateResponse
	logger.InfoWithTracking(static.HANDLER, tracking)

	return c.JSON(http.StatusOK, templateResponse)
}

func (h *BaseHandler) GenCustomResponse(response interface{}, code, message string, c echo.Context, tracking *models.Tracking) error {
	templateResponse := models.TemplateResponse{
		Code:    code,
		Message: message,
		Data:    response,
	}
	tracking.Response = templateResponse
	logger.InfoWithTracking(static.HANDLER, tracking)

	return c.JSON(http.StatusOK, templateResponse)
}

func (h *BaseHandler) GenErrorResponse(code, message string, httpStatus int, c echo.Context, tracking *models.Tracking) error {
	templateResponse := models.TemplateResponse{
		Code:    code,
		Message: message,
	}
	tracking.Response = templateResponse
	logger.ErrorWithTracking(static.HANDLER, message, tracking)

	return c.JSON(http.StatusBadRequest, templateResponse)
}

func (h *BaseHandler) PreProcess(request interface{}, tracking *models.Tracking) error {
	logger.Info("PreProcess|Request=%+v", request)
	return nil
}

func (h *BaseHandler) PostProcess(response interface{}, tracking *models.Tracking) error {
	logger.Info("PostProcess|Response=%+v", response)
	return nil
}

func (h *BaseHandler) Process(request interface{}, tracking *models.Tracking, hFunc handlerFunc) (events.APIGatewayProxyResponse, error) {
	// Pre-process
	h.PreProcess(request, tracking)

	// Process
	response, err := hFunc(request, tracking)

	// Post-process
	h.PostProcess(response, tracking)

	return response, err
}
