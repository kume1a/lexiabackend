package shared

type HttpErrorDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type OkDTO struct {
	Ok bool `json:"ok"`
}
