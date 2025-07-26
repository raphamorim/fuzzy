package fuzzy

import (
	"github.com/cespare/xxhash/v2"
	"math"
	"math/rand"
	"sort"
	"unicode"
)

type LSH struct {
	numHashTables   int
	numHashFuncs    int
	hashTables      []map[uint64][]int
	hashFunctions   [][]hashFunc
	shingleSize     int
	corpus          []string
}

type hashFunc struct {
	a uint64
	b uint64
	p uint64
}

func NewLSH(numHashTables, numHashFuncs, shingleSize int) *LSH {
	lsh := &LSH{
		numHashTables: numHashTables,
		numHashFuncs:  numHashFuncs,
		shingleSize:   shingleSize,
		hashTables:    make([]map[uint64][]int, numHashTables),
		hashFunctions: make([][]hashFunc, numHashTables),
		corpus:        make([]string, 0),
	}
	
	for i := 0; i < numHashTables; i++ {
		lsh.hashTables[i] = make(map[uint64][]int)
		lsh.hashFunctions[i] = make([]hashFunc, numHashFuncs)
		for j := 0; j < numHashFuncs; j++ {
			lsh.hashFunctions[i][j] = hashFunc{
				a: rand.Uint64(),
				b: rand.Uint64(),
				p: 4294967311,
			}
		}
	}
	
	return lsh
}

func (lsh *LSH) Add(text string) int {
	id := len(lsh.corpus)
	lsh.corpus = append(lsh.corpus, text)
	
	shingles := lsh.getShingles(text)
	
	for i := 0; i < lsh.numHashTables; i++ {
		signature := lsh.computeSignature(shingles, lsh.hashFunctions[i])
		lsh.hashTables[i][signature] = append(lsh.hashTables[i][signature], id)
	}
	
	return id
}

func (lsh *LSH) getShingles(text string) []uint64 {
	runes := []rune(text)
	if len(runes) < lsh.shingleSize {
		return []uint64{xxhash.Sum64String(text)}
	}
	
	shingles := make([]uint64, 0, len(runes)-lsh.shingleSize+1)
	for i := 0; i <= len(runes)-lsh.shingleSize; i++ {
		shingle := string(runes[i : i+lsh.shingleSize])
		shingles = append(shingles, xxhash.Sum64String(shingle))
	}
	
	return shingles
}

func (lsh *LSH) computeSignature(shingles []uint64, hashFuncs []hashFunc) uint64 {
	minHashes := make([]uint64, len(hashFuncs))
	for i := range minHashes {
		minHashes[i] = math.MaxUint64
	}
	
	for _, shingle := range shingles {
		for i, hf := range hashFuncs {
			hash := (hf.a*shingle + hf.b) % hf.p
			if hash < minHashes[i] {
				minHashes[i] = hash
			}
		}
	}
	
	var signature uint64
	for _, minHash := range minHashes {
		signature = signature*31 + minHash
	}
	
	return signature
}

func (lsh *LSH) Query(text string, threshold float64) []int {
	shingles := lsh.getShingles(text)
	candidates := make(map[int]int)
	
	for i := 0; i < lsh.numHashTables; i++ {
		signature := lsh.computeSignature(shingles, lsh.hashFunctions[i])
		if ids, exists := lsh.hashTables[i][signature]; exists {
			for _, id := range ids {
				candidates[id]++
			}
		}
	}
	
	type result struct {
		id         int
		similarity float64
	}
	
	results := make([]result, 0)
	
	for id, count := range candidates {
		if float64(count)/float64(lsh.numHashTables) >= threshold {
			similarity := lsh.jaccardSimilarity(text, lsh.corpus[id])
			if similarity >= threshold {
				results = append(results, result{id: id, similarity: similarity})
			}
		}
	}
	
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})
	
	ids := make([]int, len(results))
	for i, r := range results {
		ids[i] = r.id
	}
	
	return ids
}

func (lsh *LSH) jaccardSimilarity(text1, text2 string) float64 {
	shingles1 := lsh.getShingles(text1)
	shingles2 := lsh.getShingles(text2)
	
	set1 := make(map[uint64]bool)
	set2 := make(map[uint64]bool)
	
	for _, s := range shingles1 {
		set1[s] = true
	}
	for _, s := range shingles2 {
		set2[s] = true
	}
	
	intersection := 0
	for s := range set1 {
		if set2[s] {
			intersection++
		}
	}
	
	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0
	}
	
	return float64(intersection) / float64(union)
}

type SimHash struct {
	hashBits int
	corpus   []string
	hashes   []uint64
}

func NewSimHash(hashBits int) *SimHash {
	return &SimHash{
		hashBits: hashBits,
		corpus:   make([]string, 0),
		hashes:   make([]uint64, 0),
	}
}

func (sh *SimHash) Add(text string) int {
	id := len(sh.corpus)
	sh.corpus = append(sh.corpus, text)
	sh.hashes = append(sh.hashes, sh.computeHash(text))
	return id
}

func (sh *SimHash) computeHash(text string) uint64 {
	tokens := tokenize(text)
	vector := make([]int, sh.hashBits)
	
	for _, token := range tokens {
		hash := xxhash.Sum64String(token)
		for i := 0; i < sh.hashBits; i++ {
			if (hash>>i)&1 == 1 {
				vector[i]++
			} else {
				vector[i]--
			}
		}
	}
	
	var simhash uint64
	for i := 0; i < sh.hashBits; i++ {
		if vector[i] > 0 {
			simhash |= 1 << i
		}
	}
	
	return simhash
}

func (sh *SimHash) Query(text string, maxHammingDistance int) []int {
	queryHash := sh.computeHash(text)
	results := make([]int, 0)
	
	for i, hash := range sh.hashes {
		distance := hammingDistance(queryHash, hash)
		if distance <= maxHammingDistance {
			results = append(results, i)
		}
	}
	
	return results
}

func hammingDistance(a, b uint64) int {
	xor := a ^ b
	count := 0
	for xor != 0 {
		count++
		xor &= xor - 1
	}
	return count
}

func tokenize(text string) []string {
	tokens := make([]string, 0)
	current := make([]rune, 0)
	
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current = append(current, unicode.ToLower(r))
		} else if len(current) > 0 {
			tokens = append(tokens, string(current))
			current = current[:0]
		}
	}
	
	if len(current) > 0 {
		tokens = append(tokens, string(current))
	}
	
	return tokens
}