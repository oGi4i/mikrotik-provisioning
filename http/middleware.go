package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"io/ioutil"
	cfg "mikrotik_provisioning/config"
	"mikrotik_provisioning/pkg"
	"mikrotik_provisioning/types"
	valid "mikrotik_provisioning/validate"
	"net/http"
	"strings"
)

func EnsureAddressListExists(i *pkg.Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			addressList := new(types.AddressList)
			var err error

			ctx := context.Background()
			if addressListName := chi.URLParam(r, "addressListName"); addressListName != "" {
				addressList, err = i.Storage.GetAddressListByName(ctx, addressListName)
			} else {
				render.Render(w, r, types.ErrNotFound)
				return
			}
			if err != nil {
				render.Render(w, r, types.ErrInternalServerError(err))
				return
			}

			c := context.WithValue(ctx, "addressList", addressList)
			next.ServeHTTP(w, r.WithContext(c))
		}
		return http.HandlerFunc(fn)
	}
}

func EnsureAddressListNotExists(i *pkg.Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			data := new(types.AddressListRequest)

			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.Unmarshal(bodyBytes, data); err != nil {
				render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}
			if err := valid.Validate.Struct(data); err != nil {
				render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}

			ctx := context.Background()
			result, err := i.Storage.GetAddressListByName(ctx, data.Name)
			if err != nil {
				render.Render(w, r, types.ErrInternalServerError(err))
				return
			}
			if result != nil {
				render.Render(w, r, types.ErrInvalidRequest(errors.New(fmt.Sprintf("address list already exists: %s", result.Name))))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func EnsureAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			authValues := strings.Split(auth, ":")
			if len(authValues) == 2 {
				accessKey := authValues[0]
				secretKey := authValues[1]
				for _, v := range cfg.Config.Access.Users {
					if v.AccessKey == accessKey && v.SecretKey == secretKey {
						next.ServeHTTP(w, r)
					}
				}
			} else {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	})
}

func CheckAcceptHeader(contentTypes ...string) func(next http.Handler) http.Handler {
	cT := make([]string, 0)
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
				render.Render(w, r, types.ErrInvalidRequest(errors.New(fmt.Sprintf("invalid format parameter value: %s", format))))
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

func EnsureStaticDNSEntriesNotExist(i *pkg.Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			data := new(types.StaticDNSBatchRequest)

			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.Unmarshal(bodyBytes, data); err != nil {
				render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}
			if err := valid.Validate.Struct(data); err != nil {
				render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}

			ctx := context.Background()
			results, err := i.Storage.GetAllStaticDNS(ctx)
			if err != nil {
				render.Render(w, r, types.ErrInternalServerError(err))
				return
			}
			if results != nil {
				for _, entry := range data.Entries {
					for _, v := range results {
						if v.Name == entry.Name {
							render.Render(w, r, types.ErrInvalidRequest(errors.New(fmt.Sprintf("statis DNS entry already exists: %s", v.Name))))
							return
						}
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func EnsureStaticDNSEntriesExist(i *pkg.Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			data := new(types.StaticDNSBatchRequest)

			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.Unmarshal(bodyBytes, data); err != nil {
				render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}
			if err := valid.Validate.Struct(data); err != nil {
				render.Render(w, r, types.ErrInvalidRequest(err))
				return
			}

			ctx := context.Background()
			results, err := i.Storage.GetAllStaticDNS(ctx)
			if err != nil {
				render.Render(w, r, types.ErrInternalServerError(err))
				return
			}
			if results != nil {
				for _, entry := range data.Entries {
					var exists bool
					for _, v := range results {
						if v.Name == entry.Name {
							exists = true
						}
					}
					if !exists {
						render.Render(w, r, types.ErrInvalidRequest(errors.New(fmt.Sprintf("statis DNS entry not exists: %s", entry.Name))))
						return
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func EnsureStaticDNSEntryExists(i *pkg.Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := new(types.StaticDNSEntry)
			var err error

			ctx := context.Background()
			if staticDNSName := chi.URLParam(r, "staticDNSName"); staticDNSName != "" {
				entry, err = i.Storage.GetStaticDNSEntryByName(ctx, staticDNSName)
			} else {
				render.Render(w, r, types.ErrNotFound)
				return
			}
			if err != nil {
				render.Render(w, r, types.ErrInternalServerError(err))
				return
			}

			c := context.WithValue(ctx, "staticDNSEntry", entry)
			next.ServeHTTP(w, r.WithContext(c))
		}
		return http.HandlerFunc(fn)
	}
}

func EnsureStaticDNSEntryNotExists(i *pkg.Implementation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := new(types.StaticDNSEntry)
			var err error

			if staticDNSName := chi.URLParam(r, "staticDNSName"); staticDNSName != "" {
				entry, err = i.Storage.GetStaticDNSEntryByName(r.Context(), staticDNSName)
			} else {
				bodyBytes, _ := ioutil.ReadAll(r.Body)
				r.Body.Close()
				r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				if err := json.Unmarshal(bodyBytes, entry); err != nil {
					render.Render(w, r, types.ErrInvalidRequest(err))
					return
				}
				if err := valid.Validate.Struct(entry); err != nil {
					render.Render(w, r, types.ErrInvalidRequest(err))
					return
				}

				entries, err := i.Storage.GetAllStaticDNS(r.Context())
				if err != nil {
					render.Render(w, r, types.ErrInternalServerError(err))
					return
				}

				for _, e := range entries {
					if e.Name == entry.Name {
						render.Render(w, r, types.ErrInvalidRequest(errors.New(fmt.Sprintf("statis DNS entry already exists: %s", e.Name))))
						return
					}
				}
			}
			if err != nil {
				render.Render(w, r, types.ErrInternalServerError(err))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
