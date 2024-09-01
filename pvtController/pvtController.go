package pvtcontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// AuthUser structure (deprecated)
type AuthUser struct {
	Email    string
	UserID   string
	RegionID uint
}

// Context interface with methods similar to what was provided by Echo
type Context interface {
	Path() string
	Request() *http.Request
	String(code int, s string) error
	JSON(code int, i interface{}) error
	NoContent(code int) error
	Get(key string) interface{}
	Set(key string, val interface{})
	Authorization() (string, error)
	BindBody(interface{}) error
	BodyBytes() ([]byte, error)
	BodyString() (string, error)
	GetLogger() Logger

	// Deprecated methods
	SetAuthUser(*AuthUser)
	GetAuthUser() *AuthUser
}

// Simple context implementation
type simpleContext struct {
	request  *http.Request
	response http.ResponseWriter
	store    map[string]interface{}
	logger   Logger
	authUser *AuthUser
}

func NewContext(w http.ResponseWriter, r *http.Request, logger Logger) Context {
	return &simpleContext{
		request:  r,
		response: w,
		store:    make(map[string]interface{}),
		logger:   logger,
	}
}

// Path returns the path of the request
func (c *simpleContext) Path() string {
	return c.request.URL.Path
}

// Request returns the HTTP request
func (c *simpleContext) Request() *http.Request {
	return c.request
}

// String sends a string response
func (c *simpleContext) String(code int, s string) error {
	c.response.WriteHeader(code)
	_, err := c.response.Write([]byte(s))
	return err
}

// JSON sends a JSON response
func (c *simpleContext) JSON(code int, i interface{}) error {
	c.response.Header().Set("Content-Type", "application/json")
	c.response.WriteHeader(code)
	return json.NewEncoder(c.response).Encode(i)
}

// NoContent sends a response with no content
func (c *simpleContext) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

// Get retrieves a value from the context store
func (c *simpleContext) Get(key string) interface{} {
	return c.store[key]
}

// Set stores a value in the context store
func (c *simpleContext) Set(key string, val interface{}) {
	c.store[key] = val
}

// Authorization retrieves the Authorization header
func (c *simpleContext) Authorization() (string, error) {
	auth := c.request.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("authorization header missing")
	}
	return auth, nil
}

// BindBody reads the request body and binds it to the provided interface
func (c *simpleContext) BindBody(i interface{}) error {
	body, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return err
	}
	defer c.request.Body.Close()

	// Replacing the request body with the read bytes for further use
	c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return json.Unmarshal(body, i)
}

// BodyBytes returns the request body as bytes
func (c *simpleContext) BodyBytes() ([]byte, error) {
	body, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return nil, err
	}
	defer c.request.Body.Close()

	// Replacing the request body with the read bytes for further use
	c.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

// BodyString returns the request body as a string
func (c *simpleContext) BodyString() (string, error) {
	body, err := c.BodyBytes()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetLogger returns the logger
func (c *simpleContext) GetLogger() Logger {
	return c.logger
}

// SetAuthUser stores the AuthUser in the context (deprecated)
func (c *simpleContext) SetAuthUser(user *AuthUser) {
	c.authUser = user
}

// GetAuthUser retrieves the AuthUser from the context (deprecated)
func (c *simpleContext) GetAuthUser() *AuthUser {
	return c.authUser
}

// Logger interface with basic logging methods
type Logger interface {
	Errorf(format string, args ...interface{})
}
