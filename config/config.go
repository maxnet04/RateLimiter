package config

import (
	"fmt"
	"rate-limiter/application/controllers"
	"rate-limiter/application/middleware"
	"rate-limiter/application/repository"
	"rate-limiter/application/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type Configure struct {
	apiKey           string `mapstructure:"API_KEY"`
	rate_limit_ip    int    `mapstructure:"RATE_LIMIT_IP"`
	rate_limit_token int    `mapstructure:"RATE_LIMIT_TOKEN"`
	block_duration   int    `mapstructure:"BLOCK_DURATION"`
}

func Initialize() {

	// Configuração das variáveis de ambiente
	cfg, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	cfg.apiKey = viper.GetString("API_KEY")

	cfg.rate_limit_ip = viper.GetInt("RATE_LIMIT_IP")
	cfg.rate_limit_token = viper.GetInt("RATE_LIMIT_TOKEN")
	cfg.block_duration = viper.GetInt("BLOCK_DURATION_SECONDS")

	redisRepository := repository.NewRedisRepository()
	limiterUseCase := usecases.NewLimiterUseCase(redisRepository)

	// Configuração do middleware
	rateLimiterConfig := middleware.RateLimiterConfig{
		Token:          cfg.apiKey,           // Token registrado no config
		RequestsToken:  cfg.rate_limit_token, // Limite de requisições por token
		RequestsIP:     cfg.rate_limit_ip,    // Limite de requisições por IP
		BlockDuration:  cfg.block_duration,   // Duração do bloqueio
		LimiterUseCase: limiterUseCase,
	}

	// Aplica o middleware
	app.Use(middleware.RateLimiterMiddleware(rateLimiterConfig))

	setRoutes(app)

	// Inicializa o servidor
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}

func LoadConfig(path string) (*Configure, error) {
	var cfg *Configure
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.env")
	viper.AutomaticEnv()

	fmt.Println("Loading config from path:", path)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}

func setRoutes(app *fiber.App) {

	// Controllers
	rateLimiterController := controllers.NewRateLimiterController()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	app.Get("/", rateLimiterController.GetController)
}
