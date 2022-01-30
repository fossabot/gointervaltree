// Package gointervaltree provides functionality for indexing a set of integer intervals, e.g. [start, end)
// based on http://en.wikipedia.org/wiki/Interval_tree. Copyright 2022, Kirill Danilov. Licensed under MIT license.
package gointervaltree

import (
	"log"
	"reflect"
	"sort"
)

// IntervalTree struct defines data structure for indexing a set of integer intervals, e.g. [start, end).
type IntervalTree struct {
	min              int
	max              int
	center           int
	singleInterval   []interface{}
	leftSubtree      *IntervalTree
	rightSubtree     *IntervalTree
	midSortedByStart []interface{}
	midSortedByEnd   []interface{}
}

// NewIntervalTree method instantiates an instance of IntervalTree struct creating a node for keeping intervals.
func NewIntervalTree(min int, max int) (tree *IntervalTree) {
	tree = new(IntervalTree)
	tree.min = min
	tree.max = max
	if !(tree.min < tree.max) {
		log.Panic("AssertionError: interval tree start must be numerically less than its end")
	}
	tree.center = (min + max) / 2
	tree.singleInterval = nil
	tree.leftSubtree = nil
	tree.rightSubtree = nil
	tree.midSortedByStart = []interface{}{}
	tree.midSortedByEnd = []interface{}{}
	return tree
}

// addInterval method adds intervals to the tree without sorting them along the way.
func (tree *IntervalTree) addInterval(start int, end int, data interface{}) {
	if (end - start) <= 0 {
		return
	}
	if tree.singleInterval == nil {
		tree.singleInterval = []interface{}{start, end, data}
	} else if reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		tree.addIntervalMain(start, end, data)
	} else {
		tree.addIntervalMain(tree.singleInterval[0].(int), tree.singleInterval[1].(int), tree.singleInterval[2])
		tree.singleInterval = []interface{}{0}
		tree.addIntervalMain(start, end, data)
	}
}

// addIntervalMain method is a technical method used inside addInterval.
func (tree *IntervalTree) addIntervalMain(start int, end int, data interface{}) {

	if end <= tree.center {
		if tree.leftSubtree == nil {
			tree.leftSubtree = NewIntervalTree(tree.min, tree.center)
		}
		tree.leftSubtree.addInterval(start, end, data)
	} else if start > tree.center {
		if tree.rightSubtree == nil {
			tree.rightSubtree = NewIntervalTree(tree.center, tree.max)
		}
		tree.rightSubtree.addInterval(start, end, data)
	} else {
		tree.midSortedByStart = append(tree.midSortedByStart, []interface{}{start, end, data})
		tree.midSortedByEnd = append(tree.midSortedByEnd, []interface{}{start, end, data})
	}
}

// sort method is used to sort intervals within the tree and must be invoked after adding intervals.
func (tree *IntervalTree) sort() {
	if tree.singleInterval == nil || !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		return
	}

	sort.Slice(tree.midSortedByStart, func(i, j int) bool {
		return tree.midSortedByStart[i].([3]interface{})[0].(int) < tree.midSortedByStart[j].([3]interface{})[0].(int)
	})

	sort.Slice(tree.midSortedByEnd, func(i, j int) bool {
		return tree.midSortedByEnd[i].([3]interface{})[1].(int) > tree.midSortedByEnd[j].([3]interface{})[1].(int)
	})
}

// query method returns all intervals in the tree which overlap given point,
// i.e. all (start, end, data) records, for which (start <= x < end).
func (tree *IntervalTree) query(x int) []interface{} {
	var result []interface{}
	return tree.queryMain(x, result)
}

// queryMain method is a technical method used inside query.
func (tree *IntervalTree) queryMain(x int, result []interface{}) []interface{} {
	if tree.singleInterval == nil {
		return result
	} else if !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		if tree.singleInterval[0].(int) <= x && x < tree.singleInterval[1].(int) {
			result = append(result, tree.singleInterval)
		}
		return result
	} else if x < tree.center {
		if tree.leftSubtree != nil {
			result = append(result, tree.leftSubtree.queryMain(x, result)...)
		}
		for _, element := range tree.midSortedByStart {
			if element.([]interface{})[0].(int) <= x {
				result = append(result, element)
			} else {
				break
			}
		}
		return result
	} else {
		for _, element := range tree.midSortedByEnd {
			if element.([]interface{})[1].(int) > x {
				result = append(result, element)
			} else {
				break
			}
		}
		if tree.rightSubtree != nil {
			result = append(result, tree.rightSubtree.queryMain(x, result)...)

		}
		return result
	}
}

// len method represents the number of intervals maintained in the tree, zero- or negative-size intervals
// are not registered.
func (tree *IntervalTree) len() int {
	if tree.singleInterval == nil {
		return 0
	} else if !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		return 1
	} else {
		size := len(tree.midSortedByStart)
		if tree.leftSubtree != nil {
			size += tree.leftSubtree.len()
		}
		if tree.rightSubtree != nil {
			size += tree.rightSubtree.len()
		}
		return size
	}
}

// iter method returns a slice of all intervals maintained in the tree.
func (tree *IntervalTree) iter() []interface{} {
	var result []interface{}
	if tree.singleInterval == nil {
		return result
	} else if !reflect.DeepEqual(tree.singleInterval, []interface{}{0}) {
		result = append(result, tree.singleInterval)
		return result
	} else {
		if tree.leftSubtree != nil {
			result = append(result, tree.leftSubtree.iter()...)
		}
		if tree.rightSubtree != nil {
			result = append(result, tree.rightSubtree.iter()...)
		}
		for _, element := range tree.midSortedByStart {
			result = append(result, element)
		}
		return result
	}
}