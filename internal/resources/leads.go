package resources

import (
	"context"

	"github.com/rs/rest-layer/schema"
)

type Leads struct {
	GroupFieldname string
	ContextKey interface{}
}

func (l Leads) Hook() Hook {
	return Hook{
		name: l.GroupFieldname,
		key: l.ContextKey,
	}
}

func (l Leads) Schema() schema.Schema {

	stringField := schema.String{
		MinLen: 2,
		MaxLen: 50,
	}
	required := schema.Field{
		Required:   true,
		Filterable: true,
		Validator: &stringField,
	}
	optional := schema.Field{
		Filterable: true,
		Validator: &stringField,
	}

	onInit := func(k interface{}) func(
		context.Context,
		interface{},
	) interface{} {
		return func(
			ctx context.Context,
			value interface{},
		) interface{} {
			if value := ctx.Value(k); value != nil {
				return value
			}
			return value
		}
	}

	fields := schema.Fields{
		"id": schema.IDField,
		l.GroupFieldname: {
			Filterable: true,
			ReadOnly:   true,
			Hidden:     true,
			OnInit:     onInit(l.ContextKey),
		},
		"FirstName": required,
		"LastName": required,
		"Email": required,
		"Company": optional,
		"PostCode": optional,
		"AcceptTerms": {
			Required:   true,
			Filterable: true,
			Validator: &schema.Bool{},
		},
		"DateCreated": schema.CreatedField,
	}

	return schema.Schema{Fields: fields}
}
