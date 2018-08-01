package main

import (
	"fmt"
	"strings"

	"github.com/k0kubun/pp"
)

func Pretty(msg interface{}) {
	pp.Println(msg)
}

func RemoveSliceDuplicates(elements []string, noEmpty bool) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}
	for v := range elements {
		if noEmpty {
			if len(strings.TrimSpace(elements[v])) == 0 {
				continue
			}
		}
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func AddTopicsPrefix(topics []string, prefix string, noEmpty bool) []string {
	topics = RemoveSliceDuplicates(topics, true)
	topicsPrefixed := make([]string, 0)
	for _, topic := range topics {
		name := fmt.Sprintf("%s%s", prefix, topic)
		if len(strings.TrimSpace(name)) == 0 {
			continue
		}
		topicsPrefixed = append(topicsPrefixed, fmt.Sprintf("%s%s", prefix, topic))
	}
	return topicsPrefixed
}

func Topics2Tags(topics []string, prefix string, noEmpty bool) []Tag {
	topics = RemoveSliceDuplicates(topics, true)
	tags := make([]Tag, 0)
	for _, topic := range topics {
		name := fmt.Sprintf("%s%s", prefix, topic)
		if len(strings.TrimSpace(name)) == 0 {
			continue
		}
		tags = append(tags, Tag{
			Name: fmt.Sprintf("%s%s", prefix, topic),
		})
	}
	return tags
}
