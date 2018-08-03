package taggraph

import "sort"

//TagGrapher is an interface for graph like collection of tags
type TagGrapher interface {
	GetTag(name string) (Tagger, bool)
	SetTag(name string)
	Delete(name string)
	AddChildToTag(child, parent string)
	RemoveChildFromTag(child, parent string)
}

type tags []*tag

//NewTagGaph returns a TagGrapher
func NewTagGaph() TagGrapher {
	return &tags{}
}

func (t *tags) Len() int           { return len((*t)) }
func (t *tags) Swap(i, j int)      { (*t)[i], (*t)[j] = (*t)[j], (*t)[i] }
func (t *tags) Less(i, j int) bool { return (*t)[i].name < (*t)[j].name }

func (t *tags) searchFunc(name string) func(int) bool {
	return func(i int) bool {
		return (*t)[i].name >= name
	}
}

func (t *tags) get(name string) (*tag, int, bool) {
	i := sort.Search(t.Len(), t.searchFunc(name))
	if i < t.Len() && (*t)[i].name == name {
		return (*t)[i], i, true
	}
	return nil, 0, false
}

func (t *tags) GetTag(name string) (Tagger, bool) {
	tag, _, ok := t.get(name)
	return tag, ok
}

func (t *tags) SetTag(name string) {
	if _, _, ok := t.get(name); !ok {
		(*t) = append((*t), &tag{
			name:       name,
			childTags:  tags{},
			parentTags: tags{},
		})
	}
	sort.Sort(t)
}

func (t *tags) Delete(name string) {
	if node, i, ok := t.get(name); ok {
		(*t) = deleteIndexFromTags((*t), i)
		for _, v := range *t {
			v.parentTags = deleteTagFromTags(v.parentTags, node)
			sort.Sort(&v.parentTags)
			v.childTags = deleteTagFromTags(v.childTags, node)
			sort.Sort(&v.childTags)
		}
	}
	sort.Sort(t)
}

func (t *tags) AddChildToTag(child, parent string) {
	t.SetTag(child)
	t.SetTag(parent)
	childTag, _, _ := t.get(child)
	parentTag, _, _ := t.get(parent)

	var hasParent bool
	for _, p := range childTag.Parents() {
		if p == parentTag.Name() {
			hasParent = true
		}
	}
	if !hasParent {
		childTag.parentTags = append(childTag.parentTags, parentTag)
		sort.Sort(&childTag.parentTags)
	}

	var hasChild bool
	for _, c := range parentTag.Children() {
		if c == childTag.Name() {
			hasChild = true
		}
	}
	if !hasChild {
		parentTag.childTags = append(parentTag.childTags, childTag)
		sort.Sort(&parentTag.childTags)
	}
}

func (t *tags) RemoveChildFromTag(child, parent string) {
	childTag, _, childOK := t.get(child)
	parentTag, _, parentOK := t.get(parent)
	if childOK && parentOK {
		childTag.parentTags = deleteTagFromTags(childTag.parentTags, parentTag)
		sort.Sort(&childTag.parentTags)
		parentTag.childTags = deleteTagFromTags(parentTag.childTags, childTag)
		sort.Sort(&parentTag.childTags)
	}
}

func deleteIndexFromTags(slice tags, i int) []*tag {
	slice[i] = slice[slice.Len()-1]
	slice[slice.Len()-1] = nil
	return slice[:slice.Len()-1]
}

func deleteTagFromTags(slice tags, t *tag) []*tag {
	for i, v := range slice {
		if v == t {
			slice = deleteIndexFromTags(slice, i)
			break
		}
	}
	return slice
}
