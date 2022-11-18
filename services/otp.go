package services

import (
	"context"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Marcel-MD/rooms-go-api/rdb"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type IOtpService interface {
	Generate(email string) (string, error)
	Verify(email string, otp string) error
}

type OtpService struct {
	rdb    *redis.Client
	ctx    context.Context
	expiry time.Duration
}

var (
	otpOnce    sync.Once
	otpService IOtpService
)

func GetOtpService() IOtpService {
	otpOnce.Do(func() {
		log.Info().Msg("Initializing otp service")

		expiryStr := os.Getenv("OTP_EXPIRY")
		expiry, err := time.ParseDuration(expiryStr)
		if err != nil {
			expiry = 10 * time.Minute
		}

		rdb, ctx := rdb.GetRDB()
		otpService = &OtpService{
			rdb:    rdb,
			ctx:    ctx,
			expiry: expiry,
		}
	})
	return otpService
}

func (s *OtpService) Generate(email string) (string, error) {
	log.Debug().Msg("Generating otp")

	num := 100000 + rand.Intn(800000)
	otp := strconv.Itoa(num)

	err := s.rdb.Set(s.ctx, email, otp, s.expiry).Err()
	if err != nil {
		log.Err(err).Msg("Error setting otp in redis")
		return "", err
	}

	return otp, nil
}

func (s *OtpService) Verify(email string, otp string) error {
	log.Debug().Msg("Validating otp")

	otpFromRedis, err := s.rdb.Get(s.ctx, email).Result()
	if err != nil {
		log.Err(err).Msg("Error getting otp from redis")
		return err
	}

	if otpFromRedis != otp {
		return errors.New("otp is not valid")
	}

	return nil
}
