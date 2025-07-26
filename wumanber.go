package fuzzy

import (
	"container/heap"
	"fmt"
	"math"
)

type WuManber struct {
	pattern      string
	patternLen   int
	alphabet     map[rune]int
	alphabetSize int
}

func NewWuManber(pattern string) *WuManber {
	wm := &WuManber{
		pattern:    pattern,
		patternLen: len(pattern),
		alphabet:   make(map[rune]int),
	}
	
	for _, r := range pattern {
		if _, exists := wm.alphabet[r]; !exists {
			wm.alphabet[r] = wm.alphabetSize
			wm.alphabetSize++
		}
	}
	
	return wm
}

func (wm *WuManber) Search(text string, maxErrors int) []Match {
	if maxErrors == 0 {
		return wm.exactSearch(text)
	}
	
	matches := make([]Match, 0)
	seen := make(map[string]bool)
	
	R := make([]uint64, maxErrors+1)
	patternMask := make([]uint64, wm.alphabetSize)
	
	for i := 0; i < wm.alphabetSize; i++ {
		patternMask[i] = ^uint64(0)
	}
	
	for i, r := range wm.pattern {
		if idx, exists := wm.alphabet[r]; exists {
			patternMask[idx] &= ^(uint64(1) << i)
		}
	}
	
	for i := 0; i <= maxErrors; i++ {
		R[i] = ^((uint64(1) << i) - 1)
	}
	
	for j, r := range text {
		oldR := R[0]
		var charMask uint64
		
		if idx, exists := wm.alphabet[r]; exists {
			charMask = patternMask[idx]
		} else {
			charMask = ^uint64(0)
		}
		
		R[0] = ((R[0] << 1) | 1) & charMask
		
		for k := 1; k <= maxErrors; k++ {
			tmp := R[k]
			R[k] = ((R[k] << 1) & charMask) | oldR | ((oldR | R[k-1]) << 1) | 1
			oldR = tmp
		}
		
		for k := 0; k <= maxErrors; k++ {
			if (R[k] & (uint64(1) << (wm.patternLen - 1))) == 0 {
				start := j - wm.patternLen + 1
				if start < 0 {
					start = 0
				}
				
				key := fmt.Sprintf("%d-%d-%d", start, j+1, k)
				if !seen[key] {
					seen[key] = true
					matches = append(matches, Match{
						Start:    start,
						End:      j + 1,
						Distance: k,
					})
				}
				break
			}
		}
	}
	
	return matches
}

func (wm *WuManber) exactSearch(text string) []Match {
	matches := make([]Match, 0)
	textRunes := []rune(text)
	patternRunes := []rune(wm.pattern)
	
	for i := 0; i <= len(textRunes)-len(patternRunes); i++ {
		match := true
		for j := 0; j < len(patternRunes); j++ {
			if textRunes[i+j] != patternRunes[j] {
				match = false
				break
			}
		}
		if match {
			matches = append(matches, Match{
				Start:    i,
				End:      i + len(patternRunes),
				Distance: 0,
			})
		}
	}
	
	return matches
}

type Match struct {
	Start    int
	End      int
	Distance int
}

type UkkonenAStar struct {
	pattern string
	maxDist int
}

func NewUkkonenAStar(pattern string, maxDist int) *UkkonenAStar {
	return &UkkonenAStar{
		pattern: pattern,
		maxDist: maxDist,
	}
}

type astarNode struct {
	i, j     int
	g, h, f  int
	parent   *astarNode
	index    int
}

type priorityQueue []*astarNode

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*astarNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

func (ua *UkkonenAStar) Search(text string) []Match {
	patternRunes := []rune(ua.pattern)
	textRunes := []rune(text)
	m := len(patternRunes)
	n := len(textRunes)
	
	matches := make([]Match, 0)
	seen := make(map[string]bool)
	
	for startPos := 0; startPos < n; startPos++ {
		pq := make(priorityQueue, 0)
		heap.Init(&pq)
		
		start := &astarNode{
			i: 0,
			j: startPos,
			g: 0,
			h: m,
			f: m,
		}
		heap.Push(&pq, start)
		
		visited := make(map[string]bool)
		
		for pq.Len() > 0 {
			current := heap.Pop(&pq).(*astarNode)
			
			if current.g > ua.maxDist {
				continue
			}
			
			if current.i == m {
				key := fmt.Sprintf("%d-%d-%d", startPos, current.j, current.g)
				if !seen[key] {
					seen[key] = true
					matches = append(matches, Match{
						Start:    startPos,
						End:      current.j,
						Distance: current.g,
					})
				}
				break
			}
			
			key := fmt.Sprintf("%d,%d", current.i, current.j)
			if visited[key] {
				continue
			}
			visited[key] = true
			
			if current.j < n && current.i < m && textRunes[current.j] == patternRunes[current.i] {
				next := &astarNode{
					i:      current.i + 1,
					j:      current.j + 1,
					g:      current.g,
					parent: current,
				}
				next.h = m - next.i
				next.f = next.g + next.h
				heap.Push(&pq, next)
			}
			
			if current.g < ua.maxDist {
				if current.i < m {
					next := &astarNode{
						i:      current.i + 1,
						j:      current.j,
						g:      current.g + 1,
						parent: current,
					}
					next.h = m - next.i
					next.f = next.g + next.h
					heap.Push(&pq, next)
				}
				
				if current.j < n {
					next := &astarNode{
						i:      current.i,
						j:      current.j + 1,
						g:      current.g + 1,
						parent: current,
					}
					next.h = m - next.i
					next.f = next.g + next.h
					heap.Push(&pq, next)
				}
				
				if current.i < m && current.j < n {
					next := &astarNode{
						i:      current.i + 1,
						j:      current.j + 1,
						g:      current.g + 1,
						parent: current,
					}
					next.h = m - next.i
					next.f = next.g + next.h
					heap.Push(&pq, next)
				}
			}
		}
	}
	
	return matches
}

type BitonicSort struct {
	items []FuzzyMatch
}

type FuzzyMatch struct {
	Text     string
	Score    float64
	Distance int
}

func NewBitonicSort(matches []FuzzyMatch) *BitonicSort {
	size := 1
	for size < len(matches) {
		size *= 2
	}
	
	items := make([]FuzzyMatch, size)
	copy(items, matches)
	
	for i := len(matches); i < size; i++ {
		items[i] = FuzzyMatch{Score: -1, Distance: math.MaxInt32}
	}
	
	return &BitonicSort{items: items}
}

func (bs *BitonicSort) Sort() []FuzzyMatch {
	bs.bitonicSort(0, len(bs.items), true)
	
	result := make([]FuzzyMatch, 0)
	for _, item := range bs.items {
		if item.Score >= 0 {
			result = append(result, item)
		}
	}
	
	return result
}

func (bs *BitonicSort) bitonicSort(low, cnt int, dir bool) {
	if cnt > 1 {
		k := cnt / 2
		bs.bitonicSort(low, k, true)
		bs.bitonicSort(low+k, k, false)
		bs.bitonicMerge(low, cnt, dir)
	}
}

func (bs *BitonicSort) bitonicMerge(low, cnt int, dir bool) {
	if cnt > 1 {
		k := cnt / 2
		for i := low; i < low+k; i++ {
			bs.compareAndSwap(i, i+k, dir)
		}
		bs.bitonicMerge(low, k, dir)
		bs.bitonicMerge(low+k, k, dir)
	}
}

func (bs *BitonicSort) compareAndSwap(i, j int, dir bool) {
	if dir == (bs.items[i].Score < bs.items[j].Score) {
		bs.items[i], bs.items[j] = bs.items[j], bs.items[i]
	}
}