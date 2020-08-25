package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mikrotik_provisioning/internal/config"
)

type Index struct {
	Version    int32            `bson:"v"`
	Key        map[string]int32 `bson:"key"`
	Name       string           `bson:"name"`
	Namespace  string           `bson:"ns"`
	Unique     bool             `bson:"unique"`
	Background bool             `bson:"background"`
}

func createIndexes(ctx context.Context, collection *mongo.Collection, indexList []*config.CollectionIndexes) error {
	cur, err := collection.Indexes().List(ctx)
	defer cur.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to get list of indexes for collection: %s from mongodb with error: %q", collection.Name(), err)
	}

	indexes := make([]string, 0)
	for cur.Next(ctx) {
		index := new(Index)
		err := cur.Decode(&index)

		if err != nil {
			return fmt.Errorf("failed to decode index model for collection: %s from mongodb with error: %q", collection.Name(), err)
		}

		indexes = append(indexes, index.Name)
	}

	for _, index := range indexList {
		var exists bool
		for _, i := range indexes {
			if index.Name == i {
				exists = true
				break
			}
		}

		if exists {
			continue
		}

		res, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: index.Field, Value: 1}},
			Options: options.Index().SetBackground(true).SetName(index.Name).SetUnique(index.Unique),
		})

		if err != nil || res != index.Name {
			return fmt.Errorf("failed to create index for collection: %s in mongodb with error: %q", index.Name, err)
		}
	}

	return nil
}
