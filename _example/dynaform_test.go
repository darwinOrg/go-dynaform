package _example

import (
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-dynaform"
	"github.com/darwinOrg/go-dynaform/dal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
	"github.com/shopspring/decimal"
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	conf := &daog.DbConf{
		DbUrl:    "root:12345678@tcp(localhost:3306)/test?charset=utf8mb4&loc=Local&interpolateParams=true&parseTime=true&timeout=1s&readTimeout=2s&writeTimeout=2s",
		Size:     200,
		IdleCons: 50,
		IdleTime: 1200,
		Life:     3600,
		LogSQL:   true,
	}
	var err error
	dal.DynaFormDb, err = daog.NewDatasource(conf)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestDeploy(t *testing.T) {
	name := "Order/Eor"
	properties, _ := os.ReadFile("Eor.properties.json")
	web, _ := os.ReadFile("Eor.web.json")
	ctx := &dgctx.DgContext{TraceId: uuid.NewString()}
	schema, err := dynaform.Deploy(ctx, name, string(properties), string(web))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", schema)
}

func TestGetSchema(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: uuid.NewString()}
	schema, err := dynaform.GetSchema(ctx, 1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", schema)
}

func TestGetWeb(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: uuid.NewString()}
	web, err := dynaform.GetWeb(ctx, 1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%v\n", web)
}

func TestCreateObject(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: uuid.NewString()}
	dyna, err := dynaform.CreateObject(ctx, 1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%d\n", dyna.Get(dynaform.SchemaId))
}

func TestSetDataField(t *testing.T) {
	type TmpStruct struct {
		SchemaId int64
	}
	dyna := dynaform.NewDynamic(int64(1), &TmpStruct{SchemaId: int64(1)})
	fmt.Printf("%d\n", dyna.Get(dynaform.SchemaId))
	dyna.Set(dynaform.SchemaId, int64(2))
	fmt.Printf("%d\n", dyna.Get(dynaform.SchemaId))
}

func TestInsert(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: uuid.NewString()}
	data := map[string]any{
		"schemaId":         int64(1),
		"firstName":        "xx",
		"lastName":         "yy",
		"nationality":      "BONAIRE_SINT_EUSTATIUS_AND_SABA",
		"countryCode":      "BONAIRE_SINT_EUSTATIUS_AND_SABA",
		"issueCountryCode": "BONAIRE_SINT_EUSTATIUS_AND_SABA",
		"handleId06Status": "0",
		"serviceItemTypes": []string{"xxx"},
		"createPersonnel":  "0",
		"personnelIdentity": dynaform.PersonnelIdentity{
			Type: "PASSPORT",
			No:   "1312442345",
		},
	}
	dynamic, err := dynaform.GetObjectFromDynamicData(ctx, data)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	dynamic.Set("firstName", "飘")

	entity := &TestEntity{
		Foo:        "foo",
		Bar:        "bar",
		Data:       dynamic,
		CreatedAt:  ttypes.NormalDatetime(time.Now()),
		ModifiedAt: ttypes.NormalDatetime(time.Now()),
	}

	InsertTestEntity(ctx, entity)
}

func TestDynaSet(t *testing.T) {
	dc := &dgctx.DgContext{
		TraceId: "12141414",
	}

	dynamic, err := dynaform.GetObjectFromJson(dc, "{\"customerId\":21,\"customerName\":\"Qa-墨西哥\",\"customerCountry\":\"MEX\",\"countryCode\":\"MEX\",\"lastName\":\"吴\",\"firstName\":\"铁锤\",\"fullName\":\"吴铁锤\",\"issueCountryCode\":\"AFG\",\"personnelIdentity\":{\"type\":\"idCard\",\"no\":\"1131425345345\"},\"startEmploy\":\"2023-06-01\",\"endEmploy\":\"2025-06-01\",\"jobTitle\":\"100\",\"jobType\":\"fullTime\",\"workType\":\"remote\",\"workHoursPerWeek\":\"50\",\"preTaxSalary\":\"100\",\"workAgeContract\":\"no\",\"insuranceType\":\"111\",\"companyMainBusiness\":\"100100100\",\"professionals\":\"100100100100\",\"maritalStatus\":\"single\",\"nationality\":\"AFG\",\"dateOfBirth\":\"2023-06-01\",\"mobilePhone\":\"100100\",\"gender\":\"male\",\"residentialAddress\":\"100100100\",\"privateEmail\":\"100100@11.com\",\"emergencyContactName\":\"100100\",\"emergencyContactPhone\":\"100100\",\"emergencyContactRelationship\":\"relative\",\"previousEducation\":[{\"collegeName\":\"100\",\"levelOfDegree\":\"100\",\"from\":\"2023-01\",\"to\":\"2023-06\"}],\"previousWorkExperience\":[{\"employerName\":\"100\",\"rolePosition\":\"100100\",\"from\":\"2023-01\",\"to\":\"2023-06\"}],\"bankName\":\"100100\",\"accountNumber\":\"100100100100100100\",\"accountName\":\"100100\",\"recipientAddress\":\"100100\",\"city\":\"100100\",\"state\":\"100100\",\"postalCode\":\"100100\",\"swiftCode\":\"100100\",\"numberCLABE\":\"100100100100100100\",\"personnelContract\":[{\"fileId\":7797783,\"filename\":\"55B98BEF780C44D9A4D96FC35DA27583-6-2 (1) (1).jpg\",\"format\":\"jpg\",\"size\":24365,\"url\":\"https://qa-e.dghire.com/fgw/download/personnel/45/?id=7797783\"}],\"schemaId\":61}")
	t.Logf("dynamic: %+v,err: %v + \n", dynamic, err)
	dynamic.Set("customerName", "Qa-墨西哥~")
	t.Logf("dynamic: %s \n", dynamic.String())
}

func TestDynaGet(t *testing.T) {
	dc := &dgctx.DgContext{
		TraceId: "12141414",
	}
	decimal.MarshalJSONWithoutQuotes = true

	dynamic, err := dynaform.GetObjectFromJson(dc, "{\"customerId\":21,\"customerName\":\"Qa-墨西哥\",\"customerCountry\":\"MEX\",\"countryCode\":\"MEX\",\"lastName\":\"吴\",\"firstName\":\"铁锤\",\"fullName\":\"吴铁锤\",\"issueCountryCode\":\"AFG\",\"personnelIdentity\":{\"type\":\"idCard\",\"no\":\"1131425345345\"},\"startEmploy\":\"2023-06-01\",\"endEmploy\":\"2025-06-01\",\"jobTitle\":\"100\",\"jobType\":\"fullTime\",\"workType\":\"remote\",\"workHoursPerWeek\":\"50\",\"preTaxSalary\":100,\"workAgeContract\":\"no\",\"insuranceType\":\"111\",\"companyMainBusiness\":\"100100100\",\"professionals\":\"100100100100\",\"maritalStatus\":\"single\",\"nationality\":\"AFG\",\"dateOfBirth\":\"2023-06-01\",\"mobilePhone\":\"100100\",\"gender\":\"male\",\"residentialAddress\":\"100100100\",\"privateEmail\":\"100100@11.com\",\"emergencyContactName\":\"100100\",\"emergencyContactPhone\":\"100100\",\"emergencyContactRelationship\":\"relative\",\"previousEducation\":[{\"collegeName\":\"100\",\"levelOfDegree\":\"100\",\"from\":\"2023-01\",\"to\":\"2023-06\"}],\"previousWorkExperience\":[{\"employerName\":\"100\",\"rolePosition\":\"100100\",\"from\":\"2023-01\",\"to\":\"2023-06\"}],\"bankName\":\"100100\",\"accountNumber\":\"100100100100100100\",\"accountName\":\"100100\",\"recipientAddress\":\"100100\",\"city\":\"100100\",\"state\":\"100100\",\"postalCode\":\"100100\",\"swiftCode\":\"100100\",\"numberCLABE\":\"100100100100100100\",\"personnelContract\":[{\"fileId\":7797783,\"filename\":\"55B98BEF780C44D9A4D96FC35DA27583-6-2 (1) (1).jpg\",\"format\":\"jpg\",\"size\":24365,\"url\":\"https://qa-e.dghire.com/fgw/download/personnel/45/?id=7797783\"}],\"schemaId\":61}")
	t.Logf("dynamic: %+v,err: %v + \n", dynamic, err)
	val := dynamic.Get("preTaxSalary")
	fmt.Println(val)
	v2 := dynamic.Get("personnelContract")
	fmt.Printf("%T \n", v2)

}
