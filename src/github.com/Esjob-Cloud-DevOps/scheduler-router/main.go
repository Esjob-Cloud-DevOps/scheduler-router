package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	addr = flag.String("listen", ":8899", "The address of listen on")
	conf = flag.String("conf", "", "The config string of JSON type")
)

func getAvailableServer(servers []string) (string, error) {
	for _, s := range servers {
		if _, err := net.DialTimeout("tcp", s, time.Second*2); err == nil {
			return s, nil
		}
	}
	return "", errors.New(fmt.Sprintf("All servers %v are unavailable", servers))
}
func main() {
	flag.Parse()
	if len(*conf) < 1 {
		log.Fatal("Lack 'conf' param")
	}

	log.Printf("conf str: %s", *conf)
	proxyMaps := make(map[string][]string)
	if err := json.Unmarshal([]byte(*conf), &proxyMaps); err != nil {
		log.Panic(err)
	}
	log.Printf("conf: %v", proxyMaps)
	mux := http.NewServeMux()
	for k, v := range proxyMaps {
		servers := v
		mux.HandleFunc("/"+k+"/", func(w http.ResponseWriter, r *http.Request) {
			server, err := getAvailableServer(servers)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, err)
				log.Printf("Error: %v", err)
				return
			}
			url := r.URL
			oldUrl := url.String()
			url.Host = server
			url.Path = strings.Join(strings.Split(url.Path, "/")[2:], "/")
			log.Printf("%s is redirected to url: %s", oldUrl, url.String())
			http.Redirect(w, r, url.String(), 307)
		})
	}

	http.ListenAndServe(*addr, mux)
}
