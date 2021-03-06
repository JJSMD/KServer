package main

import (
	"KServer/manage"
	"KServer/manage/config"
	"KServer/server/oauth/services"
	"KServer/server/utils"
	"KServer/server/utils/msg"
	"time"
)

func main() {
	conf := config.NewManageConfig()
	conf.Message.Kafka = true
	conf.Server.Head = msg.OauthTopic
	conf.DB.Redis = true
	m := manage.NewManage(conf)
	// 新建一个服务管理器

	// 启动redis
	redisConf := config.NewRedisConfig(utils.RedisConFile)
	m.DB().Redis().StartMasterPool(redisConf.GetMasterAddr(), redisConf.Master.PassWord, redisConf.Master.MaxIdle, redisConf.Master.MaxActive)
	m.DB().Redis().StartSlavePool(redisConf.GetSlaveAddr(), redisConf.Slave.PassWord, redisConf.Slave.MaxIdle, redisConf.Slave.MaxActive)

	// 启动kafka
	kafkaConf := config.NewKafkaConfig(utils.KafkaConFile)
	_ = m.Message().Kafka().Send().Open([]string{kafkaConf.GetAddr()})

	oauth := services.NewOauth(m)

	m.Message().Kafka().AddRouter(msg.OauthTopic, msg.OauthId, oauth.ResponseOauth)
	m.Message().Kafka().StartListen([]string{kafkaConf.GetAddr()}, msg.OauthTopic, -1)

	// 服务中心监听
	//d:=generalService.NewIDiscovery(m)
	m.Message().Kafka().CallRegisterService(msg.OauthId, msg.OauthTopic, m.Server().GetId(), m.Server().GetHost(), m.Server().GetPort(), utils.KafkaType)

	m.Server().Start()

	// 注销服务中心
	m.Message().Kafka().CallLogoutService(msg.OauthId, msg.OauthTopic, m.Server().GetId())
	time.Sleep(5 * time.Second)
	Close(m)

}

func Close(m manage.IManage) {

	_ = m.DB().Redis().CloseMaster()
	_ = m.DB().Redis().CloseSlave()

}
