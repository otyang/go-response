package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsJsonErrorGetDetails(t *testing.T) {
	t.Run("json syntax Type Error", func(t *testing.T) {
		v1 := &json.SyntaxError{
			Offset: 1,
		}

		ok, err := IsJsonErrorGetDetails(v1)
		assert.Error(t, err)
		assert.Equal(t, true, ok)
	})

	t.Run("json unmarshal Type Error", func(t *testing.T) {
		v2 := &json.UnmarshalTypeError{
			Offset: 90, Field: "field",
		}

		wantErr := fmt.Errorf(
			"body contains incorrect JSON type [for field %q] (at character %d)",
			v2.Field,
			v2.Offset,
		)
		ok, gotErr := IsJsonErrorGetDetails(v2)
		assert.Equal(t, true, ok)
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("error is io.EOF", func(t *testing.T) {
		wantErr := errors.New("body must not be empty")
		ok, gotErr := IsJsonErrorGetDetails(io.EOF)
		assert.Equal(t, true, ok)
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("non json error", func(t *testing.T) {
		wantErr := errors.New("just a normal error")
		ok, gotErr := IsJsonErrorGetDetails(errors.New("just a normal error"))
		assert.Equal(t, false, ok)
		assert.Equal(t, wantErr, gotErr)
	})

	t.Run("no error: nil", func(t *testing.T) {
		ok, err := IsJsonErrorGetDetails(nil)
		assert.NoError(t, err)
		assert.Equal(t, false, ok)
	})
}
