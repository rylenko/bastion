module sapphire

go 1.23.4

replace (
	github.com/rylenko/sapphire/pkg/shield => ./pkg/shield
	github.com/rylenko/sapphire/pkg/shieldprovider => ./pkg/shieldprovider
)

require (
	github.com/rylenko/sapphire/pkg/shield v0.0.0-00010101000000-000000000000
	github.com/rylenko/sapphire/pkg/shieldprovider v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
