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

	body, err := json.Marshal(map[string]string{
		"msg": "Hello World!",
	})
	assert.NoError(err)

	sig := ed25519.Sign(privateKey, body)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	req.Header.Set(WebhookSignatureHeader, base64.StdEncoding.EncodeToString(sig))

	err = VerifyWebhook(req, publicKey)
	assert.NoError(err)
}
