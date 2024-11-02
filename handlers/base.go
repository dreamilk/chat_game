package handlers

type Return struct {
	Code    int         `json:"code"`
	Data    interface{} `josn:"data"`
	Message string      `json:"message"`
}
