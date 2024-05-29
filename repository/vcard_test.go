//go:build test
// +build test

package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVcardString(t *testing.T) {
	vcard := VCard{Name: "Foo", Email: "bar@baz.ai"}
	s := vcard.String()
	assert.Equal(t, `BEGIN:VCARD
VERSION:3.0
FN:Foo
EMAIL:bar@baz.ai
END:VCARD`, s)
}
