package analysis

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/globocom/husky/types"
	"gopkg.in/mgo.v2/bson"
)

// EnryStartAnalysis checks the languages of a repository, update them into mongoDB, and performs new analysis
func EnryStartAnalysis(CID string, cOutput string, RID string) {

	// step 0: get analysis based on RID
	analysisQuery := map[string]interface{}{"RID": RID}
	analysis, err := FindOneDBAnalysis(analysisQuery)
	if err != nil {
		fmt.Println("Could not find analysis by this RID:", err)
	}

	// step 1: get each language found in cOutput into repository.languages
	languagesRepository := []string{}
	reg, err := regexp.Compile(`[^a-zA-Z]+`)
	if err != nil {
		fmt.Println("Error regexp:", err)
	}
	outputWithoutUnicode := strings.Split(reg.ReplaceAllString(cOutput, " "), " ")
	for _, language := range outputWithoutUnicode {
		if language != "" {
			languagesRepository = append(languagesRepository, language)
		}
	}

	// step 2: for each language, include a default securityTest to the repository
	securityTestList := []types.SecurityTest{}
	for _, language := range languagesRepository {
		securityTestQuery := map[string]interface{}{"language": language, "default": true}
		securityTestResult, err := FindOneDBSecurityTest(securityTestQuery)
		if err == nil {
			securityTestList = append(securityTestList, securityTestResult)
		}
	}

	// step 3: update repository with new securityTest and languages taken from output
	repositoryQuery := map[string]interface{}{"URL": analysis.URL}
	updateRepositoryQuery := bson.M{
		"$set": bson.M{
			"securityTests": securityTestList,
			"languages":     languagesRepository,
		},
	}
	err = UpdateOneDBRepository(repositoryQuery, updateRepositoryQuery)
	if err != nil {
		fmt.Println("Could not update repository's securityTests:", err)
	}

	// step 4: update analysis with the new securityTests
	repositoryResult, err := FindOneDBRepository(repositoryQuery)
	if err != nil {
		fmt.Println("Error finding repository:", err)
	}
	for _, securityTest := range repositoryResult.SecurityTests {
		analysis.SecurityTests = append(analysis.SecurityTests, securityTest)
	}
	err = UpdateOneDBAnalysis(analysisQuery, analysis)
	if err != nil {
		fmt.Println("Error updating analysis:", err)
	}

	// step 5: start new securityTests
	for _, securityTest := range repositoryResult.SecurityTests {
		dockerRun(RID, &analysis, securityTest)
	}

}
