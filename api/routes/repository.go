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
		return c.String(http.StatusBadRequest, "This is not a valid repository JSON.\n")
	}

	repositoryQuery := map[string]interface{}{"URL": repository.URL}
	_, err = db.FindOneDBRepository(repositoryQuery)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Warning("CreateNewRepository", "ANALYSIS", 110, repository.URL)
			return c.String(http.StatusConflict, "This repository is already in MongoDB.\n")
		}
		log.Error("CreateNewRepository", "ANALYSIS", 1013, err)
	}

	err = db.InsertDBRepository(repository)
	if err != nil {
		log.Error("CreateNewRepository", "ANALYSIS", 2015, err)
		return c.String(http.StatusInternalServerError, "Internal error 2015.\n")
	}

	log.Info("CreateNewRepository", "ANALYSIS", 17, repository.URL)
	return c.String(http.StatusCreated, "Repository sucessfully created.\n")
}
