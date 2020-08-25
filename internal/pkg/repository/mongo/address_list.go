package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mikrotik_provisioning/internal/pkg/address_list"
)

type AddressList struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	Name      string                  `bson:"name"`
	Addresses []*address_list.Address `bson:"addresses"`
}

func (a *AddressList) ToAddressList() *address_list.AddressList {
	return &address_list.AddressList{
		ID:        a.ID.Hex(),
		Name:      a.Name,
		Addresses: a.Addresses,
	}
}

func (s *Storage) CreateAddressList(ctx context.Context, addressList *address_list.AddressList) (*address_list.AddressList, error) {
	res, err := s.collections["address-list"].InsertOne(ctx, &AddressList{
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	})
	if err != nil {
		return nil, err
	}

	addressList.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return addressList, nil
}

func (s *Storage) GetAddressLists(ctx context.Context) ([]*address_list.AddressList, error) {
	cur, err := s.collections["address-list"].Find(ctx, bson.M{})
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*address_list.AddressList, 0)
	for cur.Next(ctx) {
		data := new(AddressList)
		err := cur.Decode(data)
		if err != nil {
			return nil, err
		}

		result = append(result, data.ToAddressList())
	}

	return result, nil
}

func (s *Storage) GetAddressList(ctx context.Context, name string) (*address_list.AddressList, error) {
	res := s.collections["address-list"].FindOne(ctx, bson.M{"name": name})
	if res.Err() != nil {
		if res.Err().Error() == NoDocumentsError {
			return nil, nil
		}
		return nil, res.Err()
	}

	if b, err := res.DecodeBytes(); err != nil {
		if b.String() == "" && err.Error() == NoDocumentsError {
			return nil, nil
		}
		return nil, err
	}

	data := new(AddressList)
	err := res.Decode(data)
	if err != nil {
		return nil, err
	}

	return data.ToAddressList(), nil
}

func (s *Storage) UpdateAddressList(ctx context.Context, id string, addressList *address_list.AddressList) (*address_list.AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.collections["address-list"].FindOneAndReplace(ctx, bson.M{"_id": objectID}, &AddressList{
		ID:        objectID,
		Name:      addressList.Name,
		Addresses: addressList.Addresses,
	}, options.FindOneAndReplace().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}

	data := new(AddressList)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	return data.ToAddressList(), nil
}

func (s *Storage) UpdateEntriesInAddressList(ctx context.Context, action address_list.Action, id string, addresses []*address_list.Address) (*address_list.AddressList, error) {
	currentData, err := s.getAddressListByID(ctx, id)
	if err != nil {
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	bsonAddresses := addressesToBsonList(currentData, addresses)

	var update bson.M
	switch action {
	case address_list.AddAction:
		update = bson.M{"$push": bson.M{"addresses": bson.M{"$each": bsonAddresses}}}
	case address_list.RemoveAction:
		update = bson.M{"$pull": bson.M{"addresses": bson.M{"$in": bsonAddresses}}}
	}
	res := s.collections["address-list"].FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update)
	if res.Err() != nil {
		return nil, res.Err()
	}

	data := new(AddressList)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	newData, err := s.getAddressListByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return newData.ToAddressList(), nil
}

func (s *Storage) DeleteAddressList(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := s.collections["address-list"].DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("failed deleting mongodb object")
	}

	return nil
}

func (s *Storage) getAddressListByID(ctx context.Context, id string) (*AddressList, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res := s.collections["address-list"].FindOne(ctx, bson.M{"_id": objectID})
	if res.Err() != nil {
		if res.Err().Error() == NoDocumentsError {
			return nil, nil
		}
		return nil, res.Err()
	}

	b, err := res.DecodeBytes()
	if err != nil {
		if b.String() == "" && err.Error() == NoDocumentsError {
			return nil, nil
		}
		return nil, err
	}

	data := new(AddressList)
	err = res.Decode(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func addressesToBsonList(addressList *AddressList, addresses []*address_list.Address) primitive.A {
	bsonA := primitive.A{}
	ok := true
	for _, a := range addresses {
		for _, b := range addressList.Addresses {
			if b == a {
				ok = false
				break
			}
		}
		if ok {
			bsonA = append(bsonA, a)
		}
	}

	return bsonA
}
