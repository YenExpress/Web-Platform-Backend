package ratelimiter

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RateLimiter struct {
	maxWrongAttemptsByIPperDay      int
	maxConsecutiveFailsByEmailAndIP int
	client                          *redis.Client
	SlowBruteKeyPrefix              string
	ConsecutiveFailsKeyPrefix       string
}

func (limiter *RateLimiter) NoteFailedByIP(ip_addr string) {
	key := fmt.Sprintf("%v_%v", ip_addr, limiter.SlowBruteKeyPrefix)
	var limit LoginLimit
	err := limit.unmarshalBinary([]byte(limiter.client.Get(key).Val()))
	if err != nil {
		limit := &LoginLimit{FailedAttempts: 1, BlockDuration: time.Now().Add(time.Hour * 24)}
		jsonified, _ := limit.marshalBinary()
		limiter.client.Set(key, jsonified, time.Hour*24).Val()
	} else {
		updated_limit := &LoginLimit{FailedAttempts: limit.FailedAttempts + 1,
			BlockDuration: limit.BlockDuration}
		jsonified, _ := updated_limit.marshalBinary()
		limiter.client.Set(key, jsonified, limit.BlockDuration.Sub(time.Now())).Val()
	}
}

func (limiter *RateLimiter) MaxOutFailureByIPForToday(ip_addr string) bool {
	key := fmt.Sprintf("%v_%v", ip_addr, limiter.SlowBruteKeyPrefix)
	var limit LoginLimit
	err := limit.unmarshalBinary([]byte(limiter.client.Get(key).Val()))
	if err != nil {
		return false
	} else {
		if limit.FailedAttempts >= limiter.maxWrongAttemptsByIPperDay && time.Now().Before(limit.BlockDuration) {
			return true
		} else if limit.FailedAttempts >= limiter.maxWrongAttemptsByIPperDay && time.Now().After(limit.BlockDuration) {
			limiter.client.Del(key)
		}
		return false
	}
}

func (limiter *RateLimiter) NoteFailedByEmailAndIP(email, ip_addr string) {
	key := fmt.Sprintf("%v_%v_%v", email, ip_addr, limiter.ConsecutiveFailsKeyPrefix)
	var limit LoginLimit
	err := limit.unmarshalBinary([]byte(limiter.client.Get(key).Val()))
	if err != nil {
		limit := &LoginLimit{FailedAttempts: 1, BlockDuration: time.Now().Add(time.Hour * 1)}
		jsonified, _ := limit.marshalBinary()
		limiter.client.Set(key, jsonified, time.Hour*1)
	} else {
		updated_limit := &LoginLimit{FailedAttempts: limit.FailedAttempts + 1,
			BlockDuration: limit.BlockDuration}
		jsonified, _ := updated_limit.marshalBinary()
		limiter.client.Set(key, jsonified, limit.BlockDuration.Sub(time.Now()))
	}
}

func (limiter *RateLimiter) MaxOutFailureByEmailandIPForTheHour(email, ip_addr string) bool {
	key := fmt.Sprintf("%v_%v_%v", email, ip_addr, limiter.ConsecutiveFailsKeyPrefix)
	var limit LoginLimit
	err := limit.unmarshalBinary([]byte(limiter.client.Get(key).Val()))
	if err != nil {
		return false
	} else {
		if limit.FailedAttempts >= limiter.maxConsecutiveFailsByEmailAndIP && time.Now().Before(limit.BlockDuration) {
			return true
		} else if limit.FailedAttempts >= limiter.maxConsecutiveFailsByEmailAndIP && time.Now().After(limit.BlockDuration) {
			limiter.client.Del(key)
		}
		return false
	}
}

func (limiter *RateLimiter) NoteFailure(email, ip_addr string) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		limiter.NoteFailedByEmailAndIP(email, ip_addr)
		wg.Done()
	}()
	go func() {
		limiter.NoteFailedByIP(ip_addr)
		wg.Done()
	}()
	wg.Wait()
}

func (limiter *RateLimiter) MaxOutFailure(email, ip_addr string) bool {
	var failedbyIP, failedbyEmailandIP bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func(result *bool) {
		*result = limiter.MaxOutFailureByEmailandIPForTheHour(email, ip_addr)
		wg.Done()
	}(&failedbyEmailandIP)
	go func(result *bool) {
		*result = limiter.MaxOutFailureByIPForToday(ip_addr)
		wg.Done()
	}(&failedbyIP)
	wg.Wait()
	if failedbyEmailandIP || failedbyIP {
		return true
	}
	return false
}
