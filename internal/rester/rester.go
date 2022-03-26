package rester

import (
	"fmt"
	"io/ioutil"
	"log"
	"mangathrV2/internal/utils"
	"net/http"
	"net/url"
)

type RESTer struct {
	client http.Client
}

func New() *RESTer {
	return &RESTer{
		client: http.Client{},
	}
}

// Get a URL, returns a json string. DOES NOT UNMARSHAL
func (r RESTer) Get(urlString string, headers map[string]string, params []utils.Tuple) string {

	if len(params) > 0 {
		urlString += "?"
	}

	for _, param := range params {
		key := param.A.(string)
		val := param.B.(string)

		if param.C.(bool) {
			val = url.QueryEscape(val)
		}
		urlString += fmt.Sprintf("%s=%s&", key, val)
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
