package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(key Key, v interface{}) *ListItem
	PushBack(key Key, v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Key   Key
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	first  *ListItem
	last   *ListItem
}

func (list *list) Len() int {
	return list.length
}

func (list *list) Front() *ListItem {
	return list.first
}

func (list *list) Back() *ListItem {
	return list.last
}

func (list *list) PushFront(key Key, v interface{}) *ListItem {
	newItem := &ListItem{Key: key, Value: v}

	if list.Len() == 0 {
		list.first = newItem
		list.last = newItem
	} else {
		list.first.Prev = newItem
		newItem.Next = list.first
		list.first = newItem
	}

	list.length++
	return newItem
}

func (list *list) PushBack(key Key, v interface{}) *ListItem {
	newItem := &ListItem{Key: key, Value: v}

	if list.Len() == 0 {
		list.first = newItem
		list.last = newItem
	} else {
		list.last.Next = newItem
		newItem.Prev = list.last
		list.last = newItem
	}

	list.length++
	return newItem
}

func (list *list) Remove(item *ListItem) {
	list.length--
	if list.length < 1 {
		list.first = nil
		list.last = nil
		item.Value = nil
		item = nil
		return
	}

	if item == list.last {
		list.last = item.Prev
	}

	if item == list.first {
		list.first = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	item.Value = nil
	item = nil
}

func (list *list) MoveToFront(item *ListItem) {
	if item == list.first {
		return
	}

	if item == list.last {
		list.last = item.Prev
		list.last.Next = nil
	} else {
		if item.Next != nil {
			item.Next.Prev = item.Prev
		}

		if item.Prev != nil {
			item.Prev.Next = item.Next
		}
	}

	item.Next = list.first
	item.Prev = nil
	list.first.Prev = item
	list.first = item
}

func NewList() List {
	return new(list)
}
