package handlers

import (
	"github.com/go-chi/render"
	"mikrotik_provisioning/pkg"
	"mikrotik_provisioning/types"
	"net/http"
)

func ListStaticDNSEntries(w http.ResponseWriter, r *http.Request) {
	results, err := pkg.API.Storage.GetAllStaticDNS(r.Context())
	switch r.Context().Value("format") {
	case nil:
		if err != nil {
			render.Render(w, r, types.ErrInternalServerError(err))
		}

		if err := render.RenderList(w, r, types.ListStaticDNSJSONResponse(results)); err != nil {
			render.Render(w, r, types.ErrRender(err))
		}
	case "rsc":
		if out, err := types.ListStaticDNSTextResponse(results); err != nil {
			render.Render(w, r, types.ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}

func CreateBatchStaticDNSEntries(w http.ResponseWriter, r *http.Request) {
	data := new(types.StaticDNSBatchRequest)
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	list, err := pkg.API.Storage.NewStaticDNSBatch(r.Context(), data.Entries)
	if err != nil {
		render.Render(w, r, types.ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.RenderList(w, r, types.ListStaticDNSJSONResponse(list))
}

func UpdateBatchStaticDNSEntries(w http.ResponseWriter, r *http.Request) {
	data := new(types.StaticDNSBatchRequest)
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	results, err := pkg.API.Storage.UpdateStaticDNSBatch(r.Context(), data.Entries)
	if err != nil {
		render.Render(w, r, types.ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.RenderList(w, r, types.ListStaticDNSJSONResponse(results))
}

func GetStaticDNSEntry(w http.ResponseWriter, r *http.Request) {
	staticDNSEntry := r.Context().Value("staticDNSEntry").(*types.StaticDNSEntry)

	switch r.Context().Value("format") {
	case nil:
		if err := render.Render(w, r, types.NewStaticDNSResponse(staticDNSEntry)); err != nil {
			render.Render(w, r, types.ErrRender(err))
		}
	case "rsc":
		if out, err := types.GetStaticDNSTextResponse(staticDNSEntry); err != nil {
			render.Render(w, r, types.ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}
