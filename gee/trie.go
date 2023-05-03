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
// e.g.
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
// e.g.
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

// insert adds a new route pattern to the node and its children recursively.
// e.g.
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

// search looks for a node that matches the given parts and returns it if found.
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
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	// If no matching node is found, return nil.
	return nil
}
