package main

import (
	"context"
	"testing"

	pb "github.com/tony-yang/gcp-cloud-native-stack/catalog/genproto"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	addr := run("0")
	log.Printf("Test Server listen at address: %s", addr)
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}))
	if err != nil {
		t.Error(err)
	}
	defer conn.Close()
	client := pb.NewProductCatalogServiceClient(conn)
	t.Run("Get all products from local file", func(t *testing.T) {
		res, err := client.ListProducts(ctx, &pb.Empty{})
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(res.Products, parseCatalog(), cmp.Comparer(proto.Equal)); diff != "" {
			t.Error(diff)
		}
	})

	t.Run("Get one product", func(t *testing.T) {
		got, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: "OLJCESPC7Z"})
		if err != nil {
			t.Error(err)
		}
		if want := parseCatalog()[0]; !proto.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Get non-exist product should returns not found", func(t *testing.T) {
		_, err := client.GetProduct(ctx, &pb.GetProductRequest{Id: "N/A"})
		if got, want := status.Code(err), codes.NotFound; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("Search product", func(t *testing.T) {
		res, err := client.SearchProducts(ctx, &pb.SearchProductsRequest{Query: "typewriter"})
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(res.Results, []*pb.Product{parseCatalog()[0]}, cmp.Comparer(proto.Equal)); diff != "" {
			t.Error(diff)
		}
	})
}
