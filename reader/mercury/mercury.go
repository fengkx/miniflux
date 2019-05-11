package mercury

import (
	"encoding/json"
	"fmt"

	"miniflux.app/http/client"
	"miniflux.app/url"
)

type apiResponse struct {
	Title string`json:"title"`
	Author string`json:"author"`
	PubDate string`json:"date_published"`
	Content string`json:"content"`
	URL string`json:"url"`
}

func Fetch(entryURL, mercury_api string)  (string, error) {
	params := map[string]string{
		"url": entryURL,
	}


	clt := client.New(url.AddQueryString(mercury_api, params))
	resp, err := clt.Get()
	if err != nil {
		return "", fmt.Errorf("mercury: unable to fetch %s error: %v", entryURL, err)
	}
	var apiResp apiResponse
	decode := json.NewDecoder(resp.Body)
	if err := decode.Decode(&apiResp); err != nil {
		return "", fmt.Errorf("wallabag: unable to decode token response: %v", err)
	}

	return apiResp.Content, nil
}
