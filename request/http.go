package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/yanun0323/errors"
)

var requestPool = sync.Pool{
	New: func() any {
		return &Request{
			query:   make(map[string]string),
			header:  make(map[string]string),
			bodyMap: make(map[string]any),
		}
	},
}

// Request represents constructor for http.Request.
type Request struct {
	ctx     context.Context
	method  string
	url     string
	query   map[string]string
	header  map[string]string
	body    any
	bodyMap map[string]any
	hooks   []func(*Request) error
}

// New creates a new request creator.
func New(method, url string) *Request {
	r := requestPool.Get().(*Request)
	r.Reset()
	r.method = method
	r.url = url
	return r
}

// Reset reset all parameters inside Request
func (r *Request) Reset() *Request {
	r.ctx = nil
	r.method = ""
	r.url = ""
	r.query = map[string]string{}
	r.header = map[string]string{}
	r.body = nil
	r.bodyMap = map[string]any{}
	r.hooks = nil
	return r
}

// WithContext sets the context for the request.
func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// WithQueryParam sets the query parameters for the request.
func (r *Request) WithQueryParam(key, format string, args ...any) *Request {
	if len(args) != 0 {
		r.query[key] = fmt.Sprintf(format, args...)
	} else {
		r.query[key] = format
	}
	return r
}

// WithQueryParams sets the query parameters for the request.
func (r *Request) WithQueryParams(param map[string]any) *Request {
	for key, value := range param {
		r.query[key] = fmt.Sprintf("%v", value)
	}
	return r
}

// WithHeader sets the header for the request.
func (r *Request) WithHeader(key, value string) *Request {
	r.header[key] = value
	return r
}

// WithHeaders sets the header for the request.
func (r *Request) WithHeaders(header map[string]string) *Request {
	for key, value := range header {
		r.header[key] = value
	}
	return r
}

// WithBody sets the body for the request.
//
// It will use body map which set from WithBodyMap first.
func (r *Request) WithBodyObject(p any) *Request {
	r.body = p
	return r
}

// WithBodyMap sets the body for the request.
//
// It will use body map which set from WithBodyMap first.
func (r *Request) WithBodyMap(m map[string]any) *Request {
	r.bodyMap = m
	return r
}

// WithHook sets the hook for the request.
func (r *Request) WithHook(hook func(*Request) error) *Request {
	r.hooks = append(r.hooks, hook)
	return r
}

// Create creates a http.Request from the request parameters.
func (r *Request) Create() (*http.Request, error) {
	var reader io.Reader
	if len(r.bodyMap) != 0 {
		r.body = r.bodyMap
	}

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	for i, hook := range r.hooks {
		if err := hook(r); err != nil {
			return nil, errors.Wrapf(err, "execute hook(%d)", i)
		}
	}

	if r.body != nil {
		data, err := json.Marshal(r.body)
		if err != nil {
			return nil, errors.Wrap(err, "marshal body")
		}

		if len(r.header["Content-Type"]) == 0 {
			r.header["Content-Type"] = "application/json"
		}

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(r.ctx, r.method, r.url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	q := url.Values{}
	for k, v := range r.query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range r.header {
		req.Header.Add(k, v)
	}

	return req, nil
}

// Send sends the request and returns the response.
//
// If proxy is provided, it will be used to send the request.
func (r *Request) Send(delegator ...func(*http.Request) (*http.Response, error)) (*Response, error) {
	defer requestPool.Put(r)

	do := http.DefaultClient.Do
	if len(delegator) != 0 {
		do = delegator[0]
	}

	req, err := r.Create()
	if err != nil {
		return nil, err
	}

	resp, err := do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	return &Response{HttpResponse: resp}, nil
}

// Response is a wrapper of http.Response.
type Response struct {
	HttpResponse *http.Response

	checkStatus bool
}

// WithCheckStatus makes Response ensure the status code between 200 and 2XX, otherwise the Decode function will return an error.
func (r *Response) WithCheckStatus() *Response {
	r.checkStatus = true
	return r
}

// Decode decodes body of http.Response into the object of pointer you provided,
// then close the body.
func (r *Response) Decode(p any) error {
	defer r.HttpResponse.Body.Close()
	defer io.Copy(io.Discard, r.HttpResponse.Body)

	errTmp := errors.NewTemplate(
		"url", r.HttpResponse.Request.URL.String(),
	)

	if r.checkStatus {
		if r.HttpResponse.StatusCode < 200 || r.HttpResponse.StatusCode >= 300 {
			return errTmp.Errorf("response bad status code: %d", r.HttpResponse.StatusCode)
		}
	}

	body, err := io.ReadAll(r.HttpResponse.Body)
	if err != nil {
		return errTmp.Wrap(err, "read body")
	}

	if err := json.Unmarshal(body, p); err != nil {
		return errTmp.Wrapf(err, "unmarshal body: %s", string(body))
	}

	return nil
}
