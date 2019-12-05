package main

import (
	"context"

	pb "github.com/tony-yang/gcp-cloud-native-stack/frontend/genproto"
)

func (f *frontendServer) getProducts(ctx context.Context) ([]*pb.Product, error) {
	resp, err := pb.NewProductCatalogServiceClient(f.catalogConn).ListProducts(ctx, &pb.Empty{})
	return resp.GetProducts(), err
}
