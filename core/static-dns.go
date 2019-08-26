package core

import (
	"bytes"
	"mikrotik_provisioning/models"
	"mikrotik_provisioning/pkg"
)

func ListStaticDNSTextResponse(staticDNSList []*models.StaticDNSEntry) ([]byte, error) {
	output := bytes.Buffer{}
	err := pkg.API.Templates.ExecuteTemplate(&output, "ListStaticDNS", staticDNSList)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func GetStaticDNSTextResponse(staticDNS *models.StaticDNSEntry) ([]byte, error) {
	output := bytes.Buffer{}
	err := pkg.API.Templates.ExecuteTemplate(&output, "GetStaticDNS", staticDNS)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
