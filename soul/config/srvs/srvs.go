package srvs

import(
	pb "rest/pb"
	"rest/service"
)

var ServerMap= make(map[string]map[string]interface{})

func init(){
	ServerMap["test"] = map[string]interface{}{"Host":"localhost:10000","Server":&service.TestServer{},"Register":pb.RegisterTestServer,"HandlerFromEndpoint":pb.RegisterTestHandlerFromEndpoint}
}