package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/thefloweringash/hass_ir_adapter/emitters/irkit"
)

func main() {
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		msg := irkit.Message{}
		err = json.Unmarshal(bytes, &msg)
		if err != nil {
			w.WriteHeader(400)
			_, _ = fmt.Fprintf(w, "json decode error: %v", err)
			return
		}

		// TODO: more checks
		//  - interval length is odd
		//  - mandatory header included

		fmt.Printf("message: %v\n", msg)
	})

	var listenAddr string
	flag.StringVar(&listenAddr, "listen", ":8000", "listen address")
	flag.Parse()

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
