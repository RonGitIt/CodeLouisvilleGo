package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WebResponse struct {
	Success bool
	ErrorDetails string
	S3Url string
}

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
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error uploading file to aws: %s", err)))
		}

		response := WebResponse{
			Success: true,
			ErrorDetails: "",
			S3Url: fmt.Sprintf("Url to uploaded file: %s", uploadUrl),
		}
		responseString, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("Error marshalling web return struct: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseString)


	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	/*defer func() {fmt.Println("Exiting HandleImageUpload")}()

	outputDir := "/home/jon/tmp"
	switch r.Method {
	case http.MethodPost:
		r.ParseMultipartForm(32 << 20)

		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error getting file from request: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()
		io.Copy(&buf, file)

		outputFilename := filepath.Join(outputDir, header.Filename)
		outputFile, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0666)
		defer outputFile.Close()
		io.Copy(outputFile, &buf)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("File size: %v\n", header.Size)))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}*/
}
