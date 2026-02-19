package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	UsernameField = field.StringField(
		"username",
		field.WithDescription("Username for the Oracle EBS Database connection"),
		field.WithRequired(true),
	)

	PasswordField = field.StringField(
		"password",
		field.WithDescription("Password for the Oracle EBS Database connection"),
		field.WithIsSecret(true),
		field.WithRequired(true),
	)

	ServerField = field.StringField(
		"server",
		field.WithDescription("Server for the Oracle EBS connection"),
	)

	ServiceField = field.StringField(
		"service",
		field.WithDescription("Service for the Oracle EBS connection"),
	)

	PortField = field.IntField(
		"port",
		field.WithDescription("Port for the Oracle EBS connection"),
	)

	BaseURLField = field.StringField(
		"base-url",
		field.WithDescription("Override the Oracle EBS connection URL (for testing)"),
		field.WithHidden(true),
		field.WithExportTarget(field.ExportTargetCLIOnly),
	)

	ConfigurationFields = []field.SchemaField{
		UsernameField,
		PasswordField,
		ServerField,
		ServiceField,
		PortField,
		BaseURLField,
	}

	ConfigurationSchema = field.Configuration{
		Fields: ConfigurationFields,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(ConfigurationFields)
