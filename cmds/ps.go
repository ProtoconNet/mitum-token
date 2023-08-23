package cmds

import (
	"context"

	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"

	// "github.com/ProtoconNet/mitum-token/operation/processor"
	"github.com/ProtoconNet/mitum2/isaac"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/ps"
)

var PNameOperationProcessorsMap = ps.Name("mitum-dao-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *currencyprocessor.OperationProcessor
	var set *hint.CompatibleSet

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		currencycmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &set,
	); err != nil {
		return pctx, err
	}

	err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	if err != nil {
		return pctx, err
	}
	err = opr.SetGetNewProcessorFunc(processor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}

	// if err := opr.SetProcessor(
	// 	dao.CreateDAOHint,
	// 	dao.NewCreateDAOProcessor(),
	// ); err != nil {
	// 	return pctx, err
	// } else if err := opr.SetProcessor(
	// 	dao.UpdatePolicyHint,
	// 	dao.NewUpdatePolicyProcessor(),
	// ); err != nil {
	// 	return pctx, err
	// }

	// _ = set.Add(dao.CreateDAOHint, func(height base.Height) (base.OperationProcessor, error) {
	// 	return opr.New(
	// 		height,
	// 		db.State,
	// 		nil,
	// 		nil,
	// 	)
	// })

	var f currencycmds.ProposalOperationFactHintFunc = IsSupportedProposalOperationFactHintFunc

	pctx = context.WithValue(pctx, currencycmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, set) //revive:disable-line:modifies-parameter
	pctx = context.WithValue(pctx, currencycmds.ProposalOperationFactHintContextKey, f)

	return pctx, nil
}
