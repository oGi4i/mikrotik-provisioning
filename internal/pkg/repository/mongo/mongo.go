package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"mikrotik_provisioning/internal/config"
)

const (
	NoDocumentsError = "mongo: no documents in result"
)

type Storage struct {
	collections map[string]*mongo.Collection
}

func NewMongoStorage(ctx context.Context, dbConfig *config.Database) (*Storage, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(dbConfig.DSN),
		options.Client().SetConnectTimeout(dbConfig.Timeout*time.Second),
		options.Client().SetDirect(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %q", err)
	}

	err = client.Ping(ctx, readpref.Nearest())
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %q", err)
	}

	collections := make(map[string]*mongo.Collection)
	for _, coll := range dbConfig.Collections {
		collections[coll.Resource] = client.Database(dbConfig.Name).Collection(coll.Name)

		if err := createIndexes(ctx, collections[coll.Resource], coll.Indexes); err != nil {
			return nil, err
		}
	}

	return &Storage{collections: collections}, nil
}
