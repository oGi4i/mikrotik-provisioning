package pkg

import (
	"context"
	"mikrotik_provisioning/types"
)

var API = &Implementation{}

type Storage interface {
	CreateAddressList(ctx context.Context, addressList *types.AddressList) (*types.AddressList, error)
	GetAllAddressLists(ctx context.Context) ([]*types.AddressList, error)
	GetAddressListByName(ctx context.Context, name string) (*types.AddressList, error)
	UpdateAddressListById(ctx context.Context, id string, addressList *types.AddressList) (*types.AddressList, error)
	AddAddressesToAddressListById(ctx context.Context, id string, addresses []types.Address) (*types.AddressList, error)
	RemoveAddressListById(ctx context.Context, id string) (*types.AddressList, error)
	RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []types.Address) (*types.AddressList, error)
	GetAllStaticDNS(ctx context.Context) ([]*types.StaticDNSEntry, error)
	CreateStaticDNSEntriesFromBatch(ctx context.Context, entries []*types.StaticDNSEntry) ([]*types.StaticDNSEntry, error)
	UpdateStaticDNSEntriesFromBatch(ctx context.Context, entries []*types.StaticDNSEntry) ([]*types.StaticDNSEntry, error)
	GetStaticDNSEntryByName(ctx context.Context, name string) (*types.StaticDNSEntry, error)
	CreateStaticDNSEntry(ctx context.Context, entry *types.StaticDNSEntry) (*types.StaticDNSEntry, error)
	UpdateStaticDNSEntryById(ctx context.Context, id string, entry *types.StaticDNSEntry) (*types.StaticDNSEntry, error)
	RemoveStaticDNSEntryById(ctx context.Context, id string) (*types.StaticDNSEntry, error)
}

type Implementation struct {
	Storage Storage
}

func NewMikrotikProvisioningAPI(storage Storage) *Implementation {
	return &Implementation{Storage: storage}
}
