package gee

import "strings"

type router struct {
	roots    map[string]*node       // 使用 roots 来存储每种请求方式的Trie 树根节点。
	handlers map[string]HandlerFunc // 使用 handlers 存储每种请求方式的 HandlerFunc 。
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

// 创建一个新的路由
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

		if item != "" { //开头有一个空字符，跳过
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 函数在路由器的 trie 结构中为给定的 HTTP 方法和路径搜索匹配的路由。如果找到匹配项，它会返回匹配节点和从路径中提取的参数映射。如果没有匹配项，则返回 nil。
//
//	method（字符串）：我们要为其找到匹配路由的 HTTP 方法（例如，GET、POST、PUT）。
//	path (string): 要与 trie 中存储的路由匹配的请求路径。
//
// 考虑具有以下注册路由的路由器：
// GET /api/v1/用户
// GET /api/v1/产品
// GET /p/:lang/文档
// 我们要为请求找到一个匹配的路由：GET /p/JavaScript/doc
// 方法是“GET”，路径是“/p/JavaScript/doc”。
// searchParts 数组将为 ["p", "JavaScript", "doc"]。
// params 映射被初始化为一个空映射。
// 找到“GET”方法的根节点。
// 在具有 searchParts 和 0 作为高度的根节点上调用搜索函数。
// 找到具有模式“/p/:lang/doc”的匹配节点 n。
// 匹配节点的模式被解析为部分：["p", ":lang", "doc"]。
// 对于“:lang”部分，相应的值“JavaScript”从 searchParts 中提取并作为 {“lang”: “JavaScript”} 添加到参数映射中。
// 该函数返回匹配的节点 n 和参数映射。
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	// 这部分比较难懂
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
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
