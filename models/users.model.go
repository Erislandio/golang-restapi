package models

import (
	"context"
	"encoding/json"
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

var (
	ctx = context.Background()
)

// FetchAllUsers for users list
func FetchAllUsers() (Response, error) {

	var obj User
	var users []User
	var response Response

	conn := database.CreateMysqlConn()

	users = getOnRedis("users")

	if len(users) > 0 {
		setResponseSuccess(&response, "Success from redis", len(users), http.StatusOK, users)
		return response, nil
	}

	rows, err := conn.Query(getAll)

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

	saveOnRedis("users", users)
	return response, nil
}

func marshalBinary(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func getOnRedis(key string) []User {
	var users []User
	redis := database.CreateRedisConn()

	usersFromRedis, err := redis.Get("users").Result()

	if err != nil {
		fmt.Print(err.Error())
	}

	json.Unmarshal([]byte(usersFromRedis), &users)

	return users
}

func saveOnRedis(key string, data interface{}) {
	redis := database.CreateRedisConn()

	toRedisData, _ := marshalBinary(data)

	_, err := redis.Set("users", toRedisData, 0).Result()

	if err != nil {
		fmt.Print(err.Error())
	}

}

// FindByEmail ;
func FindByEmail(email string) (Response, error) {

	var users []User
	var response Response

	conn := database.CreateMysqlConn()

	rows, err := conn.Query(getByEmail, email)

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

	stmt, err := conn.Prepare(insertIntoUsers)

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

	newUsers, _ := FetchAllUsers()
	saveOnRedis("users", newUsers.Data)

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

	rows, err := conn.Query(getByid, id)

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

// UpdateById func data
func UpdateById(name, phone, id string) (Response, error) {

	var response Response

	conn := database.CreateMysqlConn()

	updateStatement, err := conn.Prepare(updateByID)

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	result, err := updateStatement.Exec(name, phone, id)

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	final, err := result.RowsAffected()

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	if final > 0 {
		setResponseSuccess(&response, "Success", int(final), http.StatusOK, final)
	} else {
		message := fmt.Sprintf("User with id: %s not found!", id)
		setResponse(&response, message, http.StatusBadRequest, nil)
	}

	newUsers, _ := FetchAllUsers()
	saveOnRedis("users", newUsers.Data)

	return response, nil

}

// DeleteByID .
func DeleteByID(id string) (Response, error) {
	var response Response

	conn := database.CreateMysqlConn()

	deleteStatement, err := conn.Prepare(deleteUser)

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	result, err := deleteStatement.Exec(id)

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	final, err := result.RowsAffected()

	if err != nil {
		setResponse(&response, err.Error(), http.StatusInternalServerError, nil)
		return response, err
	}

	if final > 0 {
		setResponseSuccess(&response, "Success", int(final), http.StatusOK, final)
	} else {
		message := fmt.Sprintf("User with id: %s not found!", id)
		setResponse(&response, message, http.StatusBadRequest, nil)
	}

	newUsers, _ := FetchAllUsers()
	saveOnRedis("users", newUsers.Data)

	return response, nil

}

func setResponse(response *Response, message string, status int, data interface{}) {
	response.Status = status
	response.Message = message
	response.Data = data
}
