package wikilinkui

import (
	"context"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	RedisAddress string
	RedisRole    string
	VaultAddress string
	VaultRegion  string
	VaultRole    string
	VaultClient  *vault.Client
	RedisClient  *redis.Client
	RedisTTL     time.Time
}

func MakeRedisHandler(rAddr, rRole, vAddr, vRegion, vRole string) (*RedisHandler, error) {
	client, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return nil, err
	}
	client.SetAddress("http://" + vAddr)
	return &RedisHandler{
		RedisAddress: rAddr,
		RedisRole:    rRole,
		VaultAddress: vAddr,
		VaultRegion:  vRegion,
		VaultRole:    vRole,
		VaultClient:  client,
		RedisTTL:     time.Now(),
	}, nil
}

func (r *RedisHandler) GetValue(key string) (string, error) {
	if time.Now().After(r.RedisTTL) {
		// Need new redis creds
		if err := r.getRedisAuth(); err != nil {
			return "", err
		}
	}
	return r.RedisClient.Get(context.TODO(), key).Result()
}

func (r *RedisHandler) PutValue(key, value string) error {
	if time.Now().After(r.RedisTTL) {
		// Need new redis creds
		if err := r.getRedisAuth(); err != nil {
			return err
		}
	}
	return r.RedisClient.Set(context.TODO(), key, value, 0).Err()
}

// getRedisAuth reconnects to the redis server
func (r *RedisHandler) getRedisAuth() error {
	var ctx = context.Background()
	awsAuth, err := auth.NewAWSAuth(auth.WithRole(r.VaultRole))
	if err != nil {
		return err
	}
	authInfo, err := r.VaultClient.Auth().Login(ctx, awsAuth)
	if err != nil {
		return err
	}
	if authInfo == nil {
		return fmt.Errorf("vault auth failed")
	}
	secret, err := r.VaultClient.Logical().Read("database/creds/" + r.VaultRole)
	if err != nil {
		return err
	}
	username, ok := secret.Data["username"].(string)
	if !ok {
		return fmt.Errorf("username not in vault auth")
	}
	password, ok := secret.Data["password"].(string)
	if !ok {
		return fmt.Errorf("password not in vault auth")
	}
	r.RedisClient = redis.NewClient(&redis.Options{
		Addr:     r.RedisAddress,
		Username: username,
		Password: password,
		DB:       0,
	})
	r.RedisTTL = time.Now().Add(time.Second * time.Duration(secret.LeaseDuration))
	return nil
}
