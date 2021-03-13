package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/justinas/alice"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"strings"
)

func LogToStdout(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}

const parameterPrefix = "MICROWEBDAV_"

func randStr() string {
	buffer := [32]byte{}
	rand.Read(buffer[:])
	return hex.EncodeToString(buffer[:])
}

func param(name string, defaultValue string) string {
	value, ok := os.LookupEnv(parameterPrefix + strings.ToUpper(name))
	if ok {
		return value
	} else {
		return defaultValue
	}
}

func main() {
	listen := param("listen", ":8000")
	servePath := param("path", ".")
	userName := param("user", "user")
	password := param("pass", randStr())

	logger := log.New(os.Stderr, "", 0)
	logger.Print("Listening on ", listen, ", serving \"", servePath, "\", credentials ", userName, " ", password)


	//fileSystem := webdav.NewMemFS() //new(webdav.FileSystem)
	fileSystem := webdav.Dir(servePath)

	// we use in-memory locking, so no one except this process should access the files
	lockSystem := webdav.NewMemLS()

	handler := new(webdav.Handler)

	handler.FileSystem = fileSystem
	handler.LockSystem = lockSystem

	chainedHandler := alice.New(
		LogToStdout,
		handlers.ProxyHeaders,
		httpauth.SimpleBasicAuth(userName, password)).
		Then(handler)

	http.Handle("/", chainedHandler)

	http.ListenAndServe(listen, nil)
}
