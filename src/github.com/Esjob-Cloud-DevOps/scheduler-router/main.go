package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	addr = flag.String("listen", ":8899", "The address of listen on")
	roles = flag.String("roles", "", "The roles of frameworks")
)

func main() {
	flag.Parse()
	if len(*roles) < 1 {
		log.Fatal("Lack 'roles' param")
	}

	log.Printf("roles str: %s", *roles)

	mux := http.NewServeMux()
	mux.HandleFunc("/default/api/operate/sandbox", func(w http.ResponseWriter, r *http.Request) {
		appName := r.URL.Query().Get("appName")
		for _, role := range strings.Split(*roles, ",") {

			checkUrl := "http://jobcloud.api:8899/" + role + "/api/app/" + appName
			resp, err := http.Get(checkUrl)
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintln(w, err)
				log.Printf("Get %s: %v", checkUrl, err)
				return
			}
			if resp.StatusCode == http.StatusOK {
				targetUrl := "http://jobcloud.api:8899/" + role + "/api/operate/sandbox?appName=" + appName
				log.Printf("%s is redirected to url: %s", "/default/api/operate/sandbox", targetUrl)
				http.Redirect(w, r, targetUrl, 307)
				return
			}
		}
		w.WriteHeader(404)
		fmt.Fprintln(w, "Can find app")
		log.Printf("Can find app %s", appName)
	})

	http.ListenAndServe(*addr, mux)
}
