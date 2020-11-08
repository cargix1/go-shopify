package goshopify

import (
	"fmt"
	"net/http"
	"time"
)

const inventoryLevelsBasePath = "inventory_levels"

// InventoryLevelService is an interface for interacting with the
// inventory Levels endpoints of the Shopify API
// See https://shopify.dev/docs/admin-api/rest/reference/inventory/inventorylevel#index-2020-07
type InventoryLevelService interface {
	ListWithPagination(interface{}) ([]InventoryLevel, *Pagination, error)
}

// InventoryLevelServiceOp is the default implementation of the InventoryLevelService interface
type InventoryLevelServiceOp struct {
	client *Client
}

// InventoryLevel represents a Shopify inventory Level
type InventoryLevel struct {
	Available       int       `json:"available,omitempty"`
	InventoryItemID int64     `json:"inventory_item_id,omitempty"`
	LocationID      int64     `json:"location_id,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

// InventoryLevelResource is used for handling single Level requests and responses
type InventoryLevelResource struct {
	InventoryLevel *InventoryLevel `json:"inventory_Level"`
}

// InventoryLevelsResource is used for handling multiple Level responsees
type InventoryLevelsResource struct {
	InventoryLevels []InventoryLevel `json:"inventory_levels"`
}

// ListWithPagination lists inventory levels and return pagination to retrieve next/previous results.
func (s *InventoryLevelServiceOp) ListWithPagination(options interface{}) ([]InventoryLevel, *Pagination, error) {
	path := fmt.Sprintf("%s.json", inventoryLevelsBasePath)
	resource := new(InventoryLevelsResource)
	headers := http.Header{}

	headers, err := s.client.createAndDoGetHeaders("GET", path, nil, options, resource)
	if err != nil {
		return nil, nil, err
	}

	// Extract pagination info from header
	linkHeader := headers.Get("Link")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, nil, err
	}

	return resource.InventoryLevels, pagination, nil
}
