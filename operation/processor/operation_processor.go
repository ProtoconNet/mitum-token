package processor

import (
	"github.com/ProtoconNet/mitum-currency/v3/operation/currency"
	extensioncurrency "github.com/ProtoconNet/mitum-currency/v3/operation/extension"
	currencyprocessor "github.com/ProtoconNet/mitum-currency/v3/operation/processor"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-token/operation/token"
	"github.com/ProtoconNet/mitum-token/utils"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/pkg/errors"
)

const (
	DuplicationTypeSender   currencytypes.DuplicationType = "sender"
	DuplicationTypeCurrency currencytypes.DuplicationType = "currency"
)

func CheckDuplication(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) error {
	opr.Lock()
	defer opr.Unlock()

	if err := currencyprocessor.CheckDuplication(opr, op); err != nil {
		return err
	}

	var did string
	var didtype currencytypes.DuplicationType
	var err error

	switch t := op.(type) {
	case token.RegisterToken:
		did, didtype, err = checkDuplicateSender(t)
	case token.Mint:
		did, didtype, err = checkDuplicateSender(t)
	default:
		return nil
	}

	if err != nil {
		return err
	}

	if did != "" {
		if _, found := opr.Duplicated[did]; found {
			switch didtype {
			case DuplicationTypeSender:
				return errors.Errorf("violates only one sender in proposal")
			case DuplicationTypeCurrency:
				return errors.Errorf("duplicate currency id, %q found in proposal", did)
			default:
				return errors.Errorf("violates duplication in proposal")
			}
		}

		opr.Duplicated[did] = didtype
	}

	return nil
}

func checkDuplicateSender(op mitumbase.Operation) (string, currencytypes.DuplicationType, error) {
	fact, ok := op.Fact().(token.TokenFact)
	if !ok {
		return "", "", errors.Errorf(utils.ErrStringTypeCast(token.TokenFact{}, op.Fact()))
	}
	return fact.Sender().String(), DuplicationTypeSender, nil
}

func GetNewProcessor(opr *currencyprocessor.OperationProcessor, op mitumbase.Operation) (mitumbase.OperationProcessor, bool, error) {
	switch i, err := opr.GetNewProcessorFromHintset(op); {
	case err != nil:
		return nil, false, err
	case i != nil:
		return i, true, nil
	}

	switch t := op.(type) {
	case currency.CreateAccounts,
		currency.KeyUpdater,
		currency.Transfers,
		extensioncurrency.CreateContractAccounts,
		extensioncurrency.Withdraws,
		currency.CurrencyRegister,
		currency.CurrencyPolicyUpdater,
		currency.SuffrageInflation,
		token.RegisterToken,
		token.Mint:
		return nil, false, errors.Errorf("%T needs SetProcessor", t)
	default:
		return nil, false, nil
	}
}
