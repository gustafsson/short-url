//go:build test
// +build test

package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailToString(t *testing.T) {
	mailto := MailTo{Address: "bar@baz.ai", Subject: "Foo", Body: `Hello? "Bar" & "Foo": Welcome!`}
	s := mailto.String()
	assert.Equal(t, "mailto:bar%40baz.ai?body=Hello%3F+%22Bar%22+%26+%22Foo%22%3A+Welcome%21&subject=Foo", s)
}
