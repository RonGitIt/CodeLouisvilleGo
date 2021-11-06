package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

func makeHash(passphrase string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	hash := hasher.Sum(nil)
	return hash
}

// Encrypt encrypts a string using a provided password. This
// uses an AES-128 cipher and is intended only to provide
// simple obfuscation of the protected string (and avoid
// posting plaintext tokens and secrets to a public GitHub
// repo. It is not intended to provide robust security and
// should not be relied upon for that purpose.
//
// The ciphered string is returned, Base64 encoded.
func Encrypt(plaintext string, passphrase string) string {
	tempHash := makeHash(passphrase)
	gcmBlock, err := aes.NewCipher(tempHash)
	gcm, err := cipher.NewGCM(gcmBlock)
	if err != nil {
		log.Printf("Error setting up wrapping GCM block for encryption... problems ahead: %s", err)
	}
	nonce := make([]byte,  gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Error generating encrpytion nonce... duh duh duh: %s", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt decrypts a string that was previously encrypted with
// Encrypt(). If the incorrect password is provided, an error
// is returned.
func Decrypt(ciphertext string, passphrase string) (string, error)  {
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		log.Fatalf("Error during Base64 decoding of secret: %s", err)
	}

	tempHash :=  makeHash(passphrase)
	gcmBlock, err := aes.NewCipher(tempHash)
	if err != nil {
		log.Fatalf("Error setting up GCM block for decryption... not going to work: %s", err)
	}
	gcm, err := cipher.NewGCM(gcmBlock)
	if err != nil {
		log.Fatalf("Error setting up GCM wrapped block... not going to work: %s", err)
	}
	nonceSize := gcm.NonceSize()
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", fmt.Errorf("Error during decryption: %w", err)
	}

	return string(plaintext), nil


}