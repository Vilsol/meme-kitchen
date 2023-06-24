// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"memekitchen/data"
	"memekitchen/ent/schema"
	"memekitchen/ent/template"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TemplateCreate is the builder for creating a Template entity.
type TemplateCreate struct {
	config
	mutation *TemplateMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (tc *TemplateCreate) SetName(s string) *TemplateCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetData sets the "data" field.
func (tc *TemplateCreate) SetData(dt []data.TemplateText) *TemplateCreate {
	tc.mutation.SetData(dt)
	return tc
}

// SetAvgDistance sets the "avg_distance" field.
func (tc *TemplateCreate) SetAvgDistance(si *schema.BigInt) *TemplateCreate {
	tc.mutation.SetAvgDistance(si)
	return tc
}

// SetDiffDistance sets the "diff_distance" field.
func (tc *TemplateCreate) SetDiffDistance(si *schema.BigInt) *TemplateCreate {
	tc.mutation.SetDiffDistance(si)
	return tc
}

// SetPerceptionDistance sets the "perception_distance" field.
func (tc *TemplateCreate) SetPerceptionDistance(si *schema.BigInt) *TemplateCreate {
	tc.mutation.SetPerceptionDistance(si)
	return tc
}

// Mutation returns the TemplateMutation object of the builder.
func (tc *TemplateCreate) Mutation() *TemplateMutation {
	return tc.mutation
}

// Save creates the Template in the database.
func (tc *TemplateCreate) Save(ctx context.Context) (*Template, error) {
	return withHooks(ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TemplateCreate) SaveX(ctx context.Context) *Template {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TemplateCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TemplateCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tc *TemplateCreate) check() error {
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Template.name"`)}
	}
	if v, ok := tc.mutation.Name(); ok {
		if err := template.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Template.name": %w`, err)}
		}
	}
	if _, ok := tc.mutation.Data(); !ok {
		return &ValidationError{Name: "data", err: errors.New(`ent: missing required field "Template.data"`)}
	}
	if _, ok := tc.mutation.AvgDistance(); !ok {
		return &ValidationError{Name: "avg_distance", err: errors.New(`ent: missing required field "Template.avg_distance"`)}
	}
	if _, ok := tc.mutation.DiffDistance(); !ok {
		return &ValidationError{Name: "diff_distance", err: errors.New(`ent: missing required field "Template.diff_distance"`)}
	}
	if _, ok := tc.mutation.PerceptionDistance(); !ok {
		return &ValidationError{Name: "perception_distance", err: errors.New(`ent: missing required field "Template.perception_distance"`)}
	}
	return nil
}

func (tc *TemplateCreate) sqlSave(ctx context.Context) (*Template, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TemplateCreate) createSpec() (*Template, *sqlgraph.CreateSpec) {
	var (
		_node = &Template{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(template.Table, sqlgraph.NewFieldSpec(template.FieldID, field.TypeInt))
	)
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(template.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := tc.mutation.Data(); ok {
		_spec.SetField(template.FieldData, field.TypeJSON, value)
		_node.Data = value
	}
	if value, ok := tc.mutation.AvgDistance(); ok {
		_spec.SetField(template.FieldAvgDistance, field.TypeInt, value)
		_node.AvgDistance = value
	}
	if value, ok := tc.mutation.DiffDistance(); ok {
		_spec.SetField(template.FieldDiffDistance, field.TypeInt, value)
		_node.DiffDistance = value
	}
	if value, ok := tc.mutation.PerceptionDistance(); ok {
		_spec.SetField(template.FieldPerceptionDistance, field.TypeInt, value)
		_node.PerceptionDistance = value
	}
	return _node, _spec
}

// TemplateCreateBulk is the builder for creating many Template entities in bulk.
type TemplateCreateBulk struct {
	config
	builders []*TemplateCreate
}

// Save creates the Template entities in the database.
func (tcb *TemplateCreateBulk) Save(ctx context.Context) ([]*Template, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Template, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TemplateMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TemplateCreateBulk) SaveX(ctx context.Context) []*Template {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TemplateCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TemplateCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}
