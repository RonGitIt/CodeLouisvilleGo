package server_test

import (
	"awsuploader/server"
	"testing"
)

const (
	TESTPLAINTEXT = "This is the plaintext"
	TESTPASSWORD = "SuperStrongPassword987"
	TESTID = "ThisIsTheTestId"
	TESTSECRET = "ThisIsTheTestSecret"

)

func TestEncryptAndDecrypt(t *testing.T) {
	ciphertext := server.Encrypt(TESTPLAINTEXT, TESTPASSWORD)
	t.Logf("Encrypted text: %s", ciphertext)
	plaintext, err := server.Decrypt(ciphertext, TESTPASSWORD)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Logf("Decrypted text: %s", plaintext)
	if plaintext != TESTPLAINTEXT {
		t.Logf("Decrypted text (%s) doesn match expected text (%s)", plaintext, TESTPLAINTEXT)
		t.Fail()
	}
}

func TestDecryptionWithBadPassword(t *testing.T) {
	ciphertext := server.Encrypt(TESTPLAINTEXT, TESTPASSWORD)
	//t.Logf("Encrypted text: %s", ciphertext)
	_ , err := server.Decrypt(ciphertext, "This is not the password")
	if err == nil {
		t.Logf("Attempt to decrypt with wrong password should have returned an error but didn't")
		t.Fail()
	}
}
