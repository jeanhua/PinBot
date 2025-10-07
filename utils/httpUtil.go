package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type HttpUtil struct{}

func (*HttpUtil) Request(method, url string, body io.Reader, v any) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("error status code:", resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBytes, v)
	if err != nil {
		return err
	}
	return nil
}

func (*HttpUtil) RequestWithHeader(method, url string, header http.Header, body io.Reader, v any) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	request.Header = header
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("error status code:", resp.StatusCode)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBytes, v)
	if err != nil {
		return err
	}
	return nil
}

func (*HttpUtil) RequestWithNoResponse(method, url string, body io.Reader) error {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("error status code:", resp.StatusCode)
	}
	return nil
}

func (*HttpUtil) WithJsonBody(body any) io.Reader {
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		log.Println("error when json marshal: WithJsonBody")
		bodyBytes = []byte{}
	}
	return bytes.NewReader(bodyBytes)
}

func (*HttpUtil) WithStringBody(body string) io.Reader {
	return strings.NewReader(body)
}
