package storage

import (
	"context"
	"errors"
	"mikrotik_provisioning/models"
	valid "mikrotik_provisioning/validate"
)

var (
	addressLists  = make([]*models.AddressList, 0)
	staticDNSList = make([]*models.StaticDNSEntry, 0)
)

type StructStorage struct {
	data []*models.AddressList
}

func NewStructStorage(data []*models.AddressList) *StructStorage {
	return &StructStorage{data}
}

func (s *StructStorage) NewAddressList(ctx context.Context, addressList *models.AddressList) (*models.AddressList, error) {
	if err := valid.Validate.Struct(addressList); err != nil {
		return nil, err
	}

	addressLists = append(addressLists, addressList)

	return addressList, nil
}

func (s *StructStorage) GetAllAddressLists(ctx context.Context) ([]*models.AddressList, error) {
	return addressLists, nil
}

func (s *StructStorage) GetAddressListById(ctx context.Context, id string) (*models.AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) GetAddressListByName(ctx context.Context, name string) (*models.AddressList, error) {
	for _, a := range addressLists {
		if a.Name == name {
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) UpdateAddressListById(ctx context.Context, id string, addressList *models.AddressList) (*models.AddressList, error) {
	if err := valid.Validate.Struct(addressList); err != nil {
		return nil, err
	}

	for i, a := range addressLists {
		if a.ID == id {
			addressLists[i] = addressList
			return addressList, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) AddAddressesToAddressListById(ctx context.Context, id string, addresses []models.Address) (*models.AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			for _, b := range addresses {
				if err := valid.Validate.Struct(b); err != nil {
					return nil, err
				}

				if addressListContainsAddress(b.Address, a) {
					continue
				}
				a.Addresses = append(a.Addresses, b)
			}
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) RemoveAddressListById(ctx context.Context, id string) (*models.AddressList, error) {
	for i, a := range addressLists {
		if a.ID == id {
			addressLists = append((addressLists)[:i], (addressLists)[i+1:]...)
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []models.Address) (*models.AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			for _, b := range addresses {
				if err := valid.Validate.Struct(b); err != nil {
					return nil, err
				}

				if addressListContainsAddress(b.Address, a) {
					if err := s.removeAddressFromAddressList(b, a); err != nil {
						return nil, err
					}
				}
			}
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) removeAddressFromAddressList(address models.Address, addressList *models.AddressList) error {
	if err := valid.Validate.Struct(addressList); err != nil {
		return err
	}

	for i, a := range addressList.Addresses {
		if a == address {
			addressList.Addresses = append((addressList.Addresses)[:i], (addressList.Addresses)[i+1:]...)
			return nil
		}
	}

	return errors.New("address not found")
}

func addressListContainsAddress(addr string, addressList *models.AddressList) bool {
	for _, v := range addressList.Addresses {
		if v.Address == addr {
			return true
		}
	}
	return false
}
