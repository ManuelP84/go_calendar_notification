package mongo

type MongoDbSettings struct {
	Port string `envconfig:"MONGO_PORT" required:"true"`
	Host string `envconfig:"MONGO_HOST" required:"true"`
}
