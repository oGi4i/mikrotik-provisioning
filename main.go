package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

const (
	configFile = "config.yml"
)

var (
	addressLists = []*AddressList{}
	users        = []*User{}
	cfg          = &YamlConfig{}
	api          = &Implementation{}
)

func init() {
	err := cfg.initConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config with error: %q\n", err)
	}
}

func main() {
	ctx := context.Background()
	mongoDB, coll := NewDB(ctx)
	storage := NewMongoStorage(mongoDB, coll)
	api = NewMikrotikAclAPI(storage)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(CheckAcceptHeader("*/*", "application/json", "text/plain"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	setRoutes(r)

	http.ListenAndServe(":3333", r)
}

func NewDB(ctx context.Context) (*mongo.Client, *mongo.Collection) {
	ctx, _ = context.WithTimeout(context.Background(), cfg.Database.Timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.DSN))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %q\n", err)
	}

	ctx, _ = context.WithTimeout(context.Background(), cfg.Database.Timeout*time.Second*5)
	err = client.Ping(ctx, readpref.Nearest())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %q\n", err)
	}

	collection := client.Database(cfg.Database.Name).Collection(cfg.Database.Collection)

	return client, collection
}
