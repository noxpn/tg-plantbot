package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func sendP(url *url.URL, fpath string) (string, error){
	

	fp := strings.Split(fpath, "/")
	fn := fp[len(fp)-1]

	q := url.Query()
	q.Set("caption",fn)
	url.RawQuery = q.Encode()

	//log.Println(url)

	//log.Println("url", u)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	f, err := os.Open(fpath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	w, err := writer.CreateFormFile("photo", filepath.Base(fpath))
	if err != nil {
		fmt.Println("CreateFormFile: ",err)
		return "", err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		fmt.Println("send copy:",err)
		return "", err
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("writer close:",err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url.String(), payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	log.Println(">", url.RawQuery, fpath)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(body), nil
}
