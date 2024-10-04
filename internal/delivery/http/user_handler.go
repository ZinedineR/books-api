package http

import (
	_ "books-api/internal/delivery/http/response"
	"books-api/internal/entity"
	service "books-api/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHTTPHandler struct {
	Handler
	UserService service.UserService
}

func NewUserHTTPHandler(user service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		UserService: user,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the provided username and password
// @Tags Users
// @Accept json
// @Produce json
// @Param register body entity.UserLogin true "Registration Request"
// @Success 200 {object} response.DataResponse{data=entity.UserLogin} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/register [post]
func (h UserHTTPHandler) Register(ctx *gin.Context) {
	request := entity.UserLogin{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	if errException := h.UserService.Register(ctx, &request); errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, request)
}

// Login godoc
// @Summary User login
// @Description Authenticates the user and returns an access token
// @Tags Users
// @Accept json
// @Produce json
// @Param login body entity.UserLogin true "Login Request"
// @Success 200 {object} response.DataResponse{data=service.UserLoginResponse} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/login [post]
func (h UserHTTPHandler) Login(ctx *gin.Context) {
	request := entity.UserLogin{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.UserService.Login(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}
