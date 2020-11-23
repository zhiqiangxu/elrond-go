package process

const (
	blockIndex               = "blocks"
	miniblocksIndex          = "miniblocks"
	txIndex                  = "transactions"
	tpsIndex                 = "tps"
	validatorsIndex          = "validators"
	roundIndex               = "rounds"
	ratingIndex              = "rating"
	accountsIndex            = "accounts"
	accountsHistoryIndex     = "accountshistory"
	receiptsIndex            = "receipts"
	scResultsIndex           = "scresults"
	accountsESDTHistoryIndex = "accountsesdthistory"

	txPolicy                  = "transactions_policy"
	blockPolicy               = "blocks_policy"
	miniblocksPolicy          = "miniblocks_policy"
	validatorsPolicy          = "validators_policy"
	roundPolicy               = "rounds_policy"
	ratingPolicy              = "rating_policy"
	accountsHistoryPolicy     = "accountshistory_policy"
	accountsESDTHistoryPolicy = "accountsesdthistory_policy"

	metachainTpsDocID   = "meta"
	shardTpsDocIDPrefix = "shard"

	bulkSizeThreshold = 800000 // 0.8MB
)