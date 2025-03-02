package wildberies

import (
	"github.com/lanzay/wildberries"
)

func connect() error {
	// key initialized from config
	var key string
	user := wildberries.New(key)
	user.Incomes()
	return nil
}
