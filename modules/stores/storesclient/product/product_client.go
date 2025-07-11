// Code generated by go-swagger; DO NOT EDIT.

package product

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// New creates a new product API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

// New creates a new product API client with basic auth credentials.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - user: user for basic authentication header.
// - password: password for basic authentication header.
func NewClientWithBasicAuth(host, basePath, scheme, user, password string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BasicAuth(user, password)
	return &Client{transport: transport, formats: strfmt.Default}
}

// New creates a new product API client with a bearer token for authentication.
// It takes the following parameters:
// - host: http host (github.com).
// - basePath: any base path for the API client ("/v1", "/v3").
// - scheme: http scheme ("http", "https").
// - bearerToken: bearer token for Bearer authentication header.
func NewClientWithBearerToken(host, basePath, scheme, bearerToken string) ClientService {
	transport := httptransport.New(host, basePath, []string{scheme})
	transport.DefaultAuthentication = httptransport.BearerToken(bearerToken)
	return &Client{transport: transport, formats: strfmt.Default}
}

/*
Client for product API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption may be used to customize the behavior of Client methods.
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	AddProduct(params *AddProductParams, opts ...ClientOption) (*AddProductOK, error)

	DecreaseProductPrice(params *DecreaseProductPriceParams, opts ...ClientOption) (*DecreaseProductPriceOK, error)

	GetProduct(params *GetProductParams, opts ...ClientOption) (*GetProductOK, error)

	GetStoreProduct(params *GetStoreProductParams, opts ...ClientOption) (*GetStoreProductOK, error)

	IncreaseProductPrice(params *IncreaseProductPriceParams, opts ...ClientOption) (*IncreaseProductPriceOK, error)

	RebrandProduct(params *RebrandProductParams, opts ...ClientOption) (*RebrandProductOK, error)

	RemoveProduct(params *RemoveProductParams, opts ...ClientOption) (*RemoveProductOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
AddProduct adds a store product
*/
func (a *Client) AddProduct(params *AddProductParams, opts ...ClientOption) (*AddProductOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewAddProductParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "addProduct",
		Method:             "POST",
		PathPattern:        "/api/stores/{storeId}/products",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &AddProductReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*AddProductOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*AddProductDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
DecreaseProductPrice decreases the price of a product
*/
func (a *Client) DecreaseProductPrice(params *DecreaseProductPriceParams, opts ...ClientOption) (*DecreaseProductPriceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDecreaseProductPriceParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "decreaseProductPrice",
		Method:             "PUT",
		PathPattern:        "/api/stores/products/{id}/decreasePrice",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DecreaseProductPriceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*DecreaseProductPriceOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*DecreaseProductPriceDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetProduct gets a store product
*/
func (a *Client) GetProduct(params *GetProductParams, opts ...ClientOption) (*GetProductOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetProductParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getProduct",
		Method:             "GET",
		PathPattern:        "/api/stores/products/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetProductReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetProductOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetProductDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
GetStoreProduct gets a list of store products
*/
func (a *Client) GetStoreProduct(params *GetStoreProductParams, opts ...ClientOption) (*GetStoreProductOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetStoreProductParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getStoreProduct",
		Method:             "GET",
		PathPattern:        "/api/stores/{storeId}/products",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetStoreProductReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetStoreProductOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetStoreProductDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
IncreaseProductPrice increases the price of a product
*/
func (a *Client) IncreaseProductPrice(params *IncreaseProductPriceParams, opts ...ClientOption) (*IncreaseProductPriceOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewIncreaseProductPriceParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "increaseProductPrice",
		Method:             "PUT",
		PathPattern:        "/api/stores/products/{id}/increasePrice",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &IncreaseProductPriceReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*IncreaseProductPriceOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*IncreaseProductPriceDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
RebrandProduct changes the name and description of a product
*/
func (a *Client) RebrandProduct(params *RebrandProductParams, opts ...ClientOption) (*RebrandProductOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRebrandProductParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "rebrandProduct",
		Method:             "PUT",
		PathPattern:        "/api/stores/products/{id}/rebrand",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &RebrandProductReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*RebrandProductOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*RebrandProductDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
RemoveProduct removes a store product
*/
func (a *Client) RemoveProduct(params *RemoveProductParams, opts ...ClientOption) (*RemoveProductOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewRemoveProductParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "removeProduct",
		Method:             "DELETE",
		PathPattern:        "/api/stores/products/{id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &RemoveProductReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*RemoveProductOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*RemoveProductDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
