package main

import (
	"context"
	"fmt"
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
	linkt := initLink()
	for i := 1; i < 10; i++ {
		addLink(i, linkt)
	}
	addIndexLink(1000, 2, linkt)
	addDataNameLink(100, 1000, linkt)
	showLink(linkt)
	size := unsafe.Sizeof(linkt.next)
	fmt.Printf("占用的字节数是：%d 字节\n", size)
	size2 := unsafe.Sizeof(linkt.data)
	fmt.Printf("占用的字节数是：%d 字节\n", size2)
	size3 := unsafe.Sizeof(*linkt)
	fmt.Printf("占用的字节数是：%d 字节\n", size3)
}
func main() {
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
