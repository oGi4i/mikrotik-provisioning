package main

import "context"

type Storage interface {
	NewAddressList(ctx context.Context, addressList *AddressList) (*AddressList, error)
	GetAllAddressLists(ctx context.Context) ([]*AddressList, error)
	GetAddressListById(ctx context.Context, id string) (*AddressList, error)
	GetAddressListByName(ctx context.Context, name string) (*AddressList, error)
	UpdateAddressListById(ctx context.Context, id string, addressList *AddressList) (*AddressList, error)
	AddAddressesToAddressListById(ctx context.Context, id string, addresses []Address) (*AddressList, error)
	RemoveAddressListById(ctx context.Context, id string) (*AddressList, error)
	RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []Address) (*AddressList, error)
}

type Implementation struct {
	storage Storage
}

func NewMikrotikAclAPI(storage Storage) *Implementation {
	return &Implementation{storage: storage}
}
