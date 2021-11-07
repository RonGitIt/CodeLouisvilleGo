package server_test

import (
	"awsuploader/server"
	"testing"
)

func TestThatErrorIsReturnedIfPasswordFails(t *testing.T){
	badConfig := CONFIG
	badConfig.Password = "NotTheCorrectPassword"
	if	_, err := server.NewServer(badConfig); err == nil {
		t.Errorf("NewServer did not return an erro when wrong password is provided")
	}
}