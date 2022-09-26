package http

import "WeebChat/pkg/filesystem"

func GetNotFoundResponse(req HttpRequest) (res HttpResponse) {
	res = HttpResponse{
		Request:     req,
		HttpVersion: "1.1",
		StatusCode:  "404",
		StatusMsg:   "Not Found",
	}

	res.Headers = NewHeader()

	body, err := filesystem.GetFile("./404.html")
	if err != nil {
		return res
	}

	res.SetBody(body)

	return res
}

func GetInternalErrorResponse(req HttpRequest, res HttpResponse) HttpResponse {
	res.StatusCode = "500"
	res.StatusMsg = "Internal Error"

	body, err := filesystem.GetFile("./500.html")
	if err != nil {
		return res
	}

	res.SetBody(body)

	return res
}

func GetBadRequestResponse(req HttpRequest, res HttpResponse) HttpResponse {
	res.StatusCode = "400"
	res.StatusMsg = "Bad Request"

	return res
}
