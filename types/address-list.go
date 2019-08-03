package types

import (
	"bytes"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	cfg "mikrotik_provisioning/config"
	valid "mikrotik_provisioning/validate"
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

func NewAddressListResponse(addressList *AddressList) *AddressListResponse {
	return &AddressListResponse{AddressList: addressList}
}

func (rd *AddressListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListAddressListJSONResponse(addressLists []*AddressList) []render.Renderer {
	list := make([]render.Renderer, len(addressLists))

	for i, addressList := range addressLists {
		list[i] = NewAddressListResponse(addressList)
	}
	return list
}

func ListAddressListsTextResponse(addressLists []*AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := cfg.Templates.ExecuteTemplate(&output, "ListAddressLists", addressLists)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func GetAddressListTextResponse(addressList *AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := cfg.Templates.ExecuteTemplate(&output, "GetAddressList", addressList)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
