package urlutil

type (
	Matcher struct {
		tree *node
	}

	node struct {
		kind           kind
		prefix         string
		parent         *node
		staticChildren children
		paramChild     *node
		anyChild       *node
		pristinePath   string
		paramNames     []string
	}
	kind     uint8
	children []*node
)

const (
	staticKind kind = iota
	paramKind
	anyKind

	paramLabel = byte(':')
	anyLabel   = byte('*')
)

func Match(pattern string, path string) (bool, map[string]string) {
	m := NewMatcher()
	m.Add(pattern)
	res, params := m.Match(path)
	if res != pattern {
		return false, nil
	}
	return true, params
}

func NewMatcher() *Matcher {
	return &Matcher{
		tree: &node{},
	}
}

func (m *Matcher) Add(path string) {
	var paramNames []string
	pristinePath := path

	for i, lcpIndex := 0, len(path); i < lcpIndex; i++ {
		if path[i] == ':' {
			if i > 0 && path[i-1] == '\\' {
				path = path[:i-1] + path[i:]
				i--
				lcpIndex--
				continue
			}
			j := i + 1

			m.insert(path[:i], staticKind, "", nil)
			for ; i < lcpIndex && path[i] != '/'; i++ {
			}

			paramNames = append(paramNames, path[j:i])
			path = path[:j] + path[i:]
			i, lcpIndex = j, len(path)

			if i == lcpIndex {
				m.insert(path[:i], paramKind, pristinePath, paramNames)
			} else {
				m.insert(path[:i], paramKind, "", nil)
			}
		} else if path[i] == '*' {
			m.insert(path[:i], staticKind, "", nil)
			paramNames = append(paramNames, "*")
			m.insert(path[:i+1], anyKind, pristinePath, paramNames)
		}
	}

	m.insert(path, staticKind, pristinePath, paramNames)
}

func (m *Matcher) Remove(path string) bool {
	currentNode := m.tree
	var nodeToRemove *node
	prefixLen := 0
	for {
		if currentNode.pristinePath == path {
			nodeToRemove = currentNode
			break
		}
		if currentNode.kind == staticKind {
			prefixLen = prefixLen + len(currentNode.prefix)
		} else {
			for ; prefixLen < len(path) && path[prefixLen] != '/'; prefixLen++ {
			}
		}

		if prefixLen >= len(path) {
			break
		}

		next := path[prefixLen]
		switch next {
		case paramLabel:
			currentNode = currentNode.paramChild
		case anyLabel:
			currentNode = currentNode.anyChild
		default:
			currentNode = currentNode.findStaticChild(next)
		}

		if currentNode == nil {
			break
		}
	}

	if nodeToRemove == nil {
		return false
	}

	nodeToRemove.pristinePath = ""

	if nodeToRemove.isLeaf() {
		current := nodeToRemove
		for {
			parent := current.parent
			if parent == nil {
				break
			}
			switch current.kind {
			case staticKind:
				var index int
				for i, c := range parent.staticChildren {
					if c == current {
						index = i
						break
					}
				}
				parent.staticChildren = append(parent.staticChildren[:index], parent.staticChildren[index+1:]...)
			case paramKind:
				parent.paramChild = nil
			case anyKind:
				parent.anyChild = nil
			}

			if !parent.isLeaf() {
				break
			}
			current = parent
		}
	}

	return true
}

func (m *Matcher) Match(origin string) (string, map[string]string) {
	currentNode := m.tree

	var (
		search      = origin
		searchIndex = 0
		paramValues []string
	)

	backtrackToNextNodeKind := func(fromKind kind) (nextNodeKind kind, valid bool) {
		previous := currentNode
		currentNode = previous.parent
		valid = currentNode != nil

		// Next node type by priority
		if previous.kind == anyKind {
			nextNodeKind = staticKind
		} else {
			nextNodeKind = previous.kind + 1
		}

		if fromKind == staticKind {
			// when backtracking is done from static kind block we did not change search so nothing to restore
			return
		}

		// restore search to value it was before we move to current node we are backtracking from.
		if previous.kind == staticKind {
			searchIndex -= len(previous.prefix)
		} else if len(paramValues) > 0 {
			searchIndex -= len(paramValues[len(paramValues)-1])
			paramValues = paramValues[:len(paramValues)-1]
		}
		search = origin[searchIndex:]
		return
	}

	for {
		prefixLen := 0
		lcpLen := 0

		if currentNode.kind == staticKind {
			searchLen := len(search)
			prefixLen = len(currentNode.prefix)

			// LCP - Longest Common Prefix (https://en.wikipedia.org/wiki/LCP_array)
			max := prefixLen
			if searchLen < max {
				max = searchLen
			}
			for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
			}
		}

		if lcpLen != prefixLen {
			// No matching prefix, let's backtrack to the first possible alternative node of the decision path
			nk, ok := backtrackToNextNodeKind(staticKind)
			if !ok {
				return "", nil
			} else if nk == paramKind {
				goto Param
			} else {
				// Not found (this should never be possible for static node we are looking currently)
				break
			}
		}

		// The full prefix has matched, remove the prefix from the remaining search
		search = search[lcpLen:]
		searchIndex = searchIndex + lcpLen

		// Finish routing if is no request path remaining to search
		if search == "" && currentNode.pristinePath != "" {
			break
		}

		// Static node
		if search != "" {
			if child := currentNode.findStaticChild(search[0]); child != nil {
				currentNode = child
				continue
			}
		}

	Param:
		// Param node
		if child := currentNode.paramChild; search != "" && child != nil {
			currentNode = child
			i := 0
			l := len(search)
			if currentNode.isLeaf() {
				// when param node does not have any children (path param is last piece of route path) then param node should
				// act similarly to any node - consider all remaining search as match
				i = l
			} else {
				for ; i < l && search[i] != '/'; i++ {
				}
			}

			paramValues = append(paramValues, search[:i])
			search = search[i:]
			searchIndex = searchIndex + i
			continue
		}

	Any:
		// Any node
		if child := currentNode.anyChild; child != nil {
			// If any node is found, use remaining path for paramValues
			currentNode = child
			paramValues = append(paramValues, search)

			// update indexes/search in case we need to backtrack when no handler match is found
			searchIndex += +len(search)
			search = ""

			if currentNode.pristinePath != "" {
				break
			}
		}

		// Let's backtrack to the first possible alternative node of the decision path
		nk, ok := backtrackToNextNodeKind(anyKind)
		if !ok {
			break // No other possibilities on the decision path
		} else if nk == paramKind {
			goto Param
		} else if nk == anyKind {
			goto Any
		} else {
			// Not found
			break
		}
	}

	if currentNode == nil {
		return "", nil
	}
	params := make(map[string]string)
	for i, v := range paramValues {
		params[currentNode.paramNames[i]] = v
	}
	return currentNode.pristinePath, params
}

func (m *Matcher) insert(path string, t kind, pristinePath string, paramNames []string) {
	currentNode := m.tree
	search := path

	for {
		searchLen := len(search)
		prefixLen := len(currentNode.prefix)
		lcpLen := 0

		// LCP - Longest Common Prefix (https://en.wikipedia.org/wiki/LCP_array)
		max := prefixLen
		if searchLen < max {
			max = searchLen
		}
		for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		}

		if lcpLen == 0 {
			// At root node
			currentNode.prefix = search
			if pristinePath != "" {
				currentNode.kind = t
				currentNode.paramNames = paramNames
				currentNode.pristinePath = pristinePath
			}
		} else if lcpLen < prefixLen {
			n := newNode(
				currentNode.kind,
				currentNode.prefix[lcpLen:],
				currentNode,
				currentNode.staticChildren,
				currentNode.paramChild,
				currentNode.anyChild,
				currentNode.pristinePath,
				currentNode.paramNames,
			)
			for _, child := range currentNode.staticChildren {
				child.parent = n
			}
			if currentNode.paramChild != nil {
				currentNode.paramChild.parent = n
			}
			if currentNode.anyChild != nil {
				currentNode.anyChild.parent = n
			}

			// Reset parent node
			currentNode.kind = staticKind
			currentNode.prefix = currentNode.prefix[:lcpLen]
			currentNode.staticChildren = nil
			currentNode.pristinePath = ""
			currentNode.paramNames = nil
			currentNode.paramChild = nil
			currentNode.anyChild = nil

			// Only Static children could reach here
			currentNode.addStaticChild(n)

			if lcpLen == searchLen {
				// At parent node
				currentNode.kind = t
				if pristinePath != "" {
					currentNode.paramNames = paramNames
					currentNode.pristinePath = pristinePath
				}
			} else {
				// Create child node
				n = newNode(t, search[lcpLen:], currentNode, nil, nil, nil, pristinePath, paramNames)
				// Only Static children could reach here
				currentNode.addStaticChild(n)
			}
		} else if lcpLen < searchLen {
			search = search[lcpLen:]
			c := currentNode.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				currentNode = c
				continue
			}
			// Create child node
			n := newNode(t, search, currentNode, nil, nil, nil, pristinePath, paramNames)
			switch t {
			case staticKind:
				currentNode.addStaticChild(n)
			case paramKind:
				currentNode.paramChild = n
			case anyKind:
				currentNode.anyChild = n
			}
		} else {
			// Node already exists
			if pristinePath != "" {
				currentNode.paramNames = paramNames
				currentNode.pristinePath = pristinePath
			}
		}
		return
	}
}

func newNode(
	kind kind,
	prefix string,
	parent *node,
	staticChildren children,
	paramChild *node,
	anyChild *node,
	pristinePath string,
	paramNames []string,
) *node {
	return &node{
		kind:           kind,
		prefix:         prefix,
		parent:         parent,
		staticChildren: staticChildren,
		paramChild:     paramChild,
		anyChild:       anyChild,
		pristinePath:   pristinePath,
		paramNames:     paramNames,
	}
}

func (n *node) addStaticChild(c *node) {
	n.staticChildren = append(n.staticChildren, c)
}

func (n *node) findChildWithLabel(l byte) *node {
	if c := n.findStaticChild(l); c != nil {
		return c
	}
	if l == paramLabel {
		return n.paramChild
	}
	if l == anyLabel {
		return n.anyChild
	}
	return nil
}

func (n *node) findStaticChild(l byte) *node {
	for _, c := range n.staticChildren {
		if c.label() == l {
			return c
		}
	}
	return nil
}

func (n *node) isLeaf() bool {
	return len(n.staticChildren) == 0 && n.paramChild == nil && n.anyChild == nil
}

func (n *node) label() byte {
	return n.prefix[0]
}
