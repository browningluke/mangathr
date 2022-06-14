package rester

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mangathrV2/internal/logging"
	"net"
	"net/http"
	"net/url"
	"time"
)

type RESTer struct {
	client http.Client
	job    func() (interface{}, Response, error)
}

type QueryParam struct {
	Key, Value string
	Encode     bool
}

type Response struct {
	Body       []byte
	StatusCode int
	Headers    map[string][]string
}

func New() *RESTer {
	return &RESTer{
		client: http.Client{
			Transport: &http.Transport{
				MaxConnsPerHost: 50,
				MaxIdleConns:    50 * 3,
				DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig
					conn, err := tls.Dial(network, addr, tlsConfig)
					return conn, err
				},
			},
			Timeout: time.Second * 50,
		},
	}
}

func (r RESTer) Do(retries int, timeout string) interface{} {
	if retries == 0 {
		panic(errors.New("retried too many times, giving up"))
	}
	body, res, err := r.job()

	if err != nil {
		fmt.Printf("Failed, retrying... (retries left %d)\n", retries)
		logging.Debugln("Status code: ", res.StatusCode)
		logging.Debugln("Body: ", body)
		logging.Debugln("Error: ", err)

		dur, err := time.ParseDuration(timeout)
		if err != nil {
			panic(err)
		}
		time.Sleep(dur)
		return r.Do(retries-1, timeout)
	}
	return body
}

func (r RESTer) DoWithHelperFunc(retries int, timeout string, f func(res Response, err error)) interface{} {
	if retries == 0 {
		panic(errors.New("retried too many times, giving up"))
	}
	body, res, err := r.job()

	if res.StatusCode != 200 || err != nil {
		fmt.Printf("Failed, retrying... (retries left %d)\n", retries)
		logging.Debugln("Status code: ", res.StatusCode)
		logging.Debugln("Body: ", body)
		logging.Debugln("Error: ", err)

		f(res, err)
		dur, err := time.ParseDuration(timeout)
		if err != nil {
			panic(err)
		}
		time.Sleep(dur)
		return r.DoWithHelperFunc(retries-1, timeout, f)
	}

	logging.Debugln("Status code: ", res.StatusCode)
	logging.Debugln("Error: ", err)
	return body
}

// Get a URL, returns a json string. DOES NOT UNMARSHAL
func (r RESTer) Get(urlString string, headers map[string]string, params []QueryParam) *RESTer {
	r.job = func() (interface{}, Response, error) {
		res, err := r.get(urlString, headers, params)
		return string(res.Body), res, err
	}
	return &r
}

// GetBytes from a URL, returns raw bytes. DOES NOT UNMARSHAL
func (r RESTer) GetBytes(urlString string, headers map[string]string, params []QueryParam) *RESTer {
	r.job = func() (interface{}, Response, error) {
		res, err := r.get(urlString, headers, params)
		return res.Body, res, err
	}
	return &r
}

func (r RESTer) get(urlString string, headers map[string]string, params []QueryParam) (Response, error) {

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

	logging.Debugln(urlString)

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return Response{}, err
	}

	for key, element := range headers {
		req.Header.Set(key, element)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return Response{}, err
	}

	//if res.StatusCode != 200 {
	//
	//	header := res.Header
	//	remainingRetries := header["X-Ratelimit-Remaining"][0]
	//
	//	if remainingRetries == "1" {
	//
	//	}
	//
	//	return nil, errors.New("Received code: " + strconv.Itoa(res.StatusCode))
	//}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	return Response{
		StatusCode: res.StatusCode,
		Headers:    res.Header,
		Body:       body,
	}, nil
}
