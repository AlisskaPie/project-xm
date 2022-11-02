package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	delivery "github.com/AlisskaPie/project-xm/internal/company/delivery/http"
)

type httpClient struct {
	client http.Client
	host   string
	port   string
	api    string
	schema string
}

type ClientResponse struct {
	Body       []byte
	StatusCode int
}

func NewClient(host, port string) *httpClient {
	c := http.Client{}
	return &httpClient{
		client: c,
		host:   host,
		port:   port,
		api:    "companies",
		schema: "http",
	}
}

func ClientSetup() *httpClient {
	host := "localhost"
	port := "8000"
	return NewClient(host, port)
}

func (h httpClient) Create(postRequest delivery.CompanyPostRequest, authToken string) (*ClientResponse, error) {
	b, err := json.Marshal(&postRequest)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s://%s:%s/%s", h.schema, h.host, h.port, h.api)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	bearer := "Bearer " + authToken
	req.Header.Add("Authorization", bearer)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	clientResp := &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errResp := delivery.ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return clientResp, err
		}
		return clientResp, fmt.Errorf(errResp.Message)
	}

	return clientResp, nil
}

func (h httpClient) GetByID(id delivery.IDPathRequest) (*ClientResponse, error) {
	url := fmt.Sprintf("%s://%s:%s/%s/%s", h.schema, h.host, h.port, h.api, id.ID.String())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	clientResp := &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errResp := delivery.ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return clientResp, err
		}
		return clientResp, fmt.Errorf(errResp.Message)
	}

	return clientResp, nil
}

func (h httpClient) Delete(id delivery.IDPathRequest, authToken string) (*ClientResponse, error) {
	url := fmt.Sprintf("%s://%s:%s/%s/%s", h.schema, h.host, h.port, h.api, id.ID.String())
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + authToken
	req.Header.Add("Authorization", bearer)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	clientResp := &ClientResponse{
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errResp := delivery.ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return clientResp, err
		}
		return clientResp, fmt.Errorf(errResp.Message)
	}

	return clientResp, nil
}

func (h httpClient) Patch(patchRequest delivery.CompanyPatchRequest, authToken string) (*ClientResponse, error) {
	b, err := json.Marshal(&patchRequest)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s://%s:%s/%s/%s", h.schema, h.host, h.port, h.api, patchRequest.ID)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	bearer := "Bearer " + authToken
	req.Header.Add("Authorization", bearer)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	clientResp := &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		errResp := delivery.ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return clientResp, err
		}
		return clientResp, fmt.Errorf(errResp.Message)
	}

	return clientResp, nil
}
