package analysis

import (
	"errors"
	"net/url"
	"net/http"
  "io/ioutil"
  "encoding/json"
	"fmt"
	"bytes"
)
const apiKey = os.Getenv("ALCHEMY_KEY")
const baseSentimentUrl = "https://gateway-a.watsonplatform.net/calls/url/URLGetTextSentiment?apikey=" + apiKey

//curl -X POST -d "outputMode=json" -d "maxRetrieve=3" -d "url=https://www.fsf.org/blogs/community/who-in-the-world-is-changing-it-through-free-software-nominate-them-today" "https://gateway-a.watsonplatform.net/calls/url/URLGetRankedNamedEntities?apikey=$API_KEY"

// Sentiment analysis. content can be text or url. objectType is used to determine content is text or url
func AnalyzeSentimentText(content, objectType string) (*Sentiment, error) {
	if objectType != "text" && objectType != "url" {
		return nil, errors.New("Object type can be only text or url")
	}

	form := url.Values{}
	form.Add("outputMode", "json")
	if objectType == "url" {
		form.Add("url", content)
	} else {
		form.Add("text", content)
	}
  req, err := http.NewRequest("POST", baseSentimentUrl, bytes.NewBufferString(form.Encode()))
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  client := &http.Client{}
  resp, err := client.Do(req)
  respbody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  var result Sentiment
  err = json.Unmarshal(respbody, &result)
	if err != nil {
    fmt.Println(err)
    return nil, err
	}
  return &result, nil
}

type Sentiment struct {
	DocSentiment struct {
		Mixed string `json:"mixed"`
		Score float32 `json:"score,string"`
		Type  string `json:"type"`
	} `json:"docSentiment"`
	Language          string `json:"language"`
	Status            string `json:"status"`
	TotalTransactions string `json:"totalTransactions"`
	URL               string `json:"url"`
	Usage             string `json:"usage"`
}


