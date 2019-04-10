package routes

import (
	"net/http"
	"regexp"

	"github.com/globocom/huskyCI/api/analysis"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/globocom/huskyCI/api/util"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

// GetAnalysis returns the status of a given analysis given a RID.
func GetAnalysis(c echo.Context) error {

	RID := c.Param("id")
	regexpRID := `^[a-zA-Z0-9]*$`
	valid, err := regexp.MatchString(regexpRID, RID)
	if err != nil {
		log.Error("GetAnalysis", "ANALYSIS", 1008, "RID regexp ", err)
		requestResponse := util.RequestResponse(false, "internal error")
		return c.JSON(http.StatusInternalServerError, requestResponse)
	}
	if !valid {
		log.Warning("GetAnalysis", "ANALYSIS", 107, RID)
		requestResponse := util.RequestResponse(false, "invalid RID")
		return c.JSON(http.StatusBadRequest, requestResponse)
	}

	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Warning("GetAnalysis", "ANALYSIS", 106, RID)
			requestResponse := util.RequestResponse(false, "analysis not found")
			return c.JSON(http.StatusNotFound, requestResponse)
		}
		log.Error("GetAnalysis", "ANALYSIS", 1020, err)
		requestResponse := util.RequestResponse(false, "internal error")
		return c.JSON(http.StatusInternalServerError, requestResponse)
	}
	return c.JSON(http.StatusOK, analysisResult)
}

// ReceiveRequest receives the request and performs several checks before starting a new analysis.
func ReceiveRequest(c echo.Context) error {

	RID := c.Response().Header().Get(echo.HeaderXRequestID)

	// check-00: is this a valid JSON?
	repository := types.Repository{}
	err := c.Bind(&repository)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1015, err)
		requestResponse := util.RequestResponse(false, "invalid repository JSON")
		return c.JSON(http.StatusBadRequest, requestResponse)
	}

	// check-01-a: is this a git repository URL?
	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`
	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repository.URL)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository URL regexp ", err)
		requestResponse := util.RequestResponse(false, "internal error")
		return c.JSON(http.StatusInternalServerError, requestResponse)
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1016, repository.URL)
		requestResponse := util.RequestResponse(false, "invalid repository URL")
		return c.JSON(http.StatusBadRequest, requestResponse)
	}
	matches := r.FindString(repository.URL)
	repository.URL = matches

	// check-01-b: is this a git repository branch?
	regexpBranch := `^[a-zA-Z0-9_\/.-]*$`
	valid, err = regexp.MatchString(regexpBranch, repository.Branch)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository Branch regexp ", err)
		requestResponse := util.RequestResponse(false, "internal error")
		return c.JSON(http.StatusInternalServerError, requestResponse)
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1017, repository.Branch)
		requestResponse := util.RequestResponse(false, "invalid repository branch")
		return c.JSON(http.StatusBadRequest, requestResponse)
	}

	// check-02: is this repository in MongoDB?
	repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
	repositoryResult, err := db.FindOneDBRepository(repositoryQuery)
	if err == nil {
		// check-03: repository found! does it have a running status analysis?
		analysisQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		analysisResult, err := db.FindOneDBAnalysis(analysisQuery)
		if err != nil {
			if err != mgo.ErrNotFound {
				if analysisResult.Status == "running" {
					log.Warning("ReceiveRequest", "ANALYSIS", 104, analysisResult.URL)
					requestResponse := util.RequestResponse(false, "an analysis is already in place for this URL and branch")
					return c.JSON(http.StatusConflict, requestResponse)
				}
			}
			log.Error("ReceiveRequest", "ANALYSIS", 1009, err)
		}
	} else {
		// repository not found! insert it into MongoDB with default securityTests
		err = db.InsertDBRepository(repository)
		if err != nil {
			log.Error("ReceiveRequest", "ANALYSIS", 1010, err)
			requestResponse := util.RequestResponse(false, "internal error")
			return c.JSON(http.StatusInternalServerError, requestResponse)
		}
		repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		repositoryResult, err = db.FindOneDBRepository(repositoryQuery)
		if err != nil {
			// well it was supposed to be there, after all, we just inserted it.
			log.Error("ReceiveRequest", "ANALYSIS", 1011, err)
			requestResponse := util.RequestResponse(false, "internal error")
			return c.JSON(http.StatusInternalServerError, requestResponse)
		}
	}

	log.Info("ReceiveRequest", "ANALYSIS", 16, repository.Branch, repository.URL)
	go analysis.StartAnalysis(RID, repositoryResult)
	requestResponse := util.RequestResponse(true, "")
	return c.JSON(http.StatusCreated, requestResponse)
}
