package panobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	maxInputBytes int = 1_048_576
)

// Holds information about a signature.
type SignatureInfo struct {
	S  string // the signature itself, calculated from a payload
	TS string // unix milliseconds at which it was calculated
}

// Calculates a signature for the given byte payload, using the given key
// information. The events endpoint requires that you include the calculated
// signature and timestamp when making requests.
func CalculateSignature(b []byte, ki KeyInfo, now *time.Time) (SignatureInfo, error) {
	if len(b) > maxInputBytes {
		return SignatureInfo{}, fmt.Errorf(errMaxNumberSize, "input", maxInputBytes, "bytes")
	}

	var ts string
	if now != nil {
		ts = fmt.Sprint(now.UnixMilli())
	} else {
		ts = fmt.Sprint(time.Now().UnixMilli())

	}

	message := fmt.Sprintf("%s:%s:%s", "v0", ts, b)
	mac := hmac.New(sha256.New, []byte(ki.K))
	mac.Write([]byte(message))
	signature := "v0=" + string(hex.EncodeToString(mac.Sum(nil)))

	return SignatureInfo{
		S:  signature,
		TS: ts,
	}, nil
}

// Test for equality against the given signature information.
func (si SignatureInfo) Equals(other SignatureInfo) bool {
	return si.S == other.S && si.TS == other.TS
}
