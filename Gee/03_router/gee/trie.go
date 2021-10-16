package gee

type node struct {
	pattern  string  // 待匹配路由 例如 /p/:lang
	part     string  // 路由中的一部分 录入  :lang
	children []*node // 子节点 录入 [doc, tutorial, intro]
	isWild   bool    // 是否非精确匹配 part 含有 : 或 * 时为true
}

// 第一个匹配成功的结点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.isWild || child.part == part {
			return child
		}
	}
	return nil
}

// 所有匹配成功的结点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.isWild || child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入子节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 已经匹配完了
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 除了路径的下一部分
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || n.part[0] == '*' {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
