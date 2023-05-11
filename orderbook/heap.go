package orderbook

type HeapNode struct {
	key   float64
	value interface{}
}

type Heap struct {
	keyIndexMap map[float64]int
	nodes       []*HeapNode
}

func NewHeap() *Heap {
	return &Heap{keyIndexMap: make(map[float64]int), nodes: []*HeapNode{}}
}

// Create - O(logn)
// Update - O(1)
func (h *Heap) CreateOrUpdate(key float64, val interface{}) {
	idx, ok := h.keyIndexMap[key]
	if !ok {
		node := &HeapNode{key: key, value: val}
		h.push(node)
		return
	}

	h.nodes[idx].value = val
}

// Delete - O(logn)
func (h *Heap) Delete(key float64) {
	if _, ok := h.keyIndexMap[key]; !ok {
		return
	}
	h.delete(h.keyIndexMap[key])
}

// Peek - O(1)
func (h *Heap) Peek() (*HeapNode, bool) {
	if h.len() == 0 {
		return nil, false
	}
	return h.nodes[0], true
}

func (h *Heap) len() int {
	return len(h.nodes)
}

func (h *Heap) less(i, j int) bool {
	return h.nodes[i].key < h.nodes[j].key
}

func (h *Heap) swap(i, j int) {
	iKey, jKey := h.nodes[i].key, h.nodes[j].key

	h.keyIndexMap[iKey] = j
	h.keyIndexMap[jKey] = i

	h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i]
}

func (h *Heap) push(node *HeapNode) {
	h.nodes = append(h.nodes, node)
	h.keyIndexMap[node.key] = len(h.nodes) - 1
	h.up(len(h.nodes) - 1)
}

func (h *Heap) delete(idx int) *HeapNode {
	key := h.nodes[idx].key

	n := len(h.nodes) - 1
	h.swap(idx, n)
	h.down(idx, n)
	result := h.nodes[n]

	h.nodes = h.nodes[:n]
	delete(h.keyIndexMap, key)
	return result
}

func (h *Heap) pop() *HeapNode {
	return h.delete(0)
}

func (h *Heap) up(i int) {
	for {
		if i == 0 {
			break
		}
		parent := (i - 1) / 2
		if h.less(parent, i) {
			break
		}
		h.swap(parent, i)
		i = parent
	}
}

func (h *Heap) down(i, n int) {
	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		smallest := left
		if right := left + 1; right < n && h.less(right, left) {
			smallest = right
		}
		if !h.less(smallest, i) {
			break
		}
		h.swap(smallest, i)
		i = smallest
	}
}
