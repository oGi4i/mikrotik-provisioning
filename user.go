package main

import "net/http"

type User struct {
	Name     string `json:"name",yaml:"name",bson:"name"`
	Password string `json:"password",yaml:"password",bson:"password"`
}

type UserPayload struct {
	*User
	Role string `json:"role" bson:"role"`
}

func NewUserPayloadResponse(user *User) *UserPayload {
	return &UserPayload{User: user}
}

func (u *UserPayload) Bind(r *http.Request) error {
	return nil
}

func (u *UserPayload) Render(w http.ResponseWriter, r *http.Request) error {
	u.Role = "collaborator"
	return nil
}
