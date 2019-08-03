package http

import (
	"github.com/go-chi/chi"
	"mikrotik_provisioning/handlers"
	"mikrotik_provisioning/pkg"
)

func SetRoutes(r *chi.Mux) {
	r.Route("/address-list", func(r chi.Router) {
		r.Get("/", handlers.ListAddressLists)                                                              // GET /address-list
		r.With(EnsureAuth).With(EnsureAddressListNotExists(pkg.API)).Post("/", handlers.CreateAddressList) // POST /address-list

		r.Route("/{addressListId:[a-f0-9]{24}}", func(r chi.Router) {
			r.With(EnsureAddressListExists(pkg.API)).Get("/", handlers.GetAddressList)                        // GET /address-list/123
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Put("/", handlers.UpdateAddressList)    // PUT /address-list/123
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Patch("/", handlers.PatchAddressList)   // PATCH /address-list/123
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Delete("/", handlers.DeleteAddressList) // DELETE /address-list/123
		})

		r.With(EnsureAddressListExists(pkg.API)).Get("/{addressListName:[A-Za-z0-9-]+}", handlers.GetAddressList)                        // GET /address-list/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Put("/{addressListName:[A-Za-z0-9-]+}", handlers.UpdateAddressList)    // PUT /address-list/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Patch("/{addressListName:[A-Za-z0-9-]+}", handlers.PatchAddressList)   // PATCH /address-list/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Delete("/{addressListName:[A-Za-z0-9-]+}", handlers.DeleteAddressList) // DELETE /address-list/whats-up
	})

	r.Route("/dns/static", func(r chi.Router) {
		r.Route("/list", func(r chi.Router) {
			r.Get("/", handlers.ListStaticDNSEntries)                                                                        // GET /dns/static/list
			r.With(EnsureAuth).With(EnsureStaticDNSEntriesNotExist(pkg.API)).Post("/", handlers.CreateBatchStaticDNSEntries) // POST /dns/static/list
			r.With(EnsureAuth).With(EnsureStaticDNSEntriesExist(pkg.API)).Put("/", handlers.UpdateBatchStaticDNSEntries)     // PUT /dns/static/list
		})
		// TODO: pattern conflicts with staticDNSName
		r.Route("/{staticDNSId:[a-f0-9]{24}}", func(r chi.Router) {
			r.With(EnsureStaticDNSEntryExists(pkg.API)).Get("/", handlers.GetStaticDNSEntry) // GET /dns/static/123
			/*r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Put("/", handlers.UpdateAddressList)    // PUT /dns/static/123
			r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Delete("/", handlers.DeleteAddressList)*/ // DELETE /dns/static/123
		})

		r.With(EnsureStaticDNSEntryExists(pkg.API)).Get("/{staticDNSName:[a-z0-9.-]+}", handlers.GetStaticDNSEntry) // GET /dns/static/whats-up
		/*r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Put("/{addressListName:[A-Za-z0-9-]+}", handlers.UpdateAddressList)    // PUT /dns/static/whats-up
		r.With(EnsureAuth).With(EnsureAddressListExists(pkg.API)).Delete("/{addressListName:[A-Za-z0-9-]+}", handlers.DeleteAddressList)*/ // DELETE /dns/static/whats-up
	})
}
