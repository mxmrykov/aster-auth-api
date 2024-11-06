package vault

import "github.com/mxmrykov/asterix-auth/internal/config"

type IVault interface {
	GetSecret(path, variableName string) (string, error)
}

type Vault struct {
}

func NewVault(cfg *config.Vault) (IVault, error) {
	return &Vault{}, nil
}
