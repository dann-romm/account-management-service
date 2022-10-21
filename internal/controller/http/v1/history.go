package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type historyRoutes struct {
	service.Operation
}

func newHistoryRoutes(g *echo.Group, operationService service.Operation) *historyRoutes {
	r := &historyRoutes{
		Operation: operationService,
	}

	g.GET("/history", r.getHistory)

	return r
}

type getHistoryInput struct {
	AccountId int    `json:"account_id" validate:"required"`
	SortType  string `json:"sort_type,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

func (r *historyRoutes) getHistory(c echo.Context) error {
	var input getHistoryInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	operations, err := r.Operation.OperationHistory(c.Request().Context(), service.OperationHistoryInput{
		AccountId: input.AccountId,
		SortType:  input.SortType,
		Offset:    input.Offset,
		Limit:     input.Limit,
	})
	if err != nil {
		log.Debugf("error while getting operation history: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"operations": operations,
	})
}
