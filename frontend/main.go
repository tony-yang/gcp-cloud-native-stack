package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
)

const (
	defaultCurrency = "USD"
	cookieMaxAge    = 60 * 60 * 48
	cookiePrefix    = "demoshop_"
	cookieSessionID = cookiePrefix + "session-id"
)

var (
	port = "13000"
)

type ctxKeySessionID struct{}

type frontendServer struct {
	catalogAddr string
	catalogConn *grpc.ClientConn
}

func mapEnv(target *string, envKey string) error {
	v := os.Getenv(envKey)
	if v == "" {
		return fmt.Errorf("environment variable %q not set", envKey)
	}
	*target = v
	return nil
}

func connGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) error {
	var err error
	*conn, err = grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second*3),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
	if err != nil {
		return fmt.Errorf("grpc: failed to connect to %q: %v", addr, err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	addr := os.Getenv("LISTEN_ADDR")
	svc := new(frontendServer)
	if err := mapEnv(&svc.catalogAddr, "CATALOG_ADDR"); err != nil {
		log.Fatalf("failed to map component address %q to environment variable: %v", svc.catalogAddr, err)
	}
	if err := connGRPC(ctx, &svc.catalogConn, svc.catalogAddr); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	var handler http.Handler = r
	handler = ensureSessionID(handler)
	log.Infof("starting server on %s:%s", addr, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), handler))
}
