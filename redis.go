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
	if err != nil {
		return nil, err
	}
	awsAuth, err := auth.NewAWSAuth(
		auth.WithRole(vRole),
		auth.WithRegion(vRegion),
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
	secret, err := client.KVv2("database/creds/").Get(ctx, rRole)
	if err != nil {
		return nil, err
	}
	for _, value := range secret.Data {
		fmt.Printf("%+v\n", value)
	}
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
