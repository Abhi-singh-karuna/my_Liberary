package cachehandler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Abhi-singh-karuna/my_Liberary/baselogger"
	"github.com/Abhi-singh-karuna/my_Liberary/logger"

	"github.com/redis/go-redis/v9"
)

type cacheHandler struct {
	log    logger.Logger
	client *redis.Client
}

type CacheHandler interface {
	Set(string, string, int) CacheResult
	Get(string) CacheResult
	Delete(string) CacheResult
}

type CacheResult interface {
	SetVal(val string)
	Val() string
	Result() (string, error)
	String() string
}

type Cache struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewCacheHandler(cfg Redis, log *baselogger.BaseLogger) CacheHandler {
	log.Infof("Host :-  %v   -- port  %v  ---  Pass %v  --- DB  %v", cfg.GetHost(), cfg.GetPort(), cfg.GetPassword(), cfg.GetDatabase())
	log.Info("CacheHandler created variables from Config")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.GetHost(), cfg.GetPort()),
		Password: cfg.GetPassword(),
		DB:       cfg.GetDatabase(),
	})
	return &cacheHandler{
		log:    log,
		client: redisClient,
	}
}

func (c *cacheHandler) Set(key, value string, ttl int) CacheResult {
	return c.client.Set(context.Background(), key, value, time.Hour*time.Duration(ttl))
}

func (c *cacheHandler) Get(key string) CacheResult {
	return c.client.Get(context.Background(), key)
}

func (c *cacheHandler) Delete(key string) CacheResult {
	return &DelCacheResult{cmd: c.client.Del(context.Background(), key)}
}

// DelCacheResult is a wrapper around *redis.IntCmd to implement CacheResult
type DelCacheResult struct {
	cmd *redis.IntCmd
}

func (r *DelCacheResult) SetVal(val string) {
	// This method can be implemented if needed, but for delete, it's usually not required.
}

func (r *DelCacheResult) Val() string {
	// Convert the int64 result from Del command to string
	return strconv.FormatInt(r.cmd.Val(), 10)
}

func (r *DelCacheResult) Result() (string, error) {
	// Convert the int64 result from Del command to string and return it along with the error
	val, err := r.cmd.Result()
	return strconv.FormatInt(val, 10), err
}

func (r *DelCacheResult) String() string {
	// Provide a string representation of the result
	return r.Val()
}
