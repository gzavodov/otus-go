package hw6

import (
	"errors"
)

//Node represents an implementation of the Double Linked List Node
type Node struct {
	value interface{}
	next  *Node
	prev  *Node
}

//GetValue returns the value of Node
func (n *Node) GetValue() interface{} {
	return n.value
}

//GetNext returns the Next Node
func (n *Node) GetNext() *Node {
	return n.next
}

//GetPrevious returns the Previous Node
func (n *Node) GetPrevious() *Node {
	return n.prev
}

//List represents an implementation of the Double Linked List
type List struct {
	length int
	head   *Node
	tail   *Node
}

//GetLength returns Length of the List
func (l *List) GetLength() int {
	return l.length
}

//GetHead returns the First Node of the List
func (l *List) GetHead() *Node {
	return l.head
}

//GetTail returns the Last Node of the List
func (l *List) GetTail() *Node {
	return l.tail
}

//PushToHead creates new node with value specified in parameter and inserts it at beginning of the List
func (l *List) PushToHead(value interface{}) *Node {
	node := &Node{value: value}

	if l.tail == nil {
		l.tail = node
	} else {
		l.head.prev = node
		node.next = l.head
	}

	l.head = node
	l.length++

	return node
}

//PushToTail creates new node with value specified in parameter and appends it to end of the List
func (l *List) PushToTail(value interface{}) *Node {
	node := &Node{value: value}

	if l.head == nil {
		l.head = node
	} else {
		l.tail.next = node
		node.prev = l.tail
	}

	l.tail = node
	l.length++

	return node
}

//Remove delete node specified in parameter
//Error will be returned:
// - if 'node' argument is Nil
// - if 'node' argument is not belongs to the List
func (l *List) Remove(node *Node) error {
	if node == nil {
		return errors.New("the \"node\" argument must be not null")
	}

	//Ensure node belogs to this list
	isFound := false
	for curentNode := l.head; curentNode != nil; curentNode = curentNode.next {
		if curentNode == node {
			isFound = true
			break
		}
	}

	if !isFound {
		return errors.New("node is not belongs to this list")
	}

	next := node.next
	prev := node.prev

	if next != nil {
		next.prev = prev
	}
	if prev != nil {
		prev.next = next
	}

	if l.head == node {
		l.head = next
	}

	if l.tail == node {
		l.tail = prev
	}

	l.length--

	return nil
}

//Reset sets the Length of the List to zero
func (l *List) Reset() {
	if l.head == nil && l.tail == nil && l.length == 0 {
		return
	}

	l.head = nil
	l.tail = nil

	l.length = 0
}
