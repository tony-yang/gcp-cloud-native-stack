package main

import (
	"context"
	"fmt"

	pb "github.com/tony-yang/gcp-cloud-native-stack/frontend/genproto"
)

func (f *frontendServer) getProducts(ctx context.Context) ([]*pb.Product, error) {
	resp, err := pb.NewProductCatalogServiceClient(f.catalogConn).ListProducts(ctx, &pb.Empty{})
	return resp.GetProducts(), err
}

func (f *frontendServer) getProduct(ctx context.Context, id string) (*pb.Product, error) {
	resp, err := pb.NewProductCatalogServiceClient(f.catalogConn).GetProduct(ctx, &pb.GetProductRequest{Id: id})
	return resp, err
}

func (f *frontendServer) getRecommendations(ctx context.Context, userID string, productIDs []string) ([]*pb.Product, error) {
	resp, err := pb.NewRecommendationServiceClient(f.recommendationConn).ListRecommendations(ctx, &pb.ListRecommendationsRequest{UserId: userID, ProductIds: productIDs})
	if err != nil {
		return nil, err
	}
	var out []*pb.Product
	for i, pid := range resp.GetProductIds() {
		p, err := f.getProduct(ctx, pid)
		if err != nil {
			return nil, fmt.Errorf("failed to get recommended product info (#%s): %w", pid, err)
		}
		out = append(out, p)

		// Take only first four product recommendations to fit the UI
		if i >= 3 {
			break
		}
	}
	return out, err
}
