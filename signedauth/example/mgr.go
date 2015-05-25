package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/contrib/signedauth"
	"hash"
	"io/ioutil"
	"net/http"
	"time"
)

// StrictSHA1Manager is an example definition of an AuthKeyManager struct.
type StrictSHA1Manager struct {
	prefix string
	secret string
	key    string
	value  interface{}
}

// AuthHeaderPrefix returns the prefix used in the initialization.
func (mgr StrictSHA1Manager) AuthHeaderPrefix() string {
	return mgr.prefix
}

// SecretKey returns the secret key from the provided access key.
// Here should reside additional verifications on the header, or other parts of the request, if needed.
func (mgr StrictSHA1Manager) SecretKey(access string, req *http.Request) (string, *signedauth.Error) {
	if req.ContentLength != 0 && req.Body == nil {
		// Not sure whether net/http or Gin handles these kinds of fun situations.
		return "", &signedauth.Error{400, errors.New("Received a forged packet.")}
	}
	// Grabbing the date and making sure it's in the correct format and is within fifteen minutes.
	if dateHeader := req.Header.Get("Date"); dateHeader == "" {
		return "", &signedauth.Error{406, errors.New("No Date header provided.")}
	} else {
		date, derr := time.Parse("2006-01-02T15:04:05.000Z", dateHeader)
		if derr != nil {
			return "", &signedauth.Error{408, errors.New("Could not parse date.")}
		} else if time.Since(date) > time.Minute*15 {
			return "", &signedauth.Error{410, errors.New("Request is too old.")}
		}
	}
	// The headers look good, let's check the access key.
	// If the reading the access key requires any kind of IO (database, or file reading, etc.)
	// it's quite good to only verify if that access key is valid once all the checks are done.
	if access == "my_access_key" {
		return mgr.secret, nil
	}
	return "", &signedauth.Error{418, errors.New("You are a teapot.")}
}

// ContextKey returns the key which will store the return from ContextValue() in Gin's context.
func (mgr StrictSHA1Manager) ContextKey() string {
	return mgr.key
}

// ContextValue returns the value to store in Gin's context at ContextKey().
func (mgr StrictSHA1Manager) ContextValue(access string) interface{} {
	if access == "my_access_key" {
		return "All good with my access key!"
	}
	return "All good with any access key!"
}

// AuthHeaderRequired returns true because we want to forbid any non-signed request in this group.
func (mgr StrictSHA1Manager) AuthHeaderRequired() bool {
	return true
}

// HashFunction returns sha1.New. It could return sha512.New384 for example (SHA-1 has known theoretical attacks).
func (mgr StrictSHA1Manager) HashFunction() func() hash.Hash {
	return sha1.New
}

// DataToSign returns a string representing the data which will be HMAC'd with the secret and used to check
// authenticity of the request. This function is only called once all the parameters for the request are valid.
func (mgr StrictSHA1Manager) DataToSign(req *http.Request) (string, *signedauth.Error) {
	// In this example, we'll be implementing a similar signing method to the Amazon AWS REST one.
	// We'll use the HTTP-Verb, the MD5 checksum of the Body, if any, and the Date header in ISO format.
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
	// Note: We are returning a variety of error codes which don't follow the spec only for the purpose of testing.
	serialized_data := req.Method + "\n"
	if req.ContentLength != 0 {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", &signedauth.Error{402, errors.New("Could not read the body.")}
		}
		hash := md5.New()
		hash.Write(body)
		serialized_data += hex.EncodeToString(hash.Sum(nil)) + "\n"
	} else {
		serialized_data += "\n"
	}
	// We know from SecretKey that the Date header is present and fits our time constaints.
	serialized_data += req.Header.Get("Date")

	return serialized_data, nil
}
