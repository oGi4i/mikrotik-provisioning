package storage

import (
	"context"
	"errors"
	"mikrotik_provisioning/types"
	valid "mikrotik_provisioning/validate"
)

var (
	addressLists  = []*types.AddressList{}
	staticDNSList = []*types.StaticDNSEntry{}
)

type StructStorage struct {
	data []*types.AddressList
}

func NewStructStorage(data []*types.AddressList) *StructStorage {
	return &StructStorage{data}
}

func (s *StructStorage) NewAddressList(ctx context.Context, addressList *types.AddressList) (*types.AddressList, error) {
	if err := valid.Validate.Struct(addressList); err != nil {
		return nil, err
	}

	addressLists = append(addressLists, addressList)

	return addressList, nil
}

func (s *StructStorage) GetAllAddressLists(ctx context.Context) ([]*types.AddressList, error) {
	return addressLists, nil
}

func (s *StructStorage) GetAddressListById(ctx context.Context, id string) (*types.AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) GetAddressListByName(ctx context.Context, name string) (*types.AddressList, error) {
	for _, a := range addressLists {
		if a.Name == name {
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) UpdateAddressListById(ctx context.Context, id string, addressList *types.AddressList) (*types.AddressList, error) {
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

func (s *StructStorage) AddAddressesToAddressListById(ctx context.Context, id string, addresses []types.Address) (*types.AddressList, error) {
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

func (s *StructStorage) RemoveAddressListById(ctx context.Context, id string) (*types.AddressList, error) {
	for i, a := range addressLists {
		if a.ID == id {
			addressLists = append((addressLists)[:i], (addressLists)[i+1:]...)
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []types.Address) (*types.AddressList, error) {
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

func (s *StructStorage) removeAddressFromAddressList(address types.Address, addressList *types.AddressList) error {
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

func addressListContainsAddress(addr string, addressList *types.AddressList) bool {
	for _, v := range addressList.Addresses {
		if v.Address == addr {
			return true
		}
	}
	return false
}
