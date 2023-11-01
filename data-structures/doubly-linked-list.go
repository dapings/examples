package data_structures

import (
	"fmt"
)

type doublyNode struct {
	prev, next *doublyNode
	val        int
}

type doublyLinkedList struct {
	head *doublyNode
	tail *doublyNode
}

// 在链表的头部插入新节点
// 如果链表为空，则将新节点设置为头节点和尾节点。否则，将新节点插入到头节点之前，并更新头节点的pre节点。
func (l *doublyLinkedList) insertFront(val int) {
	dn := &doublyNode{val: val}

	if l.head == nil {
		l.head = dn
		l.tail = dn
	} else {
		dn.next = l.head
		l.head.prev = dn
		l.head = dn
	}
}

// 在链表的尾部插入新节点
// 如果链表为空，则将新节点设置为头节点和尾节点。否则，将新节点插入到尾节点之后，并更新尾节点的next节点。
func (l *doublyLinkedList) insertBack(val int) {
	dn := &doublyNode{val: val}

	if l.tail == nil {
		l.head = dn
		l.tail = dn
	} else {
		dn.prev = l.tail
		l.tail.next = dn
		l.tail = dn
	}
}

// 正向遍历并打印链表的所有节点的值
func (l *doublyLinkedList) displayForward() {
	if l.head == nil {
		fmt.Print("doubly linked list empty.")
		return
	}

	current := l.head
	for current != nil {
		fmt.Printf("%d ", current.val)
		current = current.next
	}
	fmt.Println()
}

// 反向遍历并打印链表的所有节点的值
func (l *doublyLinkedList) displayBackward() {
	if l.tail == nil {
		fmt.Print("doubly linked list empty.")
		return
	}

	current := l.tail
	for current != nil {
		fmt.Printf("%d ", current.val)
		current = current.prev
	}
	fmt.Println()
}
