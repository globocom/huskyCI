package db

import (
	"fmt"
	"time"

	"github.com/globocom/glbgelf"
	config "github.com/globocom/huskyci/context"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Conn *DB

// Collections names used in MongoDB.
var (
	RepositoryCollection   = "repository"
	SecurityTestCollection = "securityTest"
	AnalysisCollection     = "analysis"
)

// DB is the struct that represents mongo session.
type DB struct {
	Session *mgo.Session
}

// mongoConfig is the struct that represents mongo configuration.
type mongoConfig struct {
	Address      string
	DatabaseName string
	UserName     string
	Password     string
}

// Database is the interface's database.
type Database interface {
	Insert(obj interface{}, collection string) error
	Search(query bson.M, selectors []string, collection string, obj interface{}) error
	Update(query bson.M, updateQuery interface{}, collection string) error
	UpdateAll(query, updateQuery bson.M, collection string) error
	Upsert(query bson.M, obj interface{}, collection string) (*mgo.ChangeInfo, error)
	SearchOne(query bson.M, selectors []string, collection string, obj interface{}) error
}

// Connect connects to mongo and returns the session.
func Connect() error {
	glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "Connect",
		"info":   "DB"}, "INFO", "Connecting to mongodb")
	mongoConfig := config.ApiConfig.MongoDBConfig
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{mongoConfig.Address},
		Timeout:   mongoConfig.Timeout,
		FailFast:  true,
		Database:  mongoConfig.DatabaseName,
		Username:  mongoConfig.Username,
		Password:  mongoConfig.Password,
		PoolLimit: mongoConfig.PoolLimit,
	}
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "Connect",
			"info":   "DB"}, "ERROR", "Error connecting to Mongo:", err); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		return err
	}
	session.SetSafe(&mgo.Safe{WMode: "majority"})

	if err := session.Ping(); err != nil {
		if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
			"action": "Connect",
			"info":   "DB"}, "ERROR", "Error pinging Mongo after connection:", err); errLog != nil {
			fmt.Println("glbgelf error: ", errLog)
		}
		return err
	}

	Conn = &DB{Session: session}
	go autoReconnect()

	return nil
}

// autoReconnect checks mongo's connection each second and, if an error is found, reconect to it.
func autoReconnect() {
	glbgelf.Logger.SendLog(map[string]interface{}{
		"action": "autoReconnect",
		"info":   "DB"}, "INFO", "Initializing mongodb auto reconnect")
	var err error
	for {
		err = Conn.Session.Ping()
		if err != nil {
			if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
				"action": "autoReconnect",
				"info":   "DB"}, "ERROR", "Error pinging Mongo in autoReconnect:", err); errLog != nil {
				fmt.Println("glbgelf error: ", errLog)
			}
			Conn.Session.Refresh()
			err = Conn.Session.Ping()
			if err == nil {
				if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
					"action": "autoReconnect",
					"info":   "DB"}, "ERROR", "Reconnect to MongoDB successful."); errLog != nil {
					fmt.Println("glbgelf error: ", errLog)
				}
			} else {
				if errLog := glbgelf.Logger.SendLog(map[string]interface{}{
					"action": "autoReconnect",
					"info":   "DB"}, "ERROR", "Reconnect to MongoDB failed:", err); errLog != nil {
					fmt.Println("glbgelf error: ", errLog)
				}
			}
		}
		time.Sleep(time.Second * 1)
	}
}

// Insert inserts a new document.
func (db *DB) Insert(obj interface{}, collection string) error {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	return c.Insert(obj)
}

// Update updates a single document.
func (db *DB) Update(query, updateQuery interface{}, collection string) error {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	err := c.Update(query, updateQuery)
	return err
}

// UpdateAll updates all documents that match the query.
func (db *DB) UpdateAll(query, updateQuery interface{}, collection string) error {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	_, err := c.UpdateAll(query, updateQuery)
	return err
}

// Search searchs all documents that match the query. If selectors are present, the return will be only the chosen fields.
func (db *DB) Search(query bson.M, selectors []string, collection string, obj interface{}) error {
	session := db.Session.Clone()
	defer session.Close()
	c := session.DB("").C(collection)

	var err error
	if selectors != nil {
		selector := bson.M{}
		for _, v := range selectors {
			selector[v] = 1
		}
		err = c.Find(query).Select(selector).All(obj)
	} else {
		err = c.Find(query).All(obj)
	}
	if err == nil && obj == nil {
		err = mgo.ErrNotFound
	}
	return err
}

// SearchOne searchs for the first element that matchs with the given query.
func (db *DB) SearchOne(query bson.M, selectors []string, collection string, obj interface{}) error {
	session := db.Session.Clone()
	defer session.Close()
	c := session.DB("").C(collection)

	var err error
	if selectors != nil {
		selector := bson.M{}
		for _, v := range selectors {
			selector[v] = 1
		}
		err = c.Find(query).Select(selector).One(obj)
	} else {
		err = c.Find(query).One(obj)
	}
	if err == nil && obj == nil {
		err = mgo.ErrNotFound
	}
	return err
}

// Upsert inserts a document or update it if it already exists.
func (db *DB) Upsert(query bson.M, obj interface{}, collection string) (*mgo.ChangeInfo, error) {
	session := db.Session.Clone()
	c := session.DB("").C(collection)
	defer session.Close()
	return c.Upsert(query, obj)
}
