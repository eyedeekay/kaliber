/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/mwat56/apachelogger"
	"github.com/mwat56/errorhandler"
	"github.com/mwat56/kaliber"
	"github.com/mwat56/sessions"
)

// `setupSinals()` configures the capture of the interrupts `SIGINT`,
// `SIGKILL`, and `SIGTERM` to terminate the program gracefully.
func setupSinals(aServer *http.Server) {
	// handle `CTRL-C` and `kill(9)` and `kill(15)`.
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		for signal := range c {
			log.Printf("%s captured '%v', stopping program and exiting ...", os.Args[0], signal)
			if err := aServer.Shutdown(context.Background()); nil != err {
				log.Fatalf("%s: %v", os.Args[0], err)
			}
		}
	}()
} // setupSinals()

func main() {
	var (
		err           error
		handler       http.Handler
		ph            *kaliber.TPageHandler
		ck, cp, fn, s string
	)
	Me, _ := filepath.Abs(os.Args[0])
	if err = kaliber.DBopen(kaliber.CalibreDatabasePath()); nil != err {
		kaliber.ShowHelp()
		s = fmt.Sprintf("%s: %v", Me, err)
		apachelogger.Log("Kaliber/main", s)
		runtime.Gosched() // let the logger write
		log.Fatalln(s)
	}

	if fn, err = kaliber.AppArguments.Get("uf"); (nil == err) && (0 < len(fn)) {
		// All the following `xxxUser()` calls terminate the program
		if s, err = kaliber.AppArguments.Get("ua"); (nil == err) && (0 < len(s)) {
			kaliber.AddUser(s, fn)
		}
		if s, err = kaliber.AppArguments.Get("uc"); (nil == err) && (0 < len(s)) {
			kaliber.CheckUser(s, fn)
		}
		if s, err = kaliber.AppArguments.Get("ud"); (nil == err) && (0 < len(s)) {
			kaliber.DeleteUser(s, fn)
		}
		if s, err = kaliber.AppArguments.Get("ul"); (nil == err) && (0 < len(s)) {
			kaliber.ListUser(fn)
		}
		if s, err = kaliber.AppArguments.Get("uu"); (nil == err) && (0 < len(s)) {
			kaliber.UpdateUser(s, fn)
		}
	}

	if ph, err = kaliber.NewPageHandler(); nil != err {
		kaliber.ShowHelp()
		s = fmt.Sprintf("%s: %v", Me, err)
		apachelogger.Log("Kaliber/main", s)
		runtime.Gosched() // let the logger write
		log.Fatalln(s)
	}
	handler = errorhandler.Wrap(ph, ph)

	// inspect `sessiondir` commandline argument and setup the session handler
	if s, err = kaliber.AppArguments.Get("sessiondir"); (nil == err) && (0 < len(s)) {
		// we assume, an error means: no automatic session handling
		handler = sessions.Wrap(handler, s)
	}

	// inspect `logfile` commandline argument and setup the `ApacheLogger`
	if s, err = kaliber.AppArguments.Get("logfile"); (nil == err) && (0 < len(s)) {
		// we assume, an error means: no logfile
		handler = apachelogger.Wrap(handler, s)
	}

	// We need a `server` reference to use it in `setupSinals()` below
	// and to set some reasonable timeouts:
	server := &http.Server{
		Addr:              ph.Address(),
		Handler:           handler,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	setupSinals(server)

	ck, _ = kaliber.AppArguments.Get("certKey")
	cp, _ = kaliber.AppArguments.Get("certPem")

	if 0 < len(ck) && (0 < len(cp)) {
		s = fmt.Sprintf("%s listening HTTPS at %s", Me, ph.Address())
		log.Println(s)
		apachelogger.Log("Kaliber/main", s)
		if err = server.ListenAndServeTLS(cp, ck); nil != err {
			s = fmt.Sprintf("%s %v", Me, err)
			apachelogger.Log("Kaliber/main", s)
			runtime.Gosched() // let the logger write
			log.Fatalln(s)
		}
		return
	}

	s = fmt.Sprintf("%s listening HTTP at %s", Me, ph.Address())
	log.Println(s)
	apachelogger.Log("Kaliber/main", s)
	if err = server.ListenAndServe(); nil != err {
		s = fmt.Sprintf("%s %v", Me, err)
		apachelogger.Log("Kaliber/main", s)
		runtime.Gosched() // let the logger write
		log.Fatalln(s)
	}
} // main()

/* _EoF_ */
