//go:build generate

package main

import (
	cfg "github.com/conductorone/baton-oracle-scm/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("oracle-scm", cfg.Config)
}
