package rester

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RESTer struct {
	client http.Client
	job    func() (interface{}, error)
}

type QueryParam struct {
	Key, Value string
	Encode     bool
}

func New() *RESTer {
	return &RESTer{
		client: http.Client{},
	}
}

func (r RESTer) Do(retries int, timeout string) interface{} {
	if retries == 0 {
		panic(errors.New("retried too many times, giving up"))
	}
	res, err := r.job()

	if err != nil {
		fmt.Printf("Failed, retrying... (retries left %d)\n", retries)
		dur, err := time.ParseDuration(timeout)
		if err != nil {
			panic(err)
		}
		time.Sleep(dur)
		return r.Do(retries-1, timeout)
	}

	return res
}

// Get a URL, returns a json string. DOES NOT UNMARSHAL
func (r RESTer) Get(urlString string, headers map[string]string, params []QueryParam) *RESTer {
	r.job = func() (interface{}, error) {
		bytes, err := r.get(urlString, headers, params)
		return string(bytes), err
	}
	return &r
}

// GetBytes from a URL, returns raw bytes. DOES NOT UNMARSHAL
func (r RESTer) GetBytes(urlString string, headers map[string]string, params []QueryParam) *RESTer {
	r.job = func() (interface{}, error) {
		return r.get(urlString, headers, params)
	}
	return &r
}

func (r RESTer) get(urlString string, headers map[string]string, params []QueryParam) ([]byte, error) {

	if len(params) > 0 {
		urlString += "?"
	}

	for _, param := range params {
		val := param.Value

		if param.Encode {
			val = url.QueryEscape(val)
		}
		urlString += fmt.Sprintf("%s=%s&", param.Key, val)
	}

	//fmt.Println(urlString)

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}

	for key, element := range headers {
		req.Header.Set(key, element)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Received code: " + strconv.Itoa(res.StatusCode))
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	return body, nil
}
