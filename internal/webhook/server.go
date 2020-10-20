package main

import (
	"encoding/json"
	"net/http"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&nested.Formatter{
		HideKeys: true, NoColors: false,
		NoFieldsColors: false,
		ShowFullLevel:  true,
		FieldsOrder:    []string{"req", "form"},
	})

	log.SetLevel(log.DebugLevel)

	//
	// Enable to show caller
	//
	// log.SetReportCaller(true)
}

func main() {
	http.HandleFunc("/api/v1/alert", alertHandler)
	log.Fatal(http.ListenAndServe(":8088", nil))
}

//
// AlertResponseMessage json response
//
type AlertResponseMessage struct {
	StatusCode int
	StatusText string
	Message    string
}

func alertHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"req": r.Method + " api/v1/alert",
		}).Info("Receive alert")

		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Errorf("ParseForm() err: %v", err)
			return
		}
	}

	log.WithFields(log.Fields{
		"req":  r.Method + " api/v1/alert",
		"form": r.PostForm,
	}).Info("Receive alert")

	//
	// Prepare response message
	//
	message := AlertResponseMessage{200, "ok", "Ok!"}

	js, err := json.Marshal(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
