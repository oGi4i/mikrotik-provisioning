package app

import (
	"context"
	"text/template"

	"mikrotik_provisioning/internal/pkg/address_list"
)

type UseCases interface {
	Storage
}

type Storage interface {
	GetAddressLists(ctx context.Context) ([]*address_list.AddressList, error)
	CreateAddressList(ctx context.Context, addressList *address_list.AddressList) (*address_list.AddressList, error)
	GetAddressList(ctx context.Context, name string) (*address_list.AddressList, error)
	UpdateAddressList(ctx context.Context, id string, addressList *address_list.AddressList) (*address_list.AddressList, error)
	DeleteAddressList(ctx context.Context, id string) error
	UpdateEntriesInAddressList(ctx context.Context, action address_list.Action, id string, addresses []*address_list.Address) (*address_list.AddressList, error)
}

type Service struct {
	storage   Storage
	templates *template.Template
}

func NewMikrotikProvisioningService(storage Storage, templates *template.Template) *Service {
	return &Service{storage: storage, templates: templates}
}

func (s *Service) GetAddressLists(ctx context.Context) ([]*address_list.AddressList, error) {
	return s.storage.GetAddressLists(ctx)
}

func (s *Service) CreateAddressList(ctx context.Context, addressList *address_list.AddressList) (*address_list.AddressList, error) {
	return s.storage.CreateAddressList(ctx, addressList)
}

func (s *Service) GetAddressList(ctx context.Context, name string) (*address_list.AddressList, error) {
	return s.storage.GetAddressList(ctx, name)
}

func (s *Service) UpdateAddressList(ctx context.Context, id string, addressList *address_list.AddressList) (*address_list.AddressList, error) {
	return s.storage.UpdateAddressList(ctx, id, addressList)
}

func (s *Service) DeleteAddressList(ctx context.Context, id string) error {
	return s.storage.DeleteAddressList(ctx, id)
}

func (s *Service) UpdateEntriesInAddressList(ctx context.Context, action address_list.Action, id string, addresses []*address_list.Address) (*address_list.AddressList, error) {
	return s.storage.UpdateEntriesInAddressList(ctx, action, id, addresses)
}
