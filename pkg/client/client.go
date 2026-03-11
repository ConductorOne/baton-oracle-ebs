package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	DefaultPageSize uint = 100

	// Oracle Fusion Cloud REST API base paths.
	// HCM REST API for user accounts and roles.
	hcmBasePath = "/hcmRestApi/resources/latest"
	// FSCM REST API for SCM-specific resources.
	fscmBasePath = "/fscmRestApi/resources/latest"

	// Endpoints.
	userAccountsPath     = "/userAccounts"
	userAccountRolesPath = "/userAccounts/{{.userID}}/child/userAccountRoles"
	rolesPath            = "/roles"
)

// FusionClient handles REST API operations for Oracle Fusion Cloud SCM.
type FusionClient struct {
	httpClient  *uhttp.BaseHttpClient
	instanceURL *url.URL
}

// NewFusionClient creates a new Oracle Fusion Cloud REST API client.
// Oracle Fusion Cloud REST APIs support Basic Authentication over SSL.
func NewFusionClient(ctx context.Context, instanceURLStr, username, password string) (*FusionClient, error) {
	instanceURL, err := url.Parse(instanceURLStr)
	if err != nil {
		return nil, fmt.Errorf("baton-oracle-scm: error parsing instance URL: %w", err)
	}

	// Create HTTP client with Basic Auth.
	transport := &basicAuthTransport{
		username: username,
		password: password,
		base:     http.DefaultTransport,
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	client := &FusionClient{
		httpClient:  uhttp.NewBaseHttpClient(httpClient),
		instanceURL: instanceURL,
	}

	return client, nil
}

// basicAuthTransport adds Basic Auth header to every request.
type basicAuthTransport struct {
	username string
	password string
	base     http.RoundTripper
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.username, t.password)
	return t.base.RoundTrip(req)
}

// Validate tests the API connection by fetching a single user account.
func (c *FusionClient) Validate(ctx context.Context) error {
	opts := []URLOption{
		WithLimit(1),
		WithOffset(0),
	}
	u, err := c.constructURL(hcmBasePath, userAccountsPath, opts...)
	if err != nil {
		return fmt.Errorf("baton-oracle-scm: failed to construct validation URL: %w", err)
	}

	var response FusionListResponse[UserAccount]
	_, err = c.get(ctx, u, &response)
	if err != nil {
		return fmt.Errorf("baton-oracle-scm: failed to validate connection: %w", err)
	}

	return nil
}

// ListUserAccounts lists user accounts from Oracle Fusion Cloud HCM REST API.
func (c *FusionClient) ListUserAccounts(ctx context.Context, offset uint) ([]*UserAccount, bool, *v2.RateLimitDescription, error) {
	logger := ctxzap.Extract(ctx)
	logger.Debug("listUserAccounts called", zap.Uint("offset", offset))

	opts := []URLOption{
		WithOffset(offset),
		WithLimit(DefaultPageSize),
		WithQueryParams(map[string]string{
			"onlyData": "true",
		}),
	}

	u, err := c.constructURL(hcmBasePath, userAccountsPath, opts...)
	if err != nil {
		return nil, false, nil, fmt.Errorf("baton-oracle-scm: error constructing user accounts URL: %w", err)
	}

	var response FusionListResponse[UserAccount]
	rateLimit, err := c.get(ctx, u, &response)
	if err != nil {
		return nil, false, rateLimit, fmt.Errorf("baton-oracle-scm: error listing user accounts: %w", err)
	}

	logger.Debug("user accounts response received",
		zap.Int("count", response.Count),
		zap.Bool("hasMore", response.HasMore),
		zap.Int("itemCount", len(response.Items)))

	return response.Items, response.HasMore, rateLimit, nil
}

// ListRoles lists roles from Oracle Fusion Cloud Common Features REST API.
func (c *FusionClient) ListRoles(ctx context.Context, offset uint) ([]*Role, bool, *v2.RateLimitDescription, error) {
	logger := ctxzap.Extract(ctx)
	logger.Debug("listRoles called", zap.Uint("offset", offset))

	opts := []URLOption{
		WithOffset(offset),
		WithLimit(DefaultPageSize),
		WithQueryParams(map[string]string{
			"onlyData": "true",
		}),
	}

	u, err := c.constructURL(hcmBasePath, rolesPath, opts...)
	if err != nil {
		return nil, false, nil, fmt.Errorf("baton-oracle-scm: error constructing roles URL: %w", err)
	}

	var response FusionListResponse[Role]
	rateLimit, err := c.get(ctx, u, &response)
	if err != nil {
		return nil, false, rateLimit, fmt.Errorf("baton-oracle-scm: error listing roles: %w", err)
	}

	logger.Debug("roles response received",
		zap.Int("count", response.Count),
		zap.Bool("hasMore", response.HasMore),
		zap.Int("itemCount", len(response.Items)))

	return response.Items, response.HasMore, rateLimit, nil
}

// ListUserRoles lists roles assigned to a specific user.
func (c *FusionClient) ListUserRoles(ctx context.Context, userID int64, offset uint) ([]*UserRole, bool, *v2.RateLimitDescription, error) {
	logger := ctxzap.Extract(ctx)
	logger.Debug("listUserRoles called", zap.Int64("userID", userID), zap.Uint("offset", offset))

	opts := []URLOption{
		WithPathParams(map[string]string{"userID": strconv.FormatInt(userID, 10)}),
		WithOffset(offset),
		WithLimit(DefaultPageSize),
		WithQueryParams(map[string]string{
			"onlyData": "true",
		}),
	}

	u, err := c.constructURL(hcmBasePath, userAccountRolesPath, opts...)
	if err != nil {
		return nil, false, nil, fmt.Errorf("baton-oracle-scm: error constructing user roles URL: %w", err)
	}

	var response FusionListResponse[UserRole]
	rateLimit, err := c.get(ctx, u, &response)
	if err != nil {
		return nil, false, rateLimit, fmt.Errorf("baton-oracle-scm: error listing user roles: %w", err)
	}

	logger.Debug("user roles response received",
		zap.Int64("userID", userID),
		zap.Int("count", response.Count),
		zap.Bool("hasMore", response.HasMore),
		zap.Int("itemCount", len(response.Items)))

	return response.Items, response.HasMore, rateLimit, nil
}
