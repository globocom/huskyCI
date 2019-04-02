package routes

import (
	"net/http"
	"regexp"

	"github.com/globocom/huskyCI/api/analysis"
	"github.com/globocom/huskyCI/api/db"
	"github.com/globocom/huskyCI/api/log"
	"github.com/globocom/huskyCI/api/types"
	"github.com/labstack/echo"
	mgo "gopkg.in/mgo.v2"
)

// GetAnalysis returns the status of a given analysis given a RID.
func GetAnalysis(c echo.Context) error {

	RID := c.Param("id")
	regexpRID := `^[a-zA-Z0-9]*$`
	valid, err := regexp.MatchString(regexpRID, RID)
	if err != nil {
		log.Error("StatusAnalysis", "ANALYSIS", 1008, "RID regexp ", err)
		return c.String(http.StatusInternalServerError, "Internal error 1008.\n")
	}
	if !valid {
		log.Warning("StatusAnalysis", "ANALYSIS", 107, RID)
		return c.String(http.StatusBadRequest, "This is not a valid RID.\n")
	}

	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Warning("StatusAnalysis", "ANALYSIS", 106, RID)
			return c.String(http.StatusNotFound, "Analysis not found.\n")
		}
		return c.String(http.StatusInternalServerError, "Internal Error.\n")
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
		return c.String(http.StatusBadRequest, "This is an invalid JSON.\n")
	}

	// check-01: is this a git repository URL and a branch?
	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`
	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repository.URL)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository URL regexp ", err)
		return c.String(http.StatusInternalServerError, "Internal error 1008.\n")
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1016, repository.URL)
		return c.String(http.StatusBadRequest, "This is not a valid repository URL.\n")
	}
	matches := r.FindString(repository.URL)
	repository.URL = matches

	regexpBranch := `^[a-zA-Z0-9_\.-]*$`
	valid, err = regexp.MatchString(regexpBranch, repository.Branch)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository Branch regexp ", err)
		return c.String(http.StatusInternalServerError, "Internal error 1008.\n")
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1017, repository.Branch)
		return c.String(http.StatusBadRequest, "This is not a valid branch.\n")
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
					return c.String(http.StatusConflict, "An analysis is already in place for this URL.\n")
				}
			}
			log.Error("ReceiveRequest", "ANALYSIS", 1009, err)
		}
	} else {
		// repository not found! insert it into MongoDB with default securityTests
		err = db.InsertDBRepository(repository)
		if err != nil {
			log.Error("ReceiveRequest", "ANALYSIS", 1010, err)
			return c.String(http.StatusInternalServerError, "Internal error 1010.\n")
		}
		repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		repositoryResult, err = db.FindOneDBRepository(repositoryQuery)
		if err != nil {
			// well it was supposed to be there, after all, we just inserted it.
			log.Error("ReceiveRequest", "ANALYSIS", 1011, err)
			return c.String(http.StatusInternalServerError, "Internal error 1011.\n")
		}
	}

	log.Info("ReceiveRequest", "ANALYSIS", 16, repository.Branch, repository.URL)
	go analysis.StartAnalysis(RID, repositoryResult)
	return c.JSON(http.StatusOK, map[string]string{"RID": RID, "result": "ok", "details": "Request received."})
}
