package policy

import "container/list"

func FIFO() func(int) Cache {
	return Build(Spec{
		Wrap: func(key string, value interface{}) Node {
			return &NodeImpl{key, value}
		},
		Evict: func(ctx *Ctx) {
			ctx.List.Remove(ctx.List.Back())
		},
		Found: func(ctx *Ctx, item *list.Element) {

		},
	})
}
