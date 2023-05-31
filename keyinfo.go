package panobi

import (
	"fmt"
	"strings"
)

const (
	idLen         int    = 22
	errInvalidKey string = "invalid key"
)

// Holds information about a signing key.
type KeyInfo struct {
	K           string // actual key
	WorkspaceID string // workspace ID
	ExternalID  string // external ID
}

// Parses the given string, and returns a KeyInfo structure holding the
// component parts.
//
// You can find your key in your Panobi workspace's integration settings.
// Keys are in the format `W-E-K`, where W is the ID of your Panobi workspace,
// E is the external ID, and K is actually the secret key generated for your
// integration.
func ParseKey(input string) (KeyInfo, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 3 {
		return KeyInfo{}, fmt.Errorf(errInvalidKey)
	}

	workspaceID := strings.TrimSpace(parts[0])
	if len(workspaceID) != idLen {
		return KeyInfo{}, fmt.Errorf(errInvalidKey)
	}

	externalID := strings.TrimSpace(parts[1])
	if len(externalID) != idLen {
		return KeyInfo{}, fmt.Errorf(errInvalidKey)
	}

	k := strings.TrimSpace(parts[2])
	if k == "" {
		return KeyInfo{}, fmt.Errorf(errInvalidKey)
	}

	return KeyInfo{
		K:           k,
		WorkspaceID: workspaceID,
		ExternalID:  externalID,
	}, nil
}

// Test for equality against the given key information.
func (ki KeyInfo) Equals(other KeyInfo) bool {
	return ki.K == other.K &&
		ki.WorkspaceID == other.WorkspaceID &&
		ki.ExternalID == other.ExternalID
}
