package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"mikrotik_provisioning/pkg"
	"mikrotik_provisioning/types"
	"net/http"
)

func ListAddressLists(w http.ResponseWriter, r *http.Request) {
	results, err := pkg.API.Storage.GetAllAddressLists(r.Context())
	switch r.Context().Value("format") {
	case nil:
		if err != nil {
			render.Render(w, r, types.ErrInternalServerError(err))
		}

		if err := render.RenderList(w, r, types.ListAddressListJSONResponse(results)); err != nil {
			render.Render(w, r, types.ErrRender(err))
		}
	case "rsc":
		if out, err := types.ListAddressListsTextResponse(results); err != nil {
			render.Render(w, r, types.ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}

func CreateAddressList(w http.ResponseWriter, r *http.Request) {
	data := &types.AddressListRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	addressList := data.AddressList
	id, err := pkg.API.Storage.CreateAddressList(r.Context(), addressList)
	if err != nil {
		render.Render(w, r, types.ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, types.NewAddressListResponse(id))
}

func GetAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value("addressList").(*types.AddressList)

	switch r.Context().Value("format") {
	case nil:
		if err := render.Render(w, r, types.NewAddressListResponse(addressList)); err != nil {
			render.Render(w, r, types.ErrRender(err))
		}
	case "rsc":
		if out, err := types.GetAddressListTextResponse(addressList); err != nil {
			render.Render(w, r, types.ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}

func UpdateAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value("addressList").(*types.AddressList)

	data := &types.AddressListRequest{AddressList: addressList}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	addressList = data.AddressList
	addressList, err := pkg.API.Storage.UpdateAddressListById(r.Context(), addressList.ID, addressList)
	if err != nil {
		render.Render(w, r, types.ErrInternalServerError(err))
		return
	}

	render.Render(w, r, types.NewAddressListResponse(addressList))
}

func DeleteAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value("addressList").(*types.AddressList)

	addressList, err := pkg.API.Storage.RemoveAddressListById(r.Context(), addressList.ID)
	if err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func PatchAddressList(w http.ResponseWriter, r *http.Request) {
	var err error
	addressList := r.Context().Value("addressList").(*types.AddressList)

	data := &types.AddressListPatchRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	switch data.Action {
	case "add":
		addressList, err = pkg.API.Storage.AddAddressesToAddressListById(r.Context(), addressList.ID, data.Addresses)
		if err != nil {
			render.Render(w, r, types.ErrInternalServerError(err))
			return
		}
	case "remove":
		addressList, err = pkg.API.Storage.RemoveAddressesFromAddressListById(r.Context(), addressList.ID, data.Addresses)
		if err != nil {
			render.Render(w, r, types.ErrInternalServerError(err))
			return
		}
	default:
		render.Render(w, r, types.ErrInvalidRequest(errors.New("invalid value of Action field")))
	}

	render.Render(w, r, types.NewAddressListResponse(addressList))
}
