package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"log"
	cfg "mikrotik_provisioning/config"
	mux "mikrotik_provisioning/http"
	"mikrotik_provisioning/pkg"
	store "mikrotik_provisioning/storage"
	valid "mikrotik_provisioning/validate"
	"net/http"
)

func init() {
	if err := valid.RegisterValidators(valid.Validate); err != nil {
		log.Fatalf("failed to register custom validation functions with error: %q", err)
	}

	if err := cfg.Config.InitConfig(); err != nil {
		log.Fatalf("failed to initialize config with error: %q", err)
	}

	ctx := context.Background()
	storage, err := store.NewMongoStorage(ctx)
	if err != nil {
		log.Fatalf("failed to initialize MongoStorage with error: %q", err)
	}
	pkg.API = pkg.NewMikrotikAclAPI(storage)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(mux.CheckAcceptHeader("*/*", "application/json", "text/plain"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	mux.SetRoutes(r)

	http.ListenAndServe(":3333", r)
}
