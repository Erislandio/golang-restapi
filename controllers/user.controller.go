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
		return setResponse(e, err, result.Status)
	}

	return e.JSON(http.StatusOK, result)

}

func setResponse(e echo.Context, err error, status int) error {
	return e.JSON(status, models.Response{
		Status:  status,
		Message: err.Error(),
		Data:    nil,
	})
}

func setResponseBad(e echo.Context, message string, status int) error {
	return e.JSON(status, models.Response{
		Status:  status,
		Message: message,
		Data:    nil,
	})
}

// StoreNewUser .
func StoreNewUser(e echo.Context) error {

	body := new(UserDTO)

	if err := e.Bind(body); err != nil {
		return setResponse(e, err, http.StatusInternalServerError)
	}

	result, err := models.StoreUser(body.Name, body.Email, body.Phone)

	if err != nil {
		return setResponse(e, err, result.Status)
	}

	return e.JSON(http.StatusOK, result)
}

// GetUserByID ;
func GetUserByID(e echo.Context) error {

	idString := e.Param("id")

	result, err := models.GetUserByID(idString)

	if err != nil {
		return setResponse(e, err, result.Status)
	}

	if result.Status != 200 {
		return setResponseBad(e, result.Message, result.Status)
	}

	return e.JSON(http.StatusOK, result)

}
