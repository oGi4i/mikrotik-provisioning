package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	valid "mikrotik_provisioning/validate"
	"net/http"
)

type Address struct {
	Address  string `json:"address" bson:"address" validate:"required,ipv4|fqdn"`
	Disabled bool   `json:"disabled,omitempty" bson:"disabled,omitempty" validate:"omitempty"`
	Comment  string `json:"comment,omitempty" bson:"comment,omitempty" validate:"omitempty,comment"`
}

type AddressListMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" validate:"required"`
	Name      string             `bson:"name" validate:"required,addresslistname"`
	Addresses []Address          `bson:"addresses" validate:"required"`
}

type AddressList struct {
	ID        string    `json:"-" validate:"omitempty"`
	Name      string    `json:"name" validate:"required,addresslistname"`
	Addresses []Address `json:"addresses" validate:"required"`
}

type AddressListRequest struct {
	*AddressList
}

type AddressListResponse struct {
	*AddressList
}

type AddressListPatchRequest struct {
	Action    string    `json:"action" validate:"required,oneof=add remove"`
	Addresses []Address `json:"addresses" validate:"required"`
}

func (a *AddressListRequest) Bind(r *http.Request) error {
	if err := valid.Validate.Struct(a); err != nil {
		return err
	}

	return nil
}

func (a *AddressListPatchRequest) Bind(r *http.Request) error {
	if err := valid.Validate.Struct(a); err != nil {
		return err
	}

	return nil
}

func (rd *AddressListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
