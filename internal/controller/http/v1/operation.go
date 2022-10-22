package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type operationRoutes struct {
	service.Operation
}

func newOperationRoutes(g *echo.Group, operationService service.Operation) *operationRoutes {
	r := &operationRoutes{
		Operation: operationService,
	}

	g.GET("/history", r.getHistory)
	g.GET("/report", r.getReport)

	return r
}

type getHistoryInput struct {
	AccountId int    `json:"account_id" validate:"required"`
	SortType  string `json:"sort_type,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

func (r *operationRoutes) getHistory(c echo.Context) error {
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

type getReportInput struct {
	Month int `json:"month" validate:"required"`
	Year  int `json:"year" validate:"required"`
}

func (r *operationRoutes) getReport(c echo.Context) error {
	var input getReportInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	link, err := r.Operation.MakeReportLink(c.Request().Context(), input.Month, input.Year)
	if err != nil {
		log.Debugf("error while getting report link: %s", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"link": link,
	})
}
