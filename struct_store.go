package main

import (
	"context"
	"errors"
)

type StructStorage struct {
	data []*AddressList
}

func NewStructStorage(data []*AddressList) *StructStorage {
	return &StructStorage{data}
}

func (s *StructStorage) NewAddressList(ctx context.Context, addressList *AddressList) (*AddressList, error) {
	if err := validate.Struct(addressList); err != nil {
		return nil, err
	}

	addressLists = append(addressLists, addressList)

	return addressList, nil
}

func (s *StructStorage) GetAllAddressLists(ctx context.Context) ([]*AddressList, error) {
	return addressLists, nil
}

func (s *StructStorage) GetAddressListById(ctx context.Context, id string) (*AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) GetAddressListByName(ctx context.Context, name string) (*AddressList, error) {
	for _, a := range addressLists {
		if a.Name == name {
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) UpdateAddressListById(ctx context.Context, id string, addressList *AddressList) (*AddressList, error) {
	if err := validate.Struct(addressList); err != nil {
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

func (s *StructStorage) AddAddressesToAddressListById(ctx context.Context, id string, addresses []Address) (*AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			for _, b := range addresses {
				if err := validate.Struct(b); err != nil {
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

func (s *StructStorage) RemoveAddressListById(ctx context.Context, id string) (*AddressList, error) {
	for i, a := range addressLists {
		if a.ID == id {
			addressLists = append((addressLists)[:i], (addressLists)[i+1:]...)
			return a, nil
		}
	}

	return nil, errors.New("address list not found")
}

func (s *StructStorage) RemoveAddressesFromAddressListById(ctx context.Context, id string, addresses []Address) (*AddressList, error) {
	for _, a := range addressLists {
		if a.ID == id {
			for _, b := range addresses {
				if err := validate.Struct(b); err != nil {
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

func (s *StructStorage) removeAddressFromAddressList(address Address, addressList *AddressList) error {
	if err := validate.Struct(addressList); err != nil {
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
