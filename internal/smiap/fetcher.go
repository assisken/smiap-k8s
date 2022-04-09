package smiap

import (
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

const defaultSmiapApiUrl = "https://smiap.ru/api/v1"

type Fetcher struct {
	url string
}

type FetcherOption func(fetcher *Fetcher)

func WithSmiapUrl(apiUrl string) FetcherOption {
	return func(fetcher *Fetcher) {
		fetcher.url = apiUrl
	}
}

func NewFetcher(options ...FetcherOption) *Fetcher {
	f := &Fetcher{
		url: defaultSmiapApiUrl,
	}

	for _, option := range options {
		option(f)
	}
	return f
}

func (f Fetcher) FetchContainers() ([]Container, error) {
	resp, err := http.Get(f.url + "/containers/")
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res []Container
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}