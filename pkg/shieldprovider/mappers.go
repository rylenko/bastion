package shieldprovider

import (
	"crypto/ecdh"
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

func (provider *Provider) mapFromForeignPrivateKey(foreignKey *ecdh.PrivateKey) (*keys.Private, error) {
	if foreignKey == nil {
		return nil, fmt.Errorf("%w: foreign key is nil", ErrInvalidValue)
	}

	key := keys.NewPrivate(foreignKey.Bytes())

	return key, nil
}

func (provider *Provider) mapToForeignPrivateKey(key *keys.Private) (*ecdh.PrivateKey, error) {
	if key == nil {
		return nil, fmt.Errorf("%w: key is nil", ErrInvalidValue)
	}

	foreignKey, err := provider.curve.NewPrivateKey(key.Bytes())
	if err != nil {
		return nil, err
	}

	return foreignKey, nil
}

func (provider *Provider) mapToForeignPublicKey(key *keys.Public) (*ecdh.PublicKey, error) {
	if key == nil {
		return nil, fmt.Errorf("%w: key is nil", ErrInvalidValue)
	}

	foreignKey, err := provider.curve.NewPublicKey(key.Bytes())
	if err != nil {
		return nil, err
	}

	return foreignKey, nil
}
