package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"micro-gin/global"
	"time"
)

func RecordLog() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		bodyByte, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))
		c.Next()
		logInfo := func(c *gin.Context, bodyByte []byte, start time.Time) string {
			param := gin.LogFormatterParams{
				Request: c.Request,
			}
			path := param.Request.URL.Path
			raw := param.Request.URL.RawQuery
			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()
			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			var statusColor, methodColor, resetColor string
			if param.IsOutputColor() {
				statusColor = param.StatusCodeColor()
				methodColor = param.MethodColor()
				resetColor = param.ResetColor()
			}

			if param.Latency > time.Minute {
				param.Latency = param.Latency.Truncate(time.Second)
			}

			return fmt.Sprintf("%v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n request: %s \n errors:%s",
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				string(bodyByte),
				param.ErrorMessage,
			)
		}(c, bodyByte, start)
		global.App.Log.Info(logInfo)
	}

}
