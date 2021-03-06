package ngerest

import (
	"crypto/rsa"
	"time"
)

// HostPublicKey Host's public key used to encrypt data transmission
type HostPublicKey struct {
	KeyString string
	PublicKey *rsa.PublicKey
	Created   time.Time
	Expired   time.Duration
}
