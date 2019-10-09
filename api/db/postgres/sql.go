package db

import ()


func (sqlConfig *SqlConfig) Connect() error {
	if err := sqlConfig.Postgres.ConfigureDB(); err != nil {
		return err
	}
	sqlConfig.Postgres.ConfigurePool()
	return nil
}

func (sqlConfig *SqlConfig) CloseDB() error {
	return sqlConfig.Postgres.CloseDB()
}


func (sqlConfig *SqlConfig) Search(query bson.M, selectors []string, collection string, obj interface{}) error {}