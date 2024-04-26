package main

import (
	"fmt"
	//"os"
	"time"
	"tg-bot-echo/lib/e"
	"io/ioutil"
)


func filesInDir(dir string) error {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return e.Wrap("files in Dir: %e", err)
	}
	c := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		c++
	}
	if c > 0 {
		return nil
	} else {
		return e.Wrap("FilesInDir, %e", fmt.Errorf("dir is empty"))
	}
	
}
func getPhoto(dir string, fallback_dir string) (string, error) {

	d := dir
	if err := filesInDir(dir); err != nil {
		if err:= filesInDir(fallback_dir); err != nil{
			return "", fmt.Errorf("no files to proceed")
		} 
		d = fallback_dir
	} 
	f, err := getLastFile(d)
	if err!= nil{
		return "", fmt.Errorf("no photo to send")
	}
	return f, nil
}


func getLastFile(d string) (string, error){


	files, err := ioutil.ReadDir(d)
	if err != nil {
		return "", fmt.Errorf("ReadDir: %e", err)
	}

	m := make(map[string]time.Time)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		m[d + file.Name()] = file.ModTime()
	}

	var timestamp time.Time
	var lastfile string 
	
	for k, v := range m {
		if v.After(timestamp) {
			timestamp = v
			lastfile = k
		}
		
	}
	return lastfile, nil
}
