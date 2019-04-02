package routes

import (
	"net/http"

	"github.com/globocom/huskyCI/api/analysis"
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
		return c.String(http.StatusBadRequest, "This is not a valid securityTest JSON.\n")
	}

	securityTestQuery := map[string]interface{}{"name": securityTest.Name}
	_, err = analysis.FindOneDBSecurityTest(securityTestQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewSecurityTest", "ANALYSIS", 109, securityTest.Name)
			return c.String(http.StatusConflict, "This securityTest is already in MongoDB.\n")
		}
		log.Error("CreateNewSecurityTest", "ANALYSIS", 1012, err)
	}

	err = analysis.InsertDBSecurityTest(securityTest)
	if err != nil {
		log.Error("CreateNewSecurityTest", "ANALYSIS", 2016, err)
		return c.String(http.StatusInternalServerError, "Internal error 2016.\n")
	}

	log.Info("CreateNewSecurityTest", "ANALYSIS", 18, securityTest.Name)
	return c.String(http.StatusCreated, "SecurityTest sucessfully created.\n")
}
