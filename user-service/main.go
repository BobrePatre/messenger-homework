package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {

	type Status struct {
		Status     string `json:"status"`
		InstanceId string `json:"instanceId"`
		Service    string `json:"service"`
	}
	instance := uuid.NewString()

	router := echo.New()
	router.Use(middleware.Recover())
	router.Use(middleware.Logger())

	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Status{
			Status:     "OK, im health probe",
			InstanceId: instance,
			Service:    "user-service",
		})
	})

	router.GET("/readyz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Status{
			Status:     "OK, im readyz probe",
			InstanceId: instance,
			Service:    "user-service",
		})
	})

	router.GET("/foo", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Status{
			Status:     "OK, im foo probe",
			InstanceId: instance,
			Service:    "user-service",
		})
	})

	router.Start(":8080")
}
