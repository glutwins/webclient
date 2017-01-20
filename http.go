package webclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// DoGet Get获取返回
func DoGet(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post %s server err=%d %s", uri, resp.StatusCode, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)

}

// DoPost Post获取返回
func DoPost(uri string, contentType string, r io.Reader) ([]byte, error) {
	resp, err := http.Post(uri, contentType, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post %s server err=%d %s", uri, resp.StatusCode, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)

}

//PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields map[string]string, files map[string]string, uri string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	w := multipart.NewWriter(bodyBuf)

	for name, value := range fields {
		if err := w.WriteField(name, value); err != nil {
			return nil, err
		}
	}

	var fileWriter io.Writer
	var fileReader *os.File
	for name, value := range files {
		var err error
		if fileWriter, err = w.CreateFormFile(name, value); err != nil {
			return nil, err
		}

		if fileReader, err = os.Open(value); err != nil {
			return nil, err
		}

		_, err = io.Copy(fileWriter, fileReader)
		fileReader.Close()
		if err != nil {
			return nil, err
		}
	}

	w.Close()
	return DoPost(uri, w.FormDataContentType(), bodyBuf)
}
