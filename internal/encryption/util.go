// Package encryption provides basic functionality for payload encryption
package encryption

import (
	ecies "github.com/ecies/go/v2"
)

// GenerateKeyPair generates asymmetric key pair given the encryption type
func GenerateKeyPair() (*ecies.PrivateKey, error) {
	return ecies.GenerateKey()
}

// Encrypt encrypts payload using public key
func Encrypt(payload []byte, pubKey *ecies.PublicKey) ([]byte, error) {
	return ecies.Encrypt(pubKey, payload)
}

// Decrypt decrypts a payload using provided private key
func Decrypt(payload []byte, privKey *ecies.PrivateKey) ([]byte, error) {
	return ecies.Decrypt(privKey, payload)
}
