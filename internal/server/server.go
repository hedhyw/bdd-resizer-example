package server

import (
	"image/jpeg"
	"net"
	"net/http"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/nfnt/resize"
)

const (
	maxProxySizeBytes = 20 * 1024 * 1024
	headerContentType = "Content-Type"
	contentTypeJPEG   = "image/jpeg"
)

type Server struct {
	handler http.Handler

	httpCl *http.Client
}

func New() *Server {
	mux := http.NewServeMux()
	s := &Server{
		handler: mux,
		httpCl:  http.DefaultClient,
	}

	mux.HandleFunc("/api/image.jpg", s.handleJPEGImage)

	return s
}

func (s *Server) Serve(l net.Listener) (err error) {
	return http.Serve(l, s.handler)
}

func (s *Server) ListenAndServer(addr string) (err error) {
	return http.ListenAndServe(addr, s.handler)
}

func (s *Server) handleJPEGImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	imageSizeRaw := q.Get("size")
	imageURL := q.Get("url")

	imageSize, err := parseSize(imageSizeRaw)
	if err != nil {
		respondErr(w, err)

		return
	}

	imgFile, err := s.fetchFile(ctx, imageURL)
	if err != nil {
		respondErr(w, err)

		return
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		respondErr(w, semerr.NewUnsupportedMediaTypeError(err))

		return
	}

	resizedImg := resize.Resize(
		uint(imageSize.Dx()),
		uint(imageSize.Dy()),
		img,
		resize.Lanczos3,
	)

	w.Header().Set(headerContentType, contentTypeJPEG)
	err = jpeg.Encode(w, resizedImg, nil)
	if err != nil {
		println("error: writing img response: " + err.Error())

		return
	}
}
