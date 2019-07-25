package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Address struct {
	Address  string `json:"address" bson:"address"`
	Disabled bool   `json:"disabled,omitempty" bson:"disabled,omitempty"`
	Comment  string `json:"comment,omitempty" bson:"comment,omitempty"`
}

type AddressListMongo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Addresses []Address          `json:"addresses" bson:"addresses"`
}

type AddressList struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Addresses []Address `json:"addresses"`
}

type AddressListRequest struct {
	*AddressList
}

type AddressListResponse struct {
	*AddressList
}

type AddressListPatchRequest struct {
	Action    string    `json:"action"`
	Addresses []Address `json:"addresses"`
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

func ListAddressListTextResponse(addressLists []*AddressList) []byte {
	output := bytes.Buffer{}
	for _, addressList := range addressLists {
		output.Write(GetAddressListTextResponse(addressList))
	}
	return output.Bytes()
}

func GetAddressListTextResponse(addressList *AddressList) []byte {
	output := bytes.Buffer{}
	var disabled string
	for _, addr := range addressList.Addresses {
		switch addr.Disabled {
		case true:
			disabled = "yes"
		case false:
			disabled = "no"
		}
		output.WriteString(fmt.Sprintf("/ip firewall address-list add list=%s address=%s disabled=%s\n", addressList.Name, addr.Address, disabled))
	}
	return output.Bytes()
}
