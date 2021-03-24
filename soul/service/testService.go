package service

import(
	pb "rest/pb"
	"golang.org/x/net/context"
	"fmt"

	"rest/pkg/gcode"
)


type TestServer struct{}

func (m *TestServer) GetTestMsg(c context.Context, s *pb.TestMessage) (*pb.TestMessage, error) {
	fmt.Printf("xxxxx(%q)\n", s.Value)
	gcode.MakeCoding(gcode.MakeCodingRequest{
		Name:"产品",
		TableName:"products",
		ServerName:"Bas",
		ModuleName:"BaslProducts",
		DatabaseName:"test",
	})
	return s, nil
}

