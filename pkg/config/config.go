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

	ConfigurationFields = []field.SchemaField{
		UsernameField,
		PasswordField,
		ServerField,
		ServiceField,
		PortField,
	}

	ConfigurationSchema = field.Configuration{
		Fields: ConfigurationFields,
	}
)

//go:generate go run ./gen
var Config = field.NewConfiguration(ConfigurationFields)
