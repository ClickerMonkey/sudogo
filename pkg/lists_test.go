package sudogo

import (
	"strings"
	"testing"
)

func TestQueue(t *testing.T) {
	tests := []struct{
		expected string
		create func() Queue[string]
	}{
		{
			expected: "abc",
			create: func() Queue[string] {
				q := NewQueue[string]()
				q.Offer("a")
				q.Offer("b")
				q.Offer("c")
				return q
			},
		},
		{
			expected: "",
			create: func() Queue[string] {
				q := NewQueue[string]()
				q.Offer("a")
				q.Offer("b")
				q.Offer("c")
				q.Poll()
				q.Poll()
				q.Poll()
				q.Poll()
				return q
			},
		},
		{
			expected: "cd",
			create: func() Queue[string] {
				q := NewQueue[string]()
				q.Offer("a")
				q.Offer("b")
				q.Offer("c")
				q.Poll()
				q.Offer("d")
				q.Poll()
				return q
			},
		},
		{
			expected: "d",
			create: func() Queue[string] {
				q := NewQueue[string]()
				q.Offer("a")
				q.Offer("b")
				q.Offer("c")
				q.Poll()
				q.Poll()
				q.Poll()
				q.Offer("d")
				return q
			},
		},
	}

	for testIndex, test := range tests {
		q := test.create()
		actual := queueString(q)

		if actual != test.expected {
			t.Errorf("TestQueue failed at %d, expected %s actual %s", testIndex, test.expected, actual)
		}
	}
}

func queueString(queue Queue[string]) string {
	sb := strings.Builder{}
	node := queue.head
	for node != nil {
		sb.WriteString(node.value)
		node = node.next
	}
	return sb.String()
}


func TestStack(t *testing.T) {
	tests := []struct{
		expected string
		create func() Stack[string]
	}{
		{
			expected: "abc",
			create: func() Stack[string] {
				q := NewStack[string](0)
				q.Push("a")
				q.Push("b")
				q.Push("c")
				return q
			},
		},
		{
			expected: "",
			create: func() Stack[string] {
				q := NewStack[string](1)
				q.Push("a")
				q.Push("b")
				q.Push("c")
				q.Pop()
				q.Pop()
				q.Pop()
				q.Pop()
				return q
			},
		},
		{
			expected: "ab",
			create: func() Stack[string] {
				q := NewStack[string](2)
				q.Push("a")
				q.Push("b")
				q.Push("c")
				q.Pop()
				q.Push("d")
				q.Pop()
				return q
			},
		},
		{
			expected: "d",
			create: func() Stack[string] {
				q := NewStack[string](3)
				q.Push("a")
				q.Push("b")
				q.Push("c")
				q.Pop()
				q.Pop()
				q.Pop()
				q.Push("d")
				return q
			},
		},
	}

	for testIndex, test := range tests {
		q := test.create()
		actual := stackString(q)

		if actual != test.expected {
			t.Errorf("TestStack failed at %d, expected %s actual %s", testIndex, test.expected, actual)
		}
	}
}

func stackString(stack Stack[string]) string {
	sb := strings.Builder{}
	for i := 0; i < stack.Size(); i++ {
		sb.WriteString(*stack.At(i))
	}
	return sb.String()
}