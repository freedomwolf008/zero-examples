package main

import (
	"flag"
	"fmt"
	"log"

	"monolithic/internal/config"
	"monolithic/internal/handler"
	"monolithic/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/file-api.yaml", "the config file")

func main() {

	//对于切片来说，使用%p格式化输出时，如果前面不加取地址符，
	//那么打印的是切片中第一个元素的地址；如果前面加上取地址符&，那么打印的是该切片的地址
	var arr = []int{1, 2, 3, 4, 5}
	fmt.Printf("arr pointer: %p\n", arr)
	fmt.Printf("arr val: %v\n", arr)
	fmt.Printf("arr val: %p\n", &arr[0])
	fmt.Printf("arr pointer: %p\n", &arr)
	arr = append(arr, 6)
	fmt.Printf("after append,arr pointer: %p\n", &arr)

	//var count []int
	//var count [1]int
	//count := [1]int{0}
	count := make([]int, 1)
	log.Printf("count: %d\n", count[0])

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
