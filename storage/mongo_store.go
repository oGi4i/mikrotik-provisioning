package storage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	cfg "mikrotik_provisioning/config"
	"mikrotik_provisioning/types"
	valid "mikrotik_provisioning/validate"
	"time"
)

type MongoIndex struct {
	Version    int32            `bson:"v"`
	Key        map[string]int32 `bson:"key"`
	Name       string           `bson:"name"`
	Namespace  string           `bson:"ns"`
	Unique     bool             `bson:"unique"`
	Background bool             `bson:"background"`
}

type MongoStorage struct {
	cl    *mongo.Client
	colls map[string]*mongo.Collection
}

func NewMongoStorage(ctx context.Context) (*MongoStorage, error) {
	ctx, _ = context.WithTimeout(context.Background(), cfg.Config.Database.Timeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Config.Database.DSN))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to connect to MongoDB: %q", err))
	}

	ctx, _ = context.WithTimeout(context.Background(), cfg.Config.Database.Timeout*time.Second*5)
	err = client.Ping(ctx, readpref.Nearest())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to ping MongoDB: %q", err))
	}

	colls := make(map[string]*mongo.Collection)
	for _, coll := range cfg.Config.Database.Collections {
		colls[coll.Resource] = client.Database(cfg.Config.Database.Name).Collection(coll.Name)

		if err := createIndexes(ctx, colls[coll.Resource], coll.Indexes); err != nil {
			return nil, err
		}
	}

	return &MongoStorage{cl: client, colls: colls}, nil
}

func createIndexes(ctx context.Context, coll *mongo.Collection, indexList []cfg.CollectionIndexesConfig) error {
	cur, err := coll.Indexes().List(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to get list of indexes for collection: %s from MongoDB with error: %q", coll.Name, err))
	}
	defer cur.Close(ctx)

	indexes := make([]string, 0)
	for cur.Next(ctx) {
		index := new(MongoIndex)
		err := cur.Decode(&index)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to decode index model for collection: %s from MongoDB with error: %q", coll.Name, err))
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
			break
		}

		res, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{index.Field, 1}},
			Options: options.Index().SetBackground(true).SetName(index.Name).SetUnique(index.Unique),
		})
		if err != nil || res != index.Name {
			return errors.New(fmt.Sprintf("failed to create index for collection: %s in MongoDB with error: %q", index.Name, err))
		}
	}

	return nil
}

func (s *MongoStorage) NewAddressList(ctx context.Context, addressList *types.AddressList) (*types.AddressList, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res, err := s.colls["address-list"].InsertOne(c, &types.AddressListMongo{
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	})
	if err != nil {
		return nil, err
	}
	addressList.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return addressList, nil
}

func (s *MongoStorage) GetAllAddressLists(ctx context.Context) ([]*types.AddressList, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	cur, err := s.colls["address-list"].Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(c)

	result := make([]*types.AddressList, 0)
	for cur.Next(ctx) {
		data := new(types.AddressListMongo)

		err := cur.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		result = append(result, &types.AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		})
	}

	return result, nil
}

func (s *MongoStorage) GetAddressListById(ctx context.Context, id string) (*types.AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res := s.colls["address-list"].FindOne(c, bson.M{"_id": objectID})

	if res.Err() == nil {
		if b, err := res.DecodeBytes(); err != nil {
			if b.String() == "" && err.Error() == "mongo: no documents in result" {
				return nil, nil
			} else {
				return nil, err
			}
		}

		data := new(types.AddressListMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &types.AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) GetAddressListByName(ctx context.Context, name string) (*types.AddressList, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res := s.colls["address-list"].FindOne(c, bson.M{"name": name})

	if res.Err() == nil {
		if b, err := res.DecodeBytes(); err != nil {
			if b.String() == "" && err.Error() == "mongo: no documents in result" {
				return nil, nil
			} else {
				return nil, err
			}
		}

		data := new(types.AddressListMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &types.AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) UpdateAddressListById(ctx context.Context, id string, addressList *types.AddressList) (*types.AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.colls["address-list"].FindOneAndReplace(ctx, bson.M{"_id": objectID}, &types.AddressListMongo{
		ID:        objectID,
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	}, options.FindOneAndReplace().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, err
	}

	data := new(types.AddressListMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	return &types.AddressList{
		ID:        data.ID.Hex(),
		Name:      data.Name,
		Addresses: data.Addresses,
	}, nil
}

func (s *MongoStorage) AddAddressesToAddressListById(ctx context.Context, id string, addresses []types.Address) (*types.AddressList, error) {
	currentData, err := s.GetAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	bsonA := primitive.A{}
	ok := true
	for _, a := range addresses {
		for _, b := range currentData.Addresses {
			if b == a {
				ok = false
				break
			}
		}
		if ok {
			bsonA = append(bsonA, a)
		}
	}

	update := bson.M{"$push": bson.M{"addresses": bson.M{"$each": bsonA}}}

	res := s.colls["address-list"].FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if res.Err() != nil {
		return nil, err
	}

	data := new(types.AddressListMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	newData, err := s.GetAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(newData); err != nil {
		return nil, err
	}

	return newData, nil
}

func (s *MongoStorage) RemoveAddressListById(ctx context.Context, id string) (*types.AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid addressListId")
	}

	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res, err := s.colls["address-list"].DeleteOne(c, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("failed deleting mongoDB object")
	}

	return nil, nil
}

func (s *MongoStorage) RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []types.Address) (*types.AddressList, error) {
	currentData, err := s.GetAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(currentData); err != nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	bsonA := primitive.A{}
	ok := false
	for _, a := range addresses {
		for _, b := range currentData.Addresses {
			if b == a {
				ok = true
				break
			}
		}
		if ok {
			bsonA = append(bsonA, a)
		}
	}

	update := bson.M{"$pull": bson.M{"addresses": bson.M{"$in": bsonA}}}

	res := s.colls["address-list"].FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if res.Err() != nil {
		return nil, err
	}

	data := new(types.AddressListMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	newData, err := s.GetAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(newData); err != nil {
		return nil, err
	}

	return newData, nil
}

func (s *MongoStorage) GetAllStaticDNS(ctx context.Context) ([]*types.StaticDNSEntry, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	cur, err := s.colls["static-dns"].Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(c)

	result := make([]*types.StaticDNSEntry, 0)
	for cur.Next(ctx) {
		data := new(types.StaticDNSEntryMongo)

		err := cur.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		result = append(result, &types.StaticDNSEntry{
			ID:       data.ID.Hex(),
			Name:     data.Name,
			Regexp:   data.Regexp,
			Address:  data.Address,
			TTL:      data.TTL,
			Disabled: data.Disabled,
			Comment:  data.Comment,
		})
	}

	return result, nil
}

func (s *MongoStorage) NewStaticDNSBatch(ctx context.Context, entries []*types.StaticDNSEntry) ([]*types.StaticDNSEntry, error) {
	mongoEntries := make([]interface{}, len(entries))
	for i, entry := range entries {
		mongoEntries[i] = &types.StaticDNSEntryMongo{
			Name:     entry.Name,
			Regexp:   entry.Regexp,
			Address:  entry.Address,
			TTL:      entry.TTL,
			Disabled: entry.Disabled,
			Comment:  entry.Comment,
		}
	}

	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res, err := s.colls["static-dns"].InsertMany(c, mongoEntries)
	if err != nil {
		return nil, err
	}

	for i, id := range res.InsertedIDs {
		entries[i].ID = id.(primitive.ObjectID).Hex()
	}

	return entries, nil
}

func (s *MongoStorage) UpdateStaticDNSBatch(ctx context.Context, entries []*types.StaticDNSEntry) ([]*types.StaticDNSEntry, error) {
	mongoEntries := make([]*types.StaticDNSEntryMongo, len(entries))
	for i, entry := range entries {
		mongoEntries[i] = &types.StaticDNSEntryMongo{
			Name:     entry.Name,
			Regexp:   entry.Regexp,
			Address:  entry.Address,
			TTL:      entry.TTL,
			Disabled: entry.Disabled,
			Comment:  entry.Comment,
		}
	}

	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	for i, entry := range mongoEntries {
		res := s.colls["static-dns"].FindOneAndReplace(c, bson.M{"name": entry.Name}, entry, options.FindOneAndReplace().SetReturnDocument(options.After))
		if res.Err() != nil {
			return nil, res.Err()
		}

		data := new(types.StaticDNSEntryMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		entries[i] = &types.StaticDNSEntry{
			ID:       data.ID.Hex(),
			Name:     data.Name,
			Regexp:   data.Regexp,
			Address:  data.Address,
			TTL:      data.TTL,
			Disabled: data.Disabled,
			Comment:  data.Comment,
		}
	}

	return entries, nil
}

func (s *MongoStorage) GetStaticDNSEntryById(ctx context.Context, id string) (*types.StaticDNSEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res := s.colls["static-dns"].FindOne(c, bson.M{"_id": objectID})

	if res.Err() == nil {
		if b, err := res.DecodeBytes(); err != nil {
			if b.String() == "" && err.Error() == "mongo: no documents in result" {
				return nil, nil
			} else {
				return nil, err
			}
		}

		data := new(types.StaticDNSEntryMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &types.StaticDNSEntry{
			ID:       data.ID.Hex(),
			Name:     data.Name,
			Regexp:   data.Regexp,
			Address:  data.Address,
			TTL:      data.TTL,
			Disabled: data.Disabled,
			Comment:  data.Comment,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) GetStaticDNSEntryByName(ctx context.Context, name string) (*types.StaticDNSEntry, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res := s.colls["static-dns"].FindOne(c, bson.M{"name": name})

	if res.Err() == nil {
		if b, err := res.DecodeBytes(); err != nil {
			if b.String() == "" && err.Error() == "mongo: no documents in result" {
				return nil, nil
			} else {
				return nil, err
			}
		}

		data := new(types.StaticDNSEntryMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &types.StaticDNSEntry{
			ID:       data.ID.Hex(),
			Name:     data.Name,
			Regexp:   data.Regexp,
			Address:  data.Address,
			TTL:      data.TTL,
			Disabled: data.Disabled,
			Comment:  data.Comment,
		}, nil
	} else {
		return nil, res.Err()
	}
}
