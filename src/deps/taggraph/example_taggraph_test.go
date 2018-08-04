package taggraph_test

import (
	"fmt"
)

func Example_tags() {
	graph := NewTagGaph()

	graph.AddChildToTag("shirts", "clothes")
	graph.AddChildToTag("pants", "clothes")
	graph.AddChildToTag("dress clothes", "clothes")
	graph.AddChildToTag("shirts", "dress clothes")
	graph.AddChildToTag("shirts", "tops")
	graph.AddChildToTag("tops", "casual")
	graph.AddChildToTag("casual", "clothes")

	shirts, ok := graph.GetTag("shirts")
	if !ok {
		panic("Shirts tag not found")
	}
	paths := shirts.PathsToAllAncestorsAsString("->")

	for _, path := range paths {
		fmt.Println(path)
	}
	// Output: clothes->shirts
	// clothes->dress clothes->shirts
	// clothes->casual->tops->shirts
}
