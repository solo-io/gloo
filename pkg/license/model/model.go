package model

import (
	"net/http"
	"time"
)

type UserInfo struct {
	Email string
}

type Request struct {
	Body    []byte
	Headers http.Header
}

type KeyInfo struct {
	Expiration time.Time
}

type RequestAndKey struct {
	Request *Request
	Key     string
}
