package brale

import (
	"encoding/json"
	"fmt"
)

type BraleClient interface {
	Mint(float64, string, string) (*APIResponse, error)
	Redeem(float64, string, string) (*APIResponse, error)
}

type braleClient struct {
	BaseURL string
	Jwt     string
}

func NewBraleClient(baseURL string, jwt string) *braleClient {
	return &braleClient{
		BaseURL: baseURL,
		Jwt:     jwt,
	}
}

func (bc *braleClient) Mint(amount float64, recipient string, idem string) (*APIResponse, error) {
	return &APIResponse{}, nil
}

func (bc *braleClient) Redeem(amount float64, recipient string, idem string) (*APIResponse, error) {
	return &APIResponse{}, nil
}

type mockBraleClient struct {
}

func NewMockBraleClient() *mockBraleClient {
	return &mockBraleClient{}
}

func (m *mockBraleClient) Mint(amount float64, recipient string, idem string) (*APIResponse, error) {
	// idem would be used here to prevent double spends since
	if recipient == "0xdeadbeef" {
		errResp, err := m.loadErrorResponse()
		if err != nil {
			return nil, err
		}
		return errResp, err
	}
	return m.loadSuccessResponse()
}

func (m *mockBraleClient) Redeem(amount float64, recipient string, idem string) (*APIResponse, error) {
	// idem would be used here to prevent double spends
	if recipient == "0xdeadbeef" {
		errResp, err := m.loadErrorResponse()
		if err != nil {
			return nil, err
		}
		return errResp, err
	}
	return m.loadSuccessResponse()
}

func (m *mockBraleClient) loadSuccessResponse() (*APIResponse, error) {
	data := `{
		"data": {
		  "attributes": {
			"created": "2020-01-01T12:00:00Z",
			"status": "pending",
			"type": "mint",
			"updated": "2020-01-01T12:00:00Z"
		  },
		  "id": "2VZvtmVc2j3gQ80CTlcuQXbGrwC",
		  "type": "order"
		}
	}
	`

	var response APIResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal success response: %v", err)
	}
	return &response, nil
}

func (m *mockBraleClient) loadErrorResponse() (*APIResponse, error) {
	data := `
	{
		"errors": [
		  {
			"code": "ValidationError",
			"detail": "An error occurred with the request data.",
			"id": "error123",
			"links": {
			  "additionalProp": "/some-resource/2VZvtmVc2j3gQ80CTlcuQXbGrwC"
			},
			"meta": {
			  "additionalProp": "string"
			},
			"source": {
			  "parameter": "page[cursor]",
			  "pointer": "/body/data/attributes"
			},
			"status": "400",
			"title": "A validation error occurred"
		  }
		],
		"links": {
		  "additionalProp": {
			"href": "/some-resource/2VZvtmVc2j3gQ80CTlcuQXbGrwC"
		  }
		},
		"meta": {
		  "additionalProp": "string"
		}
	}`

	var response APIResponse
	if err := json.Unmarshal([]byte(data), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal error response: %v", err)
	}
	return &response, nil
}
