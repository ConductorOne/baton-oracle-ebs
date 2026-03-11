package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/conductorone/baton-oracle-scm/pkg/client"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
)

type Connector struct {
	fusionClient *client.FusionClient
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from Oracle Fusion Cloud SCM.
func (c *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(c.fusionClient),
		newRoleBuilder(c.fusionClient),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated HTTP client.
func (c *Connector) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (c *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Oracle SCM",
		Description: "Connector for Oracle Fusion Cloud Supply Chain Management (SCM) that syncs users and roles via the Fusion Cloud REST APIs.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured.
func (c *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	err := c.fusionClient.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("baton-oracle-scm: failed to validate connection: %w", err)
	}
	return nil, nil
}

// New returns a new instance of the Oracle SCM connector.
func New(ctx context.Context, instanceURL, username, password string) (*Connector, error) {
	fusionClient, err := client.NewFusionClient(ctx, instanceURL, username, password)
	if err != nil {
		return nil, fmt.Errorf("baton-oracle-scm: failed to create fusion client: %w", err)
	}

	return &Connector{fusionClient: fusionClient}, nil
}
