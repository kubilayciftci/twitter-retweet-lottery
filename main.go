package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func main() {

	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}

	f, err := os.Open(".keys.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.Decode(&keys)

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(keys.Key, keys.Secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var token oauth2.Token

	dec = json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     keys.Key,
		ClientSecret: keys.Secret,
		Endpoint: oauth2.Endpoint{
			TokenURL: "",
			AuthURL:  "",
		},
	}
	tclient := conf.Client(ctx, &token)
	rt, err := tclient.Get("https://api.twitter.com/1.1/statuses/retweets/1348367126151778307.json")
	if err != nil {
		panic(err)
	}
	defer rt.Body.Close()
	io.Copy(os.Stdout, rt.Body)
}
