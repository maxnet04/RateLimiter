package middleware

import (
	"log"
	"rate-limiter/application/usecases"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RateLimiterConfig Configuração personalizada para o rate limiter
type RateLimiterConfig struct {
	Token          string
	RequestsToken  int
	RequestsIP     int
	BlockDuration  int
	LimiterUseCase usecases.LimiterUseCaseInterface
}

const (
	errorMessage string = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

func RateLimiterMiddleware(config RateLimiterConfig) fiber.Handler {
	start := time.Now()

	return func(c *fiber.Ctx) error {

		// Identifica o IP do cliente
		clientIP := c.IP()
		parameter := clientIP
		limit := config.RequestsIP

		headers := c.GetReqHeaders()

		if headers["Api_key"] != nil {
			token := headers["Api_key"]
			if !strings.EqualFold(token[0], "") && strings.EqualFold(token[0], config.Token) {
				log.Printf("Token informado: %s", token)
				parameter = token[0]
				limit = config.RequestsToken
				log.Printf("Parameter: %s | Limit: %d", parameter, limit)
			}
		}

		err := config.LimiterUseCase.ValidRateLimiter(parameter, limit, config.BlockDuration)
		if err != nil && err.Error() == errorMessage {

			// Reseta o contador de requisições
			go func(ip string) {
				time.Sleep(time.Duration(time.Second * time.Duration(config.BlockDuration)))
				config.LimiterUseCase.RemoveBlock(ip)
			}(clientIP)

			if headers["Api_key"] != nil && strings.EqualFold(headers["Api_key"][0], config.Token) {
				log.Printf("Token %s excedeu o limite de %d requisições.", parameter, limit)
			} else {
				log.Printf("IP %s excedeu o limite de %d requisições.", parameter, limit)
			}

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": errorMessage,
			})
		}

		// Registra os detalhes da requisição
		log.Printf("[%s] %s | Status: %d | Request Time: %s",
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			time.Since(start),
		)

		// Continua o fluxo padrão da chamada (controller)
		return c.Next()
	}
}
