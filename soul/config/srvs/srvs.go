package srvs

import(
	pb "rest/pb"
	"rest/service"
	rb "rest/pkg/grbmq"
)

var ServerMap= make(map[string]map[string]interface{})

func init(){
	ServerMap["test"] = map[string]interface{}{"Host":"localhost:10000","Server":&service.TestServer{},"Register":pb.RegisterTestServer,"HandlerFromEndpoint":pb.RegisterTestHandlerFromEndpoint}


	//消费者配置
	Consumers := []rb.Consumer{
		{
			QueueName:"oa.employee.ihr",
			Workers:2,
			Controllers:map[string][]interface{}{
				"oa.employee.entry": {&service.TestServer{}, "GetTestMsg", &pb.TestMessage{}},

				"oa.employee.entry2": {&service.TestServer{}, "GetTestMsg", &pb.TestMessage{}},

			},
			Config : rb.ReceiverConfig{1, 1, true, false, false, true},
		},
		{
			QueueName:"oa.employee.rms",
			Workers:1,
			Controllers:map[string][]interface{}{
				"oa.employee.entry": {&service.TestServer{}, "GetTestMsg", &pb.TestMessage{}},

				"oa.employee.entry2": {&service.TestServer{}, "GetTestMsg", &pb.TestMessage{}},

			},
			Config : rb.ReceiverConfig{1, 1, true, false, false, true},
		},
	}
	rb.Mq.RunConsumerNew(Consumers)
}