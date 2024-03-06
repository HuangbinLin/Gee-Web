package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern) // 通过/进行分割

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{} // 给GET和POST都新建一个节点
	}
	// 如果是/，则会把第一个节点赋值为"/"
	r.roots[method].insert(pattern, parts, 0) //给节点插入Node结构体
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 按照/进行字符拆分
	params := make(map[string]string) //上下文c的参数
	root, ok := r.roots[method]       //取得GET还是POST的根节点

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0) // 返回最后匹配的最后一个节点，或者None，考虑了isWild和*

	if n != nil {
		parts := parsePattern(n.pattern) //把检索的最后一个节点返回，拆分数据，n.pattern记录了整条数据
		for index, part := range parts { // 如果检索的节点包含:和*，把参数替换为真实对应的数据
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path) //GET   \hellow\geektutu
	if n != nil {                             // 有检索到结果
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
