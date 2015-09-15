package pnap

import (
	"encoding/json"
)

type Response interface {
	GetResult(r interface{}) (err error)
}

type SyncResponse struct {
	body []byte
}

func (sr *SyncResponse) GetResult(r interface{}) (err error) {
	// Unmarshal into given type now
	err = json.Unmarshal(sr.body, r)
	return
}

type AsyncResponse struct {
	ResourceURL string `json:"resourceURL"`
	response    SyncResponse
}

func (ar *AsyncResponse) GetResult(r interface{}) (err error) {
	// do polling and such here
	return
}
