package http

type ContextKey string

type Format string

type Action string

const (
	FormatKey      ContextKey = "format"
	AcceptKey      ContextKey = "Accept"
	AddressListKey ContextKey = "addressList"

	RSCFormat Format = "rsc"
)
