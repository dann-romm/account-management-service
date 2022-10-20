package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type productRoutes struct {
	productService service.Product
}

func newProductRoutes(g *echo.Group, productService service.Product) *productRoutes {
	r := &productRoutes{
		productService: productService,
	}

	g.POST("/create", r.create)
	g.GET("/", r.getById)

	return r
}

type productCreateInput struct {
	Name string `json:"name" validate:"required"`
}

func (r *productRoutes) create(c echo.Context) error {
	var input productCreateInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.productService.CreateProduct(c.Request().Context(), input.Name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type getByIdInput struct {
	Id int `json:"id" validate:"required"`
}

func (r *productRoutes) getById(c echo.Context) error {
	var input getByIdInput

	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	product, err := r.productService.GetProductById(c.Request().Context(), input.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"product": product,
	})
}
