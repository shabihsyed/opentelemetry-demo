// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"fmt"
	"strconv"

	pb "github.com/opentelemetry/opentelemetry-demo/src/product-catalog/genproto/oteldemo"
)

const (
	pageSizeMetadataKey      = "page-size"
	pageTokenMetadataKey     = "page-token"
	nextPageTokenMetadataKey = "next-page-token"
)

// paginateProducts returns the page of products selected by pageSize and
// pageToken, plus the token for the following page ("" when there is none).
// An empty pageSize means no limit and an empty pageToken means the first page,
// so with both empty every product is returned. The token is the offset of the
// first product in the page.
func paginateProducts(products []*pb.Product, pageSize, pageToken string) ([]*pb.Product, string, error) {
	if pageSize == "" && pageToken == "" {
		return products, "", nil
	}

	offset := 0
	if pageToken != "" {
		n, err := strconv.Atoi(pageToken)
		if err != nil || n < 0 || n > len(products) {
			return nil, "", fmt.Errorf("invalid %s %q", pageTokenMetadataKey, pageToken)
		}
		offset = n
	}

	end := len(products)
	if pageSize != "" {
		n, err := strconv.Atoi(pageSize)
		if err != nil || n <= 0 {
			return nil, "", fmt.Errorf("invalid %s %q: must be a positive integer", pageSizeMetadataKey, pageSize)
		}
		if n < end-offset {
			end = offset + n
		}
	}

	nextPageToken := ""
	if end < len(products) {
		nextPageToken = strconv.Itoa(end)
	}
	return products[offset:end], nextPageToken, nil
}
