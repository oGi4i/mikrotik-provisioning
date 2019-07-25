package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func setRoutes(r *chi.Mux) {
	r.Route("/address-list", func(r chi.Router) {
		r.Get("/", ListAddressLists)                                         // GET /address-list
		r.With(EnsureAddressListNotExists(api)).Post("/", CreateAddressList) // POST /address-list

		r.Route("/{addressListId:[a-f0-9]+}", func(r chi.Router) {
			r.Use(EnsureAddressListExists(api)) // Load the *AddressList on the request context
			r.Get("/", GetAddressList)          // GET /address-list/123
			r.Put("/", UpdateAddressList)       // PUT /address-list/123
			r.Patch("/", PatchAddressList)      // PATCH /address-list/123
			r.Delete("/", DeleteAddressList)    // DELETE /address-list/123
		})

		r.With(EnsureAddressListExists(api)).Get("/{addressListName:[A-Za-z-]+}", GetAddressList)       // GET /address-list/whats-up
		r.With(EnsureAddressListExists(api)).Put("/{addressListName:[A-Za-z-]+}", UpdateAddressList)    // PUT /address-list/whats-up
		r.With(EnsureAddressListExists(api)).Patch("/{addressListName:[A-Za-z-]+}", PatchAddressList)   // PATCH /address-list/whats-up
		r.With(EnsureAddressListExists(api)).Delete("/{addressListName:[A-Za-z-]+}", DeleteAddressList) // DELETE /address-list/whats-up
	})

	// Mount the admin sub-router, which btw is the same as:
	// r.Route("/admin", func(r chi.Router) { admin routes here })
	r.Mount("/admin", adminRouter())
}

func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: list accounts.."))
	})
	r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
	})
	return r
}
