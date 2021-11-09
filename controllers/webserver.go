/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type WebServer struct {
}

type Report struct {
	Name       string
	Size       int64
	CreateTime string
}

func (ws *WebServer) Start(context context.Context) error {
	http.HandleFunc("/reportlist", listReport)
	http.HandleFunc("/checksize", checkSize)
	http.ListenAndServe(":8090", nil)
	return nil
}

func checkSize(w http.ResponseWriter, req *http.Request) {
	filenames, ok := req.URL.Query()["filename"]
	if !ok || len(filenames[0]) < 1 {
		log.Println("Url Param 'filename' is missing")
		return
	}

	filename := filenames[0]

	err := checkFileSize("/report/" + filename)
	if err != nil {
		log.Println("internal error")
		return
	}

	fmt.Fprintf(w, "done\n")
}

func checkFileSize(filepath string) error {
	var size int64
	size = 0

	for size == 0 {
		fi, err := os.Stat(filepath)
		if err != nil {
			return err
		}

		size = fi.Size()
		time.Sleep(1 * time.Second)
	}

	return nil
}

func listReport(w http.ResponseWriter, req *http.Request) {
	reports, err := getFilelist("/report")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(reports)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "hello\n")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getFilelist(path string) (*[]Report, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	reports := []Report{}

	for _, file := range files {
		reports = append(reports, Report{Name: file.Name(), Size: file.Size(), CreateTime: file.ModTime().String()})
	}

	return &reports, nil
}

func (ws *WebServer) NeedLeaderElection() bool {
	return true
}
