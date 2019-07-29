package main

import (
	"github.com/go-chi/chi"
)

func setRoutes(r *chi.Mux) {
	r.Route("/address-list", func(r chi.Router) {
		r.Get("/", ListAddressLists)                                                          // GET /address-list
		r.With(EnsureAuth).With(EnsureAddressListNotExists(api)).Post("/", CreateAddressList) // POST /address-list

		r.Route("/{addressListId:[a-f0-9]{24}}", func(r chi.Router) {
			r.With(EnsureAddressListExists(api)).Get("/", GetAddressList)                        // GET /address-list/123
			r.With(EnsureAuth).With(EnsureAddressListExists(api)).Put("/", UpdateAddressList)    // PUT /address-list/123
			r.With(EnsureAuth).With(EnsureAddressListExists(api)).Patch("/", PatchAddressList)   // PATCH /address-list/123
			r.With(EnsureAuth).With(EnsureAddressListExists(api)).Delete("/", DeleteAddressList) // DELETE /address-list/123
		})

		r.With(EnsureAddressListExists(api)).Get("/{addressListName:[A-Za-z0-9-]+}", GetAddressList)                        // GET /address-list/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(api)).Put("/{addressListName:[A-Za-z0-9-]+}", UpdateAddressList)    // PUT /address-list/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(api)).Patch("/{addressListName:[A-Za-z0-9-]+}", PatchAddressList)   // PATCH /address-list/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(api)).Delete("/{addressListName:[A-Za-z0-9-]+}", DeleteAddressList) // DELETE /address-list/whats-up
	})
}
