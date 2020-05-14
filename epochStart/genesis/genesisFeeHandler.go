package genesis

import (
	"math"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/process"
)

var _ process.TransactionFeeHandler = (*feeHandler)(nil)

type feeHandler struct {
}

// NewGenesisFeeHandler returns a new genesis fee handler
func NewGenesisFeeHandler() *feeHandler {
	return &feeHandler{}
}

// DeveloperPercentage -
func (f *feeHandler) DeveloperPercentage() float64 {
	return 0
}

// MinGasPrice -
func (f *feeHandler) MinGasPrice() uint64 {
	return 0
}

// MaxGasLimitPerBlock -
func (f *feeHandler) MaxGasLimitPerBlock() uint64 {
	return math.MaxUint64
}

// ComputeGasLimit -
func (f *feeHandler) ComputeGasLimit(_ process.TransactionWithFeeHandler) uint64 {
	return 0
}

// ComputeFee -
func (f *feeHandler) ComputeFee(_ process.TransactionWithFeeHandler) *big.Int {
	return big.NewInt(0)
}

// CheckValidityTxValues -
func (f *feeHandler) CheckValidityTxValues(_ process.TransactionWithFeeHandler) error {
	return nil
}

// CreateBlockStarted -
func (f *feeHandler) CreateBlockStarted() {
}

// GetAccumulatedFees -
func (f *feeHandler) GetAccumulatedFees() *big.Int {
	return big.NewInt(0)
}

// GetDeveloperFees -
func (f *feeHandler) GetDeveloperFees() *big.Int {
	return big.NewInt(0)
}

// ProcessTransactionFee -
func (f *feeHandler) ProcessTransactionFee(_ *big.Int, _ *big.Int, _ []byte) {
}

// RevertFees -
func (f *feeHandler) RevertFees(_ [][]byte) {
}

// IsInterfaceNil -
func (f *feeHandler) IsInterfaceNil() bool {
	return f == nil
}