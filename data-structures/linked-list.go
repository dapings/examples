package data_structures

import (
	"fmt"
)

type node struct {
	val int
	// 指向下一个节点
	next *node
}

type linkedList struct {
	// 指向链表的头节点
	head *node
}

// Insert  方法用于向链表中插入新节点。
// 如果链表为空，则将新节点设置为头节点。否则，遍历链表找到最后一个节点，将新节点插入到最后一个节点的后面。
func (l *linkedList) insert(val int) {
	newNode := &node{val: val}

	if l.head == nil {
		l.head = newNode
	} else {
		current := l.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
}

func (l *linkedList) insertAsPos(val, pos int) {
	newNode := &node{val: val}

	if pos <= 0 || l.head == nil {
		newNode.next = l.head
		l.head = newNode
	} else {
		current := l.head
		count := 0
		for current.next != nil && count < pos-1 {
			current = current.next
			count++
		}
		newNode.next = current.next
		current.next = newNode
	}
}

func (l *linkedList) display() {
	if l.head == nil {
		fmt.Println("linkedlist empty.")
		return
	}

	current := l.head
	for current != nil {
		fmt.Printf("%d ", current.val)
		current = current.next
	}
	fmt.Println()
}
