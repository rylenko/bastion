module bastion

go 1.23.4

replace github.com/rylenko/bastion/pkg/ratchet => ./pkg/ratchet

require (
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
