package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	baseURL = "https://api.switch-bot.com"
)

func main() {
	token, err := os.ReadFile("./token")
	if err != nil {
		os.Exit(1)
	}
	secret, err := os.ReadFile("./secret")
	if err != nil {
		os.Exit(1)
	}

	nonce := uuid.New().String()
	time := strconv.FormatInt(time.Now().UnixMilli(), 10)
	data := string(token) + time + nonce

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	client := &http.Client{}
	url, err := url.JoinPath(baseURL, "/v1.1/devices")
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", string(token))
	req.Header.Add("sign", signature)
	req.Header.Add("nonce", nonce)
	req.Header.Add("t", time)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
