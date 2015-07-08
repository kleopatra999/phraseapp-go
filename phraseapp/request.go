package phraseapp

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func sendRequestPaginated(method, rawurl, ctype string, r io.Reader, status, page, perPage int) (io.ReadCloser, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Add("page", strconv.Itoa(page))
	query.Add("per_page", strconv.Itoa(perPage))

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}

	if ctype != "" {
		req.Header.Add("Content-Type", ctype)
	}

	resp, err := send(req, status)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func sendRequest(method, url, ctype string, r io.Reader, status int) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	if ctype != "" {
		req.Header.Add("Content-Type", ctype)
	}

	resp, err := send(req, status)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func send(req *http.Request, status int) (*http.Response, error) {
	err := authenticate(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	err = handleResponseStatus(resp, status)
	if err != nil {
		resp.Body.Close()
	}
	return resp, err
}