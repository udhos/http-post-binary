// Package main implements the tool.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const version = "0.0.1"

func main() {
	infof("http-post-binary %s", version)

	var showVersion bool
	var fullURL string
	var size int
	var contentType string

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&fullURL, "url", "http://localhost:8080/test", "url")
	flag.IntVar(&size, "size", 1000, "body size")
	flag.StringVar(&contentType, "contentType", "application/octet-stream", "content-type")
	flag.Parse()

	if showVersion {
		return
	}

	infof("request size=%d url=%s content_type=%s",
		size, fullURL, contentType)

	reqBuf := make([]byte, size, size)
	reqBodyReader := bytes.NewReader(reqBuf)

	resp, errPost := http.Post(fullURL, contentType, reqBodyReader)
	if errPost != nil {
		fatalf("post: size=%d url=%s error: %v", size, fullURL, errPost)
	}

	infof("response status: %d", resp.StatusCode)

	for k, v := range resp.Header {
		for _, vv := range v {
			infof("response header: %s: %s", k, vv)
		}
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		errorf("post: size=%d url=%s error: %v", size, fullURL, errBody)
	}

	fmt.Println(string(body))
}

func errorf(format string, a ...any) {
	slog.Error(fmt.Sprintf(format, a...))
}

func fatalf(format string, a ...any) {
	slog.Error("FATAL: " + fmt.Sprintf(format, a...))
	os.Exit(1)
}

func infof(format string, a ...any) {
	slog.Info(fmt.Sprintf(format, a...))
}
