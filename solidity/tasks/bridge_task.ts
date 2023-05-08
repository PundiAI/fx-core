import {task} from "hardhat/config";
import {boolean, string} from "hardhat/internal/core/params/argumentTypes";
import {bech32} from "bech32";
import {
    DISABLE_CONFIRM_FLAG,
    DRIVER_PATH_FLAG,
    GAS_LIMIT_FLAG,
    GAS_PRICE_FLAG,
    IS_LEDGER_FLAG,
    MAX_FEE_PER_GAS_FLAG,
    MAX_PRIORITY_FEE_PER_GAS_FLAG,
    MNEMONIC_FLAG,
    NONCE_FLAG,
    PRIVATE_KEY_FLAG,
    PROMPT_CHECK_TRANSACTION_DATA,
    SUB_CHECK_PRIVATE_KEY,
    SUB_CONFIRM_TRANSACTION,
    SUB_CREATE_TRANSACTION
} from "./subtasks";
import {BigNumber} from "ethers";

task("send-to-fx", "call bridge contract sendToFx()")
    .addParam("bridgeContract", "bridge token address", undefined, string, false)
    .addParam("bridgeToken", "bridge token address", undefined, string, false)
    .addParam("amount", "amount to bridge", undefined, string, false)
    .addParam("destination", "destination address", undefined, string, false)
    .addParam("targetIbc", "target ibc address", "", string, true)
    .addParam(NONCE_FLAG, "nonce", undefined, string, true)
    .addParam(GAS_PRICE_FLAG, "gas price", undefined, string, true)
    .addParam(MAX_FEE_PER_GAS_FLAG, "max fee per gas", undefined, string, true)
    .addParam(MAX_PRIORITY_FEE_PER_GAS_FLAG, "max priority fee per gas", undefined, string, true)
    .addParam(GAS_LIMIT_FLAG, "gas limit", undefined, string, true)
    .addParam(PRIVATE_KEY_FLAG, "send tx by private key", undefined, string, true)
    .addParam(MNEMONIC_FLAG, "send tx by mnemonic", undefined, string, true)
    .addParam(IS_LEDGER_FLAG, "ledger to send tx", false, boolean, true)
    .addParam(DRIVER_PATH_FLAG, "manual HD Path derivation (overrides BIP44 config)", "m/44'/60'/0'/0/0", string, true)
    .addParam(DISABLE_CONFIRM_FLAG, "disable confirm", false, boolean, true)
    .setAction(async (taskArgs, hre) => {
        const {wallet} = await hre.run(SUB_CHECK_PRIVATE_KEY, taskArgs);

        const bridgeTokenContract = await hre.ethers.getContractAt("ERC20TokenTest", taskArgs.bridgeToken, wallet);
        const from = await wallet.getAddress();

        const allowanceAmount = await bridgeTokenContract.allowance(from, taskArgs.bridgeContract);

        if (allowanceAmount.lt(BigNumber.from(taskArgs.amount))) {
            const erc20_factory = await hre.ethers.getContractFactory("ERC20TokenTest");
            const data = erc20_factory.interface.encodeFunctionData(
                "approve",
                [taskArgs.bridgeContract, taskArgs.amount]
            )

            const tx = await hre.run(SUB_CREATE_TRANSACTION, {
                from: from,
                to: taskArgs.bridgeToken,
                data: data,
                gasPrice: taskArgs.gasPrice,
                maxFeePerGas: taskArgs.maxFeePerGas,
                maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
                nonce: taskArgs.nonce,
                gasLimit: taskArgs.gasLimit,
            });

            const {answer} = await hre.run(SUB_CONFIRM_TRANSACTION, {
                message: `${PROMPT_CHECK_TRANSACTION_DATA}(send tx approve)`,
                disableConfirm: taskArgs.disableConfirm,
            });
            if (!answer) return;

            try {
                const approveTx = await wallet.sendTransaction(tx);
                await approveTx.wait();
                console.log(`Approve success, ${approveTx.hash}`)
            } catch (e) {
                console.log(`Approve failed, ${e}`)
                return;
            }
        }
        const bridge_factory = await hre.ethers.getContractFactory("FxBridgeLogic");

        const destination_bc = bech32.fromWords(bech32.decode(taskArgs.destination).words);
        const destination_bc_hex = ('0x' + '0'.repeat(24) + Buffer.from(destination_bc).toString('hex')).toString()

        const target = hre.ethers.utils.formatBytes32String(taskArgs.targetIbc);

        const data = bridge_factory.interface.encodeFunctionData(
            "sendToFx",
            [taskArgs.bridgeToken, destination_bc_hex, target, taskArgs.amount]
        )

        const tx = await hre.run(SUB_CREATE_TRANSACTION, {
            from: from,
            to: taskArgs.bridgeContract,
            data: data,
            gasPrice: taskArgs.gasPrice,
            maxFeePerGas: taskArgs.maxFeePerGas,
            maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
            nonce: taskArgs.nonce,
            gasLimit: taskArgs.gasLimit,
        });

        const {answer} = await hre.run(SUB_CONFIRM_TRANSACTION, {
            message: `${PROMPT_CHECK_TRANSACTION_DATA}(send tx sendToFx)`,
            disableConfirm: taskArgs.disableConfirm,
        });
        if (!answer) return;

        try {
            const sendToFxTx = await wallet.sendTransaction(tx);
            await sendToFxTx.wait();
            console.log(`SendToFx success, ${sendToFxTx.hash}`)
        } catch (e) {
            console.log(`SendToFx failed, ${e}`)
        }
    });
