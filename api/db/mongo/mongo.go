// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"time"

	config "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Conn is the MongoDB connection variable.
var Conn *DB

// Collections names used in MongoDB.
var (
	RepositoryCollection   = "repository"
	SecurityTestCollection = "securityTest"
	AnalysisCollection     = "analysis"
	UserCollection         = "user"
	AccessTokenCollection  = "accessToken"
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

	log.Info("Connect", "DB", 21)
	mongoConfig := config.APIConfiguration.MongoDBConfig
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
		log.Error("Connect", "DB", 2001, err)
		return err
	}
	session.SetSafe(&mgo.Safe{WMode: "majority"})

	if err := session.Ping(); err != nil {
		log.Error("Connect", "DB", 2002, err)
		return err
	}

	Conn = &DB{Session: session}
	go autoReconnect()

	return nil
}

// autoReconnect checks mongo's connection each second and, if an error is found, reconect to it.
func autoReconnect() {
	log.Info("autoReconnect", "DB", 22)
	var err error
	for {
		err = Conn.Session.Ping()
		if err != nil {
			log.Error("autoReconnect", "DB", 2003, err)
			Conn.Session.Refresh()
			err = Conn.Session.Ping()
			if err == nil {
				log.Info("autoReconnect", "DB", 23)
			} else {
				log.Error("autoReconnect", "DB", 2004, err)
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

// Aggregation prepares a pipeline to aggregate.
func (db *DB) Aggregation(aggregation []bson.M, collection string, obj interface{}) error {
	session := db.Session.Clone()
	defer session.Close()
	c := session.DB("").C(collection)

	pipe := c.Pipe(aggregation)
	iter := pipe.Iter()
	err := iter.All(obj)

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
