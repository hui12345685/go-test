package arraytest

import (
	"fmt"
	"sort"
)

// 课程表
// prereqs记录了每个课程的前置课程
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func TestTopSort() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// 拓扑排序
func topoSort(m map[string][]string) []string {
	fmt.Println("\033[1;32;40m  \n start topoSort--------------------------------------- \033[0m")
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			//不存在item，默认就是false
			if !seen[item] {
				//fmt.Printf("!seen[item]:item:%s\n", item)
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
				//fmt.Printf("append item:%s\n", item)
			}
		}
	}
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order
}
