package config_server

type StrategySpec interface {
	Watch(key string) <-chan string
	Stop(key string)
	Get(key string) (string, error)
	Set(key string, val string) error
}
