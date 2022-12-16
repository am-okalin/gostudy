package trie

type Trie struct {
	children map[string]*Trie
	isEnd    bool
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Insert(word string) {

}

func (t *Trie) Search(word string) bool {

}

func (t *Trie) StartsWith(prefix string) bool {

}
