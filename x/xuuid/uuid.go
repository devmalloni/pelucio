package xuuid

import (
	"bytes"

	"github.com/gofrs/uuid/v5"
)

var Empty = uuid.UUID{}

func New() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}

func ParseString(s string) (uuid.UUID, error) {
	r, err := uuid.FromString(s)
	return r, err
}

func MustParseString(s string) uuid.UUID {
	return uuid.Must(uuid.FromString(s))
}

func Equal(a, b uuid.UUID) bool {
	return bytes.Equal(a.Bytes(), b.Bytes())
}

func IsNilOrEmpty(a uuid.UUID) bool {
	return Equal(a, uuid.Nil) || Equal(a, Empty)
}

func ToStrings(u ...uuid.UUID) []string {
	var s []string
	for _, i := range u {
		s = append(s, i.String())
	}

	return s
}
