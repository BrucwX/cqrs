package memory

import (
	"cqrs/internal/conf"

	"github.com/go-kratos/kratos/v3/log"
)

// Data holds the long-lived storage clients for commerce context.
type Data struct {
	// TODO: Add database client (e.g., *ent.Client, *gorm.DB, *sql.DB)
}

// NewData opens the database client and returns it with a cleanup function.
func NewData(c *conf.Data) (*Data, func(), error) {
	// TODO: Initialize database connection based on config

	cleanup := func() {
		log.Info("closing commerce data resources")
		// TODO: Close database connection
	}
	return &Data{}, cleanup, nil
}
