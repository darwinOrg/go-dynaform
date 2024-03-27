package dynaform

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	ve "github.com/darwinOrg/go-validator-ext"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
)

var customValidator = ve.NewCustomValidator()

type DynamicData map[string]any

type Dynamic struct {
	SchemaId int64
	data     any
}

func NewDynamic(schemaId int64, data any) *Dynamic {
	return &Dynamic{SchemaId: schemaId, data: data}
}

func (dyna *Dynamic) Get(field string) any {
	dd := reflect.ValueOf(dyna.data)
	if dd.Kind() == reflect.Ptr {
		dd = dd.Elem()
	}
	val := dd.FieldByName(FirstCharToUpper(field)).Interface()
	if val == nil {
		return nil
	}
	switch val.(type) {
	case *decimal.Decimal:
		return *val.(*decimal.Decimal)
	default:
		return val
	}
}

func (dyna *Dynamic) Set(field string, value any) {
	dd := reflect.ValueOf(dyna.data)
	if dd.Kind() == reflect.Ptr {
		dd = dd.Elem()
	}
	f := dd.FieldByName(FirstCharToUpper(field))
	if f.CanSet() {
		f.Set(reflect.ValueOf(value))
	}
}

func (dyna *Dynamic) String() string {
	bytes, err := json.Marshal(dyna.data)
	if err != nil {
		return ""
	}

	mp := map[string]any{}
	if err := json.Unmarshal(bytes, &mp); err != nil {
		return ""
	}

	mp[SchemaId] = dyna.SchemaId
	mpBytes, _ := json.Marshal(mp)
	return string(mpBytes)
}

func (dyna *Dynamic) Value() (driver.Value, error) {
	str, err := json.Marshal(dyna.data)
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func (dyna *Dynamic) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if data, ok := value.(string); ok {
		dc, err := GetObjectFromJson(&dgctx.DgContext{TraceId: uuid.NewString()}, data)
		if err != nil {
			return err
		}

		dyna.SchemaId = dc.SchemaId
		dyna.data = dc.data
	}

	return nil
}

func CreateObject(ctx *dgctx.DgContext, schemaId int64) (*Dynamic, error) {
	return getObjectFromDataMap(ctx, map[string]any{SchemaId: schemaId}, false)
}

func GetObjectFromJson(ctx *dgctx.DgContext, data string) (*Dynamic, error) {
	schemaIdHolder := &struct {
		SchemaId int64 `json:"schemaId"`
	}{}

	err := json.Unmarshal([]byte(data), schemaIdHolder)
	if err != nil {
		return nil, err
	}

	schemaId := schemaIdHolder.SchemaId
	if schemaId == 0 {
		return nil, errors.New("schemaId cannot be zero")
	}

	schema, err := GetSchema(ctx, schemaId)
	if err != nil {
		dglogger.Errorf(ctx, "get schema error, schemaId: %d, err: %v", schemaId, err)
		return nil, err
	}

	typ := schema.reflect2StructType()

	typStruct := reflect.New(typ).Elem().Addr().Interface()
	err = json.Unmarshal([]byte(data), typStruct)
	if err != nil {
		dglogger.Errorf(ctx, "get schema error, schemaId: %d, err: %v", schemaId, err)
		return nil, err
	}
	marshal, _ := json.Marshal(typStruct)
	fmt.Println(string(marshal))

	err = customValidator.Struct(typStruct)
	if err != nil {
		dglogger.Errorf(ctx, "validate struct data error, schemaId: %d, err: %v", schemaId, err)
		return nil, err
	}

	return NewDynamic(schemaId, typStruct), nil
}

func GetObjectFromDynamicData(ctx *dgctx.DgContext, data DynamicData) (*Dynamic, error) {
	return getObjectFromDataMap(ctx, data, true)
}

func getObjectFromDataMap(ctx *dgctx.DgContext, dataMap map[string]any, validate bool) (*Dynamic, error) {
	schemaId := dataMap[SchemaId]
	if schemaId == nil {
		return nil, errors.New("schemaId cannot be empty")
	}

	sid, err := strconv.ParseInt(fmt.Sprint(schemaId), 10, 64)
	if err != nil {
		return nil, err
	}

	schema, err := GetSchema(ctx, sid)
	if err != nil {
		dglogger.Errorf(ctx, "get schema error, schemaId: %d, err: %v", schemaId, err)
		return nil, err
	}

	data, err := schema.reflect2StructData(ctx, dataMap)
	if err != nil {
		dglogger.Errorf(ctx, "reflect struct data error, schemaId: %d, err: %v", schemaId, err)
		return nil, err
	}

	if validate {
		err = customValidator.Struct(data)
		if err != nil {
			dglogger.Errorf(ctx, "validate struct data error, schemaId: %d, err: %v", schemaId, err)
			return nil, err
		}
	}

	return NewDynamic(sid, data), nil
}
