package rootchain

import (
	"errors"
	"reflect"
	"testing"

	"github.com/rylenko/bastion/pkg/ratchet/errlist"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type testCrypto struct{}

func (tc testCrypto) AdvanceChain(_ keys.Root, _ keys.Shared) (keys.Root, keys.MessageMaster, keys.Header, error) {
	return keys.Root{}, keys.MessageMaster{}, keys.Header{}, nil
}

func TestNewConfigDefault(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig(nil)
	if err != nil {
		t.Fatalf("newConfig() expected no error but got %v", err)
	}

	if cfg.crypto == nil {
		t.Fatal("newConfig() sets no default value for crypto")
	}
}

func TestNewConfigWithCryptoSuccess(t *testing.T) {
	t.Parallel()

	cfg, err := newConfig([]Option{WithCrypto(testCrypto{})})
	if err != nil {
		t.Fatalf("newConfig() with options expected no error but got %v", err)
	}

	if reflect.TypeOf(cfg.crypto) != reflect.TypeOf(testCrypto{}) {
		t.Fatal("WithCrypto() option did not set passed crypto")
	}
}

func TestNewConfigWithCryptoError(t *testing.T) {
	t.Parallel()

	_, err := newConfig([]Option{WithCrypto(nil)})
	if err == nil || err.Error() != "option: invalid value: crypto is nil" {
		t.Fatalf("WithCrypto(nil) expected error but got %v", err)
	}

	if !errors.Is(err, errlist.ErrOption) {
		t.Fatalf("WithCrypto(nil) error is not option error but %v", err)
	}

	if !errors.Is(err, errlist.ErrInvalidValue) {
		t.Fatalf("WithCrypto(nil) error is not invalid value error but %v", err)
	}
}
