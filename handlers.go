package main

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

func ListAddressLists(w http.ResponseWriter, r *http.Request) {
	results, err := api.storage.GetAllAddressLists(r.Context())
	switch r.Context().Value("format") {
	case nil:
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
		}

		if err := render.RenderList(w, r, ListAddressListJSONResponse(results)); err != nil {
			render.Render(w, r, ErrRender(err))
		}
	case "rsc":
		if out, err := ListAddressListsTextResponse(results); err != nil {
			render.Render(w, r, ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}

func CreateAddressList(w http.ResponseWriter, r *http.Request) {
	data := &AddressListRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	addressList := data.AddressList
	id, err := api.storage.NewAddressList(r.Context(), addressList)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewAddressListResponse(id))
}

func GetAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value("addressList").(*AddressList)

	switch r.Context().Value("format") {
	case nil:
		if err := render.Render(w, r, NewAddressListResponse(addressList)); err != nil {
			render.Render(w, r, ErrRender(err))
		}
	case "rsc":
		if out, err := GetAddressListTextResponse(addressList); err != nil {
			render.Render(w, r, ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}

func UpdateAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value("addressList").(*AddressList)

	data := &AddressListRequest{AddressList: addressList}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	addressList = data.AddressList
	addressList, err := api.storage.UpdateAddressListById(r.Context(), addressList.ID, addressList)
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Render(w, r, NewAddressListResponse(addressList))
}

func DeleteAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value("addressList").(*AddressList)

	addressList, err := api.storage.RemoveAddressListById(r.Context(), addressList.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func PatchAddressList(w http.ResponseWriter, r *http.Request) {
	var err error
	addressList := r.Context().Value("addressList").(*AddressList)

	data := &AddressListPatchRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	switch data.Action {
	case "add":
		addressList, err = api.storage.AddAddressesToAddressListById(r.Context(), addressList.ID, data.Addresses)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
	case "remove":
		addressList, err = api.storage.RemoveAddressesFromAddressListById(r.Context(), addressList.ID, data.Addresses)
		if err != nil {
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
	default:
		render.Render(w, r, ErrInvalidRequest(errors.New("Invalid value of Action field")))
	}

	render.Render(w, r, NewAddressListResponse(addressList))
}
