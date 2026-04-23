package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// TOTPConfig holds the configuration for TOTP code generation and validation.
type TOTPConfig struct {
	Digits    int
	Period    int64
	Algorithm string
}

// NewTOTP creates a new TOTP configuration with standard defaults
// compatible with Google Authenticator (6 digits, 30s period, SHA1).
func NewTOTP() *TOTPConfig {
	return &TOTPConfig{
		Digits:    6,
		Period:    30,
		Algorithm: "SHA1",
	}
}

// GenerateSecret creates a cryptographically random 20-byte secret
// encoded as a base32 string suitable for QR code provisioning.
func (t *TOTPConfig) GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	// Encode without padding for compatibility with most authenticator apps
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

// GetGoogleAuthenticatorURL generates the otpauth:// URI used for QR codes.
func (t *TOTPConfig) GetGoogleAuthenticatorURL(secret, email, issuer string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d",
		urlEncode(issuer),
		urlEncode(email),
		secret,
		urlEncode(issuer),
		t.Algorithm,
		t.Digits,
		int(t.Period),
	)
}

func urlEncode(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", "%20"), ":", "%3A")
}

// GenerateCode generates the current TOTP code for the given secret.
func (t *TOTPConfig) GenerateCode(secret string) string {
	secretBytes, err := decodeSecret(secret)
	if err != nil {
		return ""
	}

	counter := time.Now().Unix() / t.Period
	return generateHOTP(secretBytes, counter, t.Digits)
}

// GenerateCodeForTime generates a TOTP code for a specific Unix timestamp.
func (t *TOTPConfig) GenerateCodeForTime(secret string, timestamp int64) string {
	secretBytes, err := decodeSecret(secret)
	if err != nil {
		return ""
	}

	counter := timestamp / t.Period
	return generateHOTP(secretBytes, counter, t.Digits)
}

// Validate checks if the provided code matches the current TOTP code,
// allowing a 1-period window in either direction for clock skew tolerance.
func (t *TOTPConfig) Validate(secret, code string) bool {
	code = strings.TrimSpace(code)
	if len(code) != t.Digits {
		return false
	}

	now := time.Now().Unix()

	// Check current, previous, and next time windows (±30s tolerance)
	for i := int64(-1); i <= 1; i++ {
		testTime := now + (i * t.Period)
		expected := t.GenerateCodeForTime(secret, testTime)
		if constantTimeCompare(code, expected) {
			return true
		}
	}

	return false
}

// GetCurrentCode returns the current TOTP code for the given secret.
func (t *TOTPConfig) GetCurrentCode(secret string) string {
	return t.GenerateCode(secret)
}

// GetTimeRemaining returns seconds until the current TOTP code expires.
func GetTimeRemaining() int {
	return 30 - (int(time.Now().Unix()) % 30)
}

// IsTimeValid returns true if there are at least 5 seconds remaining
// in the current TOTP period (to avoid submitting an about-to-expire code).
func IsTimeValid() bool {
	remain := GetTimeRemaining()
	return remain > 5
}

// ─── Internal helpers ────────────────────────────────────────────

// decodeSecret normalizes and decodes a base32-encoded TOTP secret.
func decodeSecret(secret string) ([]byte, error) {
	// Normalize: uppercase, remove spaces
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))

	// Add padding if needed
	if m := len(secret) % 8; m != 0 {
		secret += strings.Repeat("=", 8-m)
	}

	return base32.StdEncoding.DecodeString(secret)
}

// generateHOTP implements RFC 4226 HOTP using Go's standard crypto/hmac.
func generateHOTP(secret []byte, counter int64, digits int) string {
	// Encode counter as big-endian 8 bytes
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	// HMAC-SHA1 using Go's standard library (NOT custom implementation)
	mac := hmac.New(sha1.New, secret)
	mac.Write(buf)
	hash := mac.Sum(nil)

	// Dynamic truncation per RFC 4226
	offset := hash[len(hash)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7fffffff

	otp := truncated % uint32(pow(10, digits))
	return fmt.Sprintf("%0*d", digits, otp)
}

// pow returns base^exp for integer exponents.
func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

// constantTimeCompare performs a timing-safe string comparison
// to prevent side-channel attacks on TOTP validation.
func constantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
