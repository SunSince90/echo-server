// Copyright © 2021 Elis Lulja
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// All rights reserved.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	loreimpsum "gopkg.in/loremipsum.v1"
)

// GetOutboundIP Gets the preferred outbound ip of this machine
// https://stackoverflow.com/a/37382208
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func handle(w http.ResponseWriter, r *http.Request) {
	ip := GetOutboundIP()
	name, err := os.Hostname()
	msg := ""

	if err == nil {
		msg = fmt.Sprintf("Hello from %s (%s)", name, ip)
	} else {
		msg = fmt.Sprintf("Hello from %s", ip)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, msg)
}

func lorem(w http.ResponseWriter, r *http.Request) {
	paragraphs := 1
	queryPar := r.URL.Query().Get("paragraphs")
	if queryPar != "" {
		_pars, err := strconv.ParseInt(queryPar, 10, 32)
		if err == nil {
			paragraphs = int(_pars)
		}
	}

	w.WriteHeader(http.StatusOK)
	li := loreimpsum.New()
	fmt.Fprintf(w, "%s", li.Paragraphs(paragraphs))
}

func main() {
	ip := GetOutboundIP()
	r := mux.NewRouter()
	r.HandleFunc("/hey", handle)
	r.HandleFunc("/lorem-ipsum", lorem)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Printf("serving requests %s:%d\n", ip.String(), 8080)
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	log.Println("waiting for shutdown (CTRL+C)...")
	c := make(chan os.Signal, 1)
	shutdownChan := make(chan struct{})
	signal.Notify(c, os.Interrupt)
	<-c

	log.Println()
	log.Print("exit requested")

	go func() {
		defer close(shutdownChan)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	// Create a deadline to wait for.

	<-shutdownChan
	log.Println("goodbye!")
	os.Exit(0)
}
