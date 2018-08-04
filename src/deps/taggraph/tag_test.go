package taggraph

import (
	"reflect"
	"testing"
)

func TestFooMethods(t *testing.T) {
	var (
		graph      = NewTagGaph()
		wantArrays = [][]string{
			{"clothes", "shirts"},
			{"clothes", "dress clothes", "shirts"},
			{"clothes", "casual", "tops", "shirts"},
		}
		wantStrings = []string{
			"clothes->shirts",
			"clothes->dress clothes->shirts",
			"clothes->casual->tops->shirts",
		}
	)
	graph.AddChildToTag("shirts", "clothes")
	graph.AddChildToTag("pants", "clothes")
	graph.AddChildToTag("dress clothes", "clothes")
	graph.AddChildToTag("shirts", "dress clothes")
	graph.AddChildToTag("shirts", "tops")
	graph.AddChildToTag("tops", "casual")
	graph.AddChildToTag("casual", "clothes")

	shirt, _ := graph.GetTag("shirts")
	if !reflect.DeepEqual(shirt.Parents(), []string{"clothes", "dress clothes", "tops"}) {
		t.Error("Missing parents")
	}

	tops, _ := graph.GetTag("tops")
	if !reflect.DeepEqual(tops.Children(), []string{"shirts"}) {
		t.Error("Missing children")
	}

	if !reflect.DeepEqual(shirt.PathsToAllAncestors(), wantArrays) {
		t.Error("Missing Ancestors paths")
	}

	if !reflect.DeepEqual(shirt.PathsToAllAncestorsAsString("->"), wantStrings) {
		t.Error("Missing Ancestors paths")
	}
}
