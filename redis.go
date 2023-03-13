package wikilinkui

import (
	"context"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
)

type RedisHandler struct {
	RedisAddress string
	RedisRole    string
	VaultAddress string
	VaultRegion  string
	VaultRole    string
	VaultClient  *vault.Client
	RedisTTL     time.Time
}

func MakeRedisHandler(rAddr, rRole, vAddr, vRegion, vRole string) (*RedisHandler, error) {
	var ctx = context.Background()
	client, err := vault.NewClient(vault.DefaultConfig())
	client.SetAddress("http://" + vAddr)
	if err != nil {
		return nil, err
	}
	awsAuth, err := auth.NewAWSAuth(
		auth.WithRole(vRole),
	)
	if err != nil {
		return nil, err
	}
	authInfo, err := client.Auth().Login(ctx, awsAuth)
	if err != nil {
		return nil, err
	}
	if authInfo == nil {
		return nil, fmt.Errorf("auth empty")
	}
	secret, err := client.Logical().Read("database/creds/" + vRole)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v", secret)
	return &RedisHandler{
		RedisAddress: rAddr,
		RedisRole:    rRole,
		VaultAddress: vAddr,
		VaultRegion:  vRegion,
		VaultRole:    vRole,
		VaultClient:  client,
	}, nil
}

func (r *RedisHandler) GetValue(key string) (string, bool) {
	return "", false
}

func (r *RedisHandler) PutValue(key, value string) bool {
	return false
}
