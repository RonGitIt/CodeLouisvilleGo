package server_test

import (
	"awsuploader/server"
	"testing"
)

const (
	AWSTESTID = "jbkzZ96bw4hKLlChAm8qlnGWarqm/fuSKsltedfzNg4iVor2uPE61TMAEg=="
	AWSTESTSECRET = "Z5nwcnzaKiw7l6cBTQs81dD8xvDJJ9BRvbIvi7WJ4nPIAEa4lVmkGzMMTSiHVCU="
)

func TestDecryptsSecretOnNewAws(t *testing.T) {
	config := server.AwsConfig{
		Bucket: TESTBUCKET,
		Password: TESTPASSWORD,
		Id: AWSTESTID,
		Secret: AWSTESTSECRET,
	}
	myAws, err := server.NewAws(config)
	if err != nil {
		t.Logf("Failed to create AWS struct: %v", err)
		t.Fail()
	}

	if myAws.Id != TESTID {
		t.Logf("Test ID not decrypted correctly. Was %s Expected %s", myAws.Id, TESTID)
		t.Fail()
	}
	if myAws.Secret != TESTSECRET {
		t.Logf("Test secret not decrypted correctly. Was %s Expected %s", myAws.Secret, TESTSECRET)
		t.Fail()
	}
}
