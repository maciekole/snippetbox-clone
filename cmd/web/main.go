package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Define application struct to hold the application-wide dependencies for the web app
type application struct {
	logger *slog.Logger
}

func main() {
	// Define a new command-line flag with name 'addr', a default value of ":8080"
	addr := flag.String("addr", ":8080", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":8080". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to init a new structured logger, which writes to STDOUT with default settings
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	// Init a new instance of application struct containing the dependencies
	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()

	// Create a file server which servers files from static
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	logger.Info("starting server on ", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)

	// no logger.Fatal(), closest solution is to message Error and call os.Exit(1)
	logger.Error(err.Error())
	os.Exit(1)
}
