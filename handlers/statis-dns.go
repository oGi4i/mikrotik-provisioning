package handlers

import (
	"github.com/go-chi/render"
	"mikrotik_provisioning/core"
	"mikrotik_provisioning/models"
	"mikrotik_provisioning/pkg"
	"net/http"
	"strings"
)

func ListStaticDNSEntries(w http.ResponseWriter, r *http.Request) {
	results, err := pkg.API.Storage.GetAllStaticDNS(r.Context())
	switch r.Context().Value("format") {
	case nil:
		if err != nil {
			render.Render(w, r, models.ErrInternalServerError(err))
		}

		if err := render.RenderList(w, r, models.ListStaticDNSJSONResponse(results)); err != nil {
			render.Render(w, r, models.ErrRender(err))
		}
	case "rsc":
		if len(results) != 0 {
			if out, err := core.ListStaticDNSTextResponse(results); err != nil {
				render.Render(w, r, models.ErrRender(err))
			} else {
				w.Write(out)
			}
		} else {
			render.Status(r, http.StatusOK)
		}
	}
}

func CreateBatchStaticDNSEntries(w http.ResponseWriter, r *http.Request) {
	data := new(models.StaticDNSBatchRequest)
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	list, err := pkg.API.Storage.CreateStaticDNSEntriesFromBatch(r.Context(), data.Entries)
	if err != nil {
		if strings.HasPrefix(err.Error(), "time") {
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}
		render.Render(w, r, models.ErrInternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.RenderList(w, r, models.ListStaticDNSJSONResponse(list))
}

func UpdateBatchStaticDNSEntries(w http.ResponseWriter, r *http.Request) {
	data := new(models.StaticDNSBatchRequest)
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	results, err := pkg.API.Storage.UpdateStaticDNSEntriesFromBatch(r.Context(), data.Entries)
	if err != nil {
		if strings.HasPrefix(err.Error(), "time") {
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}
		render.Render(w, r, models.ErrInternalServerError(err))
		return
	}

	render.RenderList(w, r, models.ListStaticDNSJSONResponse(results))
}

func GetStaticDNSEntry(w http.ResponseWriter, r *http.Request) {
	staticDNSEntry := r.Context().Value("staticDNSEntry").(*models.StaticDNSEntry)

	switch r.Context().Value("format") {
	case nil:
		if err := render.Render(w, r, models.NewStaticDNSResponse(staticDNSEntry)); err != nil {
			render.Render(w, r, models.ErrRender(err))
		}
	case "rsc":
		if out, err := core.GetStaticDNSTextResponse(staticDNSEntry); err != nil {
			render.Render(w, r, models.ErrRender(err))
		} else {
			w.Write(out)
		}
	}
}

func CreateStaticDNSEntry(w http.ResponseWriter, r *http.Request) {
	data := new(models.StaticDNSRequest)
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	staticDNSEntry := data.StaticDNSEntry
	staticDNSEntry, err := pkg.API.Storage.CreateStaticDNSEntry(r.Context(), staticDNSEntry)
	if err != nil {
		if strings.HasPrefix(err.Error(), "time") {
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}
		render.Render(w, r, models.ErrInternalServerError(err))
		return
	}

	render.Render(w, r, models.NewStaticDNSResponse(staticDNSEntry))
}

func UpdateStaticDNSEntry(w http.ResponseWriter, r *http.Request) {
	staticDNSEntry := r.Context().Value("staticDNSEntry").(*models.StaticDNSEntry)

	data := &models.StaticDNSRequest{StaticDNSEntry: staticDNSEntry}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}
	staticDNSEntry = data.StaticDNSEntry
	staticDNSEntry, err := pkg.API.Storage.UpdateStaticDNSEntryById(r.Context(), staticDNSEntry.ID, staticDNSEntry)
	if err != nil {
		if strings.HasPrefix(err.Error(), "time") {
			render.Render(w, r, models.ErrInvalidRequest(err))
			return
		}
		render.Render(w, r, models.ErrInternalServerError(err))
		return
	}

	render.Render(w, r, models.NewStaticDNSResponse(staticDNSEntry))
}

func DeleteStaticDNSEntry(w http.ResponseWriter, r *http.Request) {
	staticDNSEntry := r.Context().Value("staticDNSEntry").(*models.StaticDNSEntry)

	staticDNSEntry, err := pkg.API.Storage.RemoveStaticDNSEntryById(r.Context(), staticDNSEntry.ID)
	if err != nil {
		render.Render(w, r, models.ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}
