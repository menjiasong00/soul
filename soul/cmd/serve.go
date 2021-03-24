package cmd

import (
	grbmq "rest/pkg/grbmq"
	pb "rest/pb"
	service "rest/service"
	gcore "rest/pkg/gcore"
	//"rest/insecure"
	"github.com/spf13/cobra"
	"log"
)


var ConsumerSettings  = []grbmq.ConsumerSetting{}
var ServerSettings = []gcore.ServeSetting{}
var grpcPort = "50050"
var httpPort = "8080"

func init() {
	//微服务配置
	ServerSettings = []gcore.ServeSetting{
		//测试demo
		{
			Host:":50050",
			Server:&service.TestServer{},
			Register:pb.RegisterTestServer,
			HandlerFromEndpoint:pb.RegisterTestHandlerFromEndpoint,
		},
	}

	//消费者配置
	ConsumerSettings = []grbmq.ConsumerSetting{
		//ihr
		{
			QueueName:"oa.employee.ihr",
			RoutingKey:"oa.employee.entry",
			Workers:2,
			Service:&service.TestServer{},
			Controller:"GetTestMsg",
			Request:&pb.TestMessage{},
			Config : grbmq.ReceiverConfig{1, 1, true, false, false, false},
		},

	}

	RootCmd.AddCommand(serveCmd)
}


// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launches the example webserver on  "+demoAddr,
	Run: func(cmd *cobra.Command, args []string) {
		//grbmq.Mq.RunConsumers(ConsumerSettings)
		log.Printf("",args)
		log.Printf("",ConsumerSettings)
		go gcore.RunServeGRPC(ServerSettings,grpcPort)
		gcore.RunServeHTTP(ServerSettings,httpPort)
		//gcore.MakeInsecure(insecure.Key,insecure.Cert)
		//serve()

	},
}