package shared

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// NewUUID returns a random UUID v4 string (RFC 4122).
func NewUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("generate uuid: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return formatUUID(b), nil
}

// ParseUUID validates and normalizes a UUID string to lowercase canonical form.
func ParseUUID(s string) (string, error) {
	var b [16]byte
	if err := decodeUUID(s, &b); err != nil {
		return "", err
	}
	return formatUUID(b), nil
}

func formatUUID(b [16]byte) string {
	var out [36]byte
	hex.Encode(out[0:8], b[0:4])
	out[8] = '-'
	hex.Encode(out[9:13], b[4:6])
	out[13] = '-'
	hex.Encode(out[14:18], b[6:8])
	out[18] = '-'
	hex.Encode(out[19:23], b[8:10])
	out[23] = '-'
	hex.Encode(out[24:36], b[10:16])
	return string(out[:])
}

func decodeUUID(s string, out *[16]byte) error {
	if len(s) != 36 {
		return NewError(CodeInvalid, "uuid must be 36 characters")
	}
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return NewError(CodeInvalid, "uuid has invalid format")
	}
	compact := make([]byte, 0, 32)
	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			continue
		}
		compact = append(compact, s[i])
	}
	if len(compact) != 32 {
		return NewError(CodeInvalid, "uuid has invalid format")
	}
	if _, err := hex.Decode(out[:], compact); err != nil {
		return Wrap(CodeInvalid, "uuid has invalid hex", err)
	}
	return nil
}
