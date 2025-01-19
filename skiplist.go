package rank

import (
	"bytes"
	"fmt"
	"math/rand"
)

// SkipListNode 表示跳表节点
type SkipListNode struct {
	value interface{}
	next  []*SkipListNode
	span  []int
}

// SkipList 表示跳表
type SkipList struct {
	head  *SkipListNode
	level int
	cmp   func(interface{}, interface{}) int
}

// NewSkipList 初始化一个新的跳表
func NewSkipList(cmp func(interface{}, interface{}) int) *SkipList {
	return &SkipList{
		head: &SkipListNode{
			next: make([]*SkipListNode, 1),
			span: []int{0},
		},
		level: 0,
		cmp:   cmp,
	}
}

// randomLevel 生成一个随机的跳表级别
func (sl *SkipList) randomLevel() int {
	level := 0
	for rand.Float64() < 0.5 {
		level++
	}
	if level > sl.level {
		level = sl.level
	}
	return level
}

// Insert 向跳表中插入元素
func (sl *SkipList) Insert(value interface{}) {
	level := sl.randomLevel()
	if level > sl.level {
		sl.level = level
		sl.head.next = append(sl.head.next, make([]*SkipListNode, level-sl.level)...)
		sl.head.span = append(sl.head.span, make([]int, level-sl.level)...)
	}
	newNode := &SkipListNode{
		value: value,
		next:  make([]*SkipListNode, level+1),
		span:  make([]int, level+1),
	}
	current := sl.head
	for i := sl.level; i >= 0; i-- {
		spanCount := 0
		for current.next[i] != nil && sl.cmp(current.next[i].value, value) < 0 {
			spanCount += current.span[i]
			current = current.next[i]
		}
		if i <= level {
			newNode.next[i] = current.next[i]
			newNode.span[i] = current.span[i] - spanCount
			current.next[i] = newNode
			current.span[i] = spanCount + 1
		}
	}
}

// Search 搜索跳表中的元素
func (sl *SkipList) Search(value interface{}) (int, *SkipListNode) {
	current := sl.head
	rank := 0
	for i := sl.level; i >= 0; i-- {
		for current.next[i] != nil && sl.cmp(current.next[i].value, value) < 0 {
			rank += current.span[i]
			current = current.next[i]
		}
	}
	if current.next[0] != nil && sl.cmp(current.next[0].value, value) == 0 {
		return rank, current.next[0]
	}
	return -1, nil
}

// Delete 删除跳表中的元素
func (sl *SkipList) Delete(value interface{}) {
	current := sl.head
	for i := sl.level; i >= 0; i-- {
		for current.next[i] != nil && sl.cmp(current.next[i].value, value) < 0 {
			current = current.next[i]
		}
		if current.next[i] != nil && sl.cmp(current.next[i].value, value) == 0 {
			current.span[i] += current.next[i].span[i] - 1
			current.next[i] = current.next[i].next[i]
		}
	}
}

// String 实现跳表的字符串表示
func (sl *SkipList) String() string {
	var buffer bytes.Buffer
	for i := sl.level; i >= 0; i-- {
		current := sl.head.next[i]
		buffer.WriteString(fmt.Sprintf("Level %d: ", i))
		for current != nil {
			buffer.WriteString(fmt.Sprintf("%v ", current.value))
			current = current.next[i]
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// Range 实现跳表的范围遍历
func (sl *SkipList) Range(start int, count int, f func(rank int, value interface{}) bool) {
	if start < 1 {
		start = 1
	}
	if count <= 0 {
		return
	}
	end := start + count
	if f == nil {
		return
	}
	current := sl.head.next[0]

	// TODO  从上往下遍历快速定位到起点
	curRank := 0
	for i := sl.level; i >= 0; i-- {
		for {
			if current.next[i] == nil {
				break
			}
			if newRank := curRank + current.span[i]; newRank < start {
				curRank = newRank
				current = current.next[i]
			}
		}
	}

	for rank := start; rank < end; rank++ {
		if !f(rank, current.value) {
			return
		}
		current = current.next[0]
	}
}
