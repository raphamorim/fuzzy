package fuzzy

// BKTree is a metric tree data structure for fast similarity search
type BKTree struct {
	root     *BKNode
	distance DistanceFunc
}

// BKNode represents a node in the BK-tree
type BKNode struct {
	word     string
	children []childNode
}

type childNode struct {
	distance int
	node     *BKNode
}

// DistanceFunc is a function that calculates distance between two strings
type DistanceFunc func(s1, s2 string) int

// NewBKTree creates a new BK-tree with the default Levenshtein distance
func NewBKTree() *BKTree {
	return &BKTree{
		distance: LevenshteinDistance,
	}
}

// NewBKTreeWithDistance creates a new BK-tree with a custom distance function
func NewBKTreeWithDistance(distFunc DistanceFunc) *BKTree {
	return &BKTree{
		distance: distFunc,
	}
}

// Add inserts a word into the BK-tree
func (t *BKTree) Add(word string) {
	if t.root == nil {
		t.root = &BKNode{word: word}
		return
	}

	node := t.root
	for {
		dist := t.distance(node.word, word)
		if dist == 0 {
			return // Word already exists
		}

		// Find child with matching distance
		found := false
		for _, child := range node.children {
			if child.distance == dist {
				node = child.node
				found = true
				break
			}
		}

		if !found {
			// Add new child
			node.children = append(node.children, childNode{
				distance: dist,
				node:     &BKNode{word: word},
			})
			return
		}
	}
}

// Search finds all words within maxDistance edits of the query
func (t *BKTree) Search(query string, maxDistance int) []string {
	if t.root == nil {
		return nil
	}

	var results []string
	candidates := []*BKNode{t.root}

	for len(candidates) > 0 {
		// Pop from stack
		node := candidates[len(candidates)-1]
		candidates = candidates[:len(candidates)-1]

		dist := t.distance(node.word, query)
		if dist <= maxDistance {
			results = append(results, node.word)
		}

		// Calculate search bounds
		minDist := dist - maxDistance
		maxDist := dist + maxDistance

		// Add children within bounds to candidates
		for _, child := range node.children {
			if child.distance >= minDist && child.distance <= maxDist {
				candidates = append(candidates, child.node)
			}
		}
	}

	return results
}

// SearchWithScores returns words with their distances
func (t *BKTree) SearchWithScores(query string, maxDistance int) []SearchResult {
	if t.root == nil {
		return nil
	}

	var results []SearchResult
	candidates := []*BKNode{t.root}

	for len(candidates) > 0 {
		// Pop from stack
		node := candidates[len(candidates)-1]
		candidates = candidates[:len(candidates)-1]

		dist := t.distance(node.word, query)
		if dist <= maxDistance {
			results = append(results, SearchResult{
				Word:     node.word,
				Distance: dist,
			})
		}

		// Calculate search bounds
		minDist := dist - maxDistance
		maxDist := dist + maxDistance

		// Add children within bounds to candidates
		for _, child := range node.children {
			if child.distance >= minDist && child.distance <= maxDist {
				candidates = append(candidates, child.node)
			}
		}
	}

	return results
}

// SearchResult contains a word and its distance from the query
type SearchResult struct {
	Word     string
	Distance int
}

// Size returns the number of words in the tree
func (t *BKTree) Size() int {
	if t.root == nil {
		return 0
	}
	return t.sizeNode(t.root)
}

func (t *BKTree) sizeNode(node *BKNode) int {
	count := 1
	for _, child := range node.children {
		count += t.sizeNode(child.node)
	}
	return count
}

// Standard Levenshtein Distance (optimized with two-row approach)
func LevenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}
	if s1 == s2 {
		return 0
	}

	// Work with bytes for ASCII strings (faster than runes)
	b1 := []byte(s1)
	b2 := []byte(s2)

	// Make sure b1 is the shorter string
	if len(b1) > len(b2) {
		b1, b2 = b2, b1
	}

	prev := make([]int, len(b1)+1)
	curr := make([]int, len(b1)+1)

	// Initialize first row
	for i := 0; i <= len(b1); i++ {
		prev[i] = i
	}

	// Fill matrix
	for j := 1; j <= len(b2); j++ {
		curr[0] = j
		for i := 1; i <= len(b1); i++ {
			cost := 0
			if b1[i-1] != b2[j-1] {
				cost = 1
			}
			curr[i] = min3(
				prev[i]+1,      // deletion
				curr[i-1]+1,    // insertion
				prev[i-1]+cost, // substitution
			)
		}
		prev, curr = curr, prev
	}

	return prev[len(b1)]
}

// DamerauLevenshteinDistance calculates the Damerau-Levenshtein distance
// allowing insertions, deletions, substitutions, and transpositions
func DamerauLevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	len1 := len(s1)
	len2 := len(s2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Create matrix
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first row and column
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)

			// Transposition
			if i > 1 && j > 1 &&
				s1[i-1] == s2[j-2] &&
				s1[i-2] == s2[j-1] {
				matrix[i][j] = min(matrix[i][j], matrix[i-2][j-2]+cost)
			}
		}
	}

	return matrix[len1][len2]
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min3(a, b, c int) int {
	return min(min(a, b), c)
}

// MyersDistance implements Myers' algorithm for computing edit distance
// This is optimized for small edit distances
func MyersDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}

	n := len(s1)
	m := len(s2)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	max := n + m
	v := make([]int, 2*max+1)
	offset := max

	v[offset+1] = 0

	for d := 0; d <= max; d++ {
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[offset+k-1] < v[offset+k+1]) {
				x = v[offset+k+1]
			} else {
				x = v[offset+k-1] + 1
			}

			y := x - k

			for x < n && y < m && s1[x] == s2[y] {
				x++
				y++
			}

			v[offset+k] = x

			if x >= n && y >= m {
				return d
			}
		}
	}

	return max
}