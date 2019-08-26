package pkg

import (
	"context"
	cfg "mikrotik_provisioning/config"
	"mikrotik_provisioning/models"
	"text/template"
)

var API = &Implementation{}

type Storage interface {
	CreateAddressList(ctx context.Context, addressList *models.AddressList) (*models.AddressList, error)
	GetAllAddressLists(ctx context.Context) ([]*models.AddressList, error)
	GetAddressListByName(ctx context.Context, name string) (*models.AddressList, error)
	UpdateAddressListById(ctx context.Context, id string, addressList *models.AddressList) (*models.AddressList, error)
	AddAddressesToAddressListById(ctx context.Context, id string, addresses []models.Address) (*models.AddressList, error)
	RemoveAddressListById(ctx context.Context, id string) (*models.AddressList, error)
	RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []models.Address) (*models.AddressList, error)
	GetAllStaticDNS(ctx context.Context) ([]*models.StaticDNSEntry, error)
	CreateStaticDNSEntriesFromBatch(ctx context.Context, entries []*models.StaticDNSEntry) ([]*models.StaticDNSEntry, error)
	UpdateStaticDNSEntriesFromBatch(ctx context.Context, entries []*models.StaticDNSEntry) ([]*models.StaticDNSEntry, error)
	GetStaticDNSEntryByName(ctx context.Context, name string) (*models.StaticDNSEntry, error)
	CreateStaticDNSEntry(ctx context.Context, entry *models.StaticDNSEntry) (*models.StaticDNSEntry, error)
	UpdateStaticDNSEntryById(ctx context.Context, id string, entry *models.StaticDNSEntry) (*models.StaticDNSEntry, error)
	RemoveStaticDNSEntryById(ctx context.Context, id string) (*models.StaticDNSEntry, error)
}

type Implementation struct {
	Storage   Storage
	Config    *cfg.ApplicationConfig
	Templates *template.Template
}

func NewMikrotikProvisioningAPI(storage Storage, config *cfg.ApplicationConfig, templates *template.Template) *Implementation {
	return &Implementation{Storage: storage, Config: config, Templates: templates}
}
