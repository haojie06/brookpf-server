package main

//StatusResponse 注意 如果要编码的话 需要用大写开头指明导出
type StatusResponse struct {
	Code      int
	Installed bool
	Enable    bool
	Records   []string
}

type MessageResponse struct {
	Code int
	Msg  string
}

type DataResponse struct {
	Code int
	Data map[string]interface{}
}

type AddPortForwardRequest struct {
	LocalPort  uint
	RemotePort uint
	Host       string
	Enable     int
}
