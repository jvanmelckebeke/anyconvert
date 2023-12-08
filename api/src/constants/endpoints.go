package constants

const statusEndpoint = "/status"
const resultEndpoint = "/result"
const SingleStatusEndpoint = statusEndpoint + "/:id"
const SingleResultEndpoint = resultEndpoint + "/:id"
const AllStatusEndpoint = statusEndpoint

func CreateStatusEndpoint(taskID string) string {
	return statusEndpoint + "/" + taskID
}

func CreateResultEndpoint(taskID string) string {
	return resultEndpoint + "/" + taskID
}
