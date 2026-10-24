// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/opentelemetry/opentelemetry-demo/src/product-catalog/genproto/oteldemo"
)

func testProducts(n int) []*pb.Product {
	products := make([]*pb.Product, n)
	for i := range products {
		products[i] = &pb.Product{Id: fmt.Sprintf("P%d", i)}
	}
	return products
}

func productIDs(products []*pb.Product) []string {
	ids := make([]string, len(products))
	for i, p := range products {
		ids[i] = p.GetId()
	}
	return ids
}

func TestPaginateProducts(t *testing.T) {
	tests := []struct {
		name          string
		catalogSize   int
		pageSize      string
		pageToken     string
		wantIDs       []string
		wantNextToken string
	}{
		{
			name:          "first page",
			catalogSize:   5,
			pageSize:      "2",
			wantIDs:       []string{"P0", "P1"},
			wantNextToken: "2",
		},
		{
			name:          "middle page",
			catalogSize:   5,
			pageSize:      "2",
			pageToken:     "2",
			wantIDs:       []string{"P2", "P3"},
			wantNextToken: "4",
		},
		{
			name:        "last page",
			catalogSize: 5,
			pageSize:    "2",
			pageToken:   "4",
			wantIDs:     []string{"P4"},
		},
		{
			name:        "last page exactly full",
			catalogSize: 4,
			pageSize:    "2",
			pageToken:   "2",
			wantIDs:     []string{"P2", "P3"},
		},
		{
			name:        "page size larger than catalog",
			catalogSize: 3,
			pageSize:    "10",
			wantIDs:     []string{"P0", "P1", "P2"},
		},
		{
			name:        "empty catalog with page size",
			catalogSize: 0,
			pageSize:    "2",
			wantIDs:     []string{},
		},
		{
			name:        "empty catalog without paging",
			catalogSize: 0,
			wantIDs:     []string{},
		},
		{
			name:        "no paging returns everything",
			catalogSize: 5,
			wantIDs:     []string{"P0", "P1", "P2", "P3", "P4"},
		},
		{
			name:        "token without page size returns the rest",
			catalogSize: 5,
			pageToken:   "3",
			wantIDs:     []string{"P3", "P4"},
		},
		{
			name:        "token at end of catalog returns empty page",
			catalogSize: 5,
			pageSize:    "2",
			pageToken:   "5",
			wantIDs:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, next, err := paginateProducts(testProducts(tt.catalogSize), tt.pageSize, tt.pageToken)
			if err != nil {
				t.Fatalf("paginateProducts() error = %v", err)
			}
			if got := productIDs(page); fmt.Sprint(got) != fmt.Sprint(tt.wantIDs) {
				t.Errorf("paginateProducts() page = %v, want %v", got, tt.wantIDs)
			}
			if next != tt.wantNextToken {
				t.Errorf("paginateProducts() next token = %q, want %q", next, tt.wantNextToken)
			}
		})
	}
}

func TestPaginateProductsInvalid(t *testing.T) {
	tests := []struct {
		name      string
		pageSize  string
		pageToken string
	}{
		{name: "zero page size", pageSize: "0"},
		{name: "negative page size", pageSize: "-1"},
		{name: "non-numeric page size", pageSize: "ten"},
		{name: "negative token", pageSize: "2", pageToken: "-1"},
		{name: "non-numeric token", pageSize: "2", pageToken: "abc"},
		{name: "token past end of catalog", pageSize: "2", pageToken: "6"},
		{name: "invalid token without page size", pageToken: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, next, err := paginateProducts(testProducts(5), tt.pageSize, tt.pageToken)
			if err == nil {
				t.Fatalf("paginateProducts() error = nil, page = %v, next = %q", productIDs(page), next)
			}
		})
	}
}

func TestSingleMetadataValue(t *testing.T) {
	tests := []struct {
		name    string
		md      metadata.MD
		want    string
		wantErr bool
	}{
		{name: "nil metadata", md: nil, want: ""},
		{name: "key absent", md: metadata.Pairs("other", "x"), want: ""},
		{name: "single value", md: metadata.Pairs(pageSizeMetadataKey, "3"), want: "3"},
		{name: "key is case-insensitive", md: metadata.Pairs("Page-Size", "3"), want: "3"},
		{
			name:    "repeated key",
			md:      metadata.Pairs(pageSizeMetadataKey, "3", pageSizeMetadataKey, "4"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := singleMetadataValue(tt.md, pageSizeMetadataKey)
			if (err != nil) != tt.wantErr {
				t.Fatalf("singleMetadataValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("singleMetadataValue() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListProductsRejectsRepeatedPagingMetadata(t *testing.T) {
	tests := []struct {
		name string
		md   metadata.MD
	}{
		{
			name: "repeated page size",
			md:   metadata.Pairs(pageSizeMetadataKey, "2", pageSizeMetadataKey, "3"),
		},
		{
			name: "repeated page token",
			md:   metadata.Pairs(pageTokenMetadataKey, "0", pageTokenMetadataKey, "2"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), tt.md)
			_, err := (&productCatalog{}).ListProducts(ctx, &pb.Empty{})
			if got := status.Code(err); got != codes.InvalidArgument {
				t.Fatalf("ListProducts() code = %v, want %v (err = %v)", got, codes.InvalidArgument, err)
			}
		})
	}
}
