package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// 添加认证中间件
func authMiddleware() gin.HandlerFunc {
	accounts := make(gin.Accounts)

	// 读取配置文件
	configPath := "fileserver.conf" // 配置文件与可执行文件同级
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过注释和空行
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("配置文件格式错误: %s", line)
		}
		accounts[parts[0]] = parts[1]
	}

	return gin.BasicAuth(accounts)
}

func main() {
	r := gin.Default()

	// 应用认证中间件到所有路由
	r.Use(authMiddleware())

	// 添加禁用缓存中间件
	// r.Use(func(c *gin.Context) {
	// 	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	// 	c.Header("Pragma", "no-cache")
	// 	c.Header("Expires", "0")
	// })

	// 配置静态文件目录
	staticDir := "./files"
	// 自动创建目录（新增代码）
	if err := os.MkdirAll(staticDir, os.ModePerm); err != nil {
		log.Fatalf("创建目录失败: %v", err)
	}

	// 注册路由处理静态文件
	r.StaticFS("/", http.Dir(staticDir))

	// 启动服务器
	go r.Run(":1080")
	// 启动HTTPS服务器（需要证书文件）
	go r.RunTLS(":1443", "./cert.pem", "./key.pem")

	// 等待程序终止
	select {}
}
