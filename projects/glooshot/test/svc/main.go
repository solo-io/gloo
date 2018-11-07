package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func toInt(s string) int {
	n, _ := strconv.ParseInt(s, 0, 0)
	return int(n)
}

func main() {

	me := toInt(os.Getenv("SELF"))
	total := toInt(os.Getenv("TOTAL"))
	selfPort := os.Getenv("PORT")
	var ports []int
	for i := 0; i < total; i++ {
		svcPort := toInt(os.Getenv(fmt.Sprintf("SVC%d", i)))
		ports = append(ports, svcPort)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if total == 0 {
			w.Write([]byte("no total"))
			return
		}
		if me+1 == total {
			w.Write([]byte("hello"))
		} else {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d", ports[me+1]))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("can't make request"))
				return
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("can't read body"))
				return
			}

			datastr := string(data)
			datastr = datastr + datastr
			w.Write([]byte(datastr))
		}
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":"+selfPort, nil))
}
