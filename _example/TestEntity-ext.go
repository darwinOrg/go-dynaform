package _example

import (
	"errors"
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	"github.com/darwinOrg/go-dynaform/dal"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/rolandhe/daog"
	txrequest "github.com/rolandhe/daog/tx"
)

func InsertTestEntity(ctx *dgctx.DgContext, entity *TestEntity) error {
	tc, err := daog.NewTransContext(dal.DynaFormDb, txrequest.RequestWrite, ctx.TraceId)
	if err != nil {
		dglogger.Errorln(ctx, err)
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			dglogger.Errorln(ctx, "panic:", r)
			err = errors.New(fmt.Sprint("panic:", r))
		}
		tc.Complete(err)
	}()

	_, err = TestEntityDao.Insert(tc, entity)
	return err
}
