package cfg

const (
	SYNCH_DIR   = "./config/synch"
	WATCHER_DIR = "./config/watch"
)

type Config interface {
	Validate()
}
