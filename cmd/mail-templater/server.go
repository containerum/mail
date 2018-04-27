package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"context"
	"net/http"
	"os/signal"

	"git.containerum.net/ch/mail-templater/pkg/router"
	"git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func initServer(c *cli.Context) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)
	for _, f := range c.GlobalFlagNames() {
		fmt.Fprintf(w, "Flag: %s\t Value: %s\n", f, c.String(f))
	}
	w.Flush()

	setupLogs(c)

	ts, err := getTemplatesStorage(c)
	defer ts.Close()
	exitOnErr(err)
	ms, err := getMessagesStorage(c)
	defer ms.Close()
	exitOnErr(err)
	us, err := getUpstream(c, ms)
	exitOnErr(err)
	uss, err := getUpstreamSimple(c, ms)
	exitOnErr(err)
	um, err := getUserManagerClient(c)
	exitOnErr(err)

	app := router.CreateRouter(&middleware.Services{
		TemplateStorage:   ts,
		MessagesStorage:   ms,
		Upstream:          us,
		UpstreamSimple:    uss,
		UserManagerClient: um,
	})

	// graceful shutdown support

	srv := http.Server{
		Addr:    ":" + c.String(portFlag),
		Handler: app,
	}

	go exitOnErr(srv.ListenAndServe())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infoln("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
