package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tokenServiceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "token_service_contract_height"),
	},
}

var tokenBalanceIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "token_balance_contract_address_height"),
	},
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	DefaultIndexes[DefaultColNameToken] = tokenServiceIndexModels
	DefaultIndexes[DefaultColNameTokenBalance] = tokenBalanceIndexModels
}
