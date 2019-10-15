package http

import (
	"github.com/go-chi/chi"
	"mikrotik_provisioning/pkg"
)

func SetRoutes(r *chi.Mux) {
	r.Route("/address-list", func(r chi.Router) {
		r.Get("/", ListAddressLists)                                                              // GET /address-list
		r.With(EnsureAuth).With(EnsureAddressListNotExists(pkg.API)).Post("/", CreateAddressList) // POST /address-list

		r.Route("/{addressListName:[A-Za-z0-9-]+}", func(r chi.Router) {
			r.With(EnsureAddressListExists(pkg.API)).Get("/", GetAddressList)                        // GET /address-list/whats-up
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Put("/", UpdateAddressList)    // PUT /address-list/whats-up
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Patch("/", PatchAddressList)   // PATCH /address-list/whats-up
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Delete("/", DeleteAddressList) // DELETE /address-list/whats-up
		})
	})

	r.Route("/dns/static", func(r chi.Router) {
		r.With(EnsureAuth).With(EnsureStaticDNSEntryNotExists(pkg.API)).Post("/", CreateStaticDNSEntry)
		r.Route("/list", func(r chi.Router) {
			r.Get("/", ListStaticDNSEntries)                                                                        // GET /dns/static/list
			r.With(EnsureAuth).With(EnsureStaticDNSEntriesNotExist(pkg.API)).Post("/", CreateBatchStaticDNSEntries) // POST /dns/static/list
			r.With(EnsureAuth).With(EnsureStaticDNSEntriesExist(pkg.API)).Put("/", UpdateBatchStaticDNSEntries)     // PUT /dns/static/list
		})

		r.Route("/{staticDNSName:[a-z0-9.-]+}", func(r chi.Router) {
			r.With(EnsureStaticDNSEntryExists(pkg.API)).Get("/", GetStaticDNSEntry)                        // GET /dns/static/whats-up
			r.With(EnsureAuth).With(EnsureStaticDNSEntryExists(pkg.API)).Put("/", UpdateStaticDNSEntry)    // PUT /dns/static/whats-up
			r.With(EnsureAuth).With(EnsureStaticDNSEntryExists(pkg.API)).Delete("/", DeleteStaticDNSEntry) // DELETE /dns/static/whats-up
		})
	})
}
