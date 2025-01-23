module github.com/rylenko/sapphire/pkg/shieldprovider

go 1.23.4

replace github.com/rylenko/sapphire/pkg/shield => ../shield

require (
	github.com/rylenko/sapphire/pkg/shield v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.32.0
)

require golang.org/x/sys v0.29.0 // indirect
