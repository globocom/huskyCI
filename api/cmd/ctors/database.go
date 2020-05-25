package ctors

import (
	"github.com/globocom/huskyCI/api/database"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// NewDatabaseSession starts a new database session.
func NewDatabaseSession(lc fx.Lifecycle, settings *viper.Viper) (database.DBSession, error) {

	databaseType := settings.GetString("HUSKYCI_DATABASE_TYPE")

	if databaseType == "mongo" {
		mongoSession, err := database.NewMongoSession(lc, settings)
		if err != nil {
			return nil, err
		}
		return mongoSession, nil
	}

	return database.NewPostgresSession(lc, settings)
}
