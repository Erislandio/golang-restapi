package models

import (
	"fmt"
	"net/http"

	"github.com/erislandio/web/restapi/database"
)

//  User struc type user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// FetchAllUsers for users list
func FetchAllUsers() (Response, error) {

	var obj User
	var users []User
	var response Response

	conn := database.CreateMysqlConn()

	sqlStatement := "SELECT * FROM `users`"

	rows, err := conn.Query(sqlStatement)

	defer rows.Close()

	if err != nil {
		return response, err
	}

	for rows.Next() {
		err = rows.Scan(&obj.ID, &obj.Name, &obj.Email, &obj.Phone)

		if err != nil {
			return response, err
		}

		users = append(users, obj)
	}

	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = users
	response.Count = len(users)

	return response, nil
}

// FindByEmail ;
func FindByEmail(email string) (Response, error) {

	var users []User
	var response Response

	conn := database.CreateMysqlConn()

	sqlStatment := "SELECT * FROM users WHERE email = ?"

	rows, err := conn.Query(sqlStatment, email)

	defer rows.Close()

	if err != nil {
		response.Data = nil
		response.Message = err.Error()
		response.Status = http.StatusInternalServerError
		return response, err
	}

	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)

		if err != nil {
			response.Data = nil
			response.Message = err.Error()
			response.Status = http.StatusInternalServerError
			return response, err
		}

		users = append(users, user)
	}

	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = users
	response.Count = len(users)

	return response, nil
}

// StoreUser ;
func StoreUser(name, email, phone string) (Response, error) {
	var response Response

	conn := database.CreateMysqlConn()

	validation, err := FindByEmail(email)

	fmt.Println(validation.Data)

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
	}

	if validation.Count > 0 {
		var message = fmt.Sprintf("User with email: %s already exists!", email)
		setResponse(&response, message, http.StatusBadRequest, nil)
		return response, nil
	}

	sqlStatement := "INSERT users (name, email, phone) VALUES (?, ?, ?)"

	stmt, err := conn.Prepare(sqlStatement)

	defer stmt.Close()

	if err != nil {
		return response, err
	}

	result, err := stmt.Exec(name, email, phone)

	if err != nil {
		return response, err
	}

	lastInsertId, err := result.LastInsertId()

	if err != nil {
		return response, err
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = map[string]int64{
		"lastId": lastInsertId,
	}
	response.Count = int(lastInsertId)

	return response, nil

}

func setResponseSuccess(response *Response, message string, count, statusCode int, data interface{}) {
	response.Status = statusCode
	response.Message = message
	response.Data = data
	response.Count = int(count)
}

// GetUserByID func get user by id
func GetUserByID(id string) (Response, error) {

	var response Response
	var users []User

	conn := database.CreateMysqlConn()

	const sqlStatement = "SELECT * FROM users WHERE id = ?"

	rows, err := conn.Query(sqlStatement, id)

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	defer rows.Close()

	for rows.Next() {

		var user User

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone)

		if err != nil {
			setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
			return response, err
		}

		users = append(users, user)
	}

	count := len(users)

	if count > 0 {
		setResponseSuccess(&response, "Success", count, http.StatusOK, users)
	} else {
		message := fmt.Sprintf("User with id: %s not found!", id)
		setResponse(&response, message, http.StatusBadRequest, nil)
	}

	return response, nil
}

func setResponse(response *Response, message string, status int, data interface{}) {
	response.Status = status
	response.Message = message
	response.Data = data
}
