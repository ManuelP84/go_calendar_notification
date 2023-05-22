package rabbit

type RabbitSettings struct {
	User     string `envconfig:"RABBIT_USER" required:"true"`
	Password string `envconfig:"RABBIT_PASSWORD" required:"true"`
	Host     string `envconfig:"RABBIT_HOST" required:"true"`
	Port     string `envconfig:"RABBIT_PORT" required:"true"`
}
