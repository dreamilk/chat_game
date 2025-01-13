package handlers

type Resp struct {
	Code    int         `json:"code"`
	Data    interface{} `josn:"data"`
	Message string      `json:"message"`
}
