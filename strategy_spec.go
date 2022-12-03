package config_server

type WatchChannel <-chan string

type StrategySpec interface {
	Watch(key string) WatchChannel
	Stop(key string) error
	Get(key string) (string, error)
	Set(key string, val string) error
}
