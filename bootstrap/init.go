package bootstrap

import "micro-gin/global"

func InitializeAll() {
	// 初始化配置
	InitializeConfig()

	// 初始化日志
	global.App.Log = InitializeLog()

	//初始化数据库
	global.App.DB = InitializeDB()

	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()
	// 初始化验证器
	InitializeValidator()

	// 初始化Redis
	global.App.Redis = InitializeRedis()

	// 初始化文件系统
	InitializeStorage()
	////仅测试注册consul
	//c := InitializeConsul()
	//c.RegisterConsul("127.0.0.1", 8888, "test-consul", "test-consul", []string{"tttt", "abc"})
	// 初始化计划任务
	InitializeCron()
	// 启动服务器
	RunServer()
}
