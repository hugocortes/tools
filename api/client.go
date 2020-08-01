package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Error struct {
	Err      error
	Request  *http.Request
	Response *http.Response
}

func (e *Error) Error() string {
	req := fmt.Sprintf("%v %v:", e.Request.Method, e.Request.URL.Path)
	if e.Err != nil {
		return fmt.Sprintf("%s %v", req, e.Err.Error())
	}
	return fmt.Sprintf("%s %d %v", req, e.Response.StatusCode, e.Err.Error())
}

type Client struct {
	url  *url.URL
	http *http.Client
}

func NewClientFromHTTP(httpClient *http.Client, url string) (*Client, error) {
	client, err := NewClient(url)
	if err != nil {
		return nil, err
	}

	client.http = httpClient
	return client, nil
}

func NewClient(baseURL string) (*Client, error) {
	client := &Client{http: http.DefaultClient}

	var err error
	client.url, err = url.Parse(baseURL)

	return client, err
}

func (c *Client) NewRequest(method string, path string, query map[string]string) (*http.Request, error) {
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

	req, err := http.NewRequest(method, fullURL.String(), nil)
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
			json.Unmarshal(data, clientError)
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
