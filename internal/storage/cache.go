package storage

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/patrickmn/go-cache"
	"time"
)

type BasicCacheModule struct {
	ttl   time.Duration
	cache *cache.Cache
}

type Config struct {
	TTL     time.Duration `env:"STORAGE_TTL" envDefault:"0s"`
	CleanUP time.Duration `env:"STORAGE_CLEANUP_INTERVAL" envDefault:"7d"`
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.TTL),
		validation.Field(&c.CleanUP),
	)
}

func New(cfg *Config) *BasicCacheModule {
	cfg.TTL = cache.NoExpiration

	return &BasicCacheModule{
		ttl: cfg.TTL,
		cache: cache.New(
			cfg.TTL,
			cfg.CleanUP,
		),
	}
}

type Key string

func (s Key) String() string { return string(s) }

type DecryptInfo struct {
	Key []int
}

func (sc *BasicCacheModule) Set(ctx context.Context, k Key, v *DecryptInfo) {
	sc.cache.Set(k.String(), v, sc.ttl)
}

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidType = errors.New("object type is not DecryptInfo")
)

func (sc *BasicCacheModule) Get(ctx context.Context, k Key) (*DecryptInfo, error) {
	val, ok := sc.cache.Get(k.String())
	if !ok {
		return nil, ErrNotFound
	}

	info, ok := val.(*DecryptInfo)
	if !ok {
		return nil, ErrInvalidType
	}

	return info, nil
}
