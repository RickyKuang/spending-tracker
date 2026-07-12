// Package auth handles password hashing, session tokens, and Bearer-token middleware.
package auth

// HashPassword hashes a plaintext password with bcrypt (cost 12).
func HashPassword(plaintext string) (string, error) {
	// TODO: implement
	return "", nil
}

// VerifyPassword checks a plaintext password against a bcrypt hash.
func VerifyPassword(hash, plaintext string) error {
	// TODO: implement
	return nil
}
