package main

import (
	"fmt"
	"strings"
	"trie-go/trie"
)

func main() {
	trie := trie.NewTrie()

	trie.Insert("/users", func() bool { return true })
	trie.Insert("/users/:userId", func() bool { return true })

	search := trie.Search("/users/1")
	println(search.Handler())
	println(search.Params["userId"])

	// invalid path
	search = trie.Search("/users/1/approve")
	println(search.Handler())

	// invalid insert
	err := trie.Insert("/users/:id/approve", func() bool { return true })
	if err != nil {
		println(err.Error())
	}

	for _, v := range strings.Split("/user/a/b/c", "/")[1:] {
		fmt.Printf("-> [%s]\n", v)
	}

}
