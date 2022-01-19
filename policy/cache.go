package policy

import (
	"container/list"
	"fmt"
)

type Cache struct {
	IsEmpty func() bool
	Clear   func()
	Get     func(key string) (interface{}, error)
	Add     func(key string, value interface{})
}

type Ctx struct {
	List  *list.List
	Limit int
}

type Environment struct {
	*Ctx
	*Spec
}

func (ctx Ctx) clear() {
	ctx.List = list.New()
}

func (ctx Ctx) isEmpty() bool {
	return ctx.List.Len() == 0
}

type Node interface {
	Key() string
	Value() interface{}
}
type NodeImpl struct {
	key   string
	value interface{}
}

func NewNode(key string, value interface{}) *NodeImpl {
	return &NodeImpl{key, value}
}

func (node *NodeImpl) Key() string {
	return node.key
}

func (node *NodeImpl) Value() interface{} {
	return node.value
}

type Spec struct {
	Wrap  func(key string, value interface{}) Node
	Evict func(ctx *Ctx)
	Found func(ctx *Ctx, item *list.Element)
}

func (env *Environment) get(key string) (interface{}, error) {
	if env.List.Len() == 0 {
		return nil, fmt.Errorf("cache miss '%s': cache is empty", key)
	}

	var currentUntyped = env.List.Front()
	for currentUntyped != nil {
		current := currentUntyped.Value.(Node)
		if current.Key() == key {
			value := current.Value()
			env.Spec.Found(env.Ctx, currentUntyped)
			return value, nil
		}

		currentUntyped = currentUntyped.Next()
	}

	return nil, fmt.Errorf("cache miss '%s': unable to find a match", key)
}

func (env *Environment) add(key string, value interface{}) {
	var currentUntyped = env.List.Front()
	for currentUntyped != nil {
		current := currentUntyped.Value.(Node)
		if current.Key() == key {
			// Nodes are immutable, just remove the node and add it back
			env.List.Remove(currentUntyped)
		}

		currentUntyped = currentUntyped.Next()
	}

	if env.List.Len()+1 > env.Limit {
		env.Spec.Evict(env.Ctx)
	}

	env.List.PushFront(env.Spec.Wrap(key, value))
}

func Build(spec Spec) func(int) Cache {
	return func(limit int) Cache {
		ctx := Ctx{Limit: limit, List: list.New()}
		env := Environment{Ctx: &ctx, Spec: &spec}

		return Cache{
			IsEmpty: func() bool {
				return env.isEmpty()
			},
			Clear: func() {
				env.clear()
			},
			Get: func(key string) (interface{}, error) {
				return env.get(key)
			},
			Add: func(key string, value interface{}) {
				env.add(key, value)
			},
		}
	}
}
