package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"micro-gin/global"
)

func getTracer() (opentracing.Tracer, io.Closer) {
	// 第一步：初始化配置信息
	cfg := &config.Configuration{
		// 采样率暂配置，设置为1，全部采样
		// 如果每个请求都保存到jeager中，压力会大，所以可以设置采集速率
		// 如：rateLimiting:每秒spans数
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,                          // 是否打印日志
			LocalAgentHostPort: global.App.Config.Jaeger.Host, // jeager默认端口是6831
		},
		ServiceName: global.App.Config.App.AppName, // 服务名字，也可以在下面NewTracer的时候传入，不过弃用了
	}

	// 第二步：通过配置，生成链路Trace
	// 传入服务名，和日志输出位置
	//cfg.New("lqz-service", config.Logger(jaeger.StdLogger))  // 弃用了
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer

}

func TrancerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取到tracer和closer
		tracer, closer := getTracer()
		// 关闭closer
		defer closer.Close()
		// 使用当前url地址创建一个span
		startSpan := tracer.StartSpan(c.Request.URL.Path)
		// 关闭span
		defer startSpan.Finish()
		// 放入ctx中
		c.Set("tracer", tracer)
		c.Set("startSpan", startSpan)
		// 继续往下执行
		c.Next()

	}
}
