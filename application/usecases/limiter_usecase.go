package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rate-limiter/application/repository"
	"strconv"
	"time"
)

type LimiterUseCaseInterface interface {
	ValidRateLimiter(parameter string, limit int, blockDuration int) error
	RemoveBlock(parameter string)
}

type limiterUseCase struct {
	redisRepository repository.RedisRepositoryInterface
}

func NewLimiterUseCase(repository repository.RedisRepositoryInterface) LimiterUseCaseInterface {
	return &limiterUseCase{
		redisRepository: repository,
	}
}

const (
	errorMessage string = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func (l *limiterUseCase) ValidRateLimiter(parameter string, limit int, blockDuration int) error {
	repository := l.redisRepository

	log.Printf("Parameter: %s", parameter)

	// Verifica se a chave de bloqueio existe
	blockKey := fmt.Sprintf("%s:block", parameter)
	blocked, _ := repository.Exists(context.Background(), blockKey)
	if blocked {
		return errors.New(errorMessage)
	}

	resp, _ := repository.Get(context.Background(), parameter)
	quantidade, _ := strconv.Atoi(resp)
	log.Printf("Quantidade atual: %d", quantidade)

	if quantidade >= limit {

		log.Printf("Excedeu limite:  %s , horario: %s , duration: %d", blockKey, time.Now().In(time.Local).Format("2006-01-02 15:04:05"), blockDuration)

		repository.Set(context.Background(), blockKey, true, time.Duration(time.Second*time.Duration(blockDuration)))

		return errors.New(errorMessage)
	}

	err := repository.Set(context.Background(), parameter, quantidade+1, time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (l *limiterUseCase) RemoveBlock(parameter string) {

	log.Printf("Remocao Bloqueio chave:  %s , horario: %s", parameter, time.Now().In(time.Local).Format("2006-01-02 15:04:05"))

	l.redisRepository.Delete(context.Background(), fmt.Sprintf("%s:block", parameter))
}
