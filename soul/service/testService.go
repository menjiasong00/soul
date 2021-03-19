package service

import(
	pb "rest/pb"
	"golang.org/x/net/context"
	"fmt"
)


type TestServer struct{}

func (m *TestServer) GetTestMsg(c context.Context, s *pb.TestMessage) (*pb.TestMessage, error) {
	fmt.Printf("xxxxx(%q)\n", s.Value)
	return s, nil
}
