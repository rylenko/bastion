module bastion

go 1.23.4

replace (
	github.com/rylenko/bastion/pkg/ratchet => ./pkg/ratchet
	github.com/rylenko/bastion/pkg/utils => ./pkg/utils
)

require github.com/rylenko/bastion/pkg/ratchet v0.0.0-00010101000000-000000000000

require (
	github.com/rylenko/bastion/pkg/utils v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
