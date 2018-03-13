package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/christophberger/grada"
)

type measurement struct {
	Metric string  `json:"metric"`
	Value  float64 `json:"value"`
}

func logit(w http.ResponseWriter, r *http.Request) {
	userAgent := []string{""}
	ok := true
	if userAgent, ok = r.Header["User-Agent"]; !ok {
		userAgent = []string{""}
	}
	referer := []string{""}
	if referer, ok = r.Header["Referer"]; !ok {
		referer = []string{""}
	}
	userID := "-"
	fmt.Printf("%s - %s [%s] \"%s %s %s\" \"%s\" \"%s\" \"-\"\n",
		r.RemoteAddr,
		userID,
		time.Now().Format("2/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.RequestURI,
		r.Proto,
		referer[0],
		userAgent[0],
	)

}

func add(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Failed to read body")
		return
	}
	defer r.Body.Close()

	if bytes.Equal(body, nil) {
		fmt.Fprintf(w, "Must have a body.")
		return
	}

	var m measurement
	err = json.Unmarshal(body, &m)

	if _, ok := metrics[m.Metric]; !ok {
		metrics[m.Metric], err = dash.CreateMetric(m.Metric, 24*time.Hour, time.Minute)
		if err != nil {
			fmt.Fprintf(w, "Failed to create metric: %s\n", err)
			return
		}
	}

	fmt.Fprintf(w, "%#v", m.Value)
	metrics[m.Metric].Add(m.Value)
	logit(w, r)
}

var metrics map[string]*grada.Metric

var dash *grada.Dashboard

func main() {

	dash = grada.GetDashboard()

	http.HandleFunc("/add", add)

	var err error
	metrics = make(map[string]*grada.Metric)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Starting")
	select {}
}
