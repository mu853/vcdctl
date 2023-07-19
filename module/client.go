package module

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"
)

type VcdClient struct {
	token      string
	site       Site
	httpClient *http.Client
}

type Response struct {
	*http.Response
	Header map[string][]string
	Body   []byte
	Error  error
}

func NewVcdClient(site Site) *VcdClient {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: transportConfig,
		Timeout:   time.Duration(30) * time.Second,
	}
	vcdClient := &VcdClient{token: "", httpClient: httpClient}
	vcdClient.site = site
	return vcdClient
}

func (c *VcdClient) Login() error {
	header := map[string]string{"Authorization": "Basic " + c.site.GetCredential()}
	res := c.Request("POST", "/cloudapi/1.0.0/sessions/provider", header, nil)
	if token, ok := res.Header["X-Vmware-Vcloud-Access-Token"]; ok {
		c.token = token[0]
	} else {
		log.Fatal(res.Header, res.Body)
	}
	return nil
}

func (c *VcdClient) Request(method string, path string, header map[string]string, req_data []byte) *Response {
	// Make request
	req, err := http.NewRequest(method, c.site.Endpoint+path, bytes.NewBuffer(req_data))
	if err != nil {
		log.Fatal(err)
	}

	// Add headers
	req.Header.Set("Accept", "application/*;version=37.1")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// Get response
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal()
	}
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return &Response{res, res.Header, res_body, nil}
}
