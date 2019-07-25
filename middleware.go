package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"io/ioutil"
	"net/http"
	"strings"
)

func EnsureAddressListExists(i *Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var addressList *AddressList
			var err error
			ctx := context.Background()

			if addressListId := chi.URLParam(r, "addressListId"); addressListId != "" {
				addressList, err = i.storage.GetAddressListById(ctx, addressListId)
			} else if addressListName := chi.URLParam(r, "addressListName"); addressListName != "" {
				addressList, err = i.storage.GetAddressListByName(ctx, addressListName)
			} else {
				render.Render(w, r, ErrNotFound)
				return
			}
			if err != nil {
				render.Render(w, r, ErrNotFound)
				return
			}

			ctx = context.WithValue(r.Context(), "addressList", addressList)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func EnsureAddressListNotExists(i *Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var addressList *AddressList
			var err error
			ctx := context.Background()

			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			err = json.Unmarshal(bodyBytes, &addressList)
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			result, err := i.storage.GetAddressListByName(ctx, addressList.Name)
			if err != nil {
				render.Render(w, r, ErrInternalServerError(err))
				return
			}
			if result != nil {
				render.Render(w, r, ErrInvalidRequest(errors.New("address list already exists")))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckAcceptHeader(contentTypes ...string) func(next http.Handler) http.Handler {
	cT := []string{}
	for _, t := range contentTypes {
		cT = append(cT, strings.ToLower(t))
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			s := strings.ToLower(strings.TrimSpace(r.Header.Get("Accept")))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			if format := r.URL.Query().Get("format"); format == "rsc" {
				ctx := context.WithValue(r.Context(), "format", format)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else if format != "" {
				render.Render(w, r, ErrInvalidRequest(errors.New("invalid format parameter value")))
				return
			}
			ctx := context.WithValue(r.Context(), "Accept", s)
			for _, t := range cT {
				if t == s {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			w.WriteHeader(http.StatusNotAcceptable)
		}
		return http.HandlerFunc(fn)
	}
}