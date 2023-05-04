// 这个文件是实现前缀路由的字典树
package gee

import "strings"

// this is new
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// matchChild is a method on the node struct that searches for a child node that matches
// the given part string. If a matching child node is found, it returns the node; otherwise, it returns nil.
// 如果找到匹配的子节点或通配符子节点（例如，“:”或“*”），则返回该子节点；否则，它返回零
// 假设我们想从根节点的子节点中为“api”部分找到一个匹配的子节点：
// 使用根节点和部分“api”调用 matchChild 函数。
// 该函数遍历根节点的子节点。
// 它检查每个子节点的部分是否与给定部分“api”匹配，或者子节点是否为通配符（例如，“:”或“*”）。
// 当它遇到“api”子节点时，它会找到匹配项并返回该子节点。
// 返回的节点是带有“api”部分的节点。如果没有匹配的子节点或通配符节点，该函数将返回 nil。
// /
// ├── p
// │   ├── :lang
// │   │   └── doc
// └── api
//
//	└── v1
//	    ├── users
//	    └── products
func (n *node) matchChild(part string) *node {
	// Iterate through the children of the current node.
	for _, child := range n.children {
		// Check if the child node's part matches the given part, or if the child node is a wildcard (e.g., : or *).
		if child.part == part || child.isWild {
			// If there's a match, return the child node.
			return child
		}
	}
	// If no matching child node is found, return nil.
	return nil
}

// matchChildren is a method on the node struct that searches for all child nodes
// that match the given part string. It returns a slice of matching child nodes.
func (n *node) matchChildren(part string) []*node {
	// Create an empty slice to store matching child nodes.
	nodes := make([]*node, 0)

	// Iterate through the children of the current node.
	for _, child := range n.children {
		// Check if the child node's part matches the given part, or if the child node is a wildcard (e.g., : or *).
		if child.part == part || child.isWild {
			// If there's a match, append the child node to the nodes slice.
			nodes = append(nodes, child)
		}
	}

	// Return the slice of matching child nodes.
	return nodes
}

// pattern: "/p/:lang/doc" parts: ["p", ":lang", "doc"] height:  trie 结构的当前深度
// insert adds a new route pattern to the node and its children recursively.
// 我们插入 "/p/:lang/doc" (parts = ["p", ":lang", "doc"]):
// “p”部分已经作为根节点的子节点存在，因此我们不创建新节点。
// 我们移动到下一个高度并将“:lang”部分作为“p”的子节点插入。
// 然后我们再次移动到下一个高度，将“doc”部分作为“:lang”的子节点插入。
// 最终的 trie 结构如下所示：
//   - (root)
//     |- p
//     |- doc
//     |- :lang
//     |- doc
func (n *node) insert(pattern string, parts []string, height int) {
	// If we reach the end of the parts array, store the pattern and return.
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// Get the current part to process and find a matching child node.
	part := parts[height]
	child := n.matchChild(part)

	// If there is no matching child, create a new one and add it to the children slice.
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// Recursively call insert on the child node, moving to the next height.
	child.insert(pattern, parts, height+1)
}

// part: ["p", "Go", "doc"] height:  trie 结构的当前深度
// search looks for a node that matches the given parts and returns it if found.
// 假设我们有如下的 trie 结构, 搜索 /p/go/doc
// /
// ├── p
// │   ├── :lang
// │   │   └── doc
// └── api
//
//	└── v1
//	    ├── users
//	    └── products
//
// 从根节点开始，该函数搜索与第一部分“p”匹配的子节点。
// 它找到一个匹配的子节点并递归调用搜索函数，下一部分“Go”和递增的高度为 1。
// 在下一级，它搜索与“Go”部分匹配的子节点。它会找到包含“:lang”部分的匹配子节点，因为“:lang”是一个可以匹配任何值的通配符。
// 它使用下一部分“doc”和递增的高度 2 递归调用搜索函数。
// 在下一级，它搜索与“doc”部分匹配的子节点。它找到一个匹配的子节点。
func (n *node) search(parts []string, height int) *node {
	// If we reach the end of the parts array or find a wildcard, check for a stored pattern.
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// If there is no pattern, return nil.
		if n.pattern == "" {
			return nil
		}
		// If there is a pattern, return the current node.
		return n
	}

	// Get the current part to process and find all matching child nodes.
	part := parts[height]
	children := n.matchChildren(part)
	// Iterate through the matching children and search recursively.
	for _, child := range children {
		// If a matching node is found, return it.
		// 递归调用search
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	// If no matching node is found, return nil.
	return nil
}
