// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/globocom/huskyCI/api/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var collectionTest = "integration_test_collection"

type mongoTestObj struct {
	Attr1 string `bson:"attribute_1"`
	Attr2 int    `bson:"attribute_2"`
}

func TestMongoIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	log.InitLog(true, "", "", "log_test", "log_test")
	RegisterFailHandler(Fail)
	RunSpecs(t, "MongoDB Suite")
}

var _ = BeforeSuite(func() {
	mongoAddress := os.Getenv("HUSKYCI_DATABASE_DB_ADDR")
	dbName := os.Getenv("HUSKYCI_DATABASE_DB_NAME")
	username := os.Getenv("HUSKYCI_DATABASE_DB_USERNAME")
	password := os.Getenv("HUSKYCI_DATABASE_DB_PASSWORD")
	dbPort, err := strconv.Atoi(os.Getenv("HUSKYCI_DATABASE_DB_PORT"))
	Expect(err).To(BeNil())

	connectionPool := 10
	connectionTimeout := time.Duration(1 * time.Second)

	errConnect := Connect(mongoAddress, dbName, username, password, connectionPool, dbPort, connectionTimeout)
	Expect(errConnect).To(BeNil())
	Expect(Conn).To(Not(BeNil()))
})

var _ = AfterSuite(func() {
	colletction := Conn.Session.DB("").C(collectionTest)
	err := colletction.DropCollection()
	Expect(err).To(BeNil())
})

var _ = Describe("Connect", func() {
	Context("When connect to MongoDB with valid parameters", func() {
		It("Should return no error when send ping", func() {
			errPing := Conn.Session.Ping()
			Expect(errPing).To(BeNil())
		})
	})
})

var _ = Describe("Insert", func() {
	Context("When insert a object", func() {
		It("Should insert with success", func() {
			newObj := mongoTestObj{
				Attr1: "insert-1",
				Attr2: 11,
			}

			errInsert := Conn.Insert(newObj, collectionTest)
			Expect(errInsert).To(BeNil())

			result := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "insert-1"}, nil, collectionTest, &result)
			Expect(errGet).To(BeNil())

			expectedResult := []mongoTestObj{
				{
					Attr1: "insert-1",
					Attr2: 11,
				},
			}

			Expect(result).To(Equal(expectedResult))
		})
	})
})

var _ = Describe("Update", func() {
	Context("When update a object", func() {
		It("Should update with success", func() {
			err := Conn.Insert(mongoTestObj{Attr1: "update-1", Attr2: 11}, collectionTest)
			Expect(err).To(BeNil())

			updatedQuery := bson.M{
				"$set": mongoTestObj{
					Attr1: "update-2",
					Attr2: 111,
				},
			}

			errupdate := Conn.Update(bson.M{"attribute_1": "update-1"}, updatedQuery, collectionTest)
			Expect(errupdate).To(BeNil())

			result := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "update-2"}, nil, collectionTest, &result)
			Expect(errGet).To(BeNil())

			expectedResult := []mongoTestObj{
				{
					Attr1: "update-2",
					Attr2: 111,
				},
			}

			Expect(result).To(Equal(expectedResult))
		})
	})
})

var _ = Describe("UpdateAll", func() {
	Context("When update objects", func() {
		It("Should update all with success", func() {
			errInsert1 := Conn.Insert(mongoTestObj{Attr1: "update-all-1", Attr2: 11}, collectionTest)
			Expect(errInsert1).To(BeNil())

			errInsert2 := Conn.Insert(mongoTestObj{Attr1: "update-all-1", Attr2: 22}, collectionTest)
			Expect(errInsert2).To(BeNil())

			updatedQuery := bson.M{
				"$set": mongoTestObj{
					Attr1: "update-all-2",
					Attr2: 33,
				},
			}

			errupdate := Conn.UpdateAll(bson.M{"attribute_1": "update-all-1"}, updatedQuery, collectionTest)
			Expect(errupdate).To(BeNil())

			result := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "update-all-2"}, nil, collectionTest, &result)
			Expect(errGet).To(BeNil())

			expectedResult := []mongoTestObj{
				{
					Attr1: "update-all-2",
					Attr2: 33,
				},
				{
					Attr1: "update-all-2",
					Attr2: 33,
				},
			}

			Expect(result).To(Equal(expectedResult))
		})
	})
})

var _ = Describe("FindAndModify", func() {
	Context("When find and modify a object", func() {
		It("Should modify with success", func() {
			objToInsert := mongoTestObj{Attr1: "find-and-modify-1", Attr2: 11}
			err := Conn.Insert(objToInsert, collectionTest)
			Expect(err).To(BeNil())

			updatedQuery := bson.M{
				"$set": mongoTestObj{
					Attr1: "find-and-modify-2",
					Attr2: 33,
				},
			}

			resultFindModify := mongoTestObj{}
			errFindModify := Conn.FindAndModify(bson.M{"attribute_1": "find-and-modify-1"}, updatedQuery, collectionTest, &resultFindModify)
			Expect(errFindModify).To(BeNil())

			resultGet := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "find-and-modify-2"}, nil, collectionTest, &resultGet)
			Expect(errGet).To(BeNil())
			Expect(resultFindModify).To(Equal(objToInsert))

			expectedResult := []mongoTestObj{
				{
					Attr1: "find-and-modify-2",
					Attr2: 33,
				},
			}

			Expect(resultGet).To(Equal(expectedResult))
		})
	})
})

var _ = Describe("Search", func() {
	Context("When search a object", func() {
		It("Should return with success", func() {
			objToInsert := mongoTestObj{Attr1: "search-1", Attr2: 11}
			err := Conn.Insert(objToInsert, collectionTest)
			Expect(err).To(BeNil())

			resultSearch := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "search-1"}, nil, collectionTest, &resultSearch)
			Expect(errGet).To(BeNil())

			expectedResult := []mongoTestObj{
				objToInsert,
			}

			Expect(resultSearch).To(Equal(expectedResult))
		})
	})
})

var _ = Describe("Aggregation", func() {
	Context("When run aggregation with count", func() {
		It("Should return count with success", func() {
			objToInsert1 := mongoTestObj{Attr1: "aggregation-test", Attr2: 11}
			err1 := Conn.Insert(objToInsert1, collectionTest)
			Expect(err1).To(BeNil())

			objToInsert2 := mongoTestObj{Attr1: "aggregation-test", Attr2: 22}
			err2 := Conn.Insert(objToInsert2, collectionTest)
			Expect(err2).To(BeNil())

			aggregationQuery := []bson.M{
				{
					"$match": bson.M{
						"attribute_1": bson.M{
							"$eq": "aggregation-test",
						},
					},
				},
				{
					"$match": bson.M{
						"attribute_2": bson.M{
							"$gte": 1,
						},
					},
				},
				{
					"$group": bson.M{
						"_id": nil,
						"count": bson.M{
							"$sum": 1,
						},
					},
				},
			}

			resultSearch, errGet := Conn.Aggregation(aggregationQuery, collectionTest)
			Expect(errGet).To(BeNil())
			expectedResult := []bson.M{
				{"_id": nil, "count": 2},
			}

			Expect(resultSearch).To(Equal(expectedResult))
		})
	})
})

var _ = Describe("SearchOne", func() {
	Context("When search a object", func() {
		It("Should return with success", func() {
			objToInsert := mongoTestObj{Attr1: "search-one-1", Attr2: 11}
			err := Conn.Insert(objToInsert, collectionTest)
			Expect(err).To(BeNil())

			resultSearch := mongoTestObj{}
			errGet := Conn.SearchOne(bson.M{"attribute_1": "search-one-1"}, nil, collectionTest, &resultSearch)
			Expect(errGet).To(BeNil())
			Expect(resultSearch).To(Equal(objToInsert))
		})
	})
})

var _ = Describe("Upsert", func() {
	Context("When upsert a object that doesn't exists", func() {
		It("Should insert with success", func() {
			obj := mongoTestObj{
				Attr1: "upsert-1",
				Attr2: 11,
			}

			query := bson.M{"attribute_1": "upsert-1"}
			_, err := Conn.Upsert(query, obj, collectionTest)
			Expect(err).To(BeNil())

			result := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "upsert-1"}, nil, collectionTest, &result)
			Expect(errGet).To(BeNil())

			expectedResult := []mongoTestObj{
				{
					Attr1: "upsert-1",
					Attr2: 11,
				},
			}

			Expect(result).To(Equal(expectedResult))
		})
	})
	Context("When upsert a object that already exists", func() {
		It("Should update with success", func() {
			errInsert := Conn.Insert(mongoTestObj{Attr1: "upsert-2", Attr2: 11}, collectionTest)
			Expect(errInsert).To(BeNil())

			obj := mongoTestObj{
				Attr1: "upsert-2",
				Attr2: 22,
			}

			query := bson.M{"attribute_1": "upsert-2"}
			_, err := Conn.Upsert(query, obj, collectionTest)
			Expect(err).To(BeNil())

			result := []mongoTestObj{}
			errGet := Conn.Search(bson.M{"attribute_1": "upsert-2"}, nil, collectionTest, &result)
			Expect(errGet).To(BeNil())

			expectedResult := []mongoTestObj{
				{
					Attr1: "upsert-2",
					Attr2: 22,
				},
			}

			Expect(result).To(Equal(expectedResult))
		})
	})
})
