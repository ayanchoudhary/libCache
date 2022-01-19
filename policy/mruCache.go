package policy

import "container/list"

func MRU() func(int) Cache {
	return Build(Spec{
		Wrap: func(key string, value interface{}) Node {
			return &NodeImpl{key: key, value: value}
		},
		Evict: func(ctx *Ctx) {
			ctx.List.Remove(ctx.List.Front())
		},
		Found: func(ctx *Ctx, item *list.Element) {
			ctx.List.MoveBefore(item, ctx.List.Front())
		},
	})
}
