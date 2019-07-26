package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoStorage struct {
	cl   *mongo.Client
	coll *mongo.Collection
}

func NewMongoStorage(client *mongo.Client, collection *mongo.Collection) *MongoStorage {
	return &MongoStorage{cl: client, coll: collection}
}

func (s *MongoStorage) NewAddressList(ctx context.Context, addressList *AddressList) (*AddressList, error) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := s.coll.InsertOne(ctx, &AddressListMongo{
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	})
	if err != nil {
		return nil, err
	}
	addressList.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return addressList, nil
}

func (s *MongoStorage) GetAllAddressLists(ctx context.Context) ([]*AddressList, error) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []*AddressList
	for cur.Next(ctx) {
		var data AddressListMongo

		err := cur.Decode(&data)
		if err != nil {
			return nil, err
		}

		result = append(result, &AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		})
	}

	return result, nil
}

func (s *MongoStorage) GetAddressListById(ctx context.Context, id string) (*AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res := s.coll.FindOne(ctx, bson.M{"_id": objectID})

	if res.Err() == nil {
		if b, err := res.DecodeBytes(); err != nil {
			if b.String() == "" && err.Error() == "mongo: no documents in result" {
				return nil, nil
			} else {
				return nil, err
			}
		}

		var data AddressListMongo
		err := res.Decode(&data)
		if err != nil {
			return nil, err
		}

		return &AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) GetAddressListByName(ctx context.Context, name string) (*AddressList, error) {
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res := s.coll.FindOne(ctx, bson.M{"name": name})

	if res.Err() == nil {
		if b, err := res.DecodeBytes(); err != nil {
			if b.String() == "" && err.Error() == "mongo: no documents in result" {
				return nil, nil
			} else {
				return nil, err
			}
		}

		var data AddressListMongo
		err := res.Decode(&data)
		if err != nil {
			return nil, err
		}

		return &AddressList{
			ID:        data.ID.Hex(),
			Name:      data.Name,
			Addresses: data.Addresses,
		}, nil
	} else {
		return nil, res.Err()
	}
}

func (s *MongoStorage) UpdateAddressListById(ctx context.Context, id string, addressList *AddressList) (*AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.coll.FindOneAndReplace(ctx, bson.M{"_id": objectID}, &AddressListMongo{
		ID:        objectID,
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	})
	if res.Err() != nil {
		return nil, err
	}

	var data *AddressListMongo
	err = res.Decode(&data)
	if err != nil {
		return nil, err
	}

	return addressList, nil
}

func (s *MongoStorage) AddAddressesToAddressListById(ctx context.Context, id string, addresses []Address) (*AddressList, error) {
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

	res := s.coll.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if res.Err() != nil {
		return nil, err
	}

	var data *AddressListMongo
	err = res.Decode(&data)
	if err != nil {
		return nil, err
	}

	newData, err := s.GetAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	return newData, nil
}

func (s *MongoStorage) RemoveAddressListById(ctx context.Context, id string) (*AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid addressListId")
	}

	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("failed deleting mongoDB object")
	}

	return nil, nil
}

func (s *MongoStorage) RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []Address) (*AddressList, error) {
	currentData, err := s.GetAddressListById(ctx, id)
	if err != nil {
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

	res := s.coll.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if res.Err() != nil {
		return nil, err
	}

	var data *AddressListMongo
	err = res.Decode(&data)
	if err != nil {
		return nil, err
	}

	newData, err := s.GetAddressListById(ctx, id)
	if err != nil {
		return nil, err
	}

	return newData, nil
}
