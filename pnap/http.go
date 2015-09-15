package pnap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Generic call
func (r *PNAP) call(method, path, qs string, inBody interface{}) (out Future, emsg string, retriable bool, eref uint64) {
	authContext := NewAuthContext(method, path, qs, r.ApplicationKey, r.SharedSecret)

	if r.Debug {
		json, err := json.Marshal(authContext)
		if err != nil {
			log.Println("Unable to marshal AuthContext.")
		} else {
			log.Printf(`AuthContext: %s`, json)
		}
	}

	url := fmt.Sprintf("%s%s%s", r.Endpoint, path, qs)
	reqBody, err := json.Marshal(inBody)
	if err != nil {
		emsg = err.Error()
		retriable = false
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Accept", documentType)
	req.Header.Set("Content-Type", documentType)
	req.Header.Set("Authorization", authContext.Authenticator)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		rawOut, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			emsg = err.Error()
			retriable = true
			return
		}
		out = &SyncResponse{body: rawOut}
	} else if resp.StatusCode == 202 {
		out = &AsyncResponse{}
		rawOut, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			emsg = err.Error()
			retriable = true
			return
		}
		err = json.Unmarshal(rawOut, &out)
		if err != nil {
			emsg = err.Error()
			retriable = true
			return
		}
	} else {
		if resp.StatusCode == 500 {
			emsg = resp.Header.Get("X-Application-Error-Description")
			eref, _ = strconv.ParseUint(resp.Header.Get("X-Application-Error-Description"), 0, 0)
			retriable = true
		} else if resp.StatusCode == 400 {
			emsg = resp.Header.Get("X-Application-Error-Description")
			eref, _ = strconv.ParseUint(resp.Header.Get("X-Application-Error-Description"), 0, 0)
			retriable = false
		} else {
			emsg = resp.Status
			retriable = false
		}
	}

	return
}
