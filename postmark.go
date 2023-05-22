package postmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	postmarkURL = `https://api.postmarkapp.com`
)

const (
	server_token  = "server"
	account_token = "account"
)

type Client struct {
	HTTPClient   *http.Client
	ServerToken  string
	AccountToken string
	BaseURL      string
}

type parameters struct {
	Method    string
	Path      string
	Payload   interface{}
	TokenType string
}

type APIError struct {
	ErrorCode int64
	Message   string
}

type Postmark struct {
	From     string
	To       string
	Subject  string
	HtmlBody string
	TextBody string
}

func NewClient(serverToken string, accountToken string) *Client {
	return &Client{
		HTTPClient:   &http.Client{},
		ServerToken:  serverToken,
		AccountToken: accountToken,
		BaseURL:      postmarkURL,
	}
}

func (client *Client) doRequest(opts parameters, dst interface{}) error {
	url := fmt.Sprintf("%s/%s", client.BaseURL, opts.Path)

	req, err := http.NewRequest(opts.Method, url, nil)
	if err != nil {
		return err
	}

	if opts.Payload != nil {
		payloadData, err := json.Marshal(opts.Payload)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(payloadData))
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	switch opts.TokenType {
	case account_token:
		req.Header.Add("X-Postmark-Account-Token", client.AccountToken)

	default:
		req.Header.Add("X-Postmark-Server-Token", client.ServerToken)
	}

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, dst)
	return err
}

// Error returns the error message details
func (res APIError) Error() string {
	return res.Message
}
