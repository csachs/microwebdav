package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
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
	authMode := param("auth_mode", "basic")
	userName := param("user", "user")
	password := param("pass", randStr())

	possibleHost := ""

	if listen[0:1] != ":" {
		possibleHost = listen
	} else {
		possibleHost = "localhost" + listen
	}

	logger := log.New(os.Stderr, "", 0)

	logger.Print("# listening on ", listen)
	logger.Print("# serving", servePath)

	if authMode != "none" {
		authMode = "basic"
		logger.Print("# auth mode ", authMode, " credentials ", userName, " ", password)
	} else {
		logger.Print("# auth mode ", authMode)
	}

	logger.Print("# possible url ", "http://", possibleHost, "/")

	if authMode != "none" {
		logger.Print("# possible url ", "http://", userName, ":", password, "@", possibleHost, "/")
	}

	webdavHandler := &webdav.Handler{
		//fileSystem := webdav.NewMemFS() //new(webdav.FileSystem)
		FileSystem: webdav.Dir(servePath),
		// we use in-memory locking, so no one except this process should access the files
		LockSystem: webdav.NewMemLS(),
	}

	handler := httpauth.SimpleBasicAuth(userName, password)(webdavHandler)

	if authMode == "none" {
		handler = webdavHandler
	}

	chainedHandler := LogToStdout(handlers.ProxyHeaders(handler))

	http.Handle("/", chainedHandler)

	http.ListenAndServe(listen, nil)
}
