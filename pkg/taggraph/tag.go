package taggraph

import "strings"

//Tagger is an interface for interacting with a node of a TagGrapher
type Tagger interface {
	Children() []string
	Parents() []string
	PathsToAllAncestors() [][]string
	PathsToAllAncestorsAsString(delim string) []string
	Name() string
}

type tag struct {
	name       string
	childTags  tags
	parentTags tags
}

func (t *tag) Name() string {
	return t.name
}

func (t *tag) Children() []string {
	return flatTags(t.childTags)
}

func (t *tag) Parents() []string {
	return flatTags(t.parentTags)
}

func flatTags(t tags) []string {
	arr := []string{}
	for _, v := range t {
		arr = append(arr, v.name)
	}
	return arr
}

func (t *tag) PathsToAllAncestors() [][]string {
	return ancestors(t, [][]string{{t.name}})
}

func (t *tag) PathsToAllAncestorsAsString(delim string) []string {
	paths := []string{}
	for _, v := range t.PathsToAllAncestors() {
		paths = append(paths, strings.Join(v, delim))
	}
	return paths
}

func ancestors(t *tag, paths [][]string) [][]string {
	newPaths := [][]string{}
	for _, path := range paths {
		for _, v := range t.parentTags {
			pathCopy := cloneStringSlice(path)
			if pathCopy[0] == t.name {
				pathCopy = prependStringSlice(pathCopy, v.name)
			}
			newPaths = append(newPaths, pathCopy)
			if v.parentTags.Len() > 0 {
				newPaths = ancestors(v, newPaths)
			}
		}
	}
	return newPaths
}

func cloneStringSlice(slice []string) []string {
	return append([]string(nil), slice...)
}

func prependStringSlice(slice []string, val string) []string {
	return append(slice[:0], append([]string{val}, slice[0:]...)...)
}
