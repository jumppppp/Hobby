package main

import (
	"context"
	"fmt"
	"hobby/cplugin"
	"hobby/ctype"
	"hobby/utils"
	"time"
	"unsafe"
)

type link struct {
	data int
	next *link
	// prior *link
	// go1   *int
}

func initLink() *link {
	return &link{data: -1, next: nil}

}
func addLink(data int, linkt *link) {

	current := linkt
	for current.next != nil {
		current = current.next
	}
	newlink := &link{data: data, next: nil}
	current.next = newlink

}
func stackAddLink(data int, linkt *link) {
	if linkt.next == nil {
		newlink := &link{}
		newlink.data = data
		linkt.next = newlink
		newlink.next = nil
		return
	} else {
		addLink(data, linkt.next)
	}
}
func showLink(linkt *link) {
	sdata := linkt.next
	for sdata != nil {

		fmt.Printf("%p,%v\n", sdata, sdata)
		sdata = sdata.next
	}
}
func addIndexLink(data int, index int, linkt *link) {
	if index <= 0 {
		fmt.Println("Index should be a positive integer")
		return
	}
	current := linkt
	cin := 1
	for current.next != nil && cin < index {
		current = current.next
		cin++
	}
	newlink := &link{data: data, next: nil}
	newlink.next = current.next
	current.next = newlink
}
func addDataNameLink(data int, datax int, linkt *link) {
	current := linkt
	for current.next != nil && current.data != datax {
		current = current.next
	}
	newlink := &link{data: data, next: nil}
	newlink.next = current.next
	current.next = newlink
}
func run() {
	linkt := utils.InitLink()
	for i := 1; i < 10; i++ {
		data := ctype.LinkData{UUID: fmt.Sprintf("%v", i), Data: i}
		utils.AddRetLink(data, linkt)
	}
	data2 := ctype.LinkData{UUID: "1000", Data: "1000"}
	utils.AddIndexLink(data2, 2, linkt)
	data3 := ctype.LinkData{UUID: "100", Data: "100"}
	utils.AddDataNameLink(data3, "1000", linkt)
	utils.ShowLink(linkt)
	size := unsafe.Sizeof(linkt.Next)
	fmt.Printf("占用的字节数是：%d 字节\n", size)
	size2 := unsafe.Sizeof(linkt.LinkData)
	fmt.Printf("占用的字节数是：%d 字节\n", size2)
	size3 := unsafe.Sizeof(*linkt)
	fmt.Printf("占用的字节数是：%d 字节\n", size3)
}
func go1() {
	// 创建一个可取消的 Context
	ctx, cancel := context.WithCancel(context.Background())
	// 在 Context 中存储键值对
	valCtx := context.WithValue(ctx, "key1", "value1")
	defer cancel() // 函数结束时调用 cancel，清理相关资源
	size3 := unsafe.Sizeof(cancel)
	fmt.Printf("占用的字节数是：%d 字节\n", size3)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done(): // 如果 Context 被取消，立即返回
				fmt.Println("Got cancel signal:", ctx.Err())
				return
			default:
				fmt.Println("Working...")
				// 使用 ctx.Value(key) 来获取存储的值
				if val, ok := ctx.Value("key1").(string); ok {
					fmt.Println("Value from context:", val)
				}
				time.Sleep(500 * time.Millisecond) // 模拟一些工作
			}
		}
	}(valCtx)

	time.Sleep(2 * time.Second) // 让 goroutine 运行一段时间
	cancel()                    // 取消 Context，触发 ctx.Done() 信号
	time.Sleep(1 * time.Second) // 等待 goroutine 结束
}
func keyBoard() {
	OutBoardData := make(chan *ctype.KeyBoardData, 128)
	done := make(chan bool, 1)
	go cplugin.KeyBoardMain(OutBoardData, done)
	for {
		select {
		case <-done:
			return
		default:
			fmt.Println(<-OutBoardData)
		}

	}
}
func main() {
	tar := ctype.ProbeTarget{"192.168.231.137", "110"}
	ok := cplugin.Probe(tar, 10*time.Second, ctype.HTTP_("192.168.231.137", "25"))
	fmt.Println(string(ok.Response))
}
