package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

// CreateNewSecurityTest inserts the given securityTest into SecurityTestCollection.
func CreateNewSecurityTest(c echo.Context) error {
	securityTest := types.SecurityTest{}
	err := c.Bind(&securityTest)
	if err != nil {
		log.Warning("CreateNewSecurityTest", "ANALYSIS", 108)
		requestResponse := util.RequestResponse(false, "invalid security test JSON")
		return c.JSON(http.StatusBadRequest, requestResponse)
	}

	securityTestQuery := map[string]interface{}{"name": securityTest.Name}
	_, err = db.FindOneDBSecurityTest(securityTestQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewSecurityTest", "ANALYSIS", 109, securityTest.Name)
			requestResponse := util.RequestResponse(false, "this security test is already registered")
			return c.JSON(http.StatusConflict, requestResponse)
		}
		log.Error("CreateNewSecurityTest", "ANALYSIS", 1012, err)
	}

	err = db.InsertDBSecurityTest(securityTest)
	if err != nil {
		log.Error("CreateNewSecurityTest", "ANALYSIS", 2016, err)
		requestResponse := util.RequestResponse(false, "internal error")
		return c.JSON(http.StatusInternalServerError, requestResponse)
	}

	log.Info("CreateNewSecurityTest", "ANALYSIS", 18, securityTest.Name)
	requestResponse := util.RequestResponse(true, "")
	return c.JSON(http.StatusCreated, requestResponse)
}
