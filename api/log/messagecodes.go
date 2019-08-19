package log

// MsgCode holds all log messages and their respective codes.
var MsgCode = map[int]string{

	// HuskyCI API infos
	11: "Starting HuskyCI.",
	12: "Environment variables set properly.",
	13: "Docker API is up and running.",
	14: "Connection with MongoDB succeed.",
	15: "Default securityTests found on MongoDB.",
	16: "Request received to start the following branch, repository and internal dependencies URL: ",
	17: "Repository created into MongoDB: ",
	18: "SecurityTest created into MongoDB: ",
	19: "SecurityTest upserted in MondoDB: ",
	20: "Default User found in MongoDB.",
	24: "URL received",

	// HuskyCI API warnings
	104: "An analysis is already in place for this URL: ",
	105: "The following analysis timed out inside MonitorAnalysis: ",
	106: "Analysis not found using the following RID: ",
	107: "Received an invalid RID: ",
	108: "Received an invalid security Test JSON: ",
	109: "The following security test is already in MongoDB: ",
	110: "The following repository is already in MongoDB: ",

	// HuskyCI API errors
	1001: "Error(s) found when starting HuskyCI API: ",
	1002: "Could not Unmarshall the following gosecOutput: ",
	1003: "Could not Unmarshall the following enryOutput: ",
	1004: "Error mapping languages: ",
	1005: "Could not Unmarshall the following brakemanOutput: ",
	1006: "Could not Unmarshall the following banditOutput: ",
	1007: "Could not bind repository JSON: ",
	1008: "Internal error MatchString: ",
	1009: "MongoDB message in FindOneDBAnalysis: ",
	1010: "Internal error inserting repository received into MongoDB: ",
	1011: "Internal error finding repository just inserted into MongoDB: ",
	1012: "MongoDB message in FindOneDBSecurityTest: ",
	1013: "MongoDB message in FindOneDBRepository: ",
	1014: "Could not Unmarshall the following retirejsOutput: ",
	1015: "Received an invalid repository JSON: ",
	1016: "Received an invalid repository URL: ",
	1017: "Received an invalid repository branch: ",
	1018: "Could not Unmarshall the following safetyOutput: ",
	1019: "Error loading viper: ",
	1020: "Error searching for an analysis: ",
	1021: "Received an invalid internal dependency URL: ",
	1022: "Could not Unmarshall the following npmauditOutput: ",
	1023: "Could not upsert securityTest into MongoDB: ",
	1024: "Received an invalid user JSON: ",
	1025: "Received an invalid Token JSON: ",
	1026: "Error during access token generation",
	1027: "Request doesn't have permission",
	1028: "Error during access token deactivation",

	// MongoDB infos
	21: "Connecting to MongoDB.",
	22: "Initializing MongoDB auto reconnect.",
	23: "Reconnect to MongoDB successful.",

	// MongoDB warnings
	201: "Enry securityTest not found.",
	202: "Gosec securityTest not found.",
	203: "Brakeman securityTest not found.",
	204: "Bandit securityTest not found.",
	205: "RetireJS securityTest not found.",
	206: "Safety securityTest not found.",

	// MongoDB errors
	2001: "Error connecting to MongoDB: ",
	2002: "Error pinging MongoDB after connection: ",
	2003: "Error pinging MongoDB in autoReconnect: ",
	2004: "Reconnect to MongoDB failed: ",
	2005: "Could not find default securityTests: ",
	2006: "Could not find securityTestName: ",
	2007: "Could not update AnalysisCollection: ",
	2008: "Could not find an analysis using the following CID: ",
	2009: "Error finding securityTest (language=Generic and default=true): ",
	2010: "Could not update repository's securityTests: ",
	2011: "Error inserting new analysis: ",
	2012: "Could not find securityTest into MongoDB using the following name: ",
	2013: "Could not update container status to timedout of an analysis: ",
	2014: "Could not find an analysis using the following RID: ",
	2015: "Could not create a new repository: ",
	2016: "Could not create a new securityTest: ",
	2017: "Could not update repository's internaldepURL: ",

	// Docker API info
	31: "Waiting pull image...",
	32: "Container started successfully: ",
	33: "Max container count reached. huskyCI is about to kill containers. ",

	// Docker API warning
	301: "",

	// Docker API errors
	3001: "Could not set DOCKER_HOST enviroment variable.",
	3002: "Could not start a new Docker API client: ",
	3005: "Could not create a new container via d.client: ",
	3006: "Could not get containers' logs: ",
	3007: "Could not read containers' STDOUT: ",
	3008: "Could not read containers' STDERR: ",
	3009: "Could not pull image into Docker API via d.client: ",
	3010: "Could not get docker image list from Docker API: ",
	3011: "Docker API Healthcheck failed: ",
	3012: "Could not create a new docker via HuskyCI: ",
	3013: "Could not pull image via HuskyCI: ",
	3014: "Could not create a new container via HuskyCI: ",
	3015: "Could not start a new container via HuskyCI: ",
	3016: "Could not wait container via HuskyCI: ",
	3017: "Could not read container output via HuskyCI: ",
	3018: "Unexpected securityTest.Name: ",
	3019: "Could not set DOCKER_CERT_PATH enviroment variable: ",
	3020: "Could not set DOCKER_TLS_VERIFY enviroment variable: ",
	3021: "Could not list current active containers: ",
	3022: "Could not stop a container via d.client: ",
	3023: "Could not remove a container via d.client: ",
	3024: "Could not call die containers: ",
	3025: "Could not update listed containers: ",
	3026: "Could not initialize default configurations: ",

	// Util package errors
	4001: "Could not read certificate file: ",
	4002: "Could not append ceritificates: ",
}
