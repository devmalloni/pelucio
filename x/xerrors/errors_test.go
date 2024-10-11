package xerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrappedError(t *testing.T) {
	staticErr := New("Err")

	wrappedErr := New("Err2").WithError(staticErr)

	assert.True(t, errors.Is(wrappedErr, staticErr), "expected inner error on wrapped error")
	assert.True(t, errors.Is(wrappedErr, wrappedErr), "expected outer error on wrapped error")
	assert.False(t, errors.Is(staticErr, wrappedErr), "expected imutable original error")
}
