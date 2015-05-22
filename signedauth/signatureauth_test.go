package signedauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
func (mgr StrictSHA1Manager) SecretKey(access string, req *http.Request) (string, *Error) {
	if req.ContentLength != 0 && req.Body == nil {
		// Not sure whether net/http or Gin handles these kinds of fun situations.
		return "", &Error{400, errors.New("Received a forged packet.")}
	}
	// Grabbing the date and making sure it's in the correct format and is within fifteen minutes.
	if dateHeader := req.Header.Get("Date"); dateHeader == "" {
		return "", &Error{406, errors.New("No Date header provided.")}
	} else {
		date, derr := time.Parse("2006-01-02T15:04:05.000Z", dateHeader)
		if derr != nil {
			return "", &Error{408, errors.New("Could not parse date.")}
		} else if time.Since(date) > time.Minute*15 {
			return "", &Error{410, errors.New("Request is too old.")}
		}
	}
	// The headers look good, let's check the access key.
	// If the reading the access key requires any kind of IO (database, or file reading, etc.)
	// it's quite good to only verify if that access key is valid once all the checks are done.
	if access == "my_access_key" {
		return mgr.secret, nil
	}
	return "", &Error{418, errors.New("You are a teapot.")}
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
func (mgr StrictSHA1Manager) DataToSign(req *http.Request) (string, *Error) {
	// In this example, we'll be implementing a similar signing method to the Amazon AWS REST one.
	// We'll use the HTTP-Verb, the MD5 checksum of the Body, if any, and the Date header in ISO format.
	// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
	// Note: We are returning a variety of error codes which don't follow the spec only for the purpose of testing.
	serialized_data := req.Method + "\n"
	if req.ContentLength != 0 {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", &Error{402, errors.New("Could not read the body.")}
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

// TestExtractAuthInfo tests the correct extraction of information from the headers.
func TestExtractAuthInfo(t *testing.T) {
	// https://github.com/smartystreets/goconvey/wiki#get-going-in-25-seconds
	Convey("Given a static manager with prefix SAUTH", t, func() {
		mgr := StrictSHA1Manager{prefix: "SAUTH", key: "contextKey", secret: "super-secret-password", value: nil}

		Convey("When the header has an incorrect prefix", func() {
			accesskey, signature, err := extractAuthInfo(mgr, "INCORRECT Something:ThereWasASpace")
			Convey("Accesskey and signature should be empty strings", func() {
				So(accesskey, ShouldEqual, "")
				So(signature, ShouldEqual, "")
			})
			Convey("The error should be a 401 with a specific message.", func() {
				So(err.Status, ShouldEqual, 401)
				So(err.Err.Error(), ShouldEqual, "Invalid authorization header.")
			})
		})

		Convey("When the header has the correct prefix but more than one space", func() {
			accesskey, signature, err := extractAuthInfo(mgr, "SAUTH Something ThereWasASpace")
			Convey("Accesskey and signature should be empty strings", func() {
				So(accesskey, ShouldEqual, "")
				So(signature, ShouldEqual, "")
			})
			Convey("The error should be a 401 with a specific message.", func() {
				So(err.Status, ShouldEqual, 401)
				So(err.Err.Error(), ShouldEqual, "Invalid authorization header.")
			})
		})

		Convey("When the header has the correct prefix but missing the seperation colon", func() {
			accesskey, signature, err := extractAuthInfo(mgr, "SAUTH SomethingThereIsNoSepColon")
			Convey("Accesskey and signature should be empty strings", func() {
				So(accesskey, ShouldEqual, "")
				So(signature, ShouldEqual, "")
			})
			Convey("The error should be a 401 with a specific message.", func() {
				So(err.Status, ShouldEqual, 401)
				So(err.Err.Error(), ShouldEqual, "Invalid format for access key and signature.")
			})
		})

		Convey("When the header is valid", func() {
			accesskey, signature, err := extractAuthInfo(mgr, "SAUTH SomeAccessKey:SomeSignature")
			Convey("Accesskey and signature should be extracted correctly", func() {
				So(accesskey, ShouldEqual, "SomeAccessKey")
				So(signature, ShouldEqual, "SomeSignature")
			})
			Convey("The error should be nil.", func() {
				So(err, ShouldEqual, nil)
			})
		})
	})
}

// TestMiddleware tests the whole signature auth middleware behavior.
func TestMiddleware(t *testing.T) {

	Convey("Given a strict manager", t, func() {
		mgr := StrictSHA1Manager{prefix: "SAUTH", key: "contextKey", secret: "super-secret-password", value: nil}
		router := gin.Default()
		router.Use(SignatureAuth(mgr))
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
		for _, meth := range methods {
			router.Handle(meth, "/test/", []gin.HandlerFunc{func(c *gin.Context) {
				c.String(http.StatusOK, "Success.")
			}})
		}
		Convey("When there is no header", func() {
			for _, meth := range methods {
				Convey(fmt.Sprintf("and doing a %s request", meth), func() {
					req := performRequest(router, meth, "/test/", nil, nil)
					Convey("the middleware should respond forbidden", func() {
						So(req.Code, ShouldEqual, 401)
					})
				})
			}
		})

		Convey("When the header has an incorrect prefix", func() {
			headers := make(map[string][]string)
			headers["Authorization"] = []string{"INCORRECT Something:ThereWasASpace"}
			for _, meth := range methods {
				Convey(fmt.Sprintf("and doing a %s request", meth), func() {
					req := performRequest(router, meth, "/test/", headers, nil)
					Convey("the middleware should respond unauthorized", func() {
						So(req.Code, ShouldEqual, 401)
					})
				})
			}
		})

		Convey("When the header has the correct prefix and date headers but an incorrect access key", func() {
			headers := make(map[string][]string)
			headers["Authorization"] = []string{"SAUTH Something:ThereWasASpace"}
			headers["Date"] = []string{time.Now().Format("2006-01-02T15:04:05.000Z")}
			for _, meth := range methods {
				Convey(fmt.Sprintf("and doing a %s request", meth), func() {
					req := performRequest(router, meth, "/test/", headers, nil)
					Convey("the middleware should respond with the Manager's secret key provided status.", func() {
						So(req.Code, ShouldEqual, 418)
					})
				})
			}
		})

		Convey("When the header has the correct prefix and access key, but incorrect signature", func() {
			headers := make(map[string][]string)
			headers["Authorization"] = []string{"SAUTH my_access_key:InvalidSignature"}

			Convey("And missing the Date header", func() {
				for _, meth := range methods {
					Convey(fmt.Sprintf("and doing a %s request", meth), func() {
						req := performRequest(router, meth, "/test/", headers, nil)
						Convey("the middleware should respond as requested by the manager.", func() {
							So(req.Code, ShouldEqual, 406)
						})
					})
				}
			})

			Convey("And the Date header is in the incorrect format", func() {
				headers["Date"] = []string{time.Now().Format("01/02 03 04 05 06")}
				for _, meth := range methods {
					Convey(fmt.Sprintf("and doing a %s request", meth), func() {
						req := performRequest(router, meth, "/test/", headers, nil)
						Convey("the middleware should respond as requested by the manager.", func() {
							So(req.Code, ShouldEqual, 408)
						})
					})
				}
			})

			Convey("And the Date header is valid but too old", func() {
				utc, _ := time.LoadLocation("UTC")
				old_date := time.Date(2006, 05, 04, 03, 02, 01, 00, utc)
				headers["Date"] = []string{old_date.Format("2006-01-02T15:04:05.000Z")}
				for _, meth := range methods {
					Convey(fmt.Sprintf("and doing a %s request", meth), func() {
						req := performRequest(router, meth, "/test/", headers, nil)
						Convey("the middleware should respond as requested by the manager.", func() {
							So(req.Code, ShouldEqual, 410)
						})
					})
				}
			})

			Convey("And the Date header is completely valid", func() {
				headers["Date"] = []string{time.Now().Format("2006-01-02T15:04:05.000Z")}
				for _, meth := range methods {
					Convey(fmt.Sprintf("and doing a %s request", meth), func() {
						req := performRequest(router, meth, "/test/", headers, nil)
						Convey("the middleware should respond unauthorized.", func() {
							So(req.Code, ShouldEqual, 401)
						})
					})
				}
			})

		})

		Convey("When the full signature is valid with no body.", func() {
			headers := make(map[string][]string)
			now := time.Now().Format("2006-01-02T15:04:05.000Z")
			headers["Date"] = []string{now}
			for _, meth := range methods {
				Convey(fmt.Sprintf("and doing a %s request", meth), func() {
					sig_data := meth + "\n\n" + now
					hash := hmac.New(sha1.New, []byte(mgr.secret))
					hash.Write([]byte(sig_data))
					signature := hex.EncodeToString(hash.Sum(nil))
					headers["Authorization"] = []string{"SAUTH my_access_key:" + signature}
					req := performRequest(router, meth, "/test/", headers, nil)
					Convey("the middleware should respond 200 OK.", func() {
						So(req.Code, ShouldEqual, 200)
					})
				})
			}
		})

		Convey("When the full signature is valid with a body.", func() {
			headers := make(map[string][]string)
			secret := "super-secret-password"
			now := time.Now().Format("2006-01-02T15:04:05.000Z")
			headers["Date"] = []string{now}
			for _, meth := range []string{"POST", "PUT"} {
				Convey(fmt.Sprintf("and doing a %s request", meth), func() {
					body := "This is the  body of my request."
					bhash := md5.New()
					bhash.Write([]byte(body))
					sig_data := meth + "\n" + hex.EncodeToString(bhash.Sum(nil)) + "\n" + now
					hash := hmac.New(sha1.New, []byte(secret))
					hash.Write([]byte(sig_data))
					signature := hex.EncodeToString(hash.Sum(nil))
					headers["Authorization"] = []string{"SAUTH my_access_key:" + signature}
					req := performRequest(router, meth, "/test/", headers, bytes.NewBufferString(body))
					Convey("the middleware should respond 200 OK.", func() {
						So(req.Code, ShouldEqual, 200)
					})
				})
			}
		})

	})
}

// performRequest is a helper to test requests, based on https://github.com/gin-gonic/gin/blob/c467186d2004be8ade88a35f5bcf71cc2c676635/routes_test.go#L19.
func performRequest(r http.Handler, method, path string, headers map[string][]string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header = headers
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
