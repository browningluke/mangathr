package rester

import (
	"io/ioutil"
	"log"
	"net/http"
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
func (r RESTer) Get(url string, headers map[string]string) string {
	req, err := http.NewRequest("GET", url, nil)
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
