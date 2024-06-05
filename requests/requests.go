package requests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type HeaderParams map[string]string

func (HeaderParams) New() HeaderParams {
	return make(map[string]string)

}

func (h HeaderParams) Has() bool {
	return len(h) > 0
}

func (h HeaderParams) AuthorizationBearer(token string) {
	h["Authorization"] = fmt.Sprintf("Bearer %s", token)

}

func (h HeaderParams) ContentType(contentType string) {
	h["Content-Type"] = contentType
}

func (h *HeaderParams) ContentTypeJson() {
	h.ContentType("application/json")
}

type Requests interface {
	Get(ctx context.Context, url string, header HeaderParams) (buf bytes.Buffer, err error)
	Post(ctx context.Context, url string, header HeaderParams, payload *bytes.Buffer) (buf bytes.Buffer, err error)
}

type requests struct{}

func New() Requests {
	return &requests{}
}

func (r requests) Get(ctx context.Context, url string, header HeaderParams) (buf bytes.Buffer, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return bytes.Buffer{}, err
	}

	if header.Has() {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	if err = r.reqError(resp); err != nil {
		return bytes.Buffer{}, err
	}

	buf, err = r.copyBody(resp)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

func (r requests) Post(ctx context.Context, url string, header HeaderParams, payload *bytes.Buffer) (buf bytes.Buffer, err error) {
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return bytes.Buffer{}, err
	}

	for key, value := range header {
		req.Header.Set(key, value)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer resp.Body.Close()

	if err = r.reqError(resp); err != nil {
		return bytes.Buffer{}, err
	}

	buf, err = r.copyBody(resp)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

func (requests) copyBody(resp *http.Response) (buf bytes.Buffer, err error) {
	if _, err = io.Copy(&buf, resp.Body); err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

func (requests) reqError(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		if wwwAuthHeader, ok := resp.Header["Www-Authenticate"]; ok {
			return fmt.Errorf("request failed: status %d, error: %+v", resp.StatusCode, wwwAuthHeader)
		}

		return fmt.Errorf("request failed: status %d, error: %+v", resp.StatusCode, resp)
	}

	return nil
}
