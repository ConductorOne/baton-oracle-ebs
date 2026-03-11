package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	InstanceURLField = field.StringField(
		"instance-url",
		field.WithDisplayName("Instance URL"),
		field.WithDescription("The Oracle Fusion Cloud instance URL (e.g. https://servername.fa.us2.oraclecloud.com)."),
		field.WithRequired(true),
	)
	UsernameField = field.StringField(
		"username",
		field.WithDisplayName("Username"),
		field.WithDescription("The Oracle Fusion Cloud username for API access."),
		field.WithRequired(true),
	)
	PasswordField = field.StringField(
		"password",
		field.WithDisplayName("Password"),
		field.WithDescription("The Oracle Fusion Cloud password for API access."),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)

	// Add the SchemaFields for the Config.
	ConfigurationFields = []field.SchemaField{
		InstanceURLField,
		UsernameField,
		PasswordField,
	}

	// FieldRelationships defines relationships between the ConfigurationFields that can be automatically validated.
	FieldRelationships = []field.SchemaFieldRelationship{}
)

//go:generate go run -tags=generate ./gen
var Config = field.NewConfiguration(
	ConfigurationFields,
	field.WithConstraints(FieldRelationships...),
	field.WithConnectorDisplayName("Oracle SCM"),
	field.WithHelpUrl("/docs/baton/oracle-scm"),
	field.WithIconUrl("/static/app-icons/oracle.svg"),
)
