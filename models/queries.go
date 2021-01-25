package models

const getAll = "SELECT * FROM `users`"
const getByEmail = "SELECT * FROM users WHERE email = ?"
const insertIntoUsers = "INSERT users (name, email, phone) VALUES (?, ?, ?)"
const getByid = "SELECT * FROM users WHERE id = ?"
const updateByID = "UPDATE users set name = ?, phone = ? WHERE id = ?"
const deleteUser = "DELETE FROM users WHERE id = ?"
