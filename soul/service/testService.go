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
		DatabaseName:"test",
		TableName:"products",
		Name:"产品",
		ServerName:"Bas",
		ModuleName:"BaslProducts",
	})
	return s, nil
}



// DlxConsumer 死信消费者
/*

CREATE TABLE `bas_mq_dlx` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `exchange` varchar(255) NOT NULL COMMENT '交换机',
  `routing_key` json NOT NULL COMMENT '路由标识',
  `queue` varchar(255) NOT NULL COMMENT '队列名称',
  `app_id` varchar(255) NOT NULL COMMENT 'app_id',
  `body` json DEFAULT NULL COMMENT '消息体',
  `header` json DEFAULT NULL COMMENT '消息头',
  `expired_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '过期时间',
  `create_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `create_by` varchar(200) DEFAULT '',
  `update_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `update_by` varchar(200) DEFAULT '更新人',
  `status` tinyint(4) DEFAULT '1' COMMENT '## 1 未处理  2 已处理',
  PRIMARY KEY (`id`),
  KEY `EX_IDX` (`exchange`),
  KEY `QUEUE_IDX` (`queue`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='mq消息死信表';

message DlxConsumerRequest{
    string header =1;
    string body =2;
}

func (s *TestServer) DlxConsumer(ctx context.Context, in *pb.DlxConsumerRequest) (*pb.BaseResponse, error) {

	type Xheader struct {
		Count int `json:"count"`
		Exchange string `json:"exchange"`
		Queue string `json:"queue"`
		Reason string `json:"reason"`
		RoutingKeys []string `json:"routing-keys"`
		Time time.Time `json:"time"`
	}
	
	type Header struct {
		AppId string `json:"appid"`
		XDeath []Xheader `json:"x-death"`
		XFirstDeathExchange string `json:"x-first-death-exchange"`
		XFirstDeathQueue string `json:"x-first-death-queue"`
		XFirstDeathReason string `json:"x-first-death-reason"`
	}

	var dlxHeader Header
	json.Unmarshal([]byte(in.Header), &dlxHeader)

	var firstXheader Xheader

	for k,v:= range dlxHeader.XDeath{
		if (k==0){
			firstXheader = v
		}else{
			if v.Time.Unix() < firstXheader.Time.Unix() {
				firstXheader = v
			}
		}
	}
	if firstXheader.Exchange =="" {
		firstXheader.Time = time.Now()
	}

	routingkey ,_:= json.Marshal(firstXheader.RoutingKeys)
	newBasMqDlx := model.BasMqDlx{
		Exchange:firstXheader.Exchange,
		Queue:firstXheader.Queue,
		Header:in.Header,
		Body:in.Body,
		CreateAt:firstXheader.Time,
		UpdateAt:firstXheader.Time,
		ExpiredAt:firstXheader.Time,
		Status:1,
		AppId:dlxHeader.AppId,
		RoutingKey: bytes.NewBuffer(routingkey).String(),
	}


	err:= gmysql.DB.Save(&newBasMqDlx).Error

	if err !=nil {
		return &pb.BaseResponse{}, err
	}

	return &pb.BaseResponse{}, nil
}
*/
