package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type reservationRoutes struct {
	reservationService service.Reservation
}

func newReservationRoutes(g *echo.Group, reservationService service.Reservation) {
	r := &reservationRoutes{
		reservationService: reservationService,
	}

	g.POST("/create", r.create)
	g.POST("/revenue", r.revenue)
	g.POST("/refund", r.refund)
}

type reservationCreateInput struct {
	AccountId int `json:"account_id" validate:"required"`
	ProductId int `json:"product_id" validate:"required"`
	OrderId   int `json:"order_id" validate:"required"`
	Amount    int `json:"amount" validate:"required"`
}

func (r *reservationRoutes) create(c echo.Context) error {
	var input reservationCreateInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.reservationService.CreateReservation(c.Request().Context(), service.ReservationCreateInput{
		AccountId: input.AccountId,
		ProductId: input.ProductId,
		OrderId:   input.OrderId,
		Amount:    input.Amount,
	})
	if err != nil {
		if err == service.ErrCannotCreateReservation {
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

type reservationRevenueInput struct {
	AccountId int `json:"account_id" validate:"required"`
	ProductId int `json:"product_id" validate:"required"`
	OrderId   int `json:"order_id" validate:"required"`
	Amount    int `json:"amount" validate:"required"`
}

func (r *reservationRoutes) revenue(c echo.Context) error {
	var input reservationRevenueInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.reservationService.RevenueReservationByOrderId(c.Request().Context(), input.OrderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

type reservationRefundInput struct {
	OrderId int `json:"order_id" validate:"required"`
}

func (r *reservationRoutes) refund(c echo.Context) error {
	var input reservationRefundInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.reservationService.RefundReservationByOrderId(c.Request().Context(), input.OrderId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
