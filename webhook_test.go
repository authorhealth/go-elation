package elation

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	assert := assert.New(t)

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	assert.NoError(err)

	event := &Event{
		Data:          []byte(`{"foo":"bar"}`),
		Action:        "action",
		EventID:       1,
		ApplicationID: "application-id",
		Resource:      "resource",
	}
	b, err := json.Marshal(event)
	assert.NoError(err)

	sig := ed25519.Sign(privateKey, b)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	req.Header.Set(WebhookSignatureHeader, base64.StdEncoding.EncodeToString(sig))

	actualEvent, err := VerifyWebhook(req, publicKey)
	assert.Equal(event, actualEvent)
	assert.NoError(err)
}

func TestWebhook_incorrect_public_key_len(t *testing.T) {
	assert := assert.New(t)

	req := httptest.NewRequest(http.MethodPost, "/", nil)

	event, err := VerifyWebhook(req, []byte("foo"))
	assert.Nil(event)
	assert.ErrorIs(err, ErrPublicKeyLength)
}
