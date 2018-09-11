package backends

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Magicking/faktur-daemon/internal/db"
	"github.com/Magicking/faktur-daemon/merkle"
	"github.com/ethereum/go-ethereum/common"
)

type HTTPBackend struct {
	ListeningAddress string // Listening address
	outChan          chan []merkle.Hashable
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

func (b *HTTPBackend) getReceiptsByRootHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Form error: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	hash, ok := r.Form["target_hash"]
	if !ok || len(hash) == 0 {
		http.Error(w, "target_hash not found in parameters", 422)
		return
	}
	rcpts, err := db.GetReceiptsByHash(ctx, common.HexToHash(hash[0]))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if len(rcpts) == 0 {
		http.Error(w, "Receipt for hash "+hash[0]+" not found", 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rcpts)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (b *HTTPBackend) saveHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Form error: %v", err)
		return
	}
	hash, ok := r.Form["hash"]
	if !ok || len(hash) == 0 {
		return
	}
	var toHashs []merkle.Hashable
	for _, e := range hash {
		toHashs = append(toHashs, common.HexToHash(e))
	}
	b.outChan <- toHashs
}

func (b *HTTPBackend) Init(ctx context.Context, hashOut chan []merkle.Hashable) error {
	b.outChan = hashOut
	http.HandleFunc("/getreceiptsbyroot", func(w http.ResponseWriter, r *http.Request) {
		b.getReceiptsByRootHandler(ctx, w, r)
	})
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		b.saveHandler(ctx, w, r)
	})
	return nil
}

func (b *HTTPBackend) Run(ctx context.Context) error {
	log.Printf("Listening on %v", b.ListeningAddress)
	return http.ListenAndServe(b.ListeningAddress, nil)
}
