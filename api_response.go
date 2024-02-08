package response

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"
)

// APIResponse defines the standard structure for all API responses.
//
// StatusCode: HTTP status code associated with the response. Not included in JSON output.
// Success: Indicates whether the request was successful (true) or not (false).
// Message: Human-readable message describing the response outcome.
// ErrorCode: Optional Application-specific error code for internal reference.
// Data: Holds the actual response data. Its type is any to allow flexibility for different data formats.
// Meta: (Optional) Holds additional information like pagination details or other metadata.
//
// Note: This struct satisfies Go's error interface, allowing it to be directly returned from functions.
type APIResponse struct {
	StatusCode int     `json:"-"`
	Success    bool    `json:"success"`
	Message    string  `json:"message"`
	ErrorCode  *string `json:"errorCode,omitempty"`
	Data       any     `json:"data,omitempty"`
	Meta       any     `json:"meta,omitempty"` // for paginations and likes
}

// Error satisfies the `error` interface by returning the response message. This enables
// using `APIResponse` as an error type and leveraging standard error handling mechanisms.
func (a *APIResponse) Error() string {
	return a.Message
}

// ToByte() encodes the response struct as a byte slice using the `gob` package.
// This can be useful for sending binary data over network connections
func (r *APIResponse) ToByte() ([]byte, error) {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ToJson() marshals the response struct to a JSON string for human-readable output.
func (r *APIResponse) ToJson() (string, error) {
	byte, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(byte), err
}

// NewAPIResponse constructs a new APIResponse object, encapsulating
// information about the API response.
//
// This function is a convenient way to build standardized API responses with all
// necessary information and avoid code repetition.
func NewAPIResponse(statusCode int, success bool, msg, errorCode string, data any) *APIResponse {
	errCode := &errorCode
	if strings.TrimSpace(errorCode) == "" {
		errCode = nil
	}

	return &APIResponse{
		StatusCode: statusCode,
		Success:    success,
		Message:    msg,
		ErrorCode:  errCode,
		Data:       data,
	}
}

// Error generates an APIResponse representing an error.
//
// return Error(http.StatusForbidden, "Access denied", "AUTH_001")
func Error(statusCode int, msg string, errorCode string) *APIResponse {
	// Check: only http status error codes are allowed.
	if statusCode < http.StatusBadRequest {
		panic("response error: cant set an error response with a non-error http status code")
	}

	return NewAPIResponse(statusCode, false, msg, errorCode, nil)
}

// Success generates an APIResponse for a successful request.
func Success(statusCode int, msg string, data any) *APIResponse {
	// Check: only http status success codes are allowed.
	if statusCode >= http.StatusBadRequest {
		panic("response error: cant set a success response with an error http status code")
	}

	if msg == "" {
		msg = "Request was successful"
	}
	return NewAPIResponse(statusCode, true, msg, "", data)
}

// Creates a api response with (HTTP 200) code
func OK(msg string, data any) *APIResponse {
	return Success(http.StatusOK, msg, data)
}

// Creates a response with (HTTP 201) code
func Created(msg string, data any) *APIResponse {
	return Success(http.StatusCreated, msg, data)
}

// Creates a success response with a list of data and meta information.
func List(msg string, data any, meta any) *APIResponse {
	if msg == "" {
		msg = "Request was successful"
	}
	rsp := NewAPIResponse(http.StatusOK, true, msg, "", data)
	rsp.Meta = meta
	return rsp
}

// Creates a response with (HTTP 400) code
func BadRequest(msg string, errorCode string) *APIResponse {
	if msg == "" {
		msg = "Request is in a bad format"
	}
	return Error(http.StatusBadRequest, msg, errorCode)
}

// Creates a response with (HTTP 401) code
func Unauthorized(msg string, errorCode string) *APIResponse {
	if msg == "" {
		msg = "Not authenticated to perform the requested action"
	}
	return Error(http.StatusUnauthorized, msg, errorCode)
}

// Creates a response with (HTTP 403) code
func Forbidden(msg string, errorCode string) *APIResponse {
	if msg == "" {
		msg = "Not authorized to perform the requested action"
	}
	return Error(http.StatusForbidden, msg, errorCode)
}

// Creates a response with (HTTP 404) code
func NotFound(msg string, errorCode string) *APIResponse {
	if msg == "" {
		msg = "Requested resource not found"
	}
	return Error(http.StatusNotFound, msg, errorCode)
}

// Creates a response with (HTTP 409) code
func Conflict(msg string, errorCode string) *APIResponse {
	if msg == "" {
		msg = "Requested resource already exist"
	}
	return Error(http.StatusConflict, msg, errorCode)
}

// creates a response with (HTTP 500)code
func InternalServerError(msg string, errorCode string) *APIResponse {
	if msg == "" {
		msg = "Something went wrong on our end."
	}
	return Error(http.StatusInternalServerError, msg, errorCode)
}

// Decodes a byte array into an APIResponse struct.
func FromJsonToAPIResponse(dataByte []byte) (*APIResponse, error) {
	var apiResponse APIResponse

	if err := json.Unmarshal(dataByte, &apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
