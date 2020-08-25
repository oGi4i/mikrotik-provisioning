package http

import (
	"net/http"

	"github.com/go-chi/render"

	"mikrotik_provisioning/internal/pkg/address_list"
)

func (h *AddressListHandler) GetAddressLists(w http.ResponseWriter, r *http.Request) {
	results, err := h.service.GetAddressLists(r.Context())
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
	}

	var out []byte
	switch r.Context().Value(FormatKey) {
	case RSCFormat:
		if len(results) != 0 {
			out, err = h.getAddressListsTextResponse(results)
			if err != nil {
				_ = render.Render(w, r, ErrRender(err))
			}
			_, _ = w.Write(out)
		} else {
			render.Status(r, http.StatusOK)
		}
	default:
		if err != nil {
			_ = render.Render(w, r, ErrInternalServerError(err))
		}

		if err := render.RenderList(w, r, getAddressListsJSONResponse(results)); err != nil {
			_ = render.Render(w, r, ErrRender(err))
		}
	}
}

func (h *AddressListHandler) CreateAddressList(w http.ResponseWriter, r *http.Request) {
	data := &address_list.AddressListRequest{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	addressList, err := h.service.CreateAddressList(r.Context(), data.AddressList)
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, newAddressListResponse(addressList))
}

func (h *AddressListHandler) GetAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value(AddressListKey).(*address_list.AddressList)

	switch r.Context().Value(FormatKey) {
	case RSCFormat:
		if out, err := h.getAddressListTextResponse(addressList); err != nil {
			_ = render.Render(w, r, ErrRender(err))
		} else {
			_, _ = w.Write(out)
		}
	default:
		if err := render.Render(w, r, newAddressListResponse(addressList)); err != nil {
			_ = render.Render(w, r, ErrRender(err))
		}
	}
}

func (h *AddressListHandler) UpdateAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value(AddressListKey).(*address_list.AddressList)

	data := &address_list.AddressListRequest{AddressList: addressList}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	addressList, err := h.service.UpdateAddressList(r.Context(), data.AddressList.ID, data.AddressList)
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
		return
	}

	_ = render.Render(w, r, newAddressListResponse(addressList))
}

func (h *AddressListHandler) DeleteAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value(AddressListKey).(*address_list.AddressList)

	err := h.service.DeleteAddressList(r.Context(), addressList.ID)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *AddressListHandler) PatchAddressList(w http.ResponseWriter, r *http.Request) {
	addressList := r.Context().Value(AddressListKey).(*address_list.AddressList)

	data := &address_list.AddressListPatchRequest{}
	err := render.Bind(r, data)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	addressList, err = h.service.UpdateEntriesInAddressList(r.Context(), data.Action, addressList.ID, data.Addresses)
	if err != nil {
		_ = render.Render(w, r, ErrInternalServerError(err))
		return
	}

	_ = render.Render(w, r, newAddressListResponse(addressList))
}
