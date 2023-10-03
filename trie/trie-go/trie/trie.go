package trie

import (
	"errors"
	"strings"
)

type TrieNode struct {
	StaticChildren map[string]*TrieNode
	ParameterChild *TrieNode
	Handler        func() bool
	ParameterName  *string
}

type Trie struct {
	Root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		Root: &TrieNode{
			StaticChildren: make(map[string]*TrieNode),
		},
	}
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		StaticChildren: make(map[string]*TrieNode),
	}
}

func (t *Trie) Insert(path string, handler func() bool) error {
	var current *TrieNode = t.Root
	var segments []string = splitPath(path)
	var seenParameters map[string]bool = make(map[string]bool)

	// _: the index of the current segment
	// segment: the current segment
	for _, segment := range segments {
		if segment[0] == ':' {
			// segment[1:] will remove the ':' from the parameter name
			// it's essentially a substring from index 1 to the end of the string
			paramName := segment[1:]

			// Check if the parameter name has already been used
			if seenParameters[paramName] {
				return errors.New("cannot use the same parameter name more than once in a route")
			}
			seenParameters[paramName] = true

			if current.ParameterChild != nil {
				// If the parameter child is not nil, then we need to check if the parameter name is the same
				// *current will access the value of the pointer
				if *current.ParameterChild.ParameterName != paramName {
					return errors.New("Invalid path: " + path + ". Parameter " + paramName + " is already defined.")
				}
			} else {
				// If the parameter child is nil, then we need to create a new node
				// &TrieNode{} will create a new TrieNode and return a pointer to it
				//paramNode := &TrieNode{}
				paramNode := newTrieNode()
				paramNode.ParameterName = &paramName // &paramName will return a pointer to the paramName variable
				current.ParameterChild = paramNode
			}

			// Set the current node to the parameter child
			current = current.ParameterChild
		} else {
			// If the segment is not a parameter, then we need to check if the segment already exists

			// Maps in Go can return two values: the value and a boolean indicating whether or not the key exists
			// when the key does not exist, the value will be the zero value of the type (zero types are for example 0 for int, "" for string, nil for pointers, etc.)
			// here we ignore the value and only check if the key exists
			// the `;` separates the two statements, so we can use the value in the next statement
			if _, exists := current.StaticChildren[segment]; !exists {
				// Create a new node and set it as the static child
				current.StaticChildren[segment] = newTrieNode()
			}

			// Set the current node to the static child
			current = current.StaticChildren[segment]
		}
	}

	// Set the handler for the current node
	// This is done outside the loop because the current node is the last node in the path
	current.Handler = handler

	return nil
}

type SearchResult struct {
	Handler func() bool
	Params  map[string]string
}

func (t *Trie) Search(path string) SearchResult {
	current := t.Root
	segments := splitPath(path)
	capturedParams := make(map[string]string)

	for _, segment := range segments {
		// next: the next node in the trie, or the child of the current node
		// exists: whether or not the next node exists
		next, exists := current.StaticChildren[segment]

		// If the next node does not exist, then we need to check if the current node has a parameter child
		if !exists && current.ParameterChild != nil {
			next = current.ParameterChild
			capturedParams[*next.ParameterName] = segment
		}

		// If the next node does not exist and the current node does not have a parameter child, then the path does not exist
		if next == nil {
			return SearchResult{Handler: func() bool { return false }, Params: nil}
		}

		// Set the current node to the next node and continue the loop with the next segment
		current = next
	}

	return SearchResult{Handler: current.Handler, Params: capturedParams}
}

func splitPath(path string) []string {
	// The strings.Split function will return an empty string as the first element
	// if the path starts with '/', so we can simply slice it off.
	return strings.Split(path, "/")[1:]
}
