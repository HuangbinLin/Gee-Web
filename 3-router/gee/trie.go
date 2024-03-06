package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) insert(pattern string, parts []string, height int) {
	// 会分层，第一层默认是/，/hellow就放到第二层，依次类推
	if len(parts) == height {
		n.pattern = pattern //最后一层才赋值
		return
	}

	part := parts[height] // 取出字符
	// 直接跑到下一层子节点列表，子节点是否匹配，不匹配则新建一个节点，匹配不做处理，进行下一个节点的匹配
	child := n.matchChild(part)
	if child == nil {
		// 如果是：或者*开头，则标记为isWild
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child) // children是一个列表
	}
	child.insert(pattern, parts, height+1) // 在子节点后面插入
}

func (n *node) search(parts []string, height int) *node {
	//是否到最后一个节点了或者HasPrefi用于检查字符串这个节点的信息是否以 * 开头，返回当前节点
	// 遇到*就停止了，遇到：的话还是会检索下面的节点的。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part) // 查找这一层，有没有符合的，返回符合的列表

	for _, child := range children { // 正常来说列表内只有一个参数，直接拿出来，进行下一步的匹配
		result := child.search(parts, height+1) //返回最后的结果，None或者节点
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
