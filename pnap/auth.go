package pnap

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// API Endpoint
//	https://providerurl.com/api
// resource path
//	/virtualmachine/55/power
// query string
//	?powerState=off
// string-to-sign
//	Create the string-to-sign for the desired operation by concatenating the HTTP verb, the Resource URL, and the Application Key.
//	String-to-sign = PUT /virtualmachine/55?powerState=off 12424d8fa15afb4
// application key - identify the account and generated through the portal
// shared secret - manually created through the portal
// request signature (SHA256)
//	Query parameters should be sorted in alphabetical order and the Request Signature should be encoded to Base 64.
//	RequestSignature = Base64 (SHA-256(String-to-Sign,sharedsecret))
// encoded credentials
//	EncodedCredentials = Base64 (ApplicationKey + “:” + RequestSignature)
// authentication header
//	Authorization = 'PNCP ' + EncodedCredentials
//	Authorization=PNCP zWM8XSTNHzzRFTVz64oBpYHVfeAgvZN/jBw=

func GetMAC(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	result := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(result)
}

type AuthContext struct {
	Method             string
	ResourcePath       string
	QueryString        string
	ApplicationKey     string
	sharedSecret       string
	StringToSign       string
	RequestSignature   string
	EncodedCredentials string
	Authenticator      string
}

func NewAuthContext(m, rp, qs, ak, ss string) AuthContext {
	c := &AuthContext{
		Method:         m,
		ResourcePath:   rp,
		QueryString:    qs,
		ApplicationKey: ak,
		sharedSecret:   ss,
	}
	c.StringToSign = fmt.Sprintf("%s %s?%s %s", c.Method, c.ResourcePath, c.QueryString, c.ApplicationKey)
	c.RequestSignature = GetMAC(c.StringToSign, c.sharedSecret)
	c.EncodedCredentials = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.ApplicationKey, c.RequestSignature)))
	c.Authenticator = fmt.Sprintf("PNCP %s", c.EncodedCredentials)
	return *c
}
