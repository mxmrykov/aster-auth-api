package vault

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault-client-go"
	"github.com/mxmrykov/asterix-auth/internal/config"
)

type IVault interface {
	GetSecret(ctx context.Context, path, variableName string) (string, error)
}

type Vault struct {
	Client *vault.Client
	Token  string
}

func NewVault(cfg *config.Vault) (IVault, error) {
	tls := vault.TLSConfiguration{}
	tls.ServerCertificate.FromFile = "./pkg/clients/vault/ca_cert.pem"

	client, err := vault.New(
		vault.WithAddress(
			fmt.Sprintf(
				"https://%s:%d",
				cfg.Host,
				cfg.Port,
			),
		),
		vault.WithRequestTimeout(cfg.ClientTimeout),
		vault.WithTLS(tls),
	)

	if err != nil {
		return nil, err
	}

	vlt := &Vault{
		Client: client,
		Token:  cfg.AuthToken,
	}

	if err := client.SetToken(vlt.Token); err != nil {
		return nil, err
	}

	return vlt, nil
}
