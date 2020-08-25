package http

import (
	"bytes"

	"github.com/go-chi/render"

	"mikrotik_provisioning/internal/pkg/address_list"
)

func newAddressListResponse(addressList *address_list.AddressList) *address_list.AddressListResponse {
	return &address_list.AddressListResponse{AddressList: addressList}
}

func getAddressListsJSONResponse(addressLists []*address_list.AddressList) []render.Renderer {
	list := make([]render.Renderer, len(addressLists))

	for i, addressList := range addressLists {
		list[i] = newAddressListResponse(addressList)
	}
	return list
}

func (h *AddressListHandler) getAddressListsTextResponse(addressLists []*address_list.AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := h.templates.ExecuteTemplate(&output, "GetAddressLists", addressLists)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func (h *AddressListHandler) getAddressListTextResponse(addressList *address_list.AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := h.templates.ExecuteTemplate(&output, "GetAddressList", addressList)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
