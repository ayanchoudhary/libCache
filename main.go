package main

import (
	"container/list"
	"fmt"
	"strconv"

	"github.com/ayanchoudhary/libCache/policy"
)

func CustomPolicy() func(int) policy.Cache {
	return policy.Build(policy.Spec{
		Wrap: func(key string, value interface{}) policy.Node {
			return policy.NewNode(key, value)
		},
		Evict: func(ctx *policy.Ctx) {
			ctx.List.Remove(ctx.List.Front())
		},
		Found: func(ctx *policy.Ctx, item *list.Element) {
			ctx.List.MoveAfter(item, ctx.List.Front())
		},
	})
}

func main() {
	size := 128
	cache := CustomPolicy()(size)

	for i := 0; i < size; i++ {
		key := strconv.Itoa(i)
		cache.Add(key, i)
		value, _ := cache.Get(key)
		fmt.Println("Key:", key, "Value:", value)
	}
}
