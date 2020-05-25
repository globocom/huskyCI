package database

// DBSession is the interface that holds the database session
type DBSession interface {
	Close() error
}
