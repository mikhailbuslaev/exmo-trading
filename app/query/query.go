package query

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type User struct {
	Id        string
	PublicKey string
	SecretKey string
}

type Response map[string]interface{}

type Query interface {
	Do() (*Response, error)
}

type PostQuery struct {
	UserParams     User
	Sign           string
	Method         string
	Params         map[string]string
	PreparedParams string
}

type GetQuery struct {
	Method string
}

func (q *PostQuery) PrepareParams() {
	post_params := url.Values{}
	post_params.Add("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	if q.Params != nil {
		for key, value := range q.Params {
			post_params.Add(key, value)
		}
	}
	q.PreparedParams = post_params.Encode()
}

func (q *PostQuery) GetSign() {
	mac := hmac.New(sha512.New, []byte(q.UserParams.SecretKey))
	mac.Write([]byte(q.PreparedParams))
	q.Sign = fmt.Sprintf("%x", mac.Sum(nil))
}

func (q *PostQuery) Do() (*Response, error) {

	q.PrepareParams()
	q.GetSign()

	req, _ := http.NewRequest("POST", "https://api.exmo.me/v1/"+q.Method, bytes.NewBuffer([]byte(q.PreparedParams)))
	req.Header.Set("Key", q.UserParams.PublicKey)
	req.Header.Set("Sign", q.Sign)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(q.PreparedParams)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ReturnResponse(resp)
}

func (q *GetQuery) Do() (*Response, error) {

	resp, err := http.Get(q.Method)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ReturnResponse(resp)
}

func ReturnResponse(resp *http.Response) (*Response, error) {
	if resp.Status != "200 OK" {
		return nil, errors.New("http status: " + resp.Status)
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		return nil, err1
	}

	var dat map[string]interface{}
	err2 := json.Unmarshal([]byte(body), &dat)
	if err2 != nil {
		return nil, err2
	}

	if result, ok := dat["result"]; ok && !result.(bool) {
		return nil, errors.New(dat["error"].(string))
	}
	r := Response{}
	r = dat
	return &r, nil
}

func Exec(q Query) (*Response, error) {
	return q.Do()
}
