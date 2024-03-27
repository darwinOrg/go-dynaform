package dal

import (
	"errors"
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
	"github.com/rolandhe/daog/ttypes"
	txrequest "github.com/rolandhe/daog/tx"
	"time"
)

var DynaFormDb daog.Datasource

func Create(ctx *dgctx.DgContext, name string, properties string, web string) (*DynaSchemaDeployment, error) {
	tc, err := daog.NewTransContext(DynaFormDb, txrequest.RequestWrite, ctx.TraceId)
	if err != nil {
		dglogger.Errorln(ctx, err)
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			dglogger.Errorln(ctx, "panic:", r)
			err = errors.New(fmt.Sprint("panic:", r))
		}
		tc.Complete(err)
	}()

	now := time.Now()
	dsd := &DynaSchemaDeployment{
		Name:       name,
		Properties: *ttypes.FromString(properties),
		Web:        *ttypes.FromString(web),
		Version:    1,
		CreatedAt:  ttypes.NormalDatetime(now),
		ModifiedAt: ttypes.NormalDatetime(now),
	}
	_, err = DynaSchemaDeploymentDao.Insert(tc, dsd)
	return dsd, err
}

func Update(ctx *dgctx.DgContext, dsd *DynaSchemaDeployment, properties string, web string) (int64, error) {
	tc, err := daog.NewTransContext(DynaFormDb, txrequest.RequestWrite, ctx.TraceId)
	if err != nil {
		dglogger.Errorln(ctx, err)
		return 0, err
	}
	defer func() {
		if r := recover(); r != nil {
			dglogger.Errorln(ctx, "panic:", r)
			err = errors.New(fmt.Sprint("panic:", r))
		}
		tc.Complete(err)
	}()

	dsd.Properties = *ttypes.FromString(properties)
	dsd.Web = *ttypes.FromString(web)
	dsd.Version++
	dsd.ModifiedAt = ttypes.NormalDatetime(time.Now())

	return DynaSchemaDeploymentDao.Update(tc, dsd)
}

func FindById(ctx *dgctx.DgContext, schemaId int64) (*DynaSchemaDeployment, error) {
	tc, err := daog.NewTransContext(DynaFormDb, txrequest.RequestReadonly, ctx.TraceId)
	if err != nil {
		dglogger.Errorln(ctx, err)
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			dglogger.Errorln(ctx, "panic:", r)
			err = errors.New(fmt.Sprint("panic:", r))
		}
		tc.Complete(err)
	}()

	mc := daog.NewMatcher().Eq(DynaSchemaDeploymentFields.Id, schemaId)
	return DynaSchemaDeploymentDao.QueryOneMatcher(tc, mc)
}

func FindByName(ctx *dgctx.DgContext, name string) (*DynaSchemaDeployment, error) {
	tc, err := daog.NewTransContext(DynaFormDb, txrequest.RequestReadonly, ctx.TraceId)
	if err != nil {
		dglogger.Errorln(ctx, err)
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			dglogger.Errorln(ctx, "panic:", r)
			err = errors.New(fmt.Sprint("panic:", r))
		}
		tc.Complete(err)
	}()

	mc := daog.NewMatcher().Eq(DynaSchemaDeploymentFields.Name, name)
	return DynaSchemaDeploymentDao.QueryOneMatcher(tc, mc)
}
