package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type (

	// Timestamp is a helper for (un)marhalling time
	Timestamp time.Time

	// HookMessage是我们从Alertmanager收到的消息
	HookMessage struct {
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		Status            string            `json:"status"`
		Receiver          string            `json:"receiver"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Alerts            []Alert           `json:"alerts"`
	}

	// Alert 是单个警报
	Alert struct {
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
		StartsAt    string            `json:"startsAt,omitempty"`
		EndsAt      string            `json:"EndsAt,omitempty"`
	}

	// 只是一个警报示例
	alertStore struct {
		sync.Mutex
		capacity int
		alerts   []*HookMessage
	}
)
    //定义路由向以及状态码
func main() {
	addr := flag.String("addr", ":8080", "address to listen for webhook")
	capacity := flag.Int("cap", 64, "capacity of the simple alerts store")
	flag.Parse()

	s := &alertStore{
		capacity: *capacity,
	}

	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/alerts", s.alertsHandler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok\n")
}

func (s *alertStore) alertsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getHandler(w, r)
	case http.MethodPost:
		s.postHandler(w, r)
	default:
		http.Error(w, "unsupported HTTP method", 400)
	}
}

func (s *alertStore) getHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	s.Lock()
	defer s.Unlock()

	if err := enc.Encode(s.alerts); err != nil {
		log.Printf("error encoding messages: %v", err)
	}
}

func (s *alertStore) postHandler(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var m HookMessage
	if err := dec.Decode(&m); err != nil {
		log.Printf("error decoding message: %v", err)
		http.Error(w, "invalid request body", 400)
		return
	}
	//利用sync库中Locker Mutex互斥锁来检验运行时检查警报

	s.Lock()
	defer s.Unlock()

	s.alerts = append(s.alerts, &m)

	if len(s.alerts) > s.capacity {
		a := s.alerts
		_, a = a[0], a[1:]
		s.alerts = a
	}
}
