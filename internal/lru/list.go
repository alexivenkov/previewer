package lru

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(string) *ListItem
	PushBack(string) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value string
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First, Last *ListItem
	len         int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) init(v string) {
	item := ListItem{
		Value: v,
	}

	l.First = &item
	l.Last = &item
	l.len++
}

func (l *list) PushFront(v string) *ListItem {
	if l.len == 0 {
		l.init(v)
		return l.First
	}

	item := ListItem{
		Value: v,
		Prev:  nil,
		Next:  l.First,
	}
	l.First.Prev = &item
	l.First = &item
	l.len++

	return &item
}

func (l *list) PushBack(v string) *ListItem {
	if l.len == 0 {
		l.init(v)
		return l.First
	}

	item := ListItem{
		Value: v,
		Prev:  l.Last,
		Next:  nil,
	}
	l.Last.Next = &item
	l.Last = &item
	l.len++

	return &item
}

func (l *list) Remove(i *ListItem) {
	if checkNilInput(i) {
		return
	}

	if i.Prev == nil {
		l.First = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.Last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if checkNilInput(i) {
		return
	}

	if i.Prev == nil {
		return
	}

	if i.Next == nil {
		i.Next = l.First
		i.Prev.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.First.Prev = i
	i.Next = l.First
	l.First = i
	i.Prev = nil
}

func checkNilInput(value *ListItem) bool {
	return value == nil
}

func NewList() List {
	return new(list)
}
