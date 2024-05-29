package handler

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProp(t *testing.T) {
    q := url.Values{}
    q.Set("key1", "value1")

    result := prop(q, "key1")
    assert.Equal(t, "value1", result)

    result = prop(q, "key2")
    assert.Equal(t, "", result)
}

func TestIntProp(t *testing.T) {
    q := url.Values{}
    q.Set("key1", "123")

    result, err := intProp(q, "key1", 0)
    assert.NoError(t, err)
    assert.Equal(t, 123, result)

    result, err = intProp(q, "key2", 456)
    assert.NoError(t, err)
    assert.Equal(t, 456, result)

    q.Set("key3", "abc")
    _, err = intProp(q, "key3", 0)
    assert.Error(t, err)
}

func TestFloatProp(t *testing.T) {
    q := url.Values{}
    q.Set("key1", "123.45")

    result, err := floatProp(q, "key1", 0.0)
    assert.NoError(t, err)
    assert.Equal(t, 123.45, result)

    result, err = floatProp(q, "key2", 456.78)
    assert.NoError(t, err)
    assert.Equal(t, 456.78, result)

    q.Set("key3", "abc")
    _, err = floatProp(q, "key3", 0.0)
    assert.Error(t, err)
}
