package digest

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-token/state"
	"github.com/ProtoconNet/mitum-token/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultColNameToken        = "digest_token"
	DefaultColNameTokenBalance = "digest_token_bl"
)

func Token(st *cdigest.Database, contract string) (*types.Design, error) {
	filter := util.NewBSONFilter("contract", contract)

	var design *types.Design
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNameToken,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			design, err = state.StateDesignValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		return nil, utilm.ErrNotFound.Errorf("token design, contract %s", contract)
	}

	return design, nil
}

func TokenBalance(st *cdigest.Database, contract, account string) (*common.Big, error) {
	filter := util.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)

	var amount common.Big
	var sta base.State
	var err error
	if err := st.MongoClient().GetByFilter(
		DefaultColNameTokenBalance,
		filter.D(),
		func(res *mongo.SingleResult) error {
			sta, err = cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}

			amount, err = state.StateTokenBalanceValue(sta)
			if err != nil {
				return err
			}

			return nil
		},
		options.FindOne().SetSort(util.NewBSONFilter("height", -1).D()),
	); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		//return nil, mitumutil.ErrNotFound.Errorf("token balance by contract %s, account %s", contract, account)
	}

	return &amount, nil
}
