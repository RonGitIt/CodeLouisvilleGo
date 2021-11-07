package server_test

import (
	"awsuploader/server"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	TESTFILE = "./testdata/testfile.txt"
	TESTFILENAME = "testfile.txt"
	TESTBUCKET = "jkp-unit-tests"
	ENCRYPTEDID = "Dx7ZgyFXUyv4x5c5vqijXwuDuj1Ej43u2jd/gcAgzIyaPCvSQchGF9ukO0zW55sq"
	ENCRYPTEDSECRET = "cJCSnSvm9S0bDFTkcm06Iw8HXivyewQaCydsdCzXQOBIDs8nud9Ab6EBPCFzcgmT6jv7HhFAfP+VMpZiWdo8QbPUP5E="
)

var (
	CONFIG = server.AwsConfig {
		Bucket: TESTBUCKET,
		Password: TESTPASSWORD,
		Id: ENCRYPTEDID,
		Secret: ENCRYPTEDSECRET,
	}

	TESTFILEMD5 []byte
)

func getMd5(input io.Reader) ([]byte, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, input); err != nil {
		return nil, fmt.Errorf("Error copying data while computing MD5 hash: %w", err)
	}

	return hash.Sum(nil), nil
}

func setupMultipartForm(nameOfFile string) (*bytes.Buffer, string, error) {
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
	var fileNameForForm string
	if nameOfFile == "" {
		fileNameForForm = testFile.Name()
	} else {
		fileNameForForm = nameOfFile
	}
	formFileWriter, err := multipartWriter.CreateFormFile("file", fileNameForForm)
	if err != nil {
		return nil, "", fmt.Errorf("Error adding file form field: %v", err)
	}

	// Copy over test file data and close the multipart writer
	if _, err := io.Copy(formFileWriter, testFile); err != nil {
		return nil, "", fmt.Errorf("Error copying data to file form field: %x", err)
	}

	// Set the MD5 hash of the file for comparison later
	testFile.Seek(0, io.SeekStart)
	TESTFILEMD5, err = getMd5(testFile)
	if err != nil {
		return nil, "", err
	}

	// Done
	multipartWriter.Close()
	header := multipartWriter.FormDataContentType()
	return &b, header, nil
}

func setupFileUploadRequest(nameOfFile string) (*http.Request, error) {
	uploadForm, contentTypeHeader, err := setupMultipartForm(nameOfFile)
	if err != nil {
		log.Printf("Error setting up 9ZfHw94LP6P4jnXBMRCUhKFpj+5Z82x3vOajVaecsZ4PTFXcM1o5XGxonLAS4dT+GJajSwUv1zGw82LXdMR6IqiyTio=test form for upload: %s", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "/upload", uploadForm)
	req.Header.Set("Content-Type", contentTypeHeader)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	return req, nil
}

func TestUploadHandlerRejectsGetRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/upload", nil)
	if err != nil {
		t.Logf("Could not create test request: %s", err)
	}

	rr := httptest.NewRecorder()
	testServer, err := server.NewServer(CONFIG)
	if err != nil {
		t.Errorf("Error creating server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleImageUpload)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("HandleImageUpload returned wrong status code (%v) instead of expected code (%v)", status, http.StatusCreated)
	}
}

func TestUploadReturnsSuccessStruct (t *testing.T) {
	req, err := setupFileUploadRequest(TESTFILENAME)
	if err != nil {
		t.Errorf("Error setting up file upload request: %v", err)
	}


	rr := httptest.NewRecorder()
	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleImageUpload)
	handler.ServeHTTP(rr, req)

	// Check and make sure it gave a 201 status code
	if rr.Code != http.StatusCreated {
		t.Errorf("HandleImageUpload did not return a 201 status code for upload")
		t.Logf("It returned %v", rr.Code)
	}
	// Check and make sure the response object indicates success
	var response server.WebResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error deserializing response from HandleImageUplad: %v", err)
	}
	if response.Success != true {
		t.Errorf("HandleImageUpload did not return Success = true")
		t.Logf("It returned %v", response.Success)
		t.Logf("Error description: %s", response.ErrorDetails)
	}

	testServer.TesthelperDeleteFile(TESTFILENAME)
}

func TestThatDuplicatesAreIdentified(t *testing.T) {
	req, err := setupFileUploadRequest("")
	if err != nil {
		t.Errorf("Error setting up file upload request: %v", err)
	}


	// Send the file once to make sure it's there
	rr := httptest.NewRecorder()
	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleImageUpload)
	handler.ServeHTTP(rr, req)

	// Then check whether it's recognized as a duplicate
	if exists, err := testServer.TesthelperDuplicateCheck(TESTFILENAME); err != nil {
		t.Errorf("Error while checking for duplicate file in AWS: %v", err)
	} else if !exists {
		t.Errorf("AlreadyExists did not correctly detect the existing file")
	}

	// Remove file that was uploaded
	_, err = testServer.TesthelperDeleteFile(TESTFILENAME)
	if err != nil {
		t.Logf("Error during cleanup: Could not delete object (%s) from test bucket: %s", TESTFILENAME, err)
	}
}

func TestThatSingletonsDontShowAsDuplicates(t *testing.T) {
	// File has never been sent, so should not come back as
	// a duplicate when we check
	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	if exists, err := testServer.TesthelperDuplicateCheck("JibberJabber.NotARealFile"); err != nil {
		t.Errorf("Error while checking for duplicate file in AWS: %v", err)
	} else if exists {
		t.Error("AlreadyExists incorrectly said that JibberJabber.NotARealFile is in the test bucket even though it isn't")
	}
}

func TestThatDuplicateFilenameUploadIsRejected(t *testing.T) {
	filename := fmt.Sprintf("%x", rand.Int63n(100000000000))
	req, err := setupFileUploadRequest(filename)
	if err != nil {
		t.Errorf("Error setting up file upload request: %v", err)
	}
	t.Logf("Testing duplicate upload with filename %s", filename)

	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleImageUpload)

	// Upload it once
	handler.ServeHTTP(httptest.NewRecorder(), req)

	// Try uploading it a second time--it should reject it and not overwrite
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleImageUpload should have responded with Bad Request on duplicate attempt, but it didn't")
	}
	var resp *server.WebResponse
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Errorf("Error decoding json response after duplicate upload attempt: %s", err)
	}
	// Check that Webresponse.Success is false
	if resp.Success != false {
		t.Errorf("Success was not false after duplicate upload attempt")
	}
	// Check that Webresponse.ErrorDetails indicates that file was a duplicate
	if !strings.Contains(strings.ToLower(resp.ErrorDetails), "duplicate filename") {
		t.Errorf("Error message not correct after duplicate upload attempt. Error message provided: %s", resp.ErrorDetails)
	}

	// Clean up... delete uploaded file from S3
	_, err = testServer.TesthelperDeleteFile(filename)
	if err != nil {
		t.Logf("Error during cleanup... could not delete object (%s) from test bucket: %s", filename, err)
	}
}

func TestGetFile(t *testing.T) {
	// Upload file that we're going to test getting
	filename := fmt.Sprintf("%x", rand.Int63n(100000000000))
	req, err := setupFileUploadRequest(filename)
	if err != nil {
		t.Errorf("Error setting up file upload request: %v", err)
	}
	t.Logf("Testing duplicate upload with filename %s", filename)

	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleImageUpload)

	// Upload it
	handler.ServeHTTP(httptest.NewRecorder(), req)

	// Download it
	handler = http.HandlerFunc(testServer.HandleGetImage)
	req, _ = http.NewRequest(http.MethodGet, "/get/" + filename, nil)
	respRecorder := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder, req)

	// Check the downloaded file!
	result := respRecorder.Result()
	if result.StatusCode != http.StatusOK {
		t.Logf("HandleGetImage did not return expected status code (%v). Got %v",
			http.StatusOK, result.StatusCode)
		t.Fail()
	}

	if !strings.Contains(result.Header["Content-Disposition"][0], "attachment") {
		t.Logf("HandleGetImage did not set the Content-Disposition header")
		t.Fail()
	}

	returnedMd5, err := getMd5(result.Body)
	if err != nil {
		t.Logf("Error getting MD5 hash of returned file: %v", err)
		t.Fail()
	} else if !bytes.Equal(returnedMd5, TESTFILEMD5) {
		t.Logf("MD5 hash of returned file does not match MD5 of original file")
		t.Fail()
	}

	// Delete the created file
	testServer.TesthelperDeleteFile(filename)
}

func TestGetFileRejectsPaths(t *testing.T){
	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleGetImage)
	req, _ := http.NewRequest(http.MethodGet, "/get/file/with/path", nil)
	respRecorder := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder, req)

	if respRecorder.Result().StatusCode != http.StatusBadRequest {
		t.Logf("HandleGetImage did not return a BadRequest when sent a filename with a path")
		t.Fail()
	}
}

func TestGetFileRejectsPostRequest(t *testing.T) {
	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleGetImage)
	req, _ := http.NewRequest(http.MethodPost, "/get/filename.txt", nil)
	respRecorder := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder, req)

	if respRecorder.Result().StatusCode != http.StatusMethodNotAllowed{
		t.Logf("HandleGetImage did not return a BadRequest when a POST request is made")
		t.Fail()
	}
}

func TestFileNotFoundReturns404(t *testing.T) {
	testServer, err := server.NewServer(CONFIG)
	if err !=  nil {
		t.Errorf("Error setting up server: %v", err)
	}
	handler := http.HandlerFunc(testServer.HandleGetImage)
	req, _ := http.NewRequest(http.MethodGet, "/get/fileIsNotHere.Nope", nil)
	respRecorder := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder, req)

	if respRecorder.Result().StatusCode != http.StatusNotFound{
		t.Logf("HandleGetImage did not return a 404 when file is not found. Gave %v",
			respRecorder.Result().StatusCode)
		t.Fail()
	}
}