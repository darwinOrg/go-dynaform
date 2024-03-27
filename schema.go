package dynaform

import (
	"encoding/json"
	"errors"
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-dynaform/dal"
	dglogger "github.com/darwinOrg/go-logger"
	ve "github.com/darwinOrg/go-validator-ext"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"reflect"
	"strings"
	"sync"
)

const FillDataSourceKey = "fill_data_source"

type DynamicWeb map[string]any

var (
	structCache = sync.Map{}
	schemaCache = sync.Map{}
	webCache    = sync.Map{}
)

type DynamicSchema struct {
	SchemaId   int64             `json:"schemaId"`
	Name       string            `json:"name"`
	Properties []DynamicProperty `json:"properties"`
	Web        string            `json:"web"`
}

type DynamicSchemaData struct {
	FillDataSource string            `json:"fillDataSource,omitempty"`
	Content        []DynamicProperty `json:"content"`
}

func (schema *DynamicSchema) reflect2StructData(ctx *dgctx.DgContext, dataMap map[string]any) (any, error) {
	nm := map[string]any{}
	for k, v := range dataMap {
		nm[FirstCharToUpper(k)] = v
	}
	typ := schema.reflect2StructType()
	dglogger.Debugf(ctx, "struct: %+v\n", typ)
	typStruct := reflect.New(typ).Elem().Addr().Interface()
	err := mapstructure.Decode(nm, &typStruct)
	if err != nil {
		dglogger.Errorln(ctx, err)
	}
	dglogger.Debugf(ctx, "struct data: %+v\n", typStruct)
	return typStruct, nil
}

func (schema *DynamicSchema) reflect2StructType() reflect.Type {
	st, _ := structCache.Load(schema.SchemaId)
	if st != nil {
		return st.(reflect.Type)
	}

	var structFields = make([]reflect.StructField, 0, len(schema.Properties))
	for _, property := range schema.Properties {
		structFields = append(structFields, reflect.StructField{
			Name: FirstCharToUpper(property.Name),
			Type: property.ReflectType(),
			Tag:  property.StructTag(),
		})
	}

	typ := reflect.StructOf(structFields)
	structCache.Store(schema.SchemaId, typ)

	return typ
}

type DynamicProperty struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	AutoFill  bool              `json:"autoFill"`
	Validates []DynamicValidate `json:"validates"`
}

func (dp DynamicProperty) ReflectType() reflect.Type {
	var tpe reflect.Type

	switch dp.Type {
	case "string":
		tpe = reflect.TypeOf("")
	case "strings":
		tpe = reflect.TypeOf([]string{})
	case "int":
		tpe = reflect.TypeOf(0)
	case "ints":
		tpe = reflect.TypeOf([]int{})
	case "long":
		tpe = reflect.TypeOf(int64(0))
	case "longs":
		tpe = reflect.TypeOf([]int64{})
	case "double":
		tpe = reflect.TypeOf(float64(0))
	case "bigDecimal":
		tpe = reflect.TypeOf(&decimal.Zero)
	case "identity":
		tpe = reflect.TypeOf(PersonnelIdentity{})
	case "objects":
		tpe = reflect.TypeOf([]any{})
	default:
		panic(fmt.Sprintf("invalid type: %s", dp.Type))
	}

	return tpe
}

func (dp DynamicProperty) StructTag() reflect.StructTag {
	if len(dp.Validates) == 0 {
		return reflect.StructTag("json:\"" + dp.Name + "\"")
	}

	var sts = make([]string, len(dp.Validates))
	for i, dw := range dp.Validates {
		sts[i] = dw.StructTag()
	}

	return reflect.StructTag("json:\"" + dp.Name + "\" binding:\"" + strings.Join(sts, ",") + "\"")
}

type DynamicValidate struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func (dw DynamicValidate) StructTag() string {
	switch dw.Name {
	case ve.NOT_NULL:
		return "required"
	case ve.IS_EMAIL:
		return "email"
	case ve.MAX_VALUE:
		return fmt.Sprintf("max=%s", fmt.Sprint(dw.Value))
	case ve.MIN_VALUE:
		return fmt.Sprintf("min=%s", fmt.Sprint(dw.Value))
	case ve.MUST_IN:
		return fmt.Sprintf("%s=%s", ve.MUST_IN, strings.ReplaceAll(dw.Value.(string), "|", ve.MustInSeq))
	}

	if dw.Value != nil && dw.Value != true && dw.Value != "true" {
		return dw.Name + "=" + fmt.Sprint(dw.Value)
	} else {
		return dw.Name
	}
}

func Deploy(ctx *dgctx.DgContext, name string, properties string, web string) (*DynamicSchema, error) {
	dsd, err := dal.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	pp, err := readProperties(properties)
	if err != nil {
		return nil, err
	}

	if dsd == nil {
		dsd, err = dal.Create(ctx, name, pp, web)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = dal.Update(ctx, dsd, pp, web)
		if err != nil {
			return nil, err
		}
	}

	ds, err := convert2DynamicSchema(dsd)
	if err != nil {
		return nil, err
	}
	schemaCache.Store(dsd.Id, ds)

	var dw DynamicWeb
	err = json.Unmarshal([]byte(ds.Web), &dw)
	if err != nil {
		return ds, err
	}
	webCache.Store(dsd.Id, dw)

	return ds, err
}

func GetSchema(ctx *dgctx.DgContext, schemaId int64) (*DynamicSchema, error) {
	schema, _ := schemaCache.Load(schemaId)
	if schema != nil {
		return schema.(*DynamicSchema), nil
	}

	dsd, err := dal.FindById(ctx, schemaId)
	if err != nil {
		return nil, err
	}
	if dsd == nil {
		return nil, errors.New("schema not found")
	}

	ds, err := convert2DynamicSchema(dsd)
	if err != nil {
		return nil, err
	}
	schemaCache.Store(schemaId, ds)

	return ds, nil
}

func GetWeb(ctx *dgctx.DgContext, schemaId int64) (DynamicWeb, error) {
	web, _ := webCache.Load(schemaId)
	if web != nil {
		return web.(DynamicWeb), nil
	}

	schema, err := GetSchema(ctx, schemaId)
	if err != nil {
		return nil, err
	}
	if schema.Web == "" {
		return nil, errors.New("web content not found")
	}

	var dw DynamicWeb
	err = json.Unmarshal([]byte(schema.Web), &dw)
	if err != nil {
		return nil, err
	}
	webCache.Store(schemaId, dw)

	return dw, nil
}

func convert2DynamicSchema(dsd *dal.DynaSchemaDeployment) (*DynamicSchema, error) {
	dynaData := &DynamicSchemaData{}
	err := json.Unmarshal([]byte(dsd.Properties.StringNilAsEmpty()), dynaData)
	if err != nil {
		return nil, err
	}

	return &DynamicSchema{
		SchemaId:   dsd.Id,
		Name:       dsd.Name,
		Properties: dynaData.Content,
		Web:        dsd.Web.StringNilAsEmpty(),
	}, nil
}

func readProperties(properties string) (string, error) {
	propMap := map[string]map[string]any{}
	err := json.Unmarshal([]byte(properties), &propMap)
	if err != nil {
		return "", err
	}

	dynaData := &DynamicSchemaData{}
	var dynaProperties []DynamicProperty
	for k1, v1 := range propMap {
		if k1 == FillDataSourceKey {
			dynaData.FillDataSource = v1["name"].(string)
			continue
		}

		tpe := v1["type"].(string)
		vm := v1["validates"].(map[string]any)
		vs := make([]DynamicValidate, len(vm))
		j := 0
		for k2, v2 := range vm {
			vs[j] = DynamicValidate{Name: k2, Value: v2}
			j++
		}

		dynaProperties = append(dynaProperties, DynamicProperty{Name: k1, Type: tpe, Validates: vs})
	}
	dynaData.Content = dynaProperties
	bytes, err := json.Marshal(dynaData)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
