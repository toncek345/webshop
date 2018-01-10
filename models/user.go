package models

// only 1 user will be admin, no login on this simple webshop

type User struct {
	Id       int
	Username string
	Password string // TODO or []byte
}
