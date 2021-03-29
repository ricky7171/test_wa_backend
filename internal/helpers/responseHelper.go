package helper

func FormatResponse(statusResponse string, obj interface{}) interface{} {
	response := make(map[string]interface{})

	response["status"] = statusResponse
	if statusResponse == "success" || statusResponse == "fail" {
		response["data"] = obj
	} else if statusResponse == "error" {
		response["message"] = obj
	} else {
		response["message"] = "Error on format response !"
	}
	return response

}
