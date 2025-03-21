//go:build test
// +build test

package postgres

import "embed"

//go:embed "all:migrations"
var Migrations embed.FS

//go:embed "all:migrations/tenants/test"
var TestTenant embed.FS
