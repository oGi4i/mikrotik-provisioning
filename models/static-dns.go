package models

import (
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson/primitive"
	valid "mikrotik_provisioning/validate"
	"net/http"
)

type StaticDNSEntry struct {
	ID       string           `json:"-" validate:"omitempty"`
	Name     string           `json:"name" validate:"required,fqdn"`
	Regexp   string           `json:"regexp,omitempty" validate:"omitempty"`
	Address  string           `json:"address" validate:"required,ipv4"`
	TTL      RouterOSDuration `json:"ttl" validate:"required"`
	Disabled bool             `json:"disabled,omitempty" validate:"omitempty"`
	Comment  string           `json:"comment,omitempty" validate:"omitempty,comment"`
}

type StaticDNSEntryMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" validate:"required"`
	Name     string             `bson:"name" validate:"required,fqdn"`
	Regexp   string             `bson:"regexp" validate:"omitempty"`
	Address  string             `bson:"address" validate:"required,ipv4"`
	TTL      RouterOSDuration   `bson:"ttl" validate:"required"`
	Disabled bool               `bson:"disabled,omitempty" validate:"omitempty"`
	Comment  string             `bson:"comment,omitempty" validate:"omitempty,comment"`
}

type StaticDNSBatchRequest struct {
	Entries []*StaticDNSEntry `json:"entries" validate:"required"`
}

func (a *StaticDNSBatchRequest) Bind(r *http.Request) error {
	if err := valid.Validate.Struct(a); err != nil {
		return err
	}

	return nil
}

type StaticDNSRequest struct {
	*StaticDNSEntry
}

type StatisDNSResponse struct {
	*StaticDNSEntry
}

func (a *StaticDNSRequest) Bind(r *http.Request) error {
	if err := valid.Validate.Struct(a); err != nil {
		return err
	}

	return nil
}

func NewStaticDNSResponse(staticDNS *StaticDNSEntry) *StatisDNSResponse {
	return &StatisDNSResponse{StaticDNSEntry: staticDNS}
}

func (rd *StatisDNSResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListStaticDNSJSONResponse(staticDNSList []*StaticDNSEntry) []render.Renderer {
	list := make([]render.Renderer, len(staticDNSList))

	for i, staticDNS := range staticDNSList {
		list[i] = NewStaticDNSResponse(staticDNS)
	}
	return list
}
