package policy

import (
	"container/list"
	"math"
)

type NodeCounted struct {
	key       string
	value     interface{}
	frequency int
}

func (node *NodeCounted) Key() string {
	return node.key
}

func (node *NodeCounted) Value() interface{} {
	return node.value
}

func (node *NodeCounted) Frequency() int {
	return node.frequency
}

func LFU() func(int) Cache {
	return Build(Spec{
		Wrap: func(key string, value interface{}) Node {
			return &NodeCounted{key: key, value: value, frequency: 0}
		},
		Evict: func(ctx *Ctx) {
			currentFrequency := math.MaxInt
			var element *list.Element = nil
			var currentUntyped = ctx.List.Front()
			for currentUntyped != nil {
				current := currentUntyped.Value.(NodeCounted)
				if current.frequency <= currentFrequency {
					element = currentUntyped
					currentFrequency = current.frequency
				}
				currentUntyped = currentUntyped.Next()
			}

			ctx.List.Remove(element)
		},
		Found: func(ctx *Ctx, item *list.Element) {
			node := item.Value.(*NodeCounted)
			node.frequency++
		},
	})
}
