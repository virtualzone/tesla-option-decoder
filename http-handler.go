package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-playground/validator"
)

type GetOptionCodeRequest struct {
	URL string `json:"url" validate:"required,url"`
}

func GetOptionCodesHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != "POST" {
		SendMethodNotAllowed(w)
		return
	}
	var m GetOptionCodeRequest
	if UnmarshalValidateBody(r, &m) != nil {
		SendBadRequest(w)
		return
	}
	if strings.Index(m.URL, "https://www.tesla.com/") != 0 {
		SendBadRequest(w)
		return
	}
	content, err := ReadFile(m.URL)
	if err != nil {
		SendNotFound(w)
		return
	}
	refList := OptionCodeListInstance
	codes, err := ExtractCodes(content)
	if err != nil {
		SendInternalServerError(w)
		return
	}
	resp := make(OptionCodeList)
	for _, code := range codes {
		details := refList[code]
		//log.Println(code + " --> " + details.Title)
		resp[code] = details
	}
	SendJSON(w, resp)
}

func SendNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func SendBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func SendInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func SendJSON(w http.ResponseWriter, v interface{}) {
	json, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		SendInternalServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func SendMethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func UnmarshalBody(r *http.Request, o interface{}) error {
	if r.Body == nil {
		return errors.New("body is NIL")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, &o); err != nil {
		return err
	}
	return nil
}

func UnmarshalValidateBody(r *http.Request, o interface{}) error {
	err := UnmarshalBody(r, &o)
	if err != nil {
		return err
	}
	err = validator.New().Struct(o)
	if err != nil {
		return err
	}
	return nil
}

func ServeHTTP() {
	log.Println("Initializing REST services...")
	httpServer := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	http.HandleFunc("/api/optioncodes", GetOptionCodesHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()
	log.Println("HTTP Server listening")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	httpServer.Shutdown(ctx)
}
