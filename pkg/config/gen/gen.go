package main

import (
	cfg "github.com/conductorone/baton-oracle-ebs/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("oracle-ebs", cfg.Config)
}
