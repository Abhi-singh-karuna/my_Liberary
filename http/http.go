package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/Abhi-singh-karuna/my_Liberary/validator" // Adjust the import path based on your project structure

	"github.com/gin-gonic/gin"
)

// GetIPAddress extracts the IP address of the user from the Gin context
func GetIPAddress(c *gin.Context) string {
	return c.ClientIP()
}

// ErrResponseWithLog logs an error and sends an error response using Gin context
func ErrResponseWithLog(c *gin.Context, err error) {
	LogResponseError(c, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// LogResponseError logs an error using Gin context without sending a response
func LogResponseError(c *gin.Context, err error) {
	// Example simple logging to standard output
	log.Printf("Error occurred, IPAddress: %s, Error: %s", GetIPAddress(c), err.Error())
}

// ReadRequest reads and validates the request body using Gin context
func ReadRequest(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		return err
	}
	return validator.ValidateStruct(c.Request.Context(), request)
}

// SanitizeRequest reads, sanitizes, and validates the request body using Gin context
func SanitizeRequest(c *gin.Context, request interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	defer c.Request.Body.Close()

	sanBody, err := SanitizeJSON(body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return err
	}

	if err = json.Unmarshal(sanBody, request); err != nil {
		return err
	}

	return validator.ValidateStruct(c.Request.Context(), request)
}

// GetCountryID extracts the "country-id" from the request headers
func GetCountryID(c *gin.Context) string {
	return c.GetHeader("country-id")
}

// SanitizeJSON is a simple example of JSON sanitization
func SanitizeJSON(input []byte) ([]byte, error) {
	// Example sanitization logic - removes HTML tags from strings
	// In practice, you might use a more sophisticated sanitization library.
	re := regexp.MustCompile(`<.*?>`)
	sanitized := re.ReplaceAll(input, []byte(""))
	return sanitized, nil
}
