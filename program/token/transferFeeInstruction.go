package token

type TransferFeeInstruction uint8

const (
	InitializeTransferFeeConfig TransferFeeInstruction = iota
	TransferCheckedWithFee
	WithdrawWithheldTokensFromMint
	WithdrawWithheldTokensFromAccounts
	HarvestWithheldTokensToMint
	SetTransferFee
)
