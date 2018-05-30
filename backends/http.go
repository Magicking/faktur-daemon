package backends

import (
	"context"
	"log"
	"net/http"
)

type HTTPBackend struct {
	ListeningAddress string // Listening address
}

var HTTPopts struct {
	ListeningAddress string `long:"listening-address" env:"HTTP_LISTENING_ADDR" default:":8090" description:"Listening address (e.g: :8090)"`
}

func NewHTTP(ctx context.Context) (*HTTPBackend, error) {
	be := HTTPBackend{
		ListeningAddress: HTTPopts.ListeningAddress,
	}
	return &be, nil
}

func (b *HTTPBackend) saveHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	log.Println("saveHandler", rw, r)
}

func (b *HTTPBackend) Init(ctx context.Context) error {
	http.HandleFunc("/save", func(rw http.ResponseWriter, r *http.Request) {
		b.saveHandler(ctx, rw, r)
	})
	return nil
}

func (b *HTTPBackend) Run(ctx context.Context, ctl chan struct{}) error {
	return http.ListenAndServe(b.ListeningAddress, nil)
}
