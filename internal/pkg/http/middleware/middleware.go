package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gopkg.in/go-playground/validator.v9"

	"mikrotik_provisioning/internal/app"
	"mikrotik_provisioning/internal/config"
	"mikrotik_provisioning/internal/pkg/address_list"
	mux "mikrotik_provisioning/internal/pkg/http"
)

type Middleware struct {
	validator *validator.Validate
	service   app.UseCases
	config    *config.Access
}

func NewMiddleware(service app.UseCases, config *config.Access) *Middleware {
	return &Middleware{validator: validator.New(), service: service, config: config}
}

func (m *Middleware) isValidAddressListRequest(request *address_list.AddressListRequest) error {
	if err := m.validator.Struct(request); err != nil {
		return err
	}

	return nil
}

func (m *Middleware) checkAccessKeys(accessKey string, secretKey string) bool {
	for _, v := range m.config.Users {
		if v.AccessKey == accessKey && v.SecretKey == secretKey {
			return true
		}
	}

	return false
}

func (m *Middleware) EnsureAddressListExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if addressListName := chi.URLParam(r, "addressListName"); addressListName != "" {
			addressList, err := m.service.GetAddressList(r.Context(), addressListName)
			if err != nil {
				_ = render.Render(w, r, mux.ErrInternalServerError(err))
				return
			}

			if addressList == nil {
				_ = render.Render(w, r, mux.ErrNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), mux.AddressListKey, addressList)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			_ = render.Render(w, r, mux.ErrNotFound)
			return
		}
	})
}

func (m *Middleware) EnsureAddressListNotExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := new(address_list.AddressListRequest)

		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := json.Unmarshal(bodyBytes, data); err != nil {
			_ = render.Render(w, r, mux.ErrInvalidRequest(err))
			return
		}

		if err := m.isValidAddressListRequest(data); err != nil {
			_ = render.Render(w, r, mux.ErrInvalidRequest(err))
			return
		}

		result, err := m.service.GetAddressList(r.Context(), data.Name)
		if err != nil {
			_ = render.Render(w, r, mux.ErrInternalServerError(err))
			return
		}

		if result != nil {
			_ = render.Render(w, r, mux.ErrInvalidRequest(fmt.Errorf("address list already exists: %s", result.Name)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) EnsureAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "" {
			authValues := strings.Split(auth, ":")
			if len(authValues) == 2 {
				accessKey := authValues[0]
				secretKey := authValues[1]
				if m.checkAccessKeys(accessKey, secretKey) {
					next.ServeHTTP(w, r)
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

func (m *Middleware) CheckAcceptHeader(contentTypes ...string) func(next http.Handler) http.Handler {
	cT := make([]string, 0)
	for _, t := range contentTypes {
		cT = append(cT, strings.ToLower(t))
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			s := strings.ToLower(strings.TrimSpace(r.Header.Get(string(mux.AcceptKey))))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			if format := r.URL.Query().Get(string(mux.FormatKey)); format == string(mux.RSCFormat) {
				ctx := context.WithValue(r.Context(), mux.FormatKey, mux.RSCFormat)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			} else if format != "" {
				_ = render.Render(w, r, mux.ErrInvalidRequest(fmt.Errorf("invalid format parameter value: %s", format)))
				return
			}

			ctx := context.WithValue(r.Context(), mux.AcceptKey, s)
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
