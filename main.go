package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	app "gonkey-example/case-app/internal"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	db, err := app.NewDB(ctx, "postgres://service:service@localhost:6543/service?sslmode=disable")
	if err != nil {
		panic(err)
	}

	const nicksGenAddr = "https://names.drycodes.com"
	apiHandlers := app.NewAPI(db, nicksGenAddr)

	const addr = ":7700"
	fmt.Printf("Starting server listening on %s", addr)
	srv := http.Server{Addr: addr, Handler: apiHandlers}
	defer srv.Close()

	go func() {
		srv.ListenAndServe()
	}()
	<-ctx.Done()
}
