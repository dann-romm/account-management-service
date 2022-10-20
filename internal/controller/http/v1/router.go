package v1

import (
	"account-management-service/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"os"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.String(200, "OK") })

	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newAccountRoutes(v1.Group("/accounts"), services.Account)
		newReservationRoutes(v1.Group("/reservations"), services.Reservation)
		newProductRoutes(v1.Group("/products"), services.Product)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
