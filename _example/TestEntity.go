package _example

import (
	"github.com/darwinOrg/go-dynaform"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var TestEntityFields = struct {
	Id         string
	Foo        string
	Bar        string
	Data       string
	CreatedAt  string
	ModifiedAt string
}{
	"id",
	"foo",
	"bar",
	"data",
	"created_at",
	"modified_at",
}

var TestEntityMeta = &daog.TableMeta[TestEntity]{
	Table: "test_entity",
	Columns: []string{
		"id",
		"foo",
		"bar",
		"data",
		"created_at",
		"modified_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *TestEntity, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "foo" == columnName {
			if point {
				return &ins.Foo
			}
			return ins.Foo
		}
		if "bar" == columnName {
			if point {
				return &ins.Bar
			}
			return ins.Bar
		}
		if "data" == columnName {
			if point {
				return &ins.Data
			}
			return ins.Data
		}
		if "created_at" == columnName {
			if point {
				return &ins.CreatedAt
			}
			return ins.CreatedAt
		}
		if "modified_at" == columnName {
			if point {
				return &ins.ModifiedAt
			}
			return ins.ModifiedAt
		}

		return nil
	},
}

var TestEntityDao daog.QuickDao[TestEntity] = &struct {
	daog.QuickDao[TestEntity]
}{
	daog.NewBaseQuickDao(TestEntityMeta),
}

type TestEntity struct {
	Id         int64
	Foo        string
	Bar        string
	Data       *dynaform.Dynamic
	CreatedAt  ttypes.NormalDatetime
	ModifiedAt ttypes.NormalDatetime
}
