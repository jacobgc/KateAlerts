package internal

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type IFTTT struct {
	endpoint string
}

func NewIFTTT(endpoint string) *IFTTT {
	return &IFTTT{
		endpoint: endpoint,
	}
}

func (i *IFTTT) TriggerEndpoint(value1, value2, value3 string) {
	data := url.Values{
		"value1": {value1},
		"value2": {value2},
		"value3": {value3},
	}

	r, err := http.PostForm(i.endpoint, data)
	if err != nil {
		panic(err)
	}
	response, _ := ioutil.ReadAll(r.Body)
	println(string(response))

}
