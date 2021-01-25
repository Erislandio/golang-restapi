package routes

import (
	"net/http"

	"github.com/erislandio/web/restapi/controllers"
	"github.com/labstack/echo"
)

// Init .
func Init() *echo.Echo {

	app := echo.New()

	app.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	app.GET("/api/v1/users", controllers.GetAllUsers)
	app.POST("/api/v1/users", controllers.StoreNewUser)
	app.GET("/api/v1/users/:id", controllers.GetUserByID)
	app.PATCH("/api/v1/users", controllers.UpdateUserInfo)
	app.DELETE("/api/v1/users/:id", controllers.DeleteUser)

	return app

}
