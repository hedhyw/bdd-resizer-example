package server

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hedhyw/semerr/pkg/v1/httperr"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

func (s Server) fetchFile(ctx context.Context, url string) (r io.Reader, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("preparing request: %w", err)

		return nil, semerr.NewBadRequestError(err)
	}

	resp, err := s.httpCl.Do(req)
	if err != nil {
		err = fmt.Errorf("doing request: %w", err)

		return nil, semerr.NewInternalServerError(err)
	}

	defer func() { err = semerr.NewMultiError(err, resp.Body.Close()) }()

	switch {
	case resp.StatusCode < http.StatusBadRequest:
		r := io.LimitReader(resp.Body, maxProxySizeBytes)

		data, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("reading body: %w", err)
		}

		return bytes.NewReader(data), nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, semerr.NewNotFoundError(semerr.Error("image not found"))
	default:
		return nil, semerr.Error("invalid status: " + resp.Status)
	}
}

func parseSize(val string) (image.Rectangle, error) {
	const errInvalidImageSize semerr.Error = "invalid size"

	toks := strings.Split(val, "x")
	if len(toks) != 2 {
		return image.Rectangle{}, semerr.NewBadRequestError(errInvalidImageSize)
	}

	width, err := strconv.Atoi(toks[0])
	if err != nil {
		err = fmt.Errorf("parsing width: %w", err)

		return image.Rectangle{}, semerr.NewBadRequestError(err)
	}

	height, err := strconv.Atoi(toks[1])
	if err != nil {
		err = fmt.Errorf("parsing height: %w", err)

		return image.Rectangle{}, semerr.NewBadRequestError(err)
	}

	if width <= 0 || height <= 0 {
		return image.Rectangle{}, semerr.NewBadRequestError(errInvalidImageSize)
	}

	return image.Rect(0, 0, width, height), nil
}

func respondErr(w http.ResponseWriter, err error) {
	code := httperr.Code(err)

	w.WriteHeader(code)

	_, err = w.Write([]byte(http.StatusText(code)))
	if err != nil {
		println("error: writing err response: " + err.Error())

		return
	}
}
