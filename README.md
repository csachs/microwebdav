# microwebdav

A very tiny, basic WebDAV server, written in Go. Mostly gluing together a few libraries in a few lines.

**Warning**: This is code quickly written while not being very familiar with Go yet ... use at your own risk.
In particular, take security into consideration when running this on a publicly accessible IP address.
I do not know how secure this code, or the libraries used, are.

## Rationale

If one has a HTTP library ready, WebDAV is a very easy protocol to speak, and it implemented in a plethora of languages.
While client support is widespread and easy, *server* support is often a complex matter, so if one wants to quickly
allow reading/writing of a particular directory, one would need to set up one of the feature complete solutions with
a lot of configuration.

In the spirit of `python -m http.server`, `microwebdav` aims to just do one thing: serve a directory read/write via WebDAV.

## Usage

To be easily configured via Docker, `microwebdav` uses environment variables for configuration.

By default, microwebdav will serve the current directory at port 8000, with username *user* and a random generated password shown at startup.

Currently, no SSL support, this can be facilitated by a reverse proxy though.
If desired, basic auth can be *disabled* by setting `MICROWEBDAV_AUTH_MODE` to `none`.

```bash
> ./microwebdav 
Listening on :8000, serving ".", credentials user b7ead73bc0d35ff23f3c7196a9173be7d2053a246b7849468aa2d325e07269c9

> MICROWEBDAV_LISTEN=:8888 MICROWEBDAV_PATH=/tmp MICROWEBDAV_USER=anotheruser MICROWEBDAV_PASS=cookie ./microwebdav
Listening on :8888, serving "/tmp", credentials anotheruser cookie
```


## Docker

`microwebdav` and its Docker image is around 9 MB compiled.
It can be useful in Docker setups, where it can e.g. in orchestrated situations help to give external access to certain volumes.
Bear in mind that no locking is happening, though, when parallel access could happen.

```bash
# run it to serve the local directory 
> docker run --rm -v `pwd`:/data -p 8000:8000 ghcr.io/csachs/microwebdav:latest
# run it and detach
> docker run -d -p 8000:8000 ghcr.io/csachs/microwebdav:latest
```

## License

MIT
