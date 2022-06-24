package models

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type Role int

const (
	Admin Role = iota
	Boss
	Regular
)

type Employee struct {
	Role     Role
	UUID     string
	Username string
	Password string
}

func NewEmployee(username, password string, role Role) Employee {
	if role == Admin {
		return Employee{UUID: "admin", Username: username, Password: password, Role: role}
	}
	return Employee{UUID: randSeq(8), Username: username, Password: password, Role: role}
}
