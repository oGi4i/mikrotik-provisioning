package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"mikrotik_provisioning/internal/app"
	"mikrotik_provisioning/internal/config"
	mux "mikrotik_provisioning/internal/pkg/http"
	mw "mikrotik_provisioning/internal/pkg/http/middleware"
	"mikrotik_provisioning/internal/pkg/repository/mongo"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("failed to parse config file with error: %q\n", err)
	}

	templateFiles, err := filepath.Glob("templates/*")
	if err != nil {
		log.Fatalf("failed to get template filenames with glob with error: %q\n", err)
	}

	templates := new(template.Template)
	templates, err = templates.Delims("#(", ")#").ParseFiles(templateFiles...)
	if err != nil {
		log.Fatalf("failed to parse template files with error: %q\n", err)
	}

	ctx := context.Background()
	mongoStore, err := mongo.NewMongoStorage(ctx, config.DB)
	if err != nil {
		log.Fatalf("failed to initialize MongoStore with error: %q\n", err)
	}

	service := app.NewMikrotikProvisioningService(mongoStore, templates)
	mw := mw.NewMiddleware(service, config.Access)
	handler := mux.NewAddressListHandler(service, templates)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(mw.CheckAcceptHeader("*/*", "application/json", "text/plain"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/address-list", func(r chi.Router) {
		r.Get("/", handler.GetAddressLists)                                                            // GET /address-list
		r.With(mw.EnsureAuth).With(mw.EnsureAddressListNotExists).Post("/", handler.CreateAddressList) // POST /address-list

		r.Route("/{addressListName:[A-Za-z0-9-]+}", func(r chi.Router) {
			r.With(mw.EnsureAddressListExists).Get("/", handler.GetAddressList)                           // GET /address-list/whats-up
			r.With(mw.EnsureAuth).With(mw.EnsureAddressListExists).Put("/", handler.UpdateAddressList)    // PUT /address-list/whats-up
			r.With(mw.EnsureAuth).With(mw.EnsureAddressListExists).Patch("/", handler.PatchAddressList)   // PATCH /address-list/whats-up
			r.With(mw.EnsureAuth).With(mw.EnsureAddressListExists).Delete("/", handler.DeleteAddressList) // DELETE /address-list/whats-up
		})
	})

	err = http.ListenAndServe(":3333", r)
	if err != nil {
		log.Fatalf("failed to initialize http server with error: %q\n", err)
	}
}
