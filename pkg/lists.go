package sudogo

type Stack[T any] struct {
	items []T
	size  int
}

func NewStack[T any](initialCapacity int) Stack[T] {
	return Stack[T]{
		items: make([]T, initialCapacity),
		size:  0,
	}
}

func (stack *Stack[T]) Peek() *T {
	if stack.size == 0 {
		return nil
	}
	return &stack.items[stack.size-1]
}

func (stack *Stack[T]) Pop() *T {
	if stack.size == 0 {
		return nil
	}
	stack.size--
	return &stack.items[stack.size]
}

func (stack *Stack[T]) Next() *T {
	if len(stack.items) <= stack.size {
		var value T
		stack.items = append(stack.items, value)
	}
	next := &stack.items[stack.size]
	stack.size++
	return next
}

func (stack *Stack[T]) Push(value T) {
	next := stack.Next()
	*next = value
}

func (stack *Stack[T]) Empty() bool {
	return stack.size == 0
}

func (stack *Stack[T]) Size() int {
	return stack.size
}

func (stack *Stack[T]) Clear() {
	stack.size = 0
}

func (stack *Stack[T]) At(index int) *T {
	if index < 0 || index >= stack.size {
		return nil
	}
	return &stack.items[index]
}

type Queue[T any] struct {
	head *QueueNode[T]
	tail *QueueNode[T]
	size int
}

func NewQueue[T any]() Queue[T] {
	return Queue[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

type QueueNode[T any] struct {
	value T
	next  *QueueNode[T]
}

func (queue *Queue[T]) Peek() *T {
	if queue.size == 0 {
		return nil
	}
	return &queue.head.value
}

func (queue *Queue[T]) Poll() *T {
	if queue.size == 0 {
		return nil
	}

	queue.size--
	value := &queue.head.value

	if queue.size == 0 {
		queue.Clear()
	} else {
		queue.head = queue.head.next
	}
	return value
}

func (queue *Queue[T]) Next() *T {
	node := &QueueNode[T]{}
	if queue.head == nil {
		queue.head = node
	}
	if queue.tail != nil {
		queue.tail.next = node
	}
	queue.tail = node
	queue.size++
	return &node.value
}

func (queue *Queue[T]) Offer(value T) {
	next := queue.Next()
	*next = value
}

func (queue *Queue[T]) Empty() bool {
	return queue.size == 0
}

func (queue *Queue[T]) Size() int {
	return queue.size
}

func (queue *Queue[T]) Clear() {
	queue.size = 0
	queue.tail = nil
	queue.head = nil
}
