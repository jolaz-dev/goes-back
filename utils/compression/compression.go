package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

const gzipEncoding = "gzip"

func Deflate(r *http.Request, b64z string) ([]byte, error) {
	encoding := r.Header.Get("Content-Encoding")

	if encoding == "" {
		return nil, nil
	}

	if encoding == gzipEncoding {
		z, err := base64.StdEncoding.DecodeString(b64z)
		if err != nil {
			return nil, err
		}

		r, err := gzip.NewReader(bytes.NewReader(z))
		if err != nil {
			return nil, err
		}

		out, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		return out, nil
	}

	return nil, errors.New("unsupported encoding: " + encoding)
}

func Compress(w http.ResponseWriter, r *http.Request, body []byte, status int) {
	accept := r.Header.Get("Accept-Encoding")
	if !strings.Contains(accept, gzipEncoding) {
		writeUncompressed(w, body, status)
		return
	}

	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
	if err != nil {
		slog.Error("Error creating gzip writer", "error", err)
		writeUncompressed(w, body, status)
		return
	}

	if _, err := gz.Write(body); err != nil {
		slog.Error("Error writing to gzip writer", "error", err)
		writeUncompressed(w, body, status)
		return
	}

	if err := gz.Close(); err != nil {
		slog.Error("Error closing gzip writer", "error", err)
		writeUncompressed(w, body, status)
		return
	}

	w.Header().Set("Content-Encoding", gzipEncoding)
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func writeUncompressed(w http.ResponseWriter, body []byte, status int) {
	w.WriteHeader(status)
	w.Write(body)
}
