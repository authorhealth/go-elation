package elation

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	WebhookSignatureHeader = "El8-Ed25519-Signature"
)

var ErrPublicKeyLength = errors.New("incorrect length of public key")

type Event struct {
	Data          json.RawMessage `json:"data"`
	Action        string          `json:"action"`
	EventID       int64           `json:"event_id"`
	ApplicationID string          `json:"application_id"`
	Resource      string          `json:"resource"`
}

func VerifyWebhook(r *http.Request, publicKey []byte) (*Event, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, ErrPublicKeyLength
	}

	sig, err := base64.StdEncoding.DecodeString(r.Header.Get(WebhookSignatureHeader))
	if err != nil {
		return nil, fmt.Errorf("decoding Ed25519 signature: %w", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("reading request body: %w", err)
	}
	//nolint
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if !ed25519.Verify(publicKey, body, sig) {
		return nil, errors.New("verifying signature")
	}

	event := &Event{}
	err = json.Unmarshal(body, event)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling body: %w", err)
	}

	return event, nil
}
