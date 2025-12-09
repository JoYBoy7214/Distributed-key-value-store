package main

import (
	"encoding/json"
	"flag"
	"fmt"
	msg "key_value_store/msg"
	"net/http"
	"strconv"
	"sync"
)

var Rmtx sync.RWMutex

func Puthandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var msg msg.Putmsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Rmtx.Lock()
	Map[msg.Key] = msg.Value
	Rmtx.Unlock()
	w.WriteHeader(http.StatusOK)
}
func Gethandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var msg msg.Getmsg
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "unable to decode", http.StatusBadRequest)
		return
	}
	var value string
	Rmtx.RLock()
	value = Map[msg.Key]
	Rmtx.RUnlock()
	if value == "" {
		http.Error(w, "key does exists", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(value)

}

var Map map[string]string

func main() {
	var port int
	Map = make(map[string]string)
	flag.IntVar(&port, "port", 8081, "help for port ")
	flag.Parse()
	http.HandleFunc("/GET", Gethandler)
	http.HandleFunc("/PUT", Puthandler)
	portString := ":" + strconv.Itoa(port)
	fmt.Println(portString)
	http.ListenAndServe(portString, nil)
}
