package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Error struct {
	Err      error
	Request  *http.Request
	Response *http.Response
}

func (e *Error) Error() string {
	req := fmt.Sprintf("%v %v:", e.Request.Method, e.Request.URL.Path)
	if e.Err != nil && e.Response != nil {
		return fmt.Sprintf("%s %d %v", req, e.Response.StatusCode, e.Err.Error())
	}
	return fmt.Sprintf("%s %v", req, e.Err.Error())
}

type Client struct {
	url  *url.URL
	http *http.Client
}

func NewClientFromHTTP(httpClient *http.Client, baseURL string) (*Client, error) {
	client := &Client{http: httpClient}

	var err error
	client.url, err = url.Parse(baseURL)

	return client, err
}

func NewClient(baseURL string) (*Client, error) {
	defaultClient := http.DefaultClient
	return NewClientFromHTTP(defaultClient, baseURL)
}

func (c *Client) NewRequest(method string, path string, query map[string]string, body interface{}) (*http.Request, error) {
	clientErr := &Error{}

	fullURL, err := c.url.Parse(path)
	if err != nil {
		clientErr.Err = err
		return nil, clientErr
	}

	if query != nil {
		q := &url.Values{}
		for k, v := range query {
			q.Add(k, v)
		}
		fullURL.RawQuery = q.Encode()
	}

	var bodyToSend io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyToSend = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(method, fullURL.String(), bodyToSend)
	if err != nil {
		clientErr.Err = err
		return nil, clientErr
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, model interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	clientError := &Error{
		Request: req,
	}

	resp, err := c.http.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			clientError.Err = ctx.Err()
		default:
			clientError.Err = err
		}
		return nil, clientError
	}
	defer resp.Body.Close()
	clientError.Response = resp

	if resp.StatusCode >= 400 {
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			contentType := resp.Header.Get("Content-Type")
			if strings.HasPrefix(contentType, "text/html") {
				clientError.Err = errors.New(string(data[:]))
			} else {
				json.Unmarshal(data, clientError)
			}

		} else {
			clientError.Err = err
		}
		return resp, clientError
	}

	if model != nil {
		bodyDec := json.NewDecoder(resp.Body)

		decErr := bodyDec.Decode(model)
		if decErr == io.EOF {
			return nil, nil
		}
		if decErr != nil {
			clientError.Err = decErr
			return resp, clientError
		}
	}

	return resp, nil
}
