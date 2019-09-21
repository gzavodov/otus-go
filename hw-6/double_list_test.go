package hw6

import (
	"errors"
	"fmt"
	"testing"
)

//Error Constructors
func NewUnexpectedValueError(expected, received interface{}) error {
	return fmt.Errorf("the node must have value \"%s\", but received \"%s\"", expected, received)
}

func NewUnexpectedLengthError(expected, received int) error {
	return fmt.Errorf("the length of list must be %d, but received %d", expected, received)
}

func NewInvalidHeadError() error {
	return errors.New("the first node of the list is not valid")
}

func NewInvalidTailError() error {
	return errors.New("the last node of the list is not valid")
}

func TestList(t *testing.T) {
	list := &List{}

	t.Run("Checking of PushToHead",
		func(t *testing.T) {
			list.Reset()
			list.PushToHead(0)
			list.PushToHead(1)
			list.PushToHead(2)

			if !(list.head != nil && list.head.value.(int) == 2 && list.GetLength() == 3) {
				t.Error(NewUnexpectedValueError(2, list.head.value))
			}
		})

	t.Run("Checking of PushToTail",
		func(t *testing.T) {
			list.Reset()
			list.PushToTail(0)
			list.PushToTail(1)
			list.PushToTail(2)

			if !(list.tail != nil && list.tail.value.(int) == 2 && list.GetLength() == 3) {
				t.Error(NewUnexpectedValueError(2, list.tail.value))
			}
		})

	t.Run("Checking of GetHead, GetTail and GetLength",
		func(t *testing.T) {
			list.Reset()
			node0 := list.PushToTail(0)
			node1 := list.PushToTail(1)

			if list.GetHead() != node0 {
				t.Error(NewInvalidHeadError())
			}

			if list.GetTail() != node1 {
				t.Error(NewInvalidTailError())
			}

			if list.GetLength() != 2 {
				t.Error(NewUnexpectedLengthError(2, list.GetLength()))
			}
		})

	t.Run("Checking of Remove (passing of the Null)",
		func(t *testing.T) {
			list.Reset()
			if list.Remove(nil) == nil {
				t.Error("Removing of Null shold fail.\n")
			}
		})

	t.Run("Checking of Remove (passing of fake node)",
		func(t *testing.T) {
			list.Reset()
			if list.Remove(&Node{}) == nil {
				t.Error("Removing of the Fake Node shold fail.\n")
			}
		})

	t.Run("Checking of Remove (passing of the head)",
		func(t *testing.T) {
			list.Reset()
			node0 := list.PushToTail(0)
			node1 := list.PushToTail(1)
			node2 := list.PushToTail(2)

			err := list.Remove(node0)
			if err != nil {
				t.Fatal(err)
			}

			if list.head != node1 {
				t.Error(NewInvalidHeadError())
			}

			if list.tail != node2 {
				t.Error(NewInvalidTailError())
			}

			if list.GetLength() != 2 {
				t.Error(NewUnexpectedLengthError(2, list.GetLength()))
			}
		})

	t.Run("Checking of Remove (passing of the tail)",
		func(t *testing.T) {
			list.Reset()
			node0 := list.PushToTail(0)
			node1 := list.PushToTail(1)
			node2 := list.PushToTail(2)

			err := list.Remove(node2)
			if err != nil {
				t.Fatal(err)
			}

			if list.head != node0 {
				t.Error(NewInvalidHeadError())
			}

			if list.tail != node1 {
				t.Error(NewInvalidTailError())
			}

			if list.GetLength() != 2 {
				t.Error(NewUnexpectedLengthError(2, list.GetLength()))
			}
		})

	t.Run("Checking of Remove (passing of ordinary node)",
		func(t *testing.T) {
			list.Reset()
			node0 := list.PushToTail(0)
			node1 := list.PushToTail(1)
			node2 := list.PushToTail(2)

			err := list.Remove(node1)
			if err != nil {
				t.Fatal(err)
			}

			if list.head != node0 {
				t.Error(NewInvalidHeadError())
			}

			if list.tail != node2 {
				t.Error(NewInvalidTailError())
			}

			if list.GetLength() != 2 {
				t.Error(NewUnexpectedLengthError(2, list.GetLength()))
			}
		})

	t.Run("Checking of navigation (forward and backward)",
		func(t *testing.T) {
			list.Reset()
			list.PushToTail(0)
			list.PushToTail(1)
			list.PushToTail(2)

			var nodeIndex int

			//Node value must be equal to node index
			nodeIndex = 0
			for curentNode := list.head; curentNode != nil; curentNode = curentNode.next {
				currentValue := curentNode.value.(int)
				if nodeIndex != currentValue {
					t.Error(NewUnexpectedValueError(nodeIndex, currentValue))
				}
				nodeIndex++
			}

			//Node value must be equal to node index
			nodeIndex = 2
			for curentNode := list.tail; curentNode != nil; curentNode = curentNode.prev {
				currentValue := curentNode.value.(int)
				if nodeIndex != currentValue {
					t.Error(NewUnexpectedValueError(nodeIndex, currentValue))
				}
				nodeIndex--
			}
		})
}
