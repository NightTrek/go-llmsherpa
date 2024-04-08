package chipper

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type LayoutPDFReader struct {
	parserAPIURL string
}

func NewLayoutPDFReader(parserAPIURL string) *LayoutPDFReader {
	return &LayoutPDFReader{
		parserAPIURL: parserAPIURL,
	}
}

func (r *LayoutPDFReader) downloadPDF(pdfURL string) (string, []byte, error) {
	// Some servers only allow browser's user_agent to download
	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"
	// Add authorization headers if using external API (see uploadPDF for an example)
	downloadHeaders := http.Header{
		"User-Agent": []string{userAgent},
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", pdfURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header = downloadHeaders

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	fileName := filepath.Base(pdfURL)
	// Note: You can change the file name here if you'd like to something else
	if resp.StatusCode == http.StatusOK {
		pdfData, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", nil, err
		}
		return fileName, pdfData, nil
	}

	return "", nil, nil
}

func (r *LayoutPDFReader) parsePDF(pdfFile string, pdfData []byte) ([]byte, error) {
	authHeader := http.Header{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", pdfFile)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(pdfData)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", r.parserAPIURL, body)
	if err != nil {
		return nil, err
	}
	req.Header = authHeader
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (r *LayoutPDFReader) ReadPDF(pathOrURL string, contents []byte) (*Document, error) {
	var pdfFile string
	var pdfData []byte
	var err error

	if contents != nil {
		pdfFile = pathOrURL
		pdfData = contents
	} else {
		parsedURL, err := url.Parse(pathOrURL)
		if err == nil && parsedURL.Scheme != "" {
			pdfFile, pdfData, err = r.downloadPDF(pathOrURL)
			if err != nil {
				return nil, err
			}
		} else {
			pdfFile = filepath.Base(pathOrURL)
			pdfData, err = os.ReadFile(pathOrURL)
			if err != nil {
				return nil, err
			}
		}
	}

	parserResponse, err := r.parsePDF(pdfFile, pdfData)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(parserResponse, &response)
	if err != nil {
		return nil, err
	}

	blocks := response["return_dict"].(map[string]interface{})["result"].(map[string]interface{})["blocks"].([]interface{})
	return NewDocument(blocks), nil
}
