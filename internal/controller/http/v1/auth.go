package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}

	g.POST("/sign-up", r.signUp)
	g.POST("/sign-in", r.signIn)
}

type signUpInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// регистрация пользователя
func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// аутентификация пользователя
func (r *authRoutes) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusBadRequest, "invalid username or password")
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
