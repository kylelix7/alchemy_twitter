package tweets

import (
	"crypto/rand"
	"fmt"
	"encoding/base64"
	"errors"
	"sort"
	"strings"
	"time"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

func obtainBearerToken(encodedKeySecret string) (string, error) {
	urlstr := "https://api.twitter.com/oauth2/token"
	body := []byte("grant_type=client_credentials")

	req, err := http.NewRequest("POST", urlstr, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic " + encodedKeySecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Status code is %d. response: %v", resp.StatusCode, string(respbody)))
	}
	type NormalResp struct {
		Type string `json:"token_type"`
		Token string `json:"access_token"`
	}
	var jsonObj NormalResp
	json.Unmarshal(respbody, &jsonObj)
	return string(jsonObj.Token), nil
}

func encodeConsumerKeySecret(consumerKey, consumerSecret string) string {
	credentialsBytes := base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(consumerKey) + ":" + url.QueryEscape(consumerSecret)))
	return fmt.Sprintf("%v", credentialsBytes)
}

func generateNounce() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", errors.New("Cannot generate nounce")
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func generateSignature(consumerKey, method, baseUrl, secret string, kvs []KeyVal) (string, error) {
	kvs = append(kvs, KeyVal{Key: "oauth_consumer_key", Val: consumerKey})
	nounce, err := generateNounce()
	if err != nil {
		return "", err
	}
	kvs = append(kvs, KeyVal{Key: "oauth_nonce", Val: nounce})
	kvs = append(kvs, KeyVal{Key: "oauth_signature_method", Val: "HMAC-SHA1"})
	kvs = append(kvs, KeyVal{Key: "oauth_timestamp", Val: fmt.Sprintf("%v", int32(time.Now().Unix()))})
	kvs = append(kvs, KeyVal{Key: "oauth_version", Val: "1.0"})
	kvs = append(kvs, KeyVal{Key: "oauth_token", Val: "1283630658-4s0zHmNnn3cBP4QmAeLn4hamrfD5fG4oOcHpRkV"})
	sort.Sort(KeyVals(kvs))
	result := ""
	for _, kv := range kvs {
		result = result + kv.Key + "=" + url.QueryEscape(kv.Val) + "&"
	}
	result = strings.ToUpper(method) + "&" + url.QueryEscape(baseUrl) + url.QueryEscape(result[0:len(result)-1])
	return "", nil
}


type KeyVal struct {
	Key string
	Val string
}

type KeyVals []KeyVal

func (slice KeyVals) Len() int {
	return len(slice)
}

func (slice KeyVals) Less(i, j int) bool {
	return slice[i].Key < slice[j].Key;
}

func (slice KeyVals) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
