package core

import (
	"bytes"
	"mikrotik_provisioning/models"
	"mikrotik_provisioning/pkg"
)

func ListAddressListsTextResponse(addressLists []*models.AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := pkg.API.Templates.ExecuteTemplate(&output, "ListAddressLists", addressLists)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func GetAddressListTextResponse(addressList *models.AddressList) ([]byte, error) {
	output := bytes.Buffer{}
	err := pkg.API.Templates.ExecuteTemplate(&output, "GetAddressList", addressList)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
