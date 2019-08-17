// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

// CreateNewRepository inserts the given repository into RepositoryCollection.
func CreateNewRepository(c echo.Context) error {
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		log.Warning("CreateNewRepository", "ANALYSIS", 101)
		reply := map[string]interface{}{"success": false, "error": "invalid repository JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	_, err = db.FindOneDBRepository(repositoryQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewRepository", "ANALYSIS", 110, repository.URL)
			reply := map[string]interface{}{"success": false, "error": "this repository is already registered"}
			return c.JSON(http.StatusConflict, reply)
		}
		log.Error("CreateNewRepository", "ANALYSIS", 1013, err)
	}

	err = db.InsertDBRepository(repository)
	if err != nil {
		log.Error("CreateNewRepository", "ANALYSIS", 2015, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	log.Info("CreateNewRepository", "ANALYSIS", 17, repository.URL)
	reply := map[string]interface{}{"success": true, "error": ""}
	return c.JSON(http.StatusCreated, reply)
}
