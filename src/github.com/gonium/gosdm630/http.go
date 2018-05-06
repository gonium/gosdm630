package sdm630

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jcuga/golongpoll"
)

const (
	SECONDS_BETWEEN_STATUSUPDATE = 1
)

// Generate the embedded assets using https://github.com/aprice/embed
//go:generate embed -c "embed.json"

func mkIndexHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	loader := GetEmbeddedContent()
	mainTemplate, err := loader.GetContents("/index.tmpl")
	if err != nil {
		log.Fatal("Failed to load embedded template: " + err.Error())
	}
	t, err := template.New("gosdm630").Parse(string(mainTemplate))
	if err != nil {
		log.Fatal("Failed to create main page template: ", err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		data := struct {
			SoftwareVersion string
			GolangVersion   string
		}{
			SoftwareVersion: RELEASEVERSION,
			GolangVersion:   runtime.Version(),
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Fatal("Failed to render main page: ", err.Error())
		}
	})
}

func mkLastAllValuesHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		ids := hc.GetSortedIDs()
		lasts := ReadingSlice{}
		for _, id := range ids {
			reading, err := hc.GetLast(id)
			if err != nil {
				// Skip this meter, it will simply not be displayed
				continue
				//w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, err.Error())
				//return
			}
			lasts = append(lasts, *reading)
		}
		if len(lasts) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "All meters are inactive.")
			return
		}
		if err := lasts.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

func mkLastSingleValuesHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		last, err := hc.GetLast(byte(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func mkLastMinuteAvgSingleHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		avg, err := hc.GetMinuteAvg(byte(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := avg.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", avg.String())
		}
	})
}

func mkLastMinuteAvgAllHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		ids := hc.GetSortedIDs()
		avgs := ReadingSlice{}
		for _, id := range ids {
			reading, err := hc.GetMinuteAvg(id)
			if err != nil {
				// Skip this meter, it will simply not be displayed
				continue
			}
			avgs = append(avgs, *reading)
		}
		if len(avgs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "All meters are inactive.")
			return
		}
		if err := avgs.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}

	})
}

func mkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := s.UpdateAndJSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

type Firehose struct {
	lpManager  *golongpoll.LongpollManager
	in         QuerySnipChannel
	statstream chan string
}

func NewFirehose(inChannel QuerySnipChannel, status *Status, verbose bool) *Firehose {
	options := golongpoll.Options{}
	// see https://github.com/jcuga/golongpoll/blob/master/longpoll.go#L81
	//options := golongpoll.Options{
	//	LoggingEnabled:                 false,
	//	MaxLongpollTimeoutSeconds:      60,
	//	MaxEventBufferSize:             250,
	//	EventTimeToLiveSeconds:         60,
	//	DeleteEventAfterFirstRetrieval: false,
	//}
	if verbose {
		options.LoggingEnabled = true
	}
	manager, err := golongpoll.StartLongpoll(options)
	if err != nil {
		log.Fatalf("Failed to create firehose longpoll manager: %q", err)
	}
	// Attach a goroutine that will push meter status information
	// periodically
	var statusstream = make(chan string)
	go func() {
		for {
			time.Sleep(SECONDS_BETWEEN_STATUSUPDATE * time.Second)
			var buffer bytes.Buffer
			if err := status.UpdateAndJSON(&buffer); err == nil {
				statusstream <- buffer.String()
				buffer.Reset()
			}
		}
	}()
	return &Firehose{
		lpManager:  manager,
		in:         inChannel,
		statstream: statusstream,
	}
}

func (f *Firehose) Run() {
	for {
		select {
		case snip := <-f.in:
			//log.Printf("FooSnip: %+v", snip)
			//var buffer bytes.Buffer
			//err := json.NewEncoder(&buffer).Encode(snip)
			//if err != nil {
			//	log.Println("FooError: ", err.Error())
			//} else {
			//	log.Println("FooSuccess: ", buffer.String())
			//}
			f.lpManager.Publish("meterupdate", snip)
		case statupdate := <-f.statstream:
			f.lpManager.Publish("statusupdate", statupdate)
		}
	}
}

func (f *Firehose) GetHandler() func(w http.ResponseWriter, r *http.Request) {
	return f.lpManager.SubscriptionHandler
}

// RunHTTPD starts httpd server
func RunHTTPD(
	mc *MeasurementCache,
	firehose *Firehose,
	s *Status,
	url string,
) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", mkIndexHandler(mc))
	router.HandleFunc("/last", mkLastAllValuesHandler(mc))
	router.HandleFunc("/last/{id:[0-9]+}", mkLastSingleValuesHandler(mc))
	router.HandleFunc("/minuteavg", mkLastMinuteAvgAllHandler(mc))
	router.HandleFunc("/minuteavg/{id:[0-9]+}", mkLastMinuteAvgSingleHandler(mc))
	router.HandleFunc("/status", mkStatusHandler(s))
	router.HandleFunc("/firehose", firehose.GetHandler())
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		GetEmbeddedContent()))
	srv := http.Server{
		Addr:         url,
		Handler:      handlers.CompressHandler(router),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
