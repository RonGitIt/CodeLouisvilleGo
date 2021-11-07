package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type WebResponse struct {
	Success bool
	ErrorDetails string
	S3Url string
}

// HandleImageUpload handles upload POST requests and uploads the
// provided file to an AWS S3 bucket. The file should be provided
// as a file form field in the request. Uploads whose filename
// already exists in the S3 bucket are rejected with a BadRequest
// response.
//
// On success, a StatusOk is returned along with a WebResponse
// object with a URL to the uploaded file (which would require appropriate
// credentials to access).
func (s *Server) HandleImageUpload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		var buf bytes.Buffer
		io.Copy(&buf, file)
		uploadUrl, err := s.awsS3.UploadFile(header.Filename, &buf)
		if err != nil && strings.Contains(err.Error(), "Duplicate filename"){
			resp, _ := json.Marshal(WebResponse{
				Success: false,
				ErrorDetails: err.Error(),
				S3Url: "",
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(resp))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error uploading file to aws: %s", err)))
			return
		}

		resp, _ := json.Marshal(WebResponse{
			Success: true,
			ErrorDetails: "",
			S3Url: fmt.Sprintf("Url to uploaded file: %s", uploadUrl),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resp)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// HandleGetImage handles GET requests and retrieves a file from the AWS S3
// bucket. The file to retrieve is provided via the URL path and must be
// only the filename being requested. Inclusion of any path ('/') will
// result in a MethodNotAllowed response.
//
// The file is returned in the response body along with appropriate
// Content-Disposition, Content-Type, and Content-Length headers.
func (s *Server) HandleGetImage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		requestedFile := strings.TrimPrefix(r.URL.Path, "/get/")
		if strings.Contains(requestedFile, "/") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Request may only be for a filename. Paths cannot be included"))
			return
		}

		fileData, size, err := s.awsS3.GetFile(requestedFile)
		if err != nil {
			if strings.Contains(err.Error(), "does not exist"){
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write([]byte(fmt.Sprintf("Error getting file: %v", err)))
			return
		}

		var dataLength int
		if size < 512 {
			dataLength = len(*fileData)
		} else {
			dataLength = 512
		}
		contentType := http.DetectContentType((*fileData)[:dataLength])
		w.Header().Set("Content-Disposition", "attachment; filename=" + strconv.Quote(requestedFile))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		w.Write(*fileData)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}