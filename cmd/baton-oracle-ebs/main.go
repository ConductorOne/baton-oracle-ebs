package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conductorone/baton-sdk/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/connectorrunner"
	"github.com/conductorone/baton-sdk/pkg/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	ebsconfig "github.com/conductorone/baton-oracle-ebs/pkg/config"
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
		ebsconfig.ConfigurationSchema,
		connectorrunner.WithDefaultCapabilitiesConnectorBuilder(&connector.OracleEBS{}),
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

func getConnector(ctx context.Context, v *viper.Viper) (types.ConnectorServer, error) {
	l := ctxzap.Extract(ctx)

	cfg := ebs.Config{
		Username: v.GetString(ebsconfig.UsernameField.FieldName),
		Password: v.GetString(ebsconfig.PasswordField.FieldName),
		Server:   v.GetString(ebsconfig.ServerField.FieldName),
		Service:  v.GetString(ebsconfig.ServiceField.FieldName),
		Port:     v.GetInt(ebsconfig.PortField.FieldName),
	}

	cb, err := connector.New(ctx, cfg)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	c, err := connectorbuilder.NewConnector(ctx, cb)
	if err != nil {
		l.Error("error creating connector", zap.Error(err))
		return nil, err
	}

	return c, nil
}
