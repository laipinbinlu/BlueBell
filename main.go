package main

import (
	"blue_bell/controller"
	"blue_bell/dao/mysql"
	"blue_bell/dao/redis"
	"blue_bell/logger"
	"blue_bell/pkg/rabbitmq"
	"blue_bell/pkg/snowflake"
	"blue_bell/router"
	"blue_bell/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// go web 开发的通用脚手架
func main() {
	//// 通过命令行来指定配置文件地址
	//if len(os.Args) < 2 {
	//	fmt.Println("args too few,need filepath args")
	//	return
	//}

	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success")

	//3.连接mysql
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	//4.连接redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// 生成用户id
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake node, err:%v\n", err)
		return
	}

	// 注册gin框架内置的校验器中的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	// 5.初始化 RabbitMQ
	if err := rabbitmq.Init(); err != nil {
		fmt.Printf("init rabbitmq failed, err:%v\n", err)
		return
	}
	defer rabbitmq.RMQ.Close()

	// 启动投票消费者
	go rabbitmq.StartVoteConsumer()

	//6.注册路由
	r := router.Setup(settings.Conf, settings.Conf.Mode)

	//7.启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
