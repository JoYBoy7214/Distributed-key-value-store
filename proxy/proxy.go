package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	message "key_value_store/msg"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"sync"
)

type Server struct {
	url    *url.URL
	Rproxy *httputil.ReverseProxy
	Weight int
}

type ConsistantHash struct {
	NumberVirtualNodes     int
	key                    []int
	ServerMapping          map[int]*Server
	mtx                    sync.Mutex
	NumberOfRepicationNode int
}

type ServerInfo struct {
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}

func (ch *ConsistantHash) AddServer(s *Server) {
	ch.mtx.Lock()
	for i := 0; i < (ch.NumberVirtualNodes * s.Weight); i++ {
		urlString := s.url.String() + "#" + strconv.Itoa(i)
		urlKey := HashHelper(urlString)
		index := sort.SearchInts(ch.key, int(urlKey))
		ch.key = slices.Insert(ch.key, index, int(urlKey))
		ch.ServerMapping[int(urlKey)] = s
		fmt.Println(index, urlString, urlKey)
	}
	ch.mtx.Unlock()

}
func (ch *ConsistantHash) GetServers(msg string) []*Server {
	KeyHash := HashHelper(msg)
	ch.mtx.Lock()
	index := sort.SearchInts(ch.key, int(KeyHash))
	if index == len(ch.key) {
		index = 0
	}
	ans := make([]*Server, ch.NumberOfRepicationNode)
	ans[0] = ch.ServerMapping[ch.key[index]]
	for i := 1; i < ch.NumberOfRepicationNode; {
		if ans[i-1].url.String() != ch.ServerMapping[ch.key[index]].url.String() {
			ans[i] = ch.ServerMapping[ch.key[index]]
			i++
		}
		index++
		if index == len(ch.key) {
			index = 0
		}
	}
	return ans
}

func HashHelper(key string) int32 {
	hasher := fnv.New32a()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		log.Fatal("error in hash writer")
	}
	ans := hasher.Sum32()
	return int32(ans)
}

func RequestHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var msg message.Getmsg
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	proxyServers := Servers.GetServers(msg.Key)
	clinet := &http.Client{}
	for i := 0; i < Servers.NumberOfRepicationNode; i++ {
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		targetUrl := proxyServers[i].url.String() + "/GET"
		req, err := http.NewRequest(http.MethodGet, targetUrl, r.Body)
		if err != nil {
			fmt.Println("error in creating request")
			continue
		}
		req.Header.Set("content-Type", "application/json")
		resp, err := clinet.Do(req)
		if err != nil {
			fmt.Println("error in making request")
			continue
		}
		if resp.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				log.Println("error in coping response body")
				continue
			}
			resp.Body.Close()
			return
		}
	}
	http.Error(w, "key not found", http.StatusNotFound)

}
func RequestHandlerPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var msg message.Getmsg
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	proxyServers := Servers.GetServers(msg.Key)
	clinet := &http.Client{}
	flag := false
	for i := 0; i < Servers.NumberOfRepicationNode; i++ {
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		targetUrl := proxyServers[i].url.String() + "/PUT"
		req, err := http.NewRequest(http.MethodPut, targetUrl, r.Body)
		if err != nil {
			fmt.Println("error in creating request")
			continue
		}
		req.Header.Set("content-Type", "application/json")
		resp, err := clinet.Do(req)
		if err != nil {
			fmt.Println("error in making request")
			continue
		}
		if resp.StatusCode == 200 {
			flag = true
		}
	}
	if flag {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "key not found", http.StatusNotFound)
}

var Servers ConsistantHash

func ConfigServers(Port ServerInfo) {
	urlString := Port.Url
	curUrl, err := url.Parse(urlString)
	if err != nil {
		log.Fatal("unable to parse the url")
	}
	Servers.AddServer(&Server{
		url:    curUrl,
		Rproxy: httputil.NewSingleHostReverseProxy(curUrl),
		Weight: Port.Weight,
	})

}
func AddServerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var serverInfo ServerInfo
	err := json.NewDecoder(r.Body).Decode(&serverInfo)
	if err != nil {
		http.Error(w, "unable to decode the serverInfo", http.StatusBadRequest)
		return
	}
	ConfigServers(serverInfo)
	w.WriteHeader(http.StatusOK)

}
func main() {
	Servers.ServerMapping = make(map[int]*Server)
	Servers.NumberVirtualNodes = 5
	Servers.NumberOfRepicationNode = 3
	http.HandleFunc("/GET", RequestHandlerGet)
	http.HandleFunc("/PUT", RequestHandlerPut)
	http.HandleFunc("/AddServer", AddServerHandler)
	http.ListenAndServe(":8080", nil)

}
