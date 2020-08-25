package address_list

import (
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type (
	Address struct {
		Address  string `json:"address" bson:"address" validator:"required,ipv4|fqdn"`
		Disabled bool   `json:"disabled,omitempty" bson:"disabled,omitempty" validator:"omitempty"`
		Comment  string `json:"comment,omitempty" bson:"comment,omitempty" validator:"omitempty,comment"`
	}

	AddressList struct {
		ID        string     `json:"-" validator:"omitempty"`
		Name      string     `json:"name" validator:"required,address_list_name"`
		Addresses []*Address `json:"addresses" validator:"required"`
	}

	AddressListRequest struct {
		*AddressList
	}

	AddressListResponse struct {
		*AddressList
	}

	AddressListPatchRequest struct {
		Action    Action     `json:"action" validator:"required,oneof=add remove"`
		Addresses []*Address `json:"addresses" validator:"required"`
	}

	Action string
)

const (
	AddAction    Action = "add"
	RemoveAction Action = "remove"
)

func (a *AddressListRequest) Bind(r *http.Request) error {
	validator := validator.New()
	if err := validator.Struct(a); err != nil {
		return err
	}

	return nil
}

func (a *AddressListPatchRequest) Bind(r *http.Request) error {
	validator := validator.New()
	if err := validator.Struct(a); err != nil {
		return err
	}

	return nil
}

func (rd *AddressListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
