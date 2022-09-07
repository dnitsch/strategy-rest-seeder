package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/a8m/envsubst"
	log "github.com/dnitsch/simplelog"
	"github.com/spyzhov/ajson"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	// RoundTrip(*http.Request) (*http.Response, error)
	// http.RoundTripper
}

type SeederImpl struct {
	log    log.Loggeriface
	client Client
	header *http.Header
	auth   *Auth
}

// TODO: change this for an interface
func (r *SeederImpl) WithClient(c Client) *SeederImpl {
	r.client = c
	return r
}

func (r *SeederImpl) WithLogger(l log.Loggeriface) *SeederImpl {
	r.log = l
	return r
}

// WithHeader allows the overwrite of default Accept and Content-Type headers
// both default to `application/json`
func (r *SeederImpl) WithHeader(header *http.Header) *SeederImpl {
	h := &http.Header{}
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/json")

	for k, _ := range *header {
		h.Add(k, header.Get(k))
	}
	r.header = h
	return r
}

func (r *SeederImpl) WithAuth(auth *AuthConfig) *SeederImpl {
	r.auth = NewAuth(*auth)
	return r
}

// Action defines the single action to make agains an endpoint
// and selecting a strategy
// Endpoint is the base url to make the requests against
// GetEndpointSuffix can be used to specify a direct ID or query params
// PostEndpointSuffix
type Action struct {
	name                 string             `yaml:"-"`
	templatedPayload     string             `yaml:"-"`
	foundId              string             `yaml:"-"`
	Strategy             string             `yaml:"strategy"`
	Order                *int               `yaml:"order,omitempty"`
	Endpoint             string             `yaml:"endpoint"`
	GetEndpointSuffix    *string            `yaml:"getEndpointSuffix,omitempty"`
	PostEndpointSuffix   *string            `yaml:"postEndpointSuffix,omitempty"`
	PutEndpointSuffix    *string            `yaml:"putEndpointSuffix,omitempty"`
	DeleteEndpointSuffix *string            `yaml:"deleteEndpointSuffix,omitempty"`
	FindByJsonPathExpr   string             `yaml:"findByJsonPathExpr,omitempty"`
	PayloadTemplate      string             `yaml:"payloadTemplate"`
	Variables            map[string]any     `yaml:"variables"`
	RuntimeVars          *map[string]string `yaml:"runtimeVars,omitempty"`
}

func (a *Action) WithName(name string) *Action {
	a.name = name
	return a
}

type Status int

const (
	StatusFatal Status = iota
	StatusRetryable
)

// do performs all the network calls -
// each request object is with Context
// re-use same context with auth implementation
func (r *SeederImpl) do(req *http.Request) ([]byte, error) {
	r.log.Debug("starting request")
	r.log.Debugf("request: %+v", req)
	respBody := []byte{}
	req.Header = *r.header
	diag := &Diagnostic{HostPathMethod: fmt.Sprintf("%s: %s%s?%s", req.Method, req.URL.Host, req.URL.Path, req.URL.RawQuery), ProceedFallback: false, IsFatal: true}

	resp, err := r.client.Do(r.doAuth(req))
	if err != nil {
		r.log.Debugf("failed to make network call: %v", err)
		diag.WithStatus(999) // networkError
		diag.WithMessage(fmt.Sprintf("failed to make network call: %v", err.Error()))
		return nil, diag
	}
	defer resp.Body.Close()

	diag.WithStatus(resp.StatusCode)

	if resp.Body != nil {
		if respBody, err = io.ReadAll(resp.Body); err != nil {
			r.log.Debugf("failed to read body, closed or empty")
			diag.WithMessage("unable to read the body")
			return nil, diag
		}
	}
	// in case we need to follow redirects (shouldn't really for backend calls...)
	if resp.StatusCode > 299 {
		diag.WithMessage(string(respBody))
		diag.WithProceedFallback(true)
		diag.WithIsFatal(false)
		r.log.Debugf("resp status code: %d", resp.StatusCode)
		if resp.StatusCode >= 500 {
			r.log.Debugf("resp status code: %d, most likely indicates a service down. should not proceed to further strategies", resp.StatusCode)
			diag.WithProceedFallback(false)
		}
		return nil, diag
	}

	return respBody, nil
}

func (r *SeederImpl) doAuth(req *http.Request) *http.Request {
	enrichedReq := req
	switch r.auth.config.AuthStrategy {
	case Basic:
		enrichedReq.SetBasicAuth(r.auth.config.Username, r.auth.config.Password)
	case OAuth:
		token, err := r.auth.oAuthConfig.Token(enrichedReq.Context())
		if err != nil {
			r.log.Errorf("failed to obtain token: %q", err)
		}
		enrichedReq.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	case BasicToToken:
		// TODO: ensure custom flow similar to OAuth happens for basicToToken credentials
		return enrichedReq
	}
	return enrichedReq
}

// get makes a network call on caller defined client.Do
// returns the body as byte array
func (r *SeederImpl) get(ctx context.Context, action *Action) ([]byte, error) {
	endpoint := action.Endpoint
	if action.GetEndpointSuffix != nil {
		endpoint = fmt.Sprintf("%s%s", endpoint, *action.GetEndpointSuffix)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)

	if err != nil {
		r.log.Debugf("failed to build request: %v", err)
	}

	return r.do(req)
}

func (r *SeederImpl) post(ctx context.Context, action *Action) error {
	endpoint := action.Endpoint
	if action.PostEndpointSuffix != nil {
		endpoint = fmt.Sprintf("%s%s", endpoint, *action.PostEndpointSuffix)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(action.templatedPayload))

	if err != nil {
		r.log.Error(err)
	}

	if _, err := r.do(req); err != nil {
		return err
	}
	return nil
}

func (r *SeederImpl) put(ctx context.Context, action *Action) error {
	// create a local reference copy in each base call
	endpoint := action.Endpoint
	if action.PutEndpointSuffix != nil {
		endpoint = fmt.Sprintf("%s%s", endpoint, *action.PutEndpointSuffix)
	}
	if action.foundId != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, action.foundId)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, strings.NewReader(action.templatedPayload))

	if err != nil {
		r.log.Error(err)
	}

	if _, err := r.do(req); err != nil {
		return err
	}
	return nil
}

// deleteMethod calls base do and returns error if any failures
func (r *SeederImpl) delete(ctx context.Context, action *Action) error {
	endpoint := action.Endpoint
	if action.DeleteEndpointSuffix != nil {
		endpoint = fmt.Sprintf("%s%s", endpoint, *action.DeleteEndpointSuffix)
	}
	if action.foundId != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, action.foundId)
	}
	req, err := http.NewRequestWithContext(ctx, "DELETE", endpoint, nil)

	if err != nil {
		r.log.Error(err)
	}

	if _, err := r.do(req); err != nil {
		return err
	}
	return nil
}

// findPathByExpression lookup func using jsonpathexpression
func (r *SeederImpl) findPathByExpression(resp []byte, pathExpression string) (string, error) {
	unescStr := "" //string(resp)
	if pathExpression == "" {
		r.log.Info("no path expression provided returning empty")
		return "", nil
	}

	unescStr, e := strconv.Unquote(fmt.Sprintf("%v", string(resp)))
	if e != nil {
		r.log.Debug("using original string")
		unescStr = string(resp)
	}

	result, err := ajson.JSONPath([]byte(unescStr), pathExpression)
	if err != nil {
		r.log.Debug("failed to perform JSON path lookup - epxression failure")
		return "", err
	}

	for _, v := range result {
		switch v.Type() {
		// case int, int64, int32, int16, int8, float64, float32:
		case ajson.String:
			str, e := strconv.Unquote(fmt.Sprintf("%v", v))
			if e != nil {
				r.log.Debugf("unable to unquote value: %v returning as is", v)
				return fmt.Sprintf("%v", v), e
			}
			return str, nil
		case ajson.Numeric:
			return fmt.Sprintf("%v", v), nil
		default:
			return "", fmt.Errorf("cannot use type: %v in further processing - can only be a numeric or string value", v.Type())
		}
	}
	r.log.Infof("expression not yielded any results")
	return "", nil
}

// TODO: set up auth

// templatePayload parses input payload and replaces all $var ${var} with
// existing global env variable as well as injected from inside RestAction
// into the local context
func (r *SeederImpl) templatePayload(payload string, vars map[string]any) string {
	tmpl, err := templatePayload(payload, vars)
	if err != nil {
		r.log.Errorf("unable to parse template: %v", err)
	}
	return tmpl
}

// templatePayload parses input payload and replaces all $var ${var} with
// existing global env variable as well as injected from inside RestAction
// into the local context
func templatePayload(payload string, vars map[string]any) (string, error) {
	for k, v := range vars {
		os.Setenv(k, fmt.Sprintf("%v", v))
	}
	return envsubst.String(payload)
}
