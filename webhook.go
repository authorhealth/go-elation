package elation

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	WebhookSignatureHeader = "El8-Ed25519-Signature"
)

func VerifyWebhook(r *http.Request, publicKey []byte) error {
	sig, err := base64.StdEncoding.DecodeString(r.Header.Get(WebhookSignatureHeader))
	if err != nil {
		return fmt.Errorf("decoding Ed25519 signature: %w", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("reading request body: %w", err)
	}
	//nolint
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if !ed25519.Verify(publicKey, body, sig) {
		return errors.New("verifying signature")
	}

	return nil
}
