package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	targetURL   string
	proxyAPIKey string
	listenAddr  string
)

func init() {
	flag.StringVar(&targetURL, "target", "", "目标 OpenAI 格式 API 的地址")
	flag.StringVar(&proxyAPIKey, "key", "sk-dummy-key", "代理服务器将使用的 API Key")
	flag.StringVar(&listenAddr, "addr", ":8080", "代理服务器监听地址")
}

func main() {
	flag.Parse()

	if targetURL == "" {
		log.Fatal("必须提供目标 API 地址，使用 -target 参数")
	}

	r := gin.Default()

	// 处理所有 OpenAI 相关的路径
	r.Any("/v1/*path", handleProxy)

	log.Printf("代理服务器启动在 %s，目标地址：%s\n", listenAddr, targetURL)
	if err := r.Run(listenAddr); err != nil {
		log.Fatal(err)
	}
}

func handleProxy(c *gin.Context) {
	// 创建到目标服务器的请求
	targetReq, err := http.NewRequest(c.Request.Method, targetURL+c.Param("path"), c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 复制原始请求的 header
	for k, v := range c.Request.Header {
		if k != "Authorization" { // 跳过原始的认证头
			targetReq.Header[k] = v
		}
	}

	// 添加模拟的 API Key
	targetReq.Header.Set("Authorization", "Bearer "+proxyAPIKey)

	// 发送请求到目标服务器
	client := &http.Client{}
	resp, err := client.Do(targetReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 设置响应头
	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}
	c.Writer.WriteHeader(resp.StatusCode)

	// 写入响应体
	if _, err := io.Copy(c.Writer, bytes.NewReader(body)); err != nil {
		log.Printf("写入响应时发生错误: %v", err)
	}
}
