package data_structures

import (
	"testing"
)

func TestInsert(t *testing.T) {
	l := linkedList{}
	l.insert(1)
	l.insert(2)
	l.insert(3)
	l.display()
}

func TestInsertAsPos(t *testing.T) {
	l := linkedList{}
	l.insertAsPos(1, -1)
	l.insertAsPos(2, 3)
	l.insertAsPos(3, 0)
	l.insertAsPos(4, 2)
	l.display()
}
