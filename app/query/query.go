package query

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type User struct {
	Id        string
	PublicKey string
	SecretKey string
}

type Query interface {
	Do() (*http.Response, error)
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
	postParams := url.Values{}
	if q.Params != nil {
		for key, value := range q.Params {
			postParams.Add(key, value)
		}
	}
	q.PreparedParams = postParams.Encode()
}

func (q *PostQuery) GetSign() {
	mac := hmac.New(sha512.New, []byte(q.UserParams.SecretKey))
	mac.Write([]byte(q.PreparedParams))
	q.Sign = fmt.Sprintf("%x", mac.Sum(nil))
}

func (q *PostQuery) Do() (*http.Response, error) {

	q.PrepareParams()
	q.GetSign()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.exmo.me/v1.1/"+q.Method, bytes.NewBuffer([]byte(q.PreparedParams)))
	req.Header.Set("Key", q.UserParams.PublicKey)
	req.Header.Set("Sign", q.Sign)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(q.PreparedParams)))
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(resp.Body)
	return resp, nil
}

func (q *GetQuery) Do() (*http.Response, error) {
	resp, err := http.Get("https://api.exmo.me/v1.1/" + q.Method)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Exec(q Query) (*http.Response, error) {
	return q.Do()
}
