package wikilinkui

import (
	"context"
	"fmt"
	"time"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/aws"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisHandler struct {
	RedisAddress string
	VaultAddress string
	VaultRole    string
	VaultClient  *vault.Client
	RedisClient  *redis.Client
	RedisTTL     time.Time
	Logger       zerolog.Logger
}

func MakeRedisHandler(rAddr, vAddr, vRole string, logger zerolog.Logger) (*RedisHandler, error) {
	client, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		return nil, err
	}
	client.SetAddress("http://" + vAddr)
	return &RedisHandler{
		RedisAddress: rAddr,
		VaultAddress: vAddr,
		VaultRole:    vRole,
		VaultClient:  client,
		RedisTTL:     time.Now(),
		Logger:       logger,
	}, nil
}

func (r *RedisHandler) GetValue(key string) (string, error) {
	if time.Now().After(r.RedisTTL) {
		// Need new redis creds
		if err := r.getRedisAuth(); err != nil {
			r.Logger.Debug().Msgf("error getting new creds %v", err)
			return "", err
		}
		r.Logger.Debug().Msg("got new creds")
	}
	res, err := r.RedisClient.Get(context.TODO(), key).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RedisHandler) PutValue(key, value string) error {
	if time.Now().After(r.RedisTTL) {
		// Need new redis creds
		if err := r.getRedisAuth(); err != nil {
			r.Logger.Debug().Msgf("error getting new creds %v", err)
			return err
		}
		r.Logger.Debug().Msg("got new creds")
	}
	if err := r.RedisClient.Set(context.TODO(), key, value, 0).Err(); err != nil {
		return err
	}
	return nil
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
