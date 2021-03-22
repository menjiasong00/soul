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
	ConsumerSettings := []rb.ConsumerSetting{
		//ihr
		{
			QueueName:"oa.employee.ihr",
			RoutingKey:"oa.employee.entry",
			Workers:2,
			Service:&service.TestServer{},
			Controller:"GetTestMsg",
			Request:&pb.TestMessage{},
			Config : rb.ReceiverConfig{1, 1, true, false, false, false},
		},
		{
			QueueName:"oa.employee.ihr",
			RoutingKey:"oa.employee.out",
			Workers:1,
			Service:&service.TestServer{},
			Controller:"GetTestMsg",
			Request:&pb.TestMessage{},
			Config : rb.ReceiverConfig{1, 1, true, false, false, false},
		},

		//rms
		{
			QueueName:"oa.employee.rms",
			RoutingKey:"oa.employee.entry",
			Workers:1,
			Service:&service.TestServer{},
			Controller:"GetTestMsg",
			Request:&pb.TestMessage{},
			Config : rb.ReceiverConfig{1, 1, true, false, false, false},
		},
		{
			QueueName:"oa.employee.rms",
			RoutingKey:"oa.employee.out",
			Workers:1,
			Service:&service.TestServer{},
			Controller:"GetTestMsg",
			Request:&pb.TestMessage{},
			Config : rb.ReceiverConfig{1, 1, true, false, false, false},
		},
	}

	rb.Mq.RunConsumers(ConsumerSettings)

}