package backends

import (
	"context"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type HTTPBackend struct {
	ListeningAddress string // Listening address
	outChan          chan common.Hash
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
	if err := r.ParseForm(); err != nil {
		log.Println("Form error: %v", err)
		return
	}
	hash, ok := r.Form["hash"]
	if !ok || len(hash) == 0 {
		return
	}
	b.outChan <- common.HexToHash(hash[0])
}

func (b *HTTPBackend) Init(ctx context.Context, hashOut chan common.Hash) error {
	b.outChan = hashOut
	http.HandleFunc("/save", func(rw http.ResponseWriter, r *http.Request) {
		b.saveHandler(ctx, rw, r)
	})
	return nil
}

func (b *HTTPBackend) Run(ctx context.Context) error {
	log.Printf("Listening on %v", b.ListeningAddress)
	return http.ListenAndServe(b.ListeningAddress, nil)
}
