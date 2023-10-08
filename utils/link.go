package utils

import (
	"hobby/ctype"
)

func InitLink() *ctype.RetLink {
	return &ctype.RetLink{LinkData: ctype.LinkData{UUID: "HEAD"}, Next: nil, Prior: nil}

}
func AddRetLink(data ctype.LinkData, linkt *ctype.RetLink) {

	current := linkt
	for current.Next != nil {
		current = current.Next
	}
	newlink := &ctype.RetLink{LinkData: data, Next: nil, Prior: current}
	current.Next = newlink

}

func ShowLink(linkt *ctype.RetLink) {
	sdata := linkt
	for sdata != nil {

		LogPf("%p,%v\n", sdata, sdata)
		sdata = sdata.Next
	}
}
func AddIndexLink(data ctype.LinkData, index int, linkt *ctype.RetLink) {
	if index <= 0 {
		LogPf("[-]Index should be a positive integer")
		return
	}
	current := linkt
	cin := 1
	for current.Next != nil && cin < index {
		current = current.Next
		cin++
	}
	newlink := &ctype.RetLink{LinkData: data, Next: current.Next, Prior: current}
	if current.Next != nil { // 需要检查current.Next是否为nil
		current.Next.Prior = newlink // 更新新节点的Next的Prior指针
	}
	current.Next = newlink // 设置当前节点的Next为新节点

}
func AddDataNameLink(data ctype.LinkData, UUID string, linkt *ctype.RetLink) {
	current := linkt
	for current.Next != nil && current.LinkData.UUID != UUID {
		current = current.Next
	}

	newlink := &ctype.RetLink{LinkData: data, Next: current.Next, Prior: current}
	if current.Next != nil { // 需要检查current.Next是否为nil
		current.Next.Prior = newlink // 更新新节点的Next的Prior指针
	}
	current.Next = newlink // 设置当前节点的Next为新节点

}
func SelectLinkDatabyUUID(UUID string, linkt *ctype.RetLink) *ctype.RetLink {
	current := linkt
	for current.Next != nil && current.LinkData.UUID != UUID {
		current = current.Next
	}
	if current.Next == nil {
		return nil
	}
	return current
}
