package rester

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type RESTer struct {
	client http.Client
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

// Get a URL, returns a json string. DOES NOT UNMARSHAL
func (r RESTer) Get(urlString string, headers map[string]string, params []QueryParam) string {

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

	fmt.Println(urlString)

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		panic(err)
	}

	for key, element := range headers {
		req.Header.Set(key, element)
	}

	res, err := r.client.Do(req)
	if err != nil {
		panic(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	//log.Printf(sb)

	return sb
}
