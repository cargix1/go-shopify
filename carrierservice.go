package goshopify

import (
	"fmt"
)

const carrierBasePath = "admin/carrier_services"

// CarrierService is an interface for interfacing with the carrier service endpoints
// of the Shopify API.
// See: https://help.shopify.com/api/reference/carrierservice
type CarrierService interface {
	List(interface{}) ([]Carrier, error)
	Get(int64, interface{}) (*Carrier, error)
	Create(Carrier) (*Carrier, error)
	Delete(int64) (*Carrier, error)
}

// CarrierServiceOp handles communication with the product related methods of
// the Shopify API.
type CarrierServiceOp struct {
	client *Client
}

// Carrier represents a Shopify carrier service
type Carrier struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	Active             bool   `json:"active"`
	ServiceDiscovery   bool   `json:"service_discovery"`
	CarrierServiceType string `json:"carrier_service_type"`
	Format             string `json:"format"`
	CallbackUrl        string `json:"callback_url"`
}

// Represents the result from the carrier_services/X.json endpoint
type CarrierResource struct {
	Carrier *Carrier `json:"carrier_service"`
}

// Represents the result from the carrier_services.json endpoint
type CarriersResource struct {
	Carriers []Carrier `json:"carrier_services"`
}

// List carrier services
func (s *CarrierServiceOp) List(options interface{}) ([]Carrier, error) {
	path := fmt.Sprintf("%s.json", carrierBasePath)
	resource := new(CarriersResource)
	err := s.client.Get(path, resource, options)
	return resource.Carriers, err
}

// Get carrier service
func (s *CarrierServiceOp) Get(carrierServiceID int64, options interface{}) (*Carrier, error) {
	path := fmt.Sprintf("%s/%v.json", carrierBasePath, carrierServiceID)
	resource := new(CarrierResource)
	err := s.client.Get(path, resource, options)
	return resource.Carrier, err
}

// Create a new carrier service
func (s *CarrierServiceOp) Create(carrier Carrier) (*Carrier, error) {
	path := fmt.Sprintf("%s.json", carrierBasePath)
	wrappedData := CarrierResource{Carrier: &carrier}
	resource := new(CarrierResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Carrier, err
}

// Delete carrier service
func (s *CarrierServiceOp) Delete(carrierServiceID int64) error {
	path := fmt.Sprintf("%s/%v.json", carrierBasePath, carrierServiceID)
	return s.client.Delete(path)
}
