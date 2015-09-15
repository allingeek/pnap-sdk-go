package pnap

import (
	"encoding/json"
	"errors"
	"time"
)

type Future interface {
	Get(r interface{}) (err error)
	TimedGet(r interface{}, ttl time.Duration) (err error)
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\
// Synchrounous Implementation
//\\//\\//\\//\\//\\//\\//\\//\\//\\

type SyncResponse struct {
	body []byte
}

func (sr SyncResponse) Get(r interface{}) (err error) {
	// Unmarshal into given type now
	err = json.Unmarshal(sr.body, r)
	return
}
func (sr SyncResponse) TimedGet(r interface{}, ttl time.Duration) (err error) {
	return sr.Get(r)
}

//\\//\\//\\//\\//\\//\\//\\//\\//\\
// Asynchronous Implementation
//\\//\\//\\//\\//\\//\\//\\//\\//\\

type AsyncResponse struct {
	ResourceURL string `json:"resourceURL"`
	response    *SyncResponse
	api         *PNAP
}

func (ar AsyncResponse) Get(r interface{}) (err error) {
	for {
		if ar.response != nil {
			return ar.response.Get(r)
		}

		result := &Task{}
		out, emsg, retriable, _ := ar.api.call(`GET`, ar.ResourceURL, ``, ``)
		if emsg != `` && !retriable {
			err = errors.New(emsg)
		}

		out.Get(result) // Unmarshal the task response (we know its a synchronous call)
		if result.RequestStateEnum == `CLOSED_SUCCESSFUL` || result.RequestStateEnum == `CLOSED_FAILED` {
			ar.response = out.(*SyncResponse)
			return ar.response.Get(r)
		}
		time.Sleep(ar.api.Backoff)
	}
}
func (ar AsyncResponse) TimedGet(r interface{}, ttl time.Duration) (err error) {
	// TODO: Implement timeout
	return nil
}
