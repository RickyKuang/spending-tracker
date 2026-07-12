// Package crypto provides AES-256-GCM sealing for secrets at rest, such as Plaid access tokens.
package crypto

// Sealer seals and opens ciphertext using a master key (APP_ENCRYPTION_KEY).
type Sealer struct {
	key []byte
}

// NewSealer builds a Sealer from a 32-byte AES-256 key.
func NewSealer(key []byte) *Sealer {
	return &Sealer{key: key}
}

// Seal encrypts plaintext with AES-256-GCM, returning nonce||ciphertext.
func (s *Sealer) Seal(plaintext []byte) ([]byte, error) {
	// TODO: implement
	return nil, nil
}

// Open decrypts a nonce||ciphertext value produced by Seal.
func (s *Sealer) Open(ciphertext []byte) ([]byte, error) {
	// TODO: implement
	return nil, nil
}
