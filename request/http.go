package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"sync"

	"github.com/yanun0323/errors"
	"github.com/yanun0323/logs"
)

var requestPool = sync.Pool{
	New: func() any {
		return &Request{
			Query:   make(map[string]string),
			Header:  make(map[string]string),
			BodyMap: make(map[string]any),
		}
	},
}

type Request struct {
	Ctx     context.Context
	Method  string
	Url     string
	Query   map[string]string
	Header  map[string]string
	Body    any
	BodyMap map[string]any
	Hooks   []func(Request) error
}

// New creates a new request creator.
func New(method, url string) Request {
	return Request{
		Method:  method,
		Url:     url,
		Query:   make(map[string]string),
		Header:  make(map[string]string),
		BodyMap: make(map[string]any),
	}
	// r := requestPool.Get().(Request)
	// r = r.Reset()
	// r.Method = method
	// r.Url = url
	// return r
}

func (r Request) Reset() Request {
	r.Ctx = nil
	r.Method = ""
	r.Url = ""
	r.Query = map[string]string{}
	r.Header = map[string]string{}
	r.Body = nil
	r.BodyMap = map[string]any{}
	r.Hooks = nil
	return r
}

// WithContext sets the context for the request.
func (r Request) WithContext(ctx context.Context) Request {
	r.Ctx = ctx
	return r
}

// WithQueryParam sets the query parameters for the request.
func (r Request) WithQueryParam(key, format string, args ...any) Request {
	if len(args) != 0 {
		r.Query[key] = fmt.Sprintf(format, args...)
	} else {
		r.Query[key] = format
	}
	return r
}

// WithQueryParams sets the query parameters for the request.
func (r Request) WithQueryParams(param map[string]any) Request {
	for key, value := range param {
		r.Query[key] = fmt.Sprintf("%v", value)
	}
	return r
}

// WithHeader sets the header for the request.
func (r Request) WithHeader(key, value string) Request {
	r.Header[key] = value
	return r
}

// WithHeaders sets the header for the request.
func (r Request) WithHeaders(header map[string]string) Request {
	copied := make(map[string]string, len(header)+len(r.Header))
	maps.Copy(copied, r.Header)
	maps.Copy(copied, header)
	r.Header = copied
	return r
}

// WithBody sets the body for the request.
func (r Request) WithBodyObject(p any) Request {
	r.Body = p
	return r
}

// WithBody sets the body for the request.
func (r Request) WithBodyMap(m map[string]any) Request {
	r.BodyMap = m
	return r
}

// WithHook sets the hook for the request.
func (r Request) WithHook(hook func(Request) error) Request {
	r.Hooks = append(r.Hooks, hook)
	return r
}

// Create creates a new request.
func (r Request) Create() (*http.Request, error) {
	var reader io.Reader
	if len(r.BodyMap) != 0 {
		r.Body = r.BodyMap
	}

	if r.Ctx == nil {
		r.Ctx = context.Background()
	}

	for i, hook := range r.Hooks {
		if err := hook(r); err != nil {
			return nil, errors.Wrapf(err, "execute hook(%d)", i)
		}
	}

	if r.Body != nil {
		data, err := json.Marshal(r.Body)
		if err != nil {
			return nil, errors.Wrap(err, "marshal body")
		}

		if len(r.Header["Content-Type"]) == 0 {
			r.Header["Content-Type"] = "application/json"
		}

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(r.Ctx, r.Method, r.Url, reader)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	q := url.Values{}
	for k, v := range r.Query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range r.Header {
		req.Header.Add(k, v)
	}

	return req, nil
}

// Send sends the request and returns the response.
//
// If proxy is provided, it will be used to send the request.
func (r Request) Send(delegator ...func(*http.Request) (*http.Response, error)) (*Response, error) {
	// defer requestPool.Put(r)

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

type Response struct {
	HttpResponse *http.Response

	checkStatus bool
}

func (r *Response) WithCheckStatus() *Response {
	r.checkStatus = true
	return r
}

func (r *Response) Decode(p any) error {
	defer r.HttpResponse.Body.Close()
	defer io.Copy(io.Discard, r.HttpResponse.Body)

	errTmp := errors.NewTemplate(
		"url", r.HttpResponse.Request.URL.String(),
	)

	if r.checkStatus {
		if r.HttpResponse.StatusCode < 200 || r.HttpResponse.StatusCode >= 300 {
			body, _ := io.ReadAll(r.HttpResponse.Body)
			return errTmp.Errorf("response bad status code: %d, body: %s", r.HttpResponse.StatusCode, string(body))
		}
	}

	body, err := io.ReadAll(r.HttpResponse.Body)
	if err != nil {
		return errTmp.Wrap(err, "read body")
	}

	if err := json.Unmarshal(body, p); err != nil {
		logs.Errorf("decode failed, body: %s", string(body))
		return errTmp.Wrapf(err, "unmarshal body: %s", string(body))
	}

	return nil
}
