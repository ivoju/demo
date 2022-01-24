package requestapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

// Info is the http req info
type ReqInfo struct {
	URL string
	//HeaderInfo HeaderInfoSchema
	HeaderInfo map[string]interface{}
	Body       []byte
}

type HeaderInfoSchema struct {
	ContentType   string
	Auth          string
	XBRISignature string
	XBRITimestamp string
}

type ResInfo struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// POST API request
func POST(reqinf ReqInfo, timeout time.Duration) (*ResInfo, error) {
	req, err := http.NewRequest("POST", reqinf.URL, bytes.NewReader(reqinf.Body))
	if err != nil {
		return nil, err
	}

	// set header
	for key, value := range reqinf.HeaderInfo {
		req.Header.Add(key, value.(string))
	}

	// execute
	cl := &http.Client{
		Timeout: timeout,
	}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &ResInfo{
		StatusCode: res.StatusCode,
		Body:       buf,
	}, nil
}
