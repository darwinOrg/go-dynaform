package dal

import (
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
)

var DynaSchemaDeploymentFields = struct {
	Id         string
	Name       string
	Properties string
	Web        string
	Version    string
	CreatedAt  string
	ModifiedAt string
}{
	"id",
	"name",
	"properties",
	"web",
	"version",
	"created_at",
	"modified_at",
}

var DynaSchemaDeploymentMeta = &daog.TableMeta[DynaSchemaDeployment]{
	Table: "dyna_schema_deployment",
	Columns: []string{
		"id",
		"name",
		"properties",
		"web",
		"version",
		"created_at",
		"modified_at",
	},
	AutoColumn: "id",
	LookupFieldFunc: func(columnName string, ins *DynaSchemaDeployment, point bool) any {
		if "id" == columnName {
			if point {
				return &ins.Id
			}
			return ins.Id
		}
		if "name" == columnName {
			if point {
				return &ins.Name
			}
			return ins.Name
		}
		if "properties" == columnName {
			if point {
				return &ins.Properties
			}
			return ins.Properties
		}
		if "web" == columnName {
			if point {
				return &ins.Web
			}
			return ins.Web
		}
		if "version" == columnName {
			if point {
				return &ins.Version
			}
			return ins.Version
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

var DynaSchemaDeploymentDao daog.QuickDao[DynaSchemaDeployment] = &struct {
	daog.QuickDao[DynaSchemaDeployment]
}{
	daog.NewBaseQuickDao(DynaSchemaDeploymentMeta),
}

type DynaSchemaDeployment struct {
	Id         int64
	Name       string
	Properties ttypes.NilableString
	Web        ttypes.NilableString
	Version    int32
	CreatedAt  ttypes.NormalDatetime
	ModifiedAt ttypes.NormalDatetime
}
