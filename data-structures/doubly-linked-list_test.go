package data_structures

import (
	"testing"
)

func TestInsertFront(t *testing.T) {
	l := doublyLinkedList{}
	l.insertFront(3)
	l.insertFront(2)
	l.insertFront(1)
	l.displayForward()
}

func TestInsertBack(t *testing.T) {
	l := doublyLinkedList{}
	l.insertBack(4)
	l.insertBack(5)
	l.insertBack(6)
	l.displayBackward()
}

func TestInsertDoubly(t *testing.T) {
	l := doublyLinkedList{}
	l.insertFront(3)
	l.insertFront(2)
	l.insertFront(1)
	l.insertBack(4)
	l.insertBack(5)
	l.insertBack(6)
	l.displayBackward()
}
