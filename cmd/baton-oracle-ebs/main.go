package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	cfg "github.com/conductorone/baton-oracle-ebs/pkg/config"
	"github.com/conductorone/baton-oracle-ebs/pkg/connector"
	"github.com/conductorone/baton-oracle-ebs/pkg/ebs"
)

var version = "dev"

func main() {
	ctx := context.Background()

	_, cmd, err := config.DefineConfiguration(
		ctx,
		"baton-oracle-ebs",
		getConnector,
		cfg.Config,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmd.Version = version

	err = cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getConnector(ctx context.Context, c *cfg.OracleEbs) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	ebsCfg := ebs.Config{
		Username: c.Username,
		Password: c.Password,
		Server:   c.Server,
		Service:  c.Service,
		Port:     c.Port,
		BaseURL:  c.BaseUrl,
	}

	cb, err := connector.New(ctx, ebsCfg)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	conn, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
