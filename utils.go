package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// IsJsonErrorGetDetails checks if an error is related to JSON parsing or unmarshalling.
// It returns a boolean indicating whether it's a JSON error and a more specific error
// message if available.
//
// The function utilizes error unwrapping and type assertions to identify different
// types of JSON errors:
//
//   - json.SyntaxError: Indicates badly-formed JSON with the character offset.
//   - io.ErrUnexpectedEOF: Generic error for incomplete JSON data.
//   - json.UnmarshalTypeError: Occurs when JSON data doesn't match the expected type,
//     including specific field name and character offset.
//   - io.EOF: Triggered if the body is empty, which violates the expected JSON format.
//
// Any other error not related to JSON parsing is returned as-is.
//
// This function is useful for handling JSON requests in APIs and providing informative
// error messages to clients about the specific JSON issue encountered.
func IsJsonErrorGetDetails(err error) (ok bool, e error) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	if err != nil {
		switch {
		case errors.As(err, &syntaxError):
			return true, fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return true, errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			return true, fmt.Errorf(
				"body contains incorrect JSON type [for field %q] (at character %d)",
				unmarshalTypeError.Field,
				unmarshalTypeError.Offset,
			)

		case errors.Is(err, io.EOF):
			return true, errors.New("body must not be empty")

		default:
			return false, err
		}
	}
	return false, nil
}
