package server_test

import (
	"awsuploader/server"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	TESTFILE = "./testdata/testfile.txt"
	TESTFILENAME = "testfile.txt"
	TESTBUCKET = "jkp-unit-tests"
)

func setupMultipartForm() (*bytes.Buffer, string, error) {
	// Setup multi-part form data
	var b bytes.Buffer
	multipartWriter := multipart.NewWriter(&b)
	// Open the test file
	testFile, err := os.Open(TESTFILE)
	if err != nil {
		fmt.Errorf("Error opening testfile: %s", err)
		return nil, "", err
	}
	defer testFile.Close()
	// Add the file form field
	formFileWriter, err := multipartWriter.CreateFormFile("file", testFile.Name())
	if err != nil {
		fmt.Errorf("Error adding file form field: %v", err)
		return nil, "", err
	}
	// Copy over test file data and close the multipart writer
	if _, err := io.Copy(formFileWriter, testFile); err != nil {
		fmt.Errorf("Error copying data to file form field: %x", err)
		return nil, "", err
	}
	multipartWriter.Close()
	header := multipartWriter.FormDataContentType()
	return &b, header, nil
}

func TestUploadHandlerRejectsGetRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/upload", nil)
	if err != nil {
		t.Logf("Could not create test request: %s", err)
	}

	rr := httptest.NewRecorder()
	testServer := server.NewServer(TESTBUCKET)
	handler := http.HandlerFunc(testServer.HandleImageUpload)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("HandleImageUpload returned wrong status code (%v) instead of expected code (%v)", status, http.StatusCreated)
	}
}

func TestUploadReturnsSuccessStruct (t *testing.T) {
	form, contentTypeHeader, err := setupMultipartForm()
	if err != nil {
		t.Errorf("Error setting up multipart form: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/upload", form)
	req.Header.Set("Content-Type", contentTypeHeader)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	testServer := server.NewServer(TESTBUCKET)
	handler := http.HandlerFunc(testServer.HandleImageUpload)
	handler.ServeHTTP(rr, req)

	// Check and make sure it gave a 201 status code
	if rr.Code != http.StatusCreated {
		t.Errorf("HandleImageUpload did not return a 201 status code for upload")
	}
	// Check and make sure the response object indicates success
	var response server.WebResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error deserializing response from HandleImageUplad: %v", err)
	}
	if response.Success != true {
		t.Errorf("HandleImageUpload did not return Success = true")
	}
}

func TestThatDuplicatesAreIdentified(t *testing.T) {
	uploadForm, contentTypeHeader, err := setupMultipartForm()
	if err != nil {
		t.Errorf("Error setting up test form for upload: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/upload", uploadForm)
	req.Header.Set("Content-Type", contentTypeHeader)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	// Send the file once to make sure it's there
	rr := httptest.NewRecorder()
	testServer := server.NewServer(TESTBUCKET)
	handler := http.HandlerFunc(testServer.HandleImageUpload)
	handler.ServeHTTP(rr, req)

	// Then check whether it's recognized as a duplicate
	if exists, err := testServer.TesthelperDuplicateCheck(TESTFILENAME); err != nil {
		t.Errorf("Error while checking for duplicate file in AWS: %v", err)
	} else if !exists {
		t.Errorf("AlreadyExists did not correctly detect the existing file")
	}
}

func TestThatSingletonsDontShowAsDuplicates(t *testing.T) {
	// File has never been sent, so should not come back as
	// a duplicate when we check
	testServer := server.NewServer(TESTBUCKET)
	if exists, err := testServer.TesthelperDuplicateCheck("JibberJabber.NotARealFile"); err != nil {
		t.Errorf("Error while checking for duplicate file in AWS: %v", err)
	} else if exists {
		t.Error("AlreadyExists incorrectly said that JibberJabber.NotARealFile is in the test bucket even though it isn't")
	}
}
