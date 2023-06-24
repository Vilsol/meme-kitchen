package schema

import (
	"database/sql"
	"database/sql/driver"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"fmt"
	"math/big"
	"memekitchen/data"
)

type BigInt big.Int

func (b *BigInt) Scan(src any) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return err
	}
	if !i.Valid {
		return nil
	}
	if _, ok := (*big.Int)(b).SetString(i.String, 10); ok {
		return nil
	}
	return fmt.Errorf("could not scan type %T with value %v into BigInt", src, src)
}

func (b *BigInt) Value() (driver.Value, error) {
	return (*big.Int)(b).String(), nil
}

type Template struct {
	ent.Schema
}

func (Template) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.JSON("data", []data.TemplateText{}),

		field.Int("avg_distance").
			GoType(new(BigInt)).
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(78, 0)",
			}),
		field.Int("diff_distance").
			GoType(new(BigInt)).
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(78, 0)",
			}),
		field.Int("perception_distance").
			GoType(new(BigInt)).
			SchemaType(map[string]string{
				dialect.Postgres: "numeric(78, 0)",
			}),
	}
}
