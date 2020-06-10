package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"regexp"
)

var httpRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "vip_http_request_count",
		Help: "http request count",
	},
	[]string{"domain", "uri", "interval"},
)

var httpRequestDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "vip_http_request_duration",
		Help: "vip http request duration",
	},
	[]string{"domain", "uri"},
)

var httpTotalDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "vip_http_request_domain_duration",
		Help: "vip http request domain duration",
	},
	[]string{"domain"},
)

var domDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "vip_front_dom_duration",
		Help: "vip dom duration",
	},
	[]string{"pageid"},
)

var fcpDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "vip_front_fcp_duration",
		Help: "vip front fcp duration",
	},
	[]string{"pageid"},
)

var atfDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "vip_front_atf_duration",
		Help: "vip front atf duration",
	},
	[]string{"pageid"},
)

var onloadDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "vip_front_onload_duration",
		Help: "vip front onload duration",
	},
	[]string{"pageid"},
)

var frontErrorCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "vip_front_error_count",
		Help: "vip front error count",
	},
	[]string{"pageid", "event_type"},
)

var freedomTimeMap = make(map[string]*prometheus.SummaryVec)

func genFreedomTimeDuration(idx int) (vec *prometheus.SummaryVec) {
	var ok bool
	label := fmt.Sprintf("vip_front_ft_%d_duration", idx)
	if vec, ok = freedomTimeMap[label]; ok {
		return
	}
	vec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: label,
			Help: fmt.Sprintf("vip front ft_%d duration", idx),
		},
		[]string{"pageid"},
	)
	freedomTimeMap[label] = vec
	return
}

type reportData struct {
	Domain string `json:"domain"`
	URI    string `json:"uri"`
	Cost   int    `json:"cost"`
}

type frontData struct {
	PageID       string `json:"pageid"`
	Dom          int    `json:"dom"`
	Fcp          int    `json:"fcp"`
	Atf          int    `json:"atf"`
	Onload       int    `json:"Onload"`
	FreedomTime0 int    `json:"freedom_time_0"`
	FreedomTime1 int    `json:"freedom_time_1"`
	FreedomTime2 int    `json:"freedom_time_2"`
	FreedomTime3 int    `json:"freedom_time_3"`
	FreedomTime4 int    `json:"freedom_time_4"`
}

type frontErrorData struct {
	PageID    string `json:"pageid"`
	Path      string `json:"path"`
	Message   string `json:"messgae"`
	EventType string `json:"event_type"`
}

func init() {
	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpTotalDuration)
	prometheus.MustRegister(domDuration)
	prometheus.MustRegister(fcpDuration)
	prometheus.MustRegister(atfDuration)
	prometheus.MustRegister(onloadDuration)
	prometheus.MustRegister(frontErrorCount)

	//FreedomTime 0-4
	for i := 0; i < 5; i++ {
		prometheus.MustRegister(genFreedomTimeDuration(i))
	}
}

func main() {

	http.HandleFunc("/report", report)
	http.HandleFunc("/front/report", frontReport)
	http.HandleFunc("/front/errorreport", frontErrorReport)
	http.Handle("/metrics", prometheus.Handler())

	err := http.ListenAndServe(":7002", nil)

	if err != nil {
		panic(err)
	}

}

func report(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, "read body failed, err:"+err.Error())
		return
	}

	if len(data) <= 0 {
		fmt.Fprintln(w, "empty post data")
		return
	}

	var reportObj *reportData
	reportObj = &reportData{}

	err = json.Unmarshal(data, reportObj)

	if err != nil {
		fmt.Fprintln(w, "decode json failed, err:"+err.Error())
		return
	}

	var interval string
	interval = "0-100"
	switch {
	case reportObj.Cost > 0 && reportObj.Cost <= 100:
		interval = "0-100"
	case reportObj.Cost > 100 && reportObj.Cost <= 1000:
		interval = "100-1k"
	case reportObj.Cost > 1000 && reportObj.Cost <= 10000:
		interval = "1k-10k"
	case reportObj.Cost > 10000:
		interval = "10k-inf"
	}

	//record count
	httpRequestCount.WithLabelValues(reportObj.Domain, reportObj.URI, interval).Inc()

	//record duration value
	httpRequestDuration.WithLabelValues(reportObj.Domain, reportObj.URI).Observe(float64(reportObj.Cost))

	httpTotalDuration.WithLabelValues(reportObj.Domain).Observe(float64(reportObj.Cost))
	httpTotalDuration.WithLabelValues("all").Observe(float64(reportObj.Cost))

	fmt.Fprintln(w, "ok")
	return
}

func frontReport(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, "read body failed, err:"+err.Error())
		return
	}

	if len(data) <= 0 {
		fmt.Fprintln(w, "empty post data")
		return
	}

	var reportObj *frontData
	reportObj = &frontData{}

	err = json.Unmarshal(data, reportObj)

	if err != nil {
		fmt.Fprintln(w, "decode json failed, err:"+err.Error())
		return
	}

	if reportObj.Dom > 0 {
		domDuration.WithLabelValues(reportObj.PageID).Observe(float64(reportObj.Dom))
	}

	if reportObj.Fcp > 0 {
		fcpDuration.WithLabelValues(reportObj.PageID).Observe(float64(reportObj.Fcp))
	}
	if reportObj.Atf > 0 {
		atfDuration.WithLabelValues(reportObj.PageID).Observe(float64(reportObj.Atf))
	}
	if reportObj.Onload > 0 {
		onloadDuration.WithLabelValues(reportObj.PageID).Observe(float64(reportObj.Onload))
	}

	objType := reflect.TypeOf(*reportObj)
	objValue := reflect.ValueOf(*reportObj)

	regRule, err := regexp.Compile("^FreedomTime(\\d+)$")
	if err != nil {
		fmt.Println(err)
	}
	for idx := 0; idx < objType.NumField(); idx++ {
		var field string
		field = objType.Field(idx).Name
		result := regRule.FindStringSubmatch(field)
		if len(result) > 0 {
			value := objValue.FieldByName(field)
			if value.Int() != 0 {
				ftidxStr := result[1]
				ftidxInt, err := strconv.Atoi(ftidxStr)
				if err != nil {
					fmt.Println(err)
				}
				duration := genFreedomTimeDuration(ftidxInt)
				duration.WithLabelValues(reportObj.PageID).Observe(float64(value.Int()))
			}
		}
	}
	fmt.Fprintln(w, "ok")
	return
}

func frontErrorReport(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintln(w, "read body failed, err:"+err.Error())
		return
	}

	if len(data) <= 0 {
		fmt.Fprintln(w, "empty post data")
		return
	}

	var reportObj *frontErrorData
	reportObj = &frontErrorData{}

	err = json.Unmarshal(data, reportObj)

	if err != nil {
		fmt.Fprintln(w, "decode json failed, err:"+err.Error())
		return
	}

	frontErrorCount.WithLabelValues(reportObj.PageID, reportObj.EventType).Inc()

	fmt.Fprintln(w, "ok")
	return

}
