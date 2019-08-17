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

// CreateNewSecurityTest inserts the given securityTest into SecurityTestCollection.
func CreateNewSecurityTest(c echo.Context) error {
	securityTest := types.SecurityTest{}
	err := c.Bind(&securityTest)
	if err != nil {
		log.Warning("CreateNewSecurityTest", "ANALYSIS", 108)
		reply := map[string]interface{}{"success": false, "error": "invalid security test JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	securityTestQuery := map[string]interface{}{"name": securityTest.Name}
	_, err = db.FindOneDBSecurityTest(securityTestQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewSecurityTest", "ANALYSIS", 109, securityTest.Name)
			reply := map[string]interface{}{"success": false, "error": "this security test is already registered"}
			return c.JSON(http.StatusConflict, reply)
		}
		log.Error("CreateNewSecurityTest", "ANALYSIS", 1012, err)
	}

	err = db.InsertDBSecurityTest(securityTest)
	if err != nil {
		log.Error("CreateNewSecurityTest", "ANALYSIS", 2016, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}

	log.Info("CreateNewSecurityTest", "ANALYSIS", 18, securityTest.Name)
	reply := map[string]interface{}{"success": true, "error": ""}
	return c.JSON(http.StatusCreated, reply)
}
