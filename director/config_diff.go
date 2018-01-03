package director

import (
	"encoding/json"
	"net/http"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type ConfigDiff struct {
	Diff [][]interface{}
}

type ConfigDiffResponse struct {
	Diff [][]interface{} `json:"diff"`
}

func NewConfigDiff(diff [][]interface{}) ConfigDiff {
	return ConfigDiff{
		Diff: diff,
	}
}

func (c Client) postConfigDiff(path string, manifest []byte, setHeaders func(*http.Request)) (ConfigDiffResponse, error) {
	var resp ConfigDiffResponse

	respBody, response, err := c.clientRequest.RawPost(path, manifest, setHeaders)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			if strings.Contains(err.Error(), "\"code\": 440012") {
				// config couldn't be found => return only the director error description
				prefixRemoved := strings.TrimPrefix(err.Error(),
					"Director responded with non-successful status code '404' response '{\"code\": 440012, \"description\": \"")
				errorDescription := strings.TrimSuffix(prefixRemoved, "\"}'")
				return resp, bosherr.Errorf(errorDescription)
			} else {
				// endpoint couldn't be found => return empty diff, just for compatibility with directors which don't have the endpoint
				return resp, nil
			}
		} else {
			return resp, bosherr.WrapErrorf(err, "Fetching diff result")
		}
	}

	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return resp, bosherr.WrapError(err, "Unmarshaling Director response")
	}

	return resp, nil
}
