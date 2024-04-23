package main

import (
	"context"

	"github.com/conductorone/baton-oracle-ebs/pkg/ebs"
	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	Config ebs.Config `mapstructure:",squash"`

	// TODO: later add more configuration options here
	// PublicKeyPath string `mapstructure:"public-key-path"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.Config.Username == "" || cfg.Config.Password == "" {
		return status.Error(codes.InvalidArgument, "username and password are required")
	}

	return nil
}

func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("username", "", "Username for the Oracle EBS Database connection. ($BATON_USERNAME)")
	cmd.PersistentFlags().String("password", "", "Password for the Oracle EBS Database connection. ($BATON_PASSWORD)")
	cmd.PersistentFlags().String("server", "", "Server for the Oracle EBS connection. ($BATON_SERVER)")
	cmd.PersistentFlags().String("service", "", "Service for the Oracle EBS connection. ($BATON_SERVICE)")
	cmd.PersistentFlags().Int("port", 0, "Port for the Oracle EBS connection. ($BATON_PORT)")
}
