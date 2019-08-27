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
	"mikrotik_provisioning/models"
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

func createIndexes(ctx context.Context, coll *mongo.Collection, indexList []*cfg.CollectionIndexesConfig) error {
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
			continue
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

func (s *MongoStorage) CreateAddressList(ctx context.Context, addressList *models.AddressList) (*models.AddressList, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res, err := s.colls["address-list"].InsertOne(c, &models.AddressListMongo{
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	})
	if err != nil {
		return nil, err
	}
	addressList.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return addressList, nil
}

func (s *MongoStorage) GetAllAddressLists(ctx context.Context) ([]*models.AddressList, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	cur, err := s.colls["address-list"].Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(c)

	result := make([]*models.AddressList, 0)
	for cur.Next(ctx) {
		data := new(models.AddressListMongo)

		err := cur.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		result = append(result, &models.AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		})
	}

	return result, nil
}

func (s *MongoStorage) getAddressListById(ctx context.Context, id string) (*models.AddressList, error) {
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

		data := new(models.AddressListMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &models.AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) GetAddressListByName(ctx context.Context, name string) (*models.AddressList, error) {
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

		data := new(models.AddressListMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &models.AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) UpdateAddressListById(ctx context.Context, id string, addressList *models.AddressList) (*models.AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.colls["address-list"].FindOneAndReplace(ctx, bson.M{"_id": objectID}, &models.AddressListMongo{
		ID:        objectID,
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	}, options.FindOneAndReplace().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}

	data := new(models.AddressListMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	return &models.AddressList{
		ID:        data.ID.Hex(),
		Name:      data.Name,
		Addresses: data.Addresses,
	}, nil
}

func (s *MongoStorage) AddAddressesToAddressListById(ctx context.Context, id string, addresses []models.Address) (*models.AddressList, error) {
	currentData, err := s.getAddressListById(ctx, id)
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
		return nil, res.Err()
	}

	data := new(models.AddressListMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	newData, err := s.getAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(newData); err != nil {
		return nil, err
	}

	return newData, nil
}

func (s *MongoStorage) RemoveAddressListById(ctx context.Context, id string) (*models.AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
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

func (s *MongoStorage) RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []models.Address) (*models.AddressList, error) {
	currentData, err := s.getAddressListById(ctx, id)
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
		return nil, res.Err()
	}

	data := new(models.AddressListMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	newData, err := s.getAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(newData); err != nil {
		return nil, err
	}

	return newData, nil
}

func (s *MongoStorage) GetAllStaticDNS(ctx context.Context) ([]*models.StaticDNSEntry, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	cur, err := s.colls["static-dns"].Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(c)

	result := make([]*models.StaticDNSEntry, 0)
	for cur.Next(ctx) {
		data := new(models.StaticDNSEntryMongo)

		err := cur.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		result = append(result, &models.StaticDNSEntry{
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

func (s *MongoStorage) CreateStaticDNSEntriesFromBatch(ctx context.Context, entries []*models.StaticDNSEntry) ([]*models.StaticDNSEntry, error) {
	mongoEntries := make([]interface{}, len(entries))
	for i, entry := range entries {
		mongoEntries[i] = &models.StaticDNSEntryMongo{
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

func (s *MongoStorage) UpdateStaticDNSEntriesFromBatch(ctx context.Context, entries []*models.StaticDNSEntry) ([]*models.StaticDNSEntry, error) {
	mongoEntries := make([]*models.StaticDNSEntryMongo, len(entries))
	for i, entry := range entries {
		mongoEntries[i] = &models.StaticDNSEntryMongo{
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

		data := new(models.StaticDNSEntryMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		entries[i] = &models.StaticDNSEntry{
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

func (s *MongoStorage) GetStaticDNSEntryByName(ctx context.Context, name string) (*models.StaticDNSEntry, error) {
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

		data := new(models.StaticDNSEntryMongo)
		err := res.Decode(data)
		if err != nil {
			return nil, err
		}

		if err := valid.Validate.Struct(data); err != nil {
			return nil, err
		}

		return &models.StaticDNSEntry{
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

func (s *MongoStorage) CreateStaticDNSEntry(ctx context.Context, entry *models.StaticDNSEntry) (*models.StaticDNSEntry, error) {
	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res, err := s.colls["static-dns"].InsertOne(c, &models.StaticDNSEntryMongo{
		Name:     entry.Name,
		Regexp:   entry.Regexp,
		Address:  entry.Address,
		TTL:      entry.TTL,
		Disabled: entry.Disabled,
		Comment:  entry.Comment,
	})
	if err != nil {
		return nil, err
	}
	entry.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return entry, nil
}

func (s *MongoStorage) UpdateStaticDNSEntryById(ctx context.Context, id string, entry *models.StaticDNSEntry) (*models.StaticDNSEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.colls["static-dns"].FindOneAndReplace(ctx, bson.M{"_id": objectID}, &models.StaticDNSEntryMongo{
		Name:     entry.Name,
		Regexp:   entry.Regexp,
		Address:  entry.Address,
		TTL:      entry.TTL,
		Disabled: entry.Disabled,
		Comment:  entry.Comment,
	}, options.FindOneAndReplace().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}

	data := new(models.StaticDNSEntryMongo)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	if err := valid.Validate.Struct(data); err != nil {
		return nil, err
	}

	return &models.StaticDNSEntry{
		ID:       data.ID.Hex(),
		Name:     data.Name,
		Regexp:   data.Regexp,
		Address:  data.Address,
		TTL:      data.TTL,
		Disabled: data.Disabled,
		Comment:  data.Comment,
	}, nil
}

func (s *MongoStorage) RemoveStaticDNSEntryById(ctx context.Context, id string) (*models.StaticDNSEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	c, _ := context.WithTimeout(ctx, cfg.Config.Database.Timeout*time.Second)
	res, err := s.colls["static-dns"].DeleteOne(c, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("failed deleting mongoDB object")
	}

	return nil, nil
}
