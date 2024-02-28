package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/xackery/eqcp/db"
	"github.com/xackery/eqcp/item"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := db.Init(ctx)
	if err != nil {
		return fmt.Errorf("db.Init: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	mux.HandleFunc("/item/", item.View)
	fmt.Println("Listening on http://127.0.0.1:8080")
	return http.ListenAndServe(":8080", mux)
}
