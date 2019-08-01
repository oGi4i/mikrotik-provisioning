package main

import (
	"bytes"
	"errors"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Address struct {
	Address  string `json:"address" bson:"address" validate:"required,ipv4|fqdn"`
	Disabled bool   `json:"disabled,omitempty" bson:"disabled,omitempty" validate:"omitempty"`
	Comment  string `json:"comment,omitempty" bson:"comment,omitempty" validate:"omitempty,comment"`
}

type AddressListMongo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" validate:"required"`
	Name      string             `json:"name" bson:"name" validate:"required,addresslistname"`
	Addresses []Address          `json:"addresses" bson:"addresses" validate:"required"`
}

type AddressList struct {
	ID        string    `json:"id,omitempty" validate:"omitempty"`
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
	if a.AddressList == nil {
		return errors.New("missing required AddressList fields")
	}

	return nil
}

func (a *AddressListPatchRequest) Bind(r *http.Request) error {
	if a.Addresses == nil {
		return errors.New("missing required Addresses field")
	}

	if a.Action == "" {
		return errors.New("missing required Action field")
	}

	return nil
}

func NewAddressListResponse(addressList *AddressList) *AddressListResponse {
	return &AddressListResponse{AddressList: addressList}
}

func (rd *AddressListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListAddressListJSONResponse(addressLists []*AddressList) []render.Renderer {
	list := []render.Renderer{}

	for _, addressList := range addressLists {
		list = append(list, NewAddressListResponse(addressList))
	}
	return list
}

func addressListContainsAddress(addr string, addressList *AddressList) bool {
	for _, v := range addressList.Addresses {
		if v.Address == addr {
			return true
		}
	}
	return false
}

func ListAddressListsTextResponse(addressLists []*AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := templates.ExecuteTemplate(&output, "ListAddressLists", addressLists)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func GetAddressListTextResponse(addressList *AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := templates.ExecuteTemplate(&output, "GetAddressList", addressList)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
