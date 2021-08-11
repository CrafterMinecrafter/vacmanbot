package models

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	IsIgnored bool   `json:"is_ignored"`
}
