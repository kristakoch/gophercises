package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Constants for working with the Twitter API.
const (
	apiURL       = "https://api.twitter.com"
	authEndpoint = "oauth2/token"
	rtsEndpoint  = "1.1/statuses/retweets"
	tlEndpoint   = "1.1/statuses/user_timeline"
	favsEndpoint = "1.1/favorites/list"
	key          = "hGpWjxMEsZBBbkK5cjfAl0tKT"
	secretKey    = "1SnNcZeiaDJBqVzkoHotW6APmvXnemWMa1DTbs9PaZlmRFs55l"
)

// Twitter holds an auth token for interacting with Twitter.
type Twitter struct {
	authToken string
}

type authData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type retweetData struct {
	User struct {
		Name string `json:"screen_name"`
	} `json:"user"`
}

type tweetData struct {
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

// New is a factory function for creating rtcontests.
func New(ctx context.Context) (*Twitter, error) {
	var t Twitter

	authToken, err := authTokenFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	t.authToken = authToken

	return &t, nil
}

// GetAndStoreRTs gets tweet data from local and the Twitter API.
func (t *Twitter) GetAndStoreRTs(ctx context.Context, tweetID string) error {

	fileRTs, err := rtsFromFile(tweetID)
	if err != nil {
		return err
	}
	apiRTJSON, err := t.rtsFromAPI(ctx, tweetID)
	if err != nil {
		return err
	}

	var apiRTs []string
	for _, rt := range apiRTJSON {
		apiRTs = append(apiRTs, rt.User.Name)
	}

	// Merge api and local data.
	// Return the tweets in a string slice.
	rts := mergeRTs(apiRTs, fileRTs)

	diff := len(rts) - len(fileRTs)
	if diff > 0 {
		log.Printf("adding %d users to file", diff)
	}

	if err = storeRTs(rts, tweetID); err != nil {
		return err
	}

	return nil
}

// RTContestWinners logs a random selection of winners.
func (t *Twitter) RTContestWinners(ctx context.Context, tweetID string, nw int) error {
	rts, err := rtsFromFile(tweetID)
	if err != nil {
		return err
	}

	logWinners(rts, nw)

	return nil
}

// TimedRTContest gets rt data every x minutes and ends after y hours.
// If internet cuts out, continues when it finds it again.
// Might be good to check each time ticker runs: duration not longer than hours * hs. If so, get data once more * choose winners.
func (t *Twitter) TimedRTContest(
	ctx context.Context,
	tweetID string,
	nw int,
	ms time.Duration,
	hs time.Duration) {

	// Run thte ticker every specified number of minutes.
	ticker := time.NewTicker(ms * time.Minute)
	done := make(chan bool)

	log.Println("beginning timed rt contest...") // * add time running
	t.GetAndStoreRTs(ctx, tweetID)

	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				log.Printf("getting rt data for tweet with id %s\n", tweetID)
				t.GetAndStoreRTs(ctx, tweetID)
			}
		}
	}()

	// Stop after a number of hours.
	time.Sleep(hs * time.Hour)
	ticker.Stop()

	done <- true

	fmt.Printf("choosing %d winner(s)\n", nw)
	t.RTContestWinners(ctx, tweetID, nw)
}

// LogUserFavs logs a numer of favorites from a user based on their ID.
func (t *Twitter) LogUserFavs(ctx context.Context, userID string, nf int) error {
	fvs, err := t.favsFromAPI(ctx, userID, nf)
	if err != nil {
		return err
	}

	fmt.Printf("favs from user %s...\n", userID)
	for _, fv := range fvs {
		tt, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", fv.CreatedAt)
		if err != nil {
			return err
		}
		ttime := tt.Format("2006-01-02 15:04:05")
		clean := strings.ReplaceAll(fv.Text, "\n", " ")

		fmt.Printf("[%s] %s\n", ttime, clean)
	}

	return nil
}

// LogUserTweets logs a number of tweets from a user based on their ID.
func (t *Twitter) LogUserTweets(ctx context.Context, userID string, nt int) error {
	tts, err := t.tweetsFromAPI(ctx, userID, nt)
	if err != nil {
		return err
	}

	fmt.Printf("tweets from user %s...\n", userID)
	for _, twt := range tts {
		tt, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", twt.CreatedAt)
		if err != nil {
			return err
		}
		ttime := tt.Format("2006-01-02 15:04:05")
		clean := strings.ReplaceAll(twt.Text, "\n", " ")

		fmt.Printf("[%s] %s\n", ttime, clean)
	}

	return nil
}

// favsFromAPI gets a number of user favorites from the API.
func (t *Twitter) favsFromAPI(ctx context.Context, userID string, nf int) ([]tweetData, error) {

	reqURL := fmt.Sprintf("%s/%s.json?count=%d&screen_name=%s", apiURL, favsEndpoint, nf, userID)
	apiReq, err := http.NewRequest("GET", reqURL, nil)
	apiReq.Header.Add("Authorization", "Bearer "+t.authToken)

	apiRespBytes, err := sendRequest(ctx, apiReq)
	if err != nil {
		log.Fatal(err)
	}

	var apiRespJSON []tweetData
	if err = json.Unmarshal(apiRespBytes, &apiRespJSON); err != nil {
		log.Print("unexpected response from API, err:", string(apiRespBytes))
		return apiRespJSON, errors.New("no results found")
	}

	return apiRespJSON, nil
}

// favsFromAPI gets a number of user tweets from the API.
func (t *Twitter) tweetsFromAPI(ctx context.Context, userID string, nt int) ([]tweetData, error) {

	// Make an API request and get the tweets.
	reqURL := fmt.Sprintf("%s/%s/%s.json?count=%d", apiURL, tlEndpoint, userID, nt)
	apiReq, err := http.NewRequest("GET", reqURL, nil)
	apiReq.Header.Add("Authorization", "Bearer "+t.authToken)

	apiRespBytes, err := sendRequest(ctx, apiReq)
	if err != nil {
		log.Fatal(err)
	}

	var apiRespJSON []tweetData
	if err = json.Unmarshal(apiRespBytes, &apiRespJSON); err != nil {
		log.Print("unexpected response from API, err:", string(apiRespBytes))
		return apiRespJSON, errors.New("no results found")
	}

	return apiRespJSON, nil
}

// rtsFromAPI gets a number of users who rted the given tweet from the API.
func (t *Twitter) rtsFromAPI(ctx context.Context, tweetID string) ([]retweetData, error) {

	// Make an API request and get the tweets.
	reqURL := fmt.Sprintf("%s/%s/%s.json?count=100", apiURL, rtsEndpoint, tweetID)
	apiReq, err := http.NewRequest("GET", reqURL, nil)
	apiReq.Header.Add("Authorization", "Bearer "+t.authToken)

	apiRespBytes, err := sendRequest(ctx, apiReq)
	if err != nil {
		log.Fatal(err)
	}
	var apiRespJSON []retweetData
	if err = json.Unmarshal(apiRespBytes, &apiRespJSON); err != nil {
		log.Print("unexpected result from API, err:", string(apiRespBytes))
		return apiRespJSON, errors.New("no results found")
	}

	return apiRespJSON, nil
}

// sendRequest sents an HTTP request to the Twitter API.
func sendRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

// authTokenFromAPI gets an auth token from Twitter.
func authTokenFromAPI(ctx context.Context) (string, error) {

	// Build the request for an auth token.
	csKeys := url.QueryEscape(key) + ":" + url.QueryEscape(secretKey)
	encodedKeys := base64.StdEncoding.EncodeToString([]byte(csKeys))

	authReq, err := http.NewRequest("POST", apiURL+"/"+authEndpoint, strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}

	authReq.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	authReq.Header.Add("Authorization", "Basic "+encodedKeys)

	// Get the auth token from the Twitter API.
	respBytes, err := sendRequest(ctx, authReq)
	if err != nil {
		return "", err
	}
	var respJSON authData
	if err = json.Unmarshal(respBytes, &respJSON); err != nil {
		log.Print("unexpected response from API, err:", string(respBytes))
		return "", errors.New("no results found")
	}

	return respJSON.AccessToken, nil
}

// rtsFromFile gets local rt data.
func rtsFromFile(tweetID string) ([]string, error) {
	var rts []string

	rtfn := tweetID + ".txt"
	if _, err := os.Stat(rtfn); os.IsNotExist(err) {
		return rts, nil
	}

	file, err := ioutil.ReadFile(rtfn)
	if err != nil {
		return rts, err
	}
	rts = strings.Split(string(file), "\n")

	return rts, nil
}

// mergeRTs merges API and local rt data.
func mergeRTs(fileRTs []string, apiRTs []string) []string {
	var uniqueRTs []string

	seen := make(map[string]bool)
	for _, rt := range fileRTs {
		if rt = strings.TrimSpace(rt); rt != "" {
			seen[rt] = true
		}
	}
	for _, rt := range apiRTs {
		if rt = strings.TrimSpace(rt); rt != "" {
			seen[rt] = true
		}
	}
	for rt := range seen {
		uniqueRTs = append(uniqueRTs, rt)
	}

	return uniqueRTs
}

// logWinners logs a number of randomly-chosen rt users.
func logWinners(rts []string, nw int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	nrts := len(rts)

	switch {
	case nrts == 0:
		fmt.Println("no entrants")
		break
	case nw > nrts:
		fmt.Printf("number of winners (%d) cannot be above number of rts (%d)\n", nw, nrts)
		break
	case nw == 1:
		idx := r.Intn(nrts)
		fmt.Printf("winner: @%s\n", rts[idx])
		break
	case nw > 1:
		perm := r.Perm(len(rts))
		var winners []string
		for i, randIndex := range perm {
			if i == nw {
				break
			}
			winners = append(winners, "@"+rts[randIndex])
		}
		fmt.Println("winners:", strings.Join(winners, ", "))
		break
	}

}

// storeRTs writes users who rted a tweet to a file.
func storeRTs(rts []string, tweetID string) error {
	rtData := strings.Join(rts, "\n")

	fd := []byte(rtData)
	if err := ioutil.WriteFile(tweetID+".txt", fd, 0644); err != nil {
		return err
	}

	return nil
}
