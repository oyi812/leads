package resources

import (
	"context"

	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/schema/query"
)

// Hook restricts access to the resource store
// by matching the keyed payload value in the
// context to the corresponding field value 
type Hook struct {
	name string
	key interface{}
}

// OnFind implements resource.FindEventHandler interface
func (h Hook) OnFind(ctx context.Context, q *query.Query, offset, limit int) error {

	value := ctx.Value(h.key)
	if value == nil {
		return resource.ErrForbidden
	}

	q.Predicate = append(q.Predicate, &query.Equal{
		Field: h.name,
		Value: value,
	})

	return nil
}

// OnGot implements resource.GotEventHandler interface
func (h Hook) OnGot(ctx context.Context, item **resource.Item, err *error) {

	if err != nil {
		return
	}

	value := ctx.Value(h.key)
	if value == nil {
		*err = resource.ErrForbidden
		return
	}

	id, ok := (*item).Payload[h.name]
	if !ok || id != value {
		*err = resource.ErrNotFound
	}

	return
}

// OnInsert implements resource.InsertEventHandler interface
func (h Hook) OnInsert(ctx context.Context, items []*resource.Item) error {

	value := ctx.Value(h.key)
	if value == nil {
		return resource.ErrForbidden
	}

	for _, item := range items {
		if id, ok := item.Payload[h.name]; ok {
			if id != value {
				return resource.ErrForbidden
			}
		} else {
			item.Payload[h.name] = value
		}
	}

	return nil
}

// OnUpdate implements resource.UpdateEventHandler interface
func (h Hook) OnUpdate(ctx context.Context, item *resource.Item, original *resource.Item) error {

	value := ctx.Value(h.key)
	if value == nil {
		return resource.ErrForbidden
	}

	id, ok := original.Payload[h.name]
	if !ok || id != value {
		return resource.ErrForbidden
	}

	// ensure field is not altered
	if id, ok := item.Payload[h.name]; !ok || id != value {
		return resource.ErrForbidden
	}

	return nil
}

// OnDelete implements resource.DeleteEventHandler interface
func (h Hook) OnDelete(ctx context.Context, item *resource.Item) error {

	value := ctx.Value(h.key)
	if value == nil {
		return resource.ErrForbidden
	}

	if item.Payload[h.name] != value {
		return resource.ErrForbidden
	}

	return nil
}

// OnClear implements resource.ClearEventHandler interface
func (h Hook) OnClear(ctx context.Context, q *query.Query) error {

	value := ctx.Value(h.key)
	if value == nil {
		return resource.ErrForbidden
	}

	// restrict impact of the clear to items for which authorised
	q.Predicate = append(q.Predicate, &query.Equal{
		Field: h.name,
		Value: value,
	})

	return nil
}
