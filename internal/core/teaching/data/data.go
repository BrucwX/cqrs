package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"
)

// ProviderSet is teaching data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewTeacherRepo,
	NewStudentRepo,
	NewCourseRepo,
)

// Data holds the long-lived storage clients for teaching context.
type Data struct {
	// TODO: Add database client (e.g., *ent.Client, *gorm.DB, *sql.DB)
}

// NewData opens the database client and returns it with a cleanup function.
func NewData(c *conf.Data) (*Data, func(), error) {
	// TODO: Initialize database connection based on config
	// Example with ent:
	// dc := c.GetDatabase()
	// db, err := ent.Open(dc.GetDriver(), dc.GetSource())
	// if err != nil {
	//     return nil, nil, err
	// }

	cleanup := func() {
		log.Info("closing teaching data resources")
		// TODO: Close database connection
	}
	return &Data{}, cleanup, nil
}
