package ratelimiter

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type FailRateLimiter struct {
	MaxWrongAttemptsByIPperDay int
	MaxConsecutiveFails        int
	Client                     *redis.Client
	TotalFailsKeyPrefix        string
	ConsecutiveFailsKeyPrefix  string
}

func (limiter *FailRateLimiter) NoteFailureForToday(ip_addr string) {
	key := fmt.Sprintf("%v_%v", ip_addr, limiter.TotalFailsKeyPrefix)
	var limit failure
	err := limit.unmarshalBinary([]byte(limiter.Client.Get(key).Val()))
	if err != nil {
		// fmt.Println(err)
		limit := &failure{Attempts: 1, BlockTill: time.Now().Add(time.Hour * 24)}
		jsonified, _ := limit.marshalBinary()
		limiter.Client.Set(key, jsonified, time.Hour*24).Val()
	} else {
		updated_limit := &failure{Attempts: limit.Attempts + 1,
			BlockTill: limit.BlockTill}
		jsonified, _ := updated_limit.marshalBinary()
		limiter.Client.Set(key, jsonified, limit.BlockTill.Sub(time.Now())).Val()
	}
}

func (limiter *FailRateLimiter) MaxOutFailureForToday(ip_addr string) bool {
	key := fmt.Sprintf("%v_%v", ip_addr, limiter.TotalFailsKeyPrefix)
	var limit failure
	err := limit.unmarshalBinary([]byte(limiter.Client.Get(key).Val()))
	if err != nil {
		return false
	} else {
		if limit.Attempts >= limiter.MaxWrongAttemptsByIPperDay && time.Now().Before(limit.BlockTill) {
			return true
		} else if limit.Attempts >= limiter.MaxWrongAttemptsByIPperDay && time.Now().After(limit.BlockTill) {
			limiter.Client.Del(key)
		}
		return false
	}
}

func (limiter *FailRateLimiter) NoteConsecutiveFailure(user_id, ip_addr string) {
	key := fmt.Sprintf("%v_%v_%v", user_id, ip_addr, limiter.ConsecutiveFailsKeyPrefix)
	var limit failure
	err := limit.unmarshalBinary([]byte(limiter.Client.Get(key).Val()))
	if err != nil {
		// fmt.Println(err)
		limit := &failure{Attempts: 1, BlockTill: time.Now().Add(time.Hour * 1)}
		jsonified, _ := limit.marshalBinary()
		limiter.Client.Set(key, jsonified, time.Hour*1)
	} else {
		updated_limit := &failure{Attempts: limit.Attempts + 1,
			BlockTill: limit.BlockTill}
		jsonified, _ := updated_limit.marshalBinary()
		limiter.Client.Set(key, jsonified, limit.BlockTill.Sub(time.Now()))
	}
}

func (limiter *FailRateLimiter) MaxOutConsecutiveFailure(user_id, ip_addr string) bool {
	key := fmt.Sprintf("%v_%v_%v", user_id, ip_addr, limiter.ConsecutiveFailsKeyPrefix)
	var limit failure
	err := limit.unmarshalBinary([]byte(limiter.Client.Get(key).Val()))
	if err != nil {
		return false
	} else {
		if limit.Attempts >= limiter.MaxConsecutiveFails && time.Now().Before(limit.BlockTill) {
			return true
		} else if limit.Attempts >= limiter.MaxConsecutiveFails && time.Now().After(limit.BlockTill) {
			limiter.Client.Del(key)
		}
		return false
	}
}

func (limiter *FailRateLimiter) NoteFailure(user_id, ip_addr string) {
	if err := limiter.Client.Ping().Err(); err != nil {
		log.Println(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		limiter.NoteConsecutiveFailure(user_id, ip_addr)
		wg.Done()
	}()
	go func() {
		limiter.NoteFailureForToday(ip_addr)
		wg.Done()
	}()
	wg.Wait()
}

func (limiter *FailRateLimiter) MaxOutFailure(user_id, ip_addr string) bool {
	if err := limiter.Client.Ping().Err(); err != nil {
		log.Println(err)
	}
	var failedbyIP, failedbyUserIDandIP bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func(result *bool) {
		*result = limiter.MaxOutConsecutiveFailure(user_id, ip_addr)
		wg.Done()
	}(&failedbyUserIDandIP)
	go func(result *bool) {
		*result = limiter.MaxOutFailureForToday(ip_addr)
		wg.Done()
	}(&failedbyIP)
	wg.Wait()
	if failedbyUserIDandIP || failedbyIP {
		return true
	}
	return false
}
