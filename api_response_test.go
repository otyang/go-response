package response

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIResponse_ToByte(t *testing.T) {
	wantStruct := &APIResponse{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "its bytes",
		ErrorCode:  nil,
		Data:       nil,
	}

	var want bytes.Buffer
	err := gob.NewEncoder(&want).Encode(wantStruct)
	assert.NoError(t, err)

	got := NewAPIResponse(http.StatusOK, true, "its bytes", "", nil)
	gotBytes, err := got.ToByte()
	assert.NoError(t, err)

	assert.Equal(t, want.Bytes(), gotBytes)
}

func TestAPIResponse_ToJson(t *testing.T) {
	resp := &APIResponse{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "api message",
		ErrorCode:  nil,
		Data:       "{data",
	}

	v, err := json.Marshal(resp)
	want := string(v)
	assert.NoError(t, err)

	got, err := resp.ToJson()
	assert.NoError(t, err)
	assert.JSONEq(t, want, got)
}

func TestNewAPIResponse(t *testing.T) {
	expected := &APIResponse{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "new",
		ErrorCode:  nil,
		Data:       nil,
	}

	got := NewAPIResponse(http.StatusOK, true, "new", "", nil)
	assert.Equal(t, expected, got)
}

func TestNewError(t *testing.T) {
	errorCode := "bad_request"

	t.Run("when error type is not empty", func(t *testing.T) {
		want := &APIResponse{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "new-error",
			ErrorCode:  &errorCode,
			Data:       nil,
		}
		got := Error(http.StatusBadRequest, "new-error", "bad_request")
		assert.Equal(t, want, got)
	})

	t.Run("when error type is empty", func(t *testing.T) {
		want := &APIResponse{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "new-error",
			ErrorCode:  nil,
			Data:       nil,
		}
		got := Error(http.StatusBadRequest, "new-error", "")
		assert.Equal(t, want, got)
	})

	t.Run("when status code isnt an error http status code", func(t *testing.T) {
		assert.Panics(t, func() {
			Error(http.StatusOK, "new-error", "")
		})
	})
}

func TestNewSuccess(t *testing.T) {
	t.Run("when status code isnt a http successful status code", func(t *testing.T) {
		data := struct{ ID string }{ID: "hello world"}
		expected := &APIResponse{
			StatusCode: http.StatusOK,
			Success:    true,
			Message:    "new-success",
			ErrorCode:  nil,
			Data:       data,
		}

		got := Success(http.StatusOK, "new-success", data)
		assert.Equal(t, expected, got)
	})

	t.Run("when status code isnt a http successful status code", func(t *testing.T) {
		assert.Panics(t, func() {
			Success(http.StatusBadRequest, "", "")
		})
	})
}
