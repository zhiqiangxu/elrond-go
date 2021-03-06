package txsFee

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm/txsFee/utils"
	"github.com/ElrondNetwork/elrond-go/process"
	"github.com/stretchr/testify/require"
)

func TestRelayedBuildInFunctionChangeOwnerCallShouldWork(t *testing.T) {
	testContext := vm.CreatePreparedTxProcessorWithVMs(
		t,
		vm.ArgEnableEpoch{
			PenalizedTooMuchGasEnableEpoch: 100,
		})
	defer testContext.Close()

	scAddress, owner := utils.DoDeploy(t, &testContext, "../arwen/testdata/counter/output/counter.wasm")
	testContext.TxFeeHandler.CreateBlockStarted()

	relayerAddr := []byte("12345678901234567890123456789033")
	newOwner := []byte("12345678901234567890123456789112")
	gasPrice := uint64(10)
	gasLimit := uint64(1000)

	txData := []byte(core.BuiltInFunctionChangeOwnerAddress + "@" + hex.EncodeToString(newOwner))
	innerTx := vm.CreateTransaction(1, big.NewInt(0), owner, scAddress, gasPrice, gasLimit, txData)

	_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, big.NewInt(30000))

	rtxData := utils.PrepareRelayerTxData(innerTx)
	rTxGasLimit := 1 + gasLimit + uint64(len(rtxData))
	rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, owner, gasPrice, rTxGasLimit, rtxData)

	retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
	require.Equal(t, vmcommon.Ok, retCode)
	require.Nil(t, err)

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckOwnerAddr(t, &testContext, scAddress, newOwner)

	expectedBalanceRelayer := big.NewInt(25760)
	vm.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedBalanceRelayer)

	expectedBalance := big.NewInt(89030)
	vm.TestAccount(t, testContext.Accounts, owner, 2, expectedBalance)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(4240), accumulatedFees)

	developerFees := testContext.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(0), developerFees)
}

func TestRelayedBuildInFunctionChangeOwnerCallWrongOwnerShouldConsumeGas(t *testing.T) {
	testContext := vm.CreatePreparedTxProcessorWithVMs(t, vm.ArgEnableEpoch{})
	defer testContext.Close()

	scAddress, owner := utils.DoDeploy(t, &testContext, "../arwen/testdata/counter/output/counter.wasm")
	testContext.TxFeeHandler.CreateBlockStarted()

	relayerAddr := []byte("12345678901234567890123456789033")
	sndAddr := []byte("12345678901234567890123456789113")
	newOwner := []byte("12345678901234567890123456789112")
	gasPrice := uint64(10)
	gasLimit := uint64(1000)

	txData := []byte(core.BuiltInFunctionChangeOwnerAddress + "@" + hex.EncodeToString(newOwner))
	innerTx := vm.CreateTransaction(1, big.NewInt(0), sndAddr, scAddress, gasPrice, gasLimit, txData)

	_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, big.NewInt(30000))

	rtxData := utils.PrepareRelayerTxData(innerTx)
	rTxGasLimit := 1 + gasLimit + uint64(len(rtxData))
	rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, owner, gasPrice, rTxGasLimit, rtxData)

	retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Equal(t, process.ErrFailedTransaction, err)
	require.Nil(t, testContext.GetLatestError())

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckOwnerAddr(t, &testContext, scAddress, owner)

	expectedBalanceRelayer := big.NewInt(16610)
	vm.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedBalanceRelayer)

	expectedBalance := big.NewInt(89030)
	vm.TestAccount(t, testContext.Accounts, owner, 1, expectedBalance)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(13390), accumulatedFees)

	developerFees := testContext.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(0), developerFees)
}

func TestRelayedBuildInFunctionChangeOwnerInvalidAddressShouldConsumeGas(t *testing.T) {
	testContext := vm.CreatePreparedTxProcessorWithVMs(t, vm.ArgEnableEpoch{})
	defer testContext.Close()

	scAddress, owner := utils.DoDeploy(t, &testContext, "../arwen/testdata/counter/output/counter.wasm")
	testContext.TxFeeHandler.CreateBlockStarted()

	relayerAddr := []byte("12345678901234567890123456789033")
	newOwner := []byte("invalidAddress")
	gasPrice := uint64(10)
	gasLimit := uint64(1000)

	txData := []byte(core.BuiltInFunctionChangeOwnerAddress + "@" + hex.EncodeToString(newOwner))
	innerTx := vm.CreateTransaction(1, big.NewInt(0), owner, scAddress, gasPrice, gasLimit, txData)

	_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, big.NewInt(30000))

	rtxData := utils.PrepareRelayerTxData(innerTx)
	rTxGasLimit := 1 + gasLimit + uint64(len(rtxData))
	rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, owner, gasPrice, rTxGasLimit, rtxData)

	retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
	require.Equal(t, vmcommon.ExecutionFailed, retCode)
	require.Equal(t, process.ErrInvalidAddressLength, testContext.GetLatestError())

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckOwnerAddr(t, &testContext, scAddress, owner)

	expectedBalanceRelayer := big.NewInt(17330)
	vm.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedBalanceRelayer)

	expectedBalance := big.NewInt(89030)
	vm.TestAccount(t, testContext.Accounts, owner, 2, expectedBalance)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(12670), accumulatedFees)

	developerFees := testContext.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(0), developerFees)
}

func TestRelayedBuildInFunctionChangeOwnerCallInsufficientGasLimitShouldConsumeGas(t *testing.T) {
	testContext := vm.CreatePreparedTxProcessorWithVMs(t, vm.ArgEnableEpoch{})
	defer testContext.Close()

	scAddress, owner := utils.DoDeploy(t, &testContext, "../arwen/testdata/counter/output/counter.wasm")
	testContext.TxFeeHandler.CreateBlockStarted()

	relayerAddr := []byte("12345678901234567890123456789033")
	newOwner := []byte("12345678901234567890123456789112")
	gasPrice := uint64(10)

	txData := []byte(core.BuiltInFunctionChangeOwnerAddress + "@" + hex.EncodeToString(newOwner))
	gasLimit := uint64(len(txData) - 1)
	innerTx := vm.CreateTransaction(1, big.NewInt(0), owner, scAddress, gasPrice, gasLimit, txData)

	_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, big.NewInt(30000))

	rtxData := utils.PrepareRelayerTxData(innerTx)
	rTxGasLimit := 1 + gasLimit + uint64(len(rtxData))
	rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, owner, gasPrice, rTxGasLimit, rtxData)

	retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
	require.Equal(t, vmcommon.UserError, retCode)
	require.Nil(t, testContext.GetLatestError())

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckOwnerAddr(t, &testContext, scAddress, owner)

	expectedBalanceRelayer := big.NewInt(25810)
	vm.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedBalanceRelayer)

	expectedBalance := big.NewInt(89030)
	vm.TestAccount(t, testContext.Accounts, owner, 2, expectedBalance)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(4190), accumulatedFees)

	developerFees := testContext.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(0), developerFees)
}

func TestRelayedBuildInFunctionChangeOwnerCallOutOfGasShouldConsumeGas(t *testing.T) {
	testContext := vm.CreatePreparedTxProcessorWithVMs(t, vm.ArgEnableEpoch{})
	defer testContext.Close()

	scAddress, owner := utils.DoDeploy(t, &testContext, "../arwen/testdata/counter/output/counter.wasm")
	testContext.TxFeeHandler.CreateBlockStarted()

	relayerAddr := []byte("12345678901234567890123456789033")
	newOwner := []byte("12345678901234567890123456789112")
	gasPrice := uint64(10)

	txData := []byte(core.BuiltInFunctionChangeOwnerAddress + "@" + hex.EncodeToString(newOwner))
	gasLimit := uint64(len(txData) + 1)
	innerTx := vm.CreateTransaction(1, big.NewInt(0), owner, scAddress, gasPrice, gasLimit, txData)

	_, _ = vm.CreateAccount(testContext.Accounts, relayerAddr, 0, big.NewInt(30000))

	rtxData := utils.PrepareRelayerTxData(innerTx)
	rTxGasLimit := 1 + gasLimit + uint64(len(rtxData))
	rtx := vm.CreateTransaction(0, innerTx.Value, relayerAddr, owner, gasPrice, rTxGasLimit, rtxData)

	retCode, err := testContext.TxProcessor.ProcessTransaction(rtx)
	require.Equal(t, vmcommon.ExecutionFailed, retCode)
	require.Equal(t, process.ErrNotEnoughGas, testContext.GetLatestError())

	_, err = testContext.Accounts.Commit()
	require.Nil(t, err)

	utils.CheckOwnerAddr(t, &testContext, scAddress, owner)

	expectedBalanceRelayer := big.NewInt(25790)
	vm.TestAccount(t, testContext.Accounts, relayerAddr, 1, expectedBalanceRelayer)

	expectedBalance := big.NewInt(89030)
	vm.TestAccount(t, testContext.Accounts, owner, 2, expectedBalance)

	// check accumulated fees
	accumulatedFees := testContext.TxFeeHandler.GetAccumulatedFees()
	require.Equal(t, big.NewInt(4210), accumulatedFees)

	developerFees := testContext.TxFeeHandler.GetDeveloperFees()
	require.Equal(t, big.NewInt(0), developerFees)
}
