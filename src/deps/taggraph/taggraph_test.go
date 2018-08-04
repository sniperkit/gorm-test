package taggraph

import (
	"testing"
)

func TestAddTag(t *testing.T) {
	graph := NewTagGaph()
	tests := []string{"foo", "bar", "baz", "quxx"}

	for _, test := range tests {
		graph.SetTag(test)
		node, ok := graph.GetTag(test)
		if !ok {
			t.Errorf("Node %v not found", test)
			continue
		}
		if node.Name() != test {
			t.Errorf("Wanted %v but got %v", node.Name(), test)
		}
	}
}

func TestDeleteNode(t *testing.T) {
	graph := &tags{}
	tests := []string{"foo", "bar", "baz", "quxx"}
	control := "buxx"

	graph.SetTag(control)
	for _, test := range tests {
		graph.SetTag(test)
	}
	for _, test := range tests {
		graph.Delete(test)
		if _, ok := graph.GetTag(control); !ok {
			t.Errorf("Control tag %v was not found", control)
		}
		if node, ok := graph.GetTag(test); ok {
			t.Errorf(
				"Found node %v when it should have been deleted",
				node.Name(),
			)
		}
	}
}

func TestEdgeMethods(t *testing.T) {
	graph := &tags{}
	tests := [][2]string{
		{"foo", "bar"},
		{"baz", "quxx"},
	}

	for _, test := range tests {
		graph.AddChildToTag(test[0], test[1])
		child, _, ok := graph.get(test[0])
		if !ok {
			t.Errorf("Node %v not found", test)
			continue
		}
		if child.Name() != test[0] {
			t.Errorf("Wanted %v but got %v", child.Name(), test[0])
		}
		parent, _, ok := graph.get(test[1])
		if !ok {
			t.Errorf("Node %v not found", test)
			continue
		}
		if parent.Name() != test[1] {
			t.Errorf("Wanted %v but got %v", parent.Name(), test[1])
		}
		if tag, ok := parent.childTags.GetTag(test[0]); !ok || tag != child {
			t.Errorf("Parent tag %v has no child %v", parent.Name(), test[0])
		}
		if tag, ok := child.parentTags.GetTag(test[1]); !ok || tag != parent {
			t.Errorf("Child tag %v has no parent %v", child.Name(), test[1])
		}
		graph.RemoveChildFromTag(test[0], test[1])

		if _, ok := child.parentTags.GetTag(test[1]); ok {
			t.Errorf("Child tag %v still has parent %v", child.Name(), test[1])
		}
		if _, ok := parent.childTags.GetTag(test[0]); ok {
			t.Errorf("Parent tag %v  still has child %v", parent.Name(), test[0])
		}
	}
}
