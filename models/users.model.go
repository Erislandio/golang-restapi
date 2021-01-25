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

	return response, nil
}

// FindByEmail ;
func FindByEmail(email string) (Response, error) {

	var users []User
	var response Response

	conn := database.CreateMysqlConn()

	sqlStatment := "SELECT * FROM users WHERE email = ?"

	rows, err := conn.Query(sqlStatment, email)

	if err != nil {
		response.Data = nil
		return response, err
	}

	for rows.Next() {
		var user User

		err := rows.Scan(user.ID, user.Name, user.Email, user.Phone)

		if err != nil {
			response.Data = nil
			return response, err
		}

		users = append(users, user)
	}

	response.Status = http.StatusOK
	response.Message = "Success"
	response.Data = users

	return response, nil
}

// StoreUser ;
func StoreUser(name, email, phone string) (Response, error) {
	var response Response

	conn := database.CreateMysqlConn()

	validation, err := FindByEmail(email)

	fmt.Println(validation.Data)

	if validation.Data != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("User with email: %s already exists!", email)
		response.Data = nil
		return response, nil
	}

	sqlStatement := "INSERT users (name, email, phone) VALUES (?, ?, ?)"

	stmt, err := conn.Prepare(sqlStatement)

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

	return response, nil

}
