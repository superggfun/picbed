package picbed

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const (
	url = "https://api.y-alpha.com/authorization/upload"
)

func PostFile(filename string) ([]byte, error) {
	var b []byte
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		log.Println("error writing to buffer")
		return b, err
	}

	fh, err := os.Open(filename)
	if err != nil {
		log.Println("error opening file")
		return b, err
	}

	boundary := body_writer.Boundary()
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		log.Printf("Error Stating file: %s", filename)
		return b, err
	}

	req, err := http.NewRequest("POST", url, request_reader)

	if err != nil {
		log.Printf("Error post file: %s", err)
		return b, err
	}

	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return b, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}
	return body, nil
}

func Weibo(cookie string, base64 string) ([]byte, error) {
	client := &http.Client{}
	d := strings.NewReader(fmt.Sprintf("---\r\nContent-Disposition: form-data; name=\"b64_data\"\r\n\r\n%v\r\n-----", base64))
	url := `https://picupload.weibo.com/interface/pic_upload.php?ori=1&mime=image%2Fjpeg&data=base64&url=0&markpos=1&logo=&nick=0&marks=1&app=miniblog`
	req, err := http.NewRequest("POST", url, d)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "multipart/form-data; boundary=-")
	req.Header.Set("cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyText[149:], nil
}
