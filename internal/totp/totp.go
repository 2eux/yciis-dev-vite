package totp

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type TOTPConfig struct {
	Digits    int
	Period   int64
	Algorithm string
}

func NewTOTP() *TOTPConfig {
	return &TOTPConfig{
		Digits:    6,
		Period:   30,
		Algorithm: "SHA1",
	}
}

func (t *TOTPConfig) GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(secret), nil
}

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

func (t *TOTPConfig) GenerateCode(secret string) string {
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}

	counter := time.Now().Unix() / t.Period
	return generateHOTP(secretBytes, counter, t.Digits)
}

func generateHOTP(secret []byte, counter int64, digits int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	hmac := hmacSHA1(secret, buf)
	offset := hmac[len(hmac)-1] & 0x0f

	truncated := binary.BigEndian.Uint32(hmac[offset : offset+4])
	truncated &= 0x7fffffff

	otp := truncated % pow(10, digits)
	return fmt.Sprintf("%0*d", digits, otp)
}

func hmacSHA1(key, message []byte) []byte {
	blockSize := 64
	if len(key) > blockSize {
		h := sha1.New()
		h.Write(key)
		key = h.Sum(nil)
	}

	keyPad := make([]byte, blockSize)
	for i := range key {
		keyPad[i] = key[i] ^ 0x5c
	}

	msgPad := make([]byte, blockSize)
	for i := range key {
		msgPad[i] = key[i] ^ 0x36
	}

	if len(message) > blockSize {
		h := sha1.New()
		h.Write(message)
		message = h.Sum(nil)
	}
	message = append(msgPad[:blockSize], message...)

	h := sha1.New()
	h.Write(keyPad[:blockSize])
	h.Write(message)
	return h.Sum(nil)
}

func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

func (t *TOTPConfig) Validate(secret, code string) bool {
	code = strings.TrimSpace(code)
	if len(code) != t.Digits {
		return false
	}

	now := time.Now().Unix()
	window := t.Period

	for i := int64(-1); i <= 1; i++ {
		testTime := now + (i * window)
		expected := t.GenerateCodeForTime(secret, testTime)
		if code == expected {
			return true
		}
	}

	return false
}

func (t *TOTPConfig) GenerateCodeForTime(secret string, timestamp int64) string {
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))
	secretBytes, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}

	counter := timestamp / t.Period
	return generateHOTP(secretBytes, counter, t.Digits)
}

func (t *TOTPConfig) GetCurrentCode(secret string) string {
	return t.GenerateCode(secret)
}

func GenerateQRCodeAsBase64(otpauthURL string) (string, error) {
	return "", fmt.Errorf("QR generation requires external library like github.com/skip2/go-qrcode")
}

type TOTPUtils struct{}

func GetTimeRemaining() int {
	return 30 - (int(time.Now().Unix()) % 30)
}

func IsTimeValid() bool {
	remain := GetTimeRemaining()
	return remain > 5
}