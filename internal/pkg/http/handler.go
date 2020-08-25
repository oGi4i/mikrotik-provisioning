package http

import (
	"net/http"
	"text/template"

	"mikrotik_provisioning/internal/app"
)

type Middleware interface {
	EnsureAddressListExists(next http.Handler) http.Handler
	EnsureAddressListNotExists(next http.Handler) http.Handler
	EnsureAuth(next http.Handler) http.Handler
	CheckAcceptHeader(contentTypes ...string) func(next http.Handler) http.Handler
}

type Handler interface {
	GetAddressLists(w http.ResponseWriter, r *http.Request)
	CreateAddressList(w http.ResponseWriter, r *http.Request)
	GetAddressList(w http.ResponseWriter, r *http.Request)
	UpdateAddressList(w http.ResponseWriter, r *http.Request)
	PatchAddressList(w http.ResponseWriter, r *http.Request)
	DeleteAddressList(w http.ResponseWriter, r *http.Request)
}

type AddressListHandler struct {
	service   app.UseCases
	templates *template.Template
}

func NewAddressListHandler(service app.UseCases, templates *template.Template) *AddressListHandler {
	return &AddressListHandler{service: service, templates: templates}
}
