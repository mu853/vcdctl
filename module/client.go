package module

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

func initClient() {
	initConfig()

	site, err := config.GetCurrentSite()
	if err != nil {
		Fatal(err)
	}
	client = *newVcdClient(site)

	if err := client.Login(); err != nil {
		Fatal(err)
	}
}

func newVcdClient(site Site) *VcdClient {
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

func (c *VcdClient) Login() error {
	header := map[string]string{"Authorization": "Basic " + c.site.GetCredential()}
	res := c.Request("POST", "/cloudapi/1.0.0/sessions/provider", header, nil)
	if token, ok := res.Header["X-Vmware-Vcloud-Access-Token"]; ok {
		c.token = token[0]
	} else {
		Fatal(res.Header, res.Body)
	}
	return nil
}

func (c *VcdClient) Request(method string, path string, header map[string]string, req_data []byte) *Response {
	// Make request
	req, err := http.NewRequest(method, c.site.Endpoint+path, bytes.NewBuffer(req_data))
	if err != nil {
		Fatal(err)
	}

	// Add headers
	req.Header.Set("Accept", "application/*;version=37.1")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	if isDebugMode {
		fmt.Printf("Method: %s\n", method)
		fmt.Printf("Path: %s\n", path)
		for key, value := range(header){
			fmt.Printf("Header: %s: %s\n", key, value)
		}
		fmt.Printf("Data: %s\n", bytes.NewBuffer(req_data))
	}

	// Get response
	res, err := c.httpClient.Do(req)
	if err != nil {
		Fatal(err)
	}
	res_body, err := io.ReadAll(res.Body)
	if err != nil {
		Fatal(err)
	}
	return &Response{res, res.Header, res_body, nil}
}
