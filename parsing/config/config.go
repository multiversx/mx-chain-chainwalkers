package config

type ElasticSearchConfig struct {
	URL      string
	Username string
	Password string
}

type Config struct {
	ElasticSearchConnector ElasticSearchConfig
}
