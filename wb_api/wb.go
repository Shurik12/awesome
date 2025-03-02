package wildberies

import (
	"github.com/lanzay/wildberries"
)

func Connect() error {
	// key initialized from config
	var key string
	user := wildberries.New(key)
	user.Incomes()
	return nil
}
