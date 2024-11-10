package vault

import (
	"context"
	"errors"
)

func (v *Vault) GetSecret(ctx context.Context, path, variableName string) (string, error) {
	response, err := v.Client.Read(ctx, path)
	if err != nil {
		return "", err
	}

	data, ok := response.Data["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("cannot assert type from vault response")
	}

	return data[variableName].(string), nil
}
