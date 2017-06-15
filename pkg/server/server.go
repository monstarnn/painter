package server

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	//"github.com/emicklei/go-restful"
	//"github.com/gorilla/mux"
	"strconv"
	"io/ioutil"
	"encoding/base64"
	"bytes"
	"github.com/gorilla/mux"
)

const LISTEN_PORT = 8084

func Start() {
	logrus.Info("Starting painter...")
	r := mux.NewRouter()
	r.HandleFunc("/", indexAction)
	r.HandleFunc("/image", sendImage)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", r)
	logrus.Infoln("Listen in *:", LISTEN_PORT)
	err := http.ListenAndServe(":"+strconv.Itoa(LISTEN_PORT), nil) // set listen port
	if err != nil {
		logrus.Fatal("ListenAndServe error: ", err)
	}
}

func indexAction(w http.ResponseWriter, _ *http.Request) {
	contents, err := ioutil.ReadFile("./templates/index.html");
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	w.Write(contents)
}

func sendImage(w http.ResponseWriter, r *http.Request) {

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: use std functions for POST parsing
	d = bytes.TrimPrefix(d, []byte("img=data:image/octet-stream;base64,"))
	png, err := base64.StdEncoding.DecodeString(string(d))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ioutil.WriteFile("./test.png", png, 0644)

	w.Write([]byte("ok"))

}
