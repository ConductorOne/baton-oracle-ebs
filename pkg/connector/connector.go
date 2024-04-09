package connector

import (
	"context"
	"fmt"
	"io"

	"github.com/conductorone/baton-oracle-ebs/pkg/ebs"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"

	ora "github.com/sijms/go-ora/v2"
)

type OracleEBS struct {
	client *ebs.Client
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (o *OracleEBS) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newUserBuilder(o.client),
		newRoleBuilder(o.client),
	}
}

// Asset takes an input AssetRef and attempts to fetch it using the connector's authenticated http client
// It streams a response, always starting with a metadata object, following by chunked payloads for the asset.
func (o *OracleEBS) Asset(ctx context.Context, asset *v2.AssetRef) (string, io.ReadCloser, error) {
	return "", nil, nil
}

// Metadata returns metadata about the connector.
func (o *OracleEBS) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "OracleEBS",
		Description: "Connector syncing EBS users and roles to Baton",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (o *OracleEBS) Validate(ctx context.Context) (annotations.Annotations, error) {
	err := o.client.Conn.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	defer o.client.Conn.Close()

	err = o.client.Conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping connection: %w", err)
	}

	return nil, nil
}

const (
	DefaultPort    = 1521
	DefaultServer  = "apps.example.com"
	DefaultService = "EBSDB"
)

// New returns a new instance of the connector.
func New(ctx context.Context, cfg ebs.Config) (*OracleEBS, error) {
	var port int
	var server, service string

	if cfg.Port == 0 {
		port = DefaultPort
	} else {
		port = cfg.Port
	}

	if cfg.Server == "" {
		server = DefaultServer
	} else {
		server = cfg.Server
	}

	if cfg.Service == "" {
		service = DefaultService
	} else {
		service = cfg.Service
	}

	connString := ora.BuildUrl(
		server,
		port,
		service,
		cfg.Username,
		cfg.Password,
		nil, // TODO: add trace file for debugging?
	)

	conn, err := ora.NewConnection(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return &OracleEBS{
		client: ebs.NewEBSClient(conn),
	}, nil
}
