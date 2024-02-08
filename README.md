# Go-response

## API Response Library for Go

**Overview:**

This library defines the `APIResponse` struct and functions for building and handling standardized responses for API applications.

**Key Features:**

* **Structured Response Format:** Encapsulates status code, success flag, message, error code (optional), data, and meta information in a single struct.
* **Error Handling:** Supports error handling with built-in functions for generating error responses and identifying JSON-related errors.
* **Multiple Encoding Options:** Responds with JSON (default) or binary data (using `gob` encoding) depending on the context.
* **Convenience Functions:** Provides helper functions for creating common response types like `OK`, `Created`, `List`, `BadRequest`, etc.

**Components:**

* `APIResponse`: The core struct for representing API responses.
* `NewAPIResponse`: Constructs a new APIResponse with specific details.
* `Error`: Creates an error response with provided status code, message, and error code.
* `Success`: Creates a success response with status code, message, and data.
* `OK`, `Created`, `List`, etc.: Convenient functions for specific response types.
* `FromJsonToAPIResponse`: Decodes a JSON byte array into an APIResponse object.
* `IsJsonErrorGetDetails`: Checks if an error is related to JSON parsing and provides details.

**Benefits:**

* Standardized response format improves clarity and consistency across different API endpoints.
* Error handling functions simplify error response construction and provide informative messages.
* Response encoding options offer flexibility for different data types and transmission formats.
* Convenience functions reduce boilerplate code and improve code efficiency.

**Usage:**

```go
 package main

import (
	"fmt"
	"net/http"

	response "github.com/otyang/go-response"
)

func main() {
	// Create an API response
	apiResponse := response.NewAPIResponse(http.StatusOK, true, "Request processed successfully", "", nil)

	// Success response with data
	apiResponse = response.Success(http.StatusOK, "Data retrieved successfully", "mo")

	// Create an error response
	apiResponse = response.Error(http.StatusBadRequest, "Invalid request parameters", "INVALID_PARAMS_001")

	// Send JSON response
	json, err := apiResponse.ToJson()
	if err != nil {
		// handle error
	}

	// Send binary response (e.g., for large data)
	data, err := apiResponse.ToByte()
	// handle error
	w.WriteHeader(apiResponse.StatusCode)
	w.Write(data)

	// Check if an error is related to JSON parsing
	isJsonError, details := response.IsJsonErrorGetDetails(err)
	if isJsonError {
		fmt.Printf("Error: %s\n", details)
	}

	// Shortcut methods
	_ = response.BadRequest("message", "INVALID_PHONE_NUMBER")
	_ = response.Unauthorized("message", "ERROR_CODE")
	_ = response.Forbidden("message", "ERROR_CODE")
	_ = response.NotFound("message", "ERROR_CODE")
	_ = response.Conflict("message", "ERROR_CODE")
	_ = response.InternalServerError("message", "ERROR_CODE")
	_ = response.OK("message", "data")
	_ = response.List(
		"Data fetched successfully",
		[]map[string]any{{"name": "John", "age": 30}},
		map[string]any{"total": 10, "page": 1},
	)
	_ = response.Success(200, "message", "data")
}
```

# go-response
