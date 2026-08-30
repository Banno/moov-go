package mhooks

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSignature(t *testing.T) {
	secret := "my-webhook-signing-secret"

	newHeaders := func(signature string) http.Header {
		headers := http.Header{}
		headers.Set("x-timestamp", "2024-04-26T21:20:55Z")
		headers.Set("x-nonce", "LwxF1Uk7QOeDq2nzB3theslHbtAo7y3uuncB1PoijwCZZaRZsd8DOtffBT7p")
		headers.Set("x-webhook-id", "dff0a709-f982-4475-81e8-214b435c74ab")
		headers.Set("x-signature", signature)
		return headers
	}

	validSignature := "6231d03752de6963087e6aea1c78a27a0617b6df1c071195f30ed85defe34e02fd0bf3995949fe12dafd747c42de9cfae03b8aafcf69cceba5495f4c7b719d82"

	isValid, err := checkSignature(newHeaders(validSignature), secret)
	require.NoError(t, err)
	require.True(t, isValid)

	for _, invalidSignature := range []string{
		"",
		"deadbeef",
		"6231d03752de6963087e6aea1c78a27a0617b6df1c071195f30ed85defe34e02fd0bf3995949fe12dafd747c42de9cfae03b8aafcf69cceba5495f4c7b719d83",
	} {
		isValid, err := checkSignature(newHeaders(invalidSignature), secret)
		require.NoError(t, err)
		require.False(t, isValid)
	}
}
