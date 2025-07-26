package fuzzy

import (
	"sort"
	"strings"
)

type SuffixArray struct {
	text     string
	suffixes []int
}

func NewSuffixArray(text string) *SuffixArray {
	sa := &SuffixArray{text: text}
	sa.build()
	return sa
}

func (sa *SuffixArray) build() {
	n := len(sa.text)
	sa.suffixes = make([]int, n)
	
	for i := 0; i < n; i++ {
		sa.suffixes[i] = i
	}
	
	sort.Slice(sa.suffixes, func(i, j int) bool {
		return sa.text[sa.suffixes[i]:] < sa.text[sa.suffixes[j]:]
	})
}

func (sa *SuffixArray) Search(pattern string) []int {
	n := len(sa.suffixes)
	left := sort.Search(n, func(i int) bool {
		suffix := sa.text[sa.suffixes[i]:]
		return strings.HasPrefix(suffix, pattern) || suffix >= pattern
	})
	
	if left >= n || !strings.HasPrefix(sa.text[sa.suffixes[left]:], pattern) {
		return nil
	}
	
	right := left + 1
	for right < n && strings.HasPrefix(sa.text[sa.suffixes[right]:], pattern) {
		right++
	}
	
	results := make([]int, right-left)
	for i := left; i < right; i++ {
		results[i-left] = sa.suffixes[i]
	}
	
	return results
}

func (sa *SuffixArray) FuzzySearch(pattern string, maxErrors int) []int {
	var results []int
	seen := make(map[int]bool)
	
	for _, suffix := range sa.suffixes {
		remaining := sa.text[suffix:]
		if len(remaining) < len(pattern)-maxErrors {
			continue
		}
		
		dist := approximateMatch(pattern, remaining, maxErrors)
		if dist <= maxErrors {
			if !seen[suffix] {
				results = append(results, suffix)
				seen[suffix] = true
			}
		}
	}
	
	return results
}

func approximateMatch(pattern, text string, maxErrors int) int {
	m := len(pattern)
	n := len(text)
	
	if m == 0 {
		return 0
	}
	
	prev := make([]int, m+1)
	curr := make([]int, m+1)
	
	for i := 0; i <= m; i++ {
		prev[i] = i
	}
	
	for j := 1; j <= n; j++ {
		curr[0] = 0
		minDist := m
		
		for i := 1; i <= m; i++ {
			cost := 0
			if pattern[i-1] != text[j-1] {
				cost = 1
			}
			
			curr[i] = min(min(prev[i]+1, curr[i-1]+1), prev[i-1]+cost)
			
			if curr[i] < minDist {
				minDist = curr[i]
			}
		}
		
		if minDist > maxErrors {
			return maxErrors + 1
		}
		
		prev, curr = curr, prev
	}
	
	minDist := m
	for i := 0; i <= m; i++ {
		if prev[i] < minDist {
			minDist = prev[i]
		}
	}
	
	return minDist
}

type FMIndex struct {
	bwt        string
	firstOcc   map[rune]int
	occTable   map[rune][]int
	suffixArr  []int
	sampleRate int
}

func NewFMIndex(text string, sampleRate int) *FMIndex {
	fm := &FMIndex{sampleRate: sampleRate}
	fm.build(text)
	return fm
}

func (fm *FMIndex) build(text string) {
	text = text + "$"
	n := len(text)
	
	sa := NewSuffixArray(text)
	fm.suffixArr = make([]int, 0, (n+fm.sampleRate-1)/fm.sampleRate)
	
	bwtBytes := make([]byte, n)
	for i, suffix := range sa.suffixes {
		if suffix == 0 {
			bwtBytes[i] = text[n-1]
		} else {
			bwtBytes[i] = text[suffix-1]
		}
		
		if i%fm.sampleRate == 0 {
			fm.suffixArr = append(fm.suffixArr, suffix)
		}
	}
	fm.bwt = string(bwtBytes)
	
	fm.firstOcc = make(map[rune]int)
	fm.occTable = make(map[rune][]int)
	
	counts := make(map[rune]int)
	for _, r := range fm.bwt {
		counts[r]++
	}
	
	chars := make([]rune, 0, len(counts))
	for r := range counts {
		chars = append(chars, r)
	}
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})
	
	total := 0
	for _, r := range chars {
		fm.firstOcc[r] = total
		total += counts[r]
		fm.occTable[r] = make([]int, n+1)
	}
	
	for r := range fm.occTable {
		count := 0
		for i := 0; i < n; i++ {
			fm.occTable[r][i] = count
			if rune(fm.bwt[i]) == r {
				count++
			}
		}
		fm.occTable[r][n] = count
	}
}

func (fm *FMIndex) Count(pattern string) int {
	if len(pattern) == 0 {
		return 0
	}
	
	runes := []rune(pattern)
	m := len(runes)
	
	c := runes[m-1]
	first, ok := fm.firstOcc[c]
	if !ok {
		return 0
	}
	
	occ, ok := fm.occTable[c]
	if !ok {
		return 0
	}
	
	sp := first
	ep := first + occ[len(fm.bwt)] - 1
	
	for i := m - 2; i >= 0 && sp <= ep; i-- {
		c = runes[i]
		first, ok = fm.firstOcc[c]
		if !ok {
			return 0
		}
		
		occ, ok = fm.occTable[c]
		if !ok {
			return 0
		}
		
		sp = first + occ[sp]
		ep = first + occ[ep+1] - 1
	}
	
	if sp > ep {
		return 0
	}
	
	return ep - sp + 1
}

func (fm *FMIndex) Locate(pattern string) []int {
	count := fm.Count(pattern)
	if count == 0 {
		return nil
	}
	
	return []int{}
}