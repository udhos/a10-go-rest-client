package a10go

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func tlsConfig() *tls.Config {
	return &tls.Config{
		//CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA, tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		//MaxVersion:               tls.VersionTLS11,
		//MinVersion:               tls.VersionTLS11,
	}
}

func httpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    tlsConfig(),
		DisableCompression: true,
		DisableKeepAlives:  true,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
}

func httpPostString(url, contentType, s string) ([]byte, error) {
	return httpPost(url, contentType, bytes.NewBufferString(s))
}

func httpDeleteString(url, contentType, s string) ([]byte, error) {
	return httpDelete(url, contentType, bytes.NewBufferString(s))
}

func httpPost(url, contentType string, body io.Reader) ([]byte, error) {
	c := httpClient()
	return clientPost(c, url, contentType, body)
}

func httpGet(url string) ([]byte, error) {
	c := httpClient()
	return clientGet(c, url)
}

func httpDelete(url string, contentType string, body io.Reader) ([]byte, error) {
	c := httpClient()
	return clientDelete(c, url, contentType, body)
}

func clientDelete(c *http.Client, url, bodyContentType string, body io.Reader) ([]byte, error) {
	return clientMethod(c, "DELETE", url, bodyContentType, body)
}

func clientMethod(c *http.Client, method, url, bodyContentType string, body io.Reader) ([]byte, error) {

	req, errNew := http.NewRequest(method, url, body)
	if errNew != nil {
		return nil, errNew
	}
	req.Header.Set("Content-Type", bodyContentType)

	resp, errDel := c.Do(req)
	if errDel != nil {
		return nil, errDel
	}

	defer resp.Body.Close()

	info, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return info, fmt.Errorf("http method=%s: read all: url=%v: %v", method, url, errRead)
	}

	if resp.StatusCode != 200 {
		return info, fmt.Errorf("http method=%s: bad status: %d", method, resp.StatusCode)
	}

	return info, nil
}

func clientPost(c *http.Client, url string, contentType string, r io.Reader) ([]byte, error) {

	resp, errPost := c.Post(url, contentType, r)
	if errPost != nil {
		return nil, errPost
	}

	defer resp.Body.Close()

	info, errBody := ioutil.ReadAll(resp.Body)
	if errBody != nil {
		return info, fmt.Errorf("httpPost: read: url=%v: %v", url, errBody)
	}

	if resp.StatusCode != 200 {
		return info, fmt.Errorf("httpPost: bad status: %d", resp.StatusCode)
	}

	return info, nil
}

func clientGet(c *http.Client, url string) ([]byte, error) {
	resp, errGet := c.Get(url)
	if errGet != nil {
		return nil, fmt.Errorf("httpGet: get url=%v: %v", url, errGet)
	}

	defer resp.Body.Close()

	info, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return info, fmt.Errorf("httpGet: read all: url=%v: %v", url, errRead)
	}

	if resp.StatusCode != 200 {
		return info, fmt.Errorf("httpGet: bad status: %d", resp.StatusCode)
	}

	return info, nil
}
