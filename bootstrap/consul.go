package bootstrap

import (
	"fmt"
	"micro-gin/global"
)
import consulapi "github.com/hashicorp/consul/api"

func InitializeConsul() *Consul {
	return &Consul{
		consulAddress: global.App.Config.Consul.Host,
	}
}

type Consul struct {
	consulAddress string
}

func (c *Consul) RegisterConsul(localIP string, localPort int, name string, id string, tags []string) error {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = id
	registration.Name = name //根据这个名称来找这个服务
	registration.Port = localPort
	//registration.Tags = []string{"lqz", "web"} //这个就是一个标签，可以根据这个来找这个服务，相当于V1.1这种
	registration.Tags = tags //这个就是一个标签，可以根据这个来找这个服务，相当于V1.1这种
	registration.Address = localIP

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d/health", registration.Address, registration.Port)
	check.Timeout = "5s"                         //超时
	check.Interval = "5s"                        //健康检查频率
	check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check
	// 注册服务到consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil

}

func (c *Consul) DeleteService(serviceId string) error {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}

	err = client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		return err
	}
	return nil

}

func (c *Consul) GetAllService() (map[string]*consulapi.AgentService, error) {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}
	res, err := client.Agent().Services()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Consul) FilterServiceByName(serviceName string) (map[string]*consulapi.AgentService, error) {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}
	filter := fmt.Sprintf(`Service=="%s"`, serviceName)
	res, err := client.Agent().ServicesWithFilter(filter) // 按服务名字过滤
	if err != nil {
		return nil, err
	}
	return res, nil

}
func (c *Consul) FilterServiceByPort(port int) (map[string]*consulapi.AgentService, error) {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = c.consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("consul client error : ", err)
	}
	filter := fmt.Sprintf(`Port=="%d"`, port)
	res, err := client.Agent().ServicesWithFilter(filter) // 按端口过滤
	if err != nil {
		return nil, err
	}
	return res, nil

}
