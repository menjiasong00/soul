package srvs

import(
	pb "github.com/philips/soul/pb"
	"github.com/philips/soul/service"
)

var ServerMap= make(map[string]map[string]interface{})

func init(){
	ServerMap["test"] = map[string]interface{}{"Host":"localhost:10000","Server":&service.TestServer{},"Register":pb.RegisterTestServer,"HandlerFromEndpoint":pb.RegisterTestHandlerFromEndpoint}
}