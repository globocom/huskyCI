// Copyright 2018 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"
	"time"

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

const logActionConnect = "Connect"
const logActionReconnect = "autoReconnect"
const logInfoMongo = "DB"

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
func Connect(address, dbName, username, password string, poolLimit, port int, timeout time.Duration) error {

	log.Info(logActionConnect, logInfoMongo, 21)
	dbAddress := fmt.Sprintf("%s:%d", address, port)
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{dbAddress},
		Timeout:   timeout,
		FailFast:  true,
		Database:  dbName,
		Username:  username,
		Password:  password,
		PoolLimit: poolLimit,
	}
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Error(logActionConnect, logInfoMongo, 2001, err)
		return err
	}
	session.SetSafe(&mgo.Safe{WMode: "majority"})

	if err := session.Ping(); err != nil {
		log.Error(logActionConnect, logInfoMongo, 2002, err)
		return err
	}

	Conn = &DB{Session: session}
	go autoReconnect()

	return nil
}

// autoReconnect checks mongo's connection each second and, if an error is found, reconect to it.
func autoReconnect() {
	log.Info(logActionReconnect, logInfoMongo, 22)
	var err error
	for {
		err = Conn.Session.Ping()
		if err != nil {
			log.Error(logActionReconnect, logInfoMongo, 2003, err)
			Conn.Session.Refresh()
			err = Conn.Session.Ping()
			if err == nil {
				log.Info(logActionReconnect, logInfoMongo, 23)
			} else {
				log.Error(logActionReconnect, logInfoMongo, 2004, err)
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
func (db *DB) Aggregation(aggregation []bson.M, collection string) (interface{}, error) {
	session := db.Session.Clone()
	defer session.Close()
	c := session.DB("").C(collection)

	pipe := c.Pipe(aggregation)
	resp := []bson.M{}
	iter := pipe.Iter()
	err := iter.All(&resp)

	return resp, err
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
