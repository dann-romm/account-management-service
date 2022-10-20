package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type accountRoutes struct {
	accountService service.Account
}

func newAccountRoutes(g *echo.Group, accountService service.Account) {
	r := &accountRoutes{
		accountService: accountService,
	}

	g.POST("/create", r.create)
	g.POST("/deposit", r.deposit) // POST, а не PUT, потому что неидемпотентно
	g.POST("/withdraw", r.withdraw)
	g.POST("/transfer", r.transfer)
}

func (r *accountRoutes) create(c echo.Context) error {
	id, err := r.accountService.CreateAccount(c.Request().Context())
	if err != nil {
		if err == service.ErrAccountAlreadyExists {
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

type accountDepositInput struct {
	Id     int `json:"id" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

func (r *accountRoutes) deposit(c echo.Context) error {
	var input accountDepositInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Deposit(c.Request().Context(), service.AccountDepositInput{
		Id:     input.Id,
		Amount: input.Amount,
	})
	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

type accountWithdrawInput struct {
	Id     int `json:"id" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

func (r *accountRoutes) withdraw(c echo.Context) error {
	var input accountWithdrawInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Withdraw(c.Request().Context(), service.AccountWithdrawInput{
		Id:     input.Id,
		Amount: input.Amount,
	})
	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

type accountTransferInput struct {
	From   int `json:"from" validate:"required"`
	To     int `json:"to" validate:"required"`
	Amount int `json:"amount" validate:"required"`
}

func (r *accountRoutes) transfer(c echo.Context) error {
	var input accountTransferInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.accountService.Transfer(c.Request().Context(), service.AccountTransferInput{
		From:   input.From,
		To:     input.To,
		Amount: input.Amount,
	})
	if err != nil {
		if err == service.ErrAccountNotFound {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
