package handler

import (
	"amass/internal/models"
	"amass/internal/service"
	"amass/internal/utils"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	BaseHandler
	service service.IUserService
}

func NewUserHandler(debug bool, s service.AllService) *userHandler {
	return &userHandler{
		BaseHandler: BaseHandler{Debug: debug},
		service:     s.IUserService,
	}
}

func (h *userHandler) Create(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	defer c.Request().Body.Close()

	tracking := &models.Tracking{
		Track:   utils.GenUUID(),
		Request: utils.BytesToString(body),
		URI:     c.Request().RequestURI,
		Method:  c.Request().Method,
	}

	req := models.UserLogin{}
	err = utils.BytesToStruct(body, &req)
	if err != nil {
		return h.GenErrorResponse("200", "GetPaperlessWithCustomerID| INVALID PARAM", http.StatusBadRequest, c, tracking)
	}

	err = h.ValidatorParam(req)
	if err != nil {
		return h.GenErrorResponse("200", "GetPaperlessWithCustomerID| INVALID PARAM", http.StatusBadRequest, c, tracking)
	}

	err = h.service.Create(&req)
	if err != nil {
		return h.GenErrorResponse("ERROR_CODE", "CreateUser| "+err.Error(), http.StatusInternalServerError, c, tracking)
	}

	return h.GenSuccessResponse(nil, c, tracking)

}

func (h *userHandler) Login(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	defer c.Request().Body.Close()

	tracking := &models.Tracking{
		Track:   utils.GenUUID(),
		Request: utils.BytesToString(body),
		URI:     c.Request().RequestURI,
		Method:  c.Request().Method,
	}

	req := models.UserLogin{}
	err = utils.BytesToStruct(body, &req)
	if err != nil {
		return h.GenErrorResponse("200", "Login| INVALID PARAM", http.StatusBadRequest, c, tracking)
	}

	err = h.ValidatorParam(req)

	if err != nil {
		return h.GenErrorResponse("200", "Login| INVALID PARAM", http.StatusBadRequest, c, tracking)
	}

	user, err := h.service.Login(&req)
	if err != nil {
		return h.GenErrorResponse("ERROR_CODE", "Login| "+err.Error(), http.StatusInternalServerError, c, tracking)
	}

	return h.GenSuccessResponse(user, c, tracking)
}

func (h *userHandler) GetUser(c echo.Context) error {
	tracking := &models.Tracking{
		Track:   utils.GenUUID(),
		Request: "from_jwt",
		URI:     c.Request().RequestURI,
		Method:  c.Request().Method,
	}

	username, _ := c.Get("username").(string)
	username = strings.TrimSpace(username)
	if username == "" {
		return h.GenErrorResponse("200", "GetUser| INVALID TOKEN CLAIM", http.StatusBadRequest, c, tracking)
	}

	user, err := h.service.GetUser(username)
	if err != nil {
		return h.GenErrorResponse("ERROR_CODE", fmt.Sprintf("GetUser| %v", err), http.StatusInternalServerError, c, tracking)
	}

	return h.GenSuccessResponse(user, c, tracking)
}
