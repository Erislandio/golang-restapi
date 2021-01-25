package controllers

import (
	"net/http"

	"github.com/erislandio/web/restapi/models"
	"github.com/labstack/echo"
)

// UserDTO ;
type UserDTO struct {
	Name  string `form:"name"`
	Email string `form:"email"`
	Phone string `form:"phone"`
}

// GetAllUsers ;
func GetAllUsers(e echo.Context) error {
	result, err := models.FetchAllUsers()

	if err != nil {
		return statusInternalServerError(e, err)
	}

	return e.JSON(http.StatusInternalServerError, result)

}

func statusInternalServerError(e echo.Context, err error) error {
	return e.JSON(http.StatusInternalServerError, models.Response{
		Status:  500,
		Message: err.Error(),
		Data:    nil,
	})
}

// StoreNewUser .
func StoreNewUser(e echo.Context) error {

	body := new(UserDTO)

	if err := e.Bind(body); err != nil {
		return err
	}

	result, err := models.StoreUser(body.Name, body.Email, body.Phone)

	if err != nil {
		return statusInternalServerError(e, err)
	}

	return e.JSON(http.StatusOK, result)

}
