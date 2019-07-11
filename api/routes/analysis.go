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
		log.Error("GetAnalysis", "ANALYSIS", 1008, "RID regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	if !valid {
		log.Warning("GetAnalysis", "ANALYSIS", 107, RID)
		reply := map[string]interface{}{"success": false, "error": "invalid RID"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	analysisQuery := map[string]interface{}{"RID": RID}
	analysisResult, err := db.FindOneDBAnalysis(analysisQuery)
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Warning("GetAnalysis", "ANALYSIS", 106, RID)
			reply := map[string]interface{}{"success": false, "error": "analysis not found"}
			return c.JSON(http.StatusNotFound, reply)
		}
		log.Error("GetAnalysis", "ANALYSIS", 1020, err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
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
		reply := map[string]interface{}{"success": false, "error": "invalid repository JSON"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// check-01-a: is this a git repository URL?
	regexpGit := `((git|ssh|http(s)?)|((git@|gitlab@)[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`
	r := regexp.MustCompile(regexpGit)
	valid, err := regexp.MatchString(regexpGit, repository.URL)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository URL regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1016, repository.URL)
		reply := map[string]interface{}{"success": false, "error": "invalid repository URL"}
		return c.JSON(http.StatusBadRequest, reply)
	}
	matches := r.FindString(repository.URL)
	repository.URL = matches

	// check-01-b: is this a git repository branch?
	regexpBranch := `^[a-zA-Z0-9_\/.-]*$`
	valid, err = regexp.MatchString(regexpBranch, repository.Branch)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository Branch regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1017, repository.Branch)
		reply := map[string]interface{}{"success": false, "error": "invalid repository branch"}
		return c.JSON(http.StatusBadRequest, reply)
	}

	// check-01-c: is this a valid dependency URL?
	regexpInternalDepURL := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
	valid, err = regexp.MatchString(regexpInternalDepURL, repository.InternalDepURL)
	if err != nil {
		log.Error("ReceiveRequest", "ANALYSIS", 1008, "Repository Branch regexp ", err)
		reply := map[string]interface{}{"success": false, "error": "internal error"}
		return c.JSON(http.StatusInternalServerError, reply)
	}
	if !valid {
		log.Error("ReceiveRequest", "ANALYSIS", 1021, repository.Branch)
		reply := map[string]interface{}{"success": false, "error": "invalid internal dependency URL"}
		return c.JSON(http.StatusBadRequest, reply)
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
					reply := map[string]interface{}{"success": false, "error": "an analysis is already in place for this URL and branch"}
					return c.JSON(http.StatusConflict, reply)
				}
			}
			log.Error("ReceiveRequest", "ANALYSIS", 1009, err)
		}
	} else {
		// repository not found! insert it into MongoDB with default securityTests
		err = db.InsertDBRepository(repository)
		if err != nil {
			log.Error("ReceiveRequest", "ANALYSIS", 1010, err)
			reply := map[string]interface{}{"success": false, "error": "internal error"}
			return c.JSON(http.StatusInternalServerError, reply)
		}
		repositoryQuery := map[string]interface{}{"repositoryURL": repository.URL, "repositoryBranch": repository.Branch}
		repositoryResult, err = db.FindOneDBRepository(repositoryQuery)
		if err != nil {
			// well it was supposed to be there, after all, we just inserted it.
			log.Error("ReceiveRequest", "ANALYSIS", 1011, err)
			reply := map[string]interface{}{"success": false, "error": "internal error"}
			return c.JSON(http.StatusInternalServerError, reply)
		}
	}

	log.Info("ReceiveRequest", "ANALYSIS", 16, repository.Branch, repository.URL)
	go analysis.StartAnalysis(RID, repositoryResult)
	reply := map[string]interface{}{"success": true, "error": ""}
	return c.JSON(http.StatusCreated, reply)
}
