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

// CreateNewRepository inserts the given repository into RepositoryCollection.
func CreateNewRepository(c echo.Context) error {
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		log.Warning("CreateNewRepository", "ANALYSIS", 101)
		requestResponse := util.RequestResponse(false, "invalid repository JSON")
		return c.JSON(http.StatusBadRequest, requestResponse)
	}

	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	_, err = db.FindOneDBRepository(repositoryQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewRepository", "ANALYSIS", 110, repository.URL)
			requestResponse := util.RequestResponse(false, "this repository is already registered")
			return c.JSON(http.StatusConflict, requestResponse)
		}
		log.Error("CreateNewRepository", "ANALYSIS", 1013, err)
	}

	err = db.InsertDBRepository(repository)
	if err != nil {
		log.Error("CreateNewRepository", "ANALYSIS", 2015, err)
		requestResponse := util.RequestResponse(false, "internal error")
		return c.JSON(http.StatusInternalServerError, requestResponse)
	}

	log.Info("CreateNewRepository", "ANALYSIS", 17, repository.URL)
	requestResponse := util.RequestResponse(true, "")
	return c.JSON(http.StatusCreated, requestResponse)
}
