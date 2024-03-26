import {task} from "hardhat/config";
import {string} from "hardhat/internal/core/params/argumentTypes";
import {
    AddTxParam,
    BridgeStateInfoToJson,
    SUB_CHECK_PRIVATE_KEY,
    SUB_CONFIRM_TRANSACTION,
    SUB_CREATE_TRANSACTION,
    SUB_GET_CONTRACT_ADDR,
    TransactionToJson
} from "./subtasks";
import {FxBridgeLogic, IFxBridgeLogic} from "../typechain-types";

const deploy = task("deploy-contract", "deploy contract")
    .addParam("contractName", "deploy contract name", undefined, string, false)
    .addVariadicPositionalParam("params", "deploy contract params", undefined, string, true)
    .setAction(async (taskArgs, hre) => {
        const {wallet} = await hre.run(SUB_CHECK_PRIVATE_KEY, taskArgs);
        const from = await wallet.getAddress();

        const contractAddress = await hre.run(SUB_GET_CONTRACT_ADDR, {from: from})
        const contractFactory = await hre.ethers.getContractFactory(taskArgs.contractName);

        const paramData = contractFactory.interface.encodeDeploy(taskArgs.params);
        const data = contractFactory.bytecode + paramData.slice(2);

        const tx = await hre.run(SUB_CREATE_TRANSACTION, {
            from: from, data: data, value: taskArgs.value,
            gasPrice: taskArgs.gasPrice,
            maxFeePerGas: taskArgs.maxFeePerGas,
            maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
            nonce: taskArgs.nonce,
            gasLimit: taskArgs.gasLimit,
        });

        const {answer} = await hre.run(SUB_CONFIRM_TRANSACTION, {
            message: `\n${TransactionToJson(tx)}\n`,
            disableConfirm: taskArgs.disableConfirm,
        });
        if (!answer) return;

        try {
            const deployTx = await wallet.sendTransaction(tx);
            await deployTx.wait();
            console.log(`${contractAddress}`)
        } catch (e) {
            console.log(`Deploy failed, ${e}`)
            return;
        }
    });

const migrateBridge = task("migrate-bridge", "migrate bridge")
    .addParam("oldBridge", "old bridge address", undefined, string, false)
    .addParam("oldRpc", "old rpc url", undefined, string, false)
    .addParam("newBridgeProxy", "new bridge proxy address", undefined, string, true)
    .addParam("newBridgeLogic", "new bridge logic address", undefined, string, true)
    .addParam("admin", "admin address", undefined, string, true)
    .setAction(async (taskArgs, hre) => {
        let {oldBridge, oldRpc, newBridgeProxy, newBridgeLogic, admin} = taskArgs
        const provider = new hre.ethers.JsonRpcProvider(oldRpc);
        const oldBridgeContract = await hre.ethers.getContractAt("FxBridgeLogic", oldBridge) as FxBridgeLogic;
        const fxBridgeId = await oldBridgeContract.connect(provider).state_fxBridgeId()
        const powerThreshold = await oldBridgeContract.connect(provider).state_powerThreshold()
        const lastEventNonce = await oldBridgeContract.connect(provider).state_lastEventNonce()
        const lastOracleSetNonce = await oldBridgeContract.connect(provider).state_lastOracleSetNonce()
        const lastOracleSetCheckpoint = await oldBridgeContract.connect(provider).state_lastOracleSetCheckpoint()
        const bridgeTokens = await oldBridgeContract.connect(provider).getBridgeTokenList()

        let bridgeTokenAddress: string[] = [];
        let lastBatchNonce: bigint[] = [];
        let tokenStatus: IFxBridgeLogic.TokenStatusStruct[] = [];

        for (let i = 0; i < bridgeTokens.length; i++) {
            const token = bridgeTokens[i];
            bridgeTokenAddress.push(token.addr);
            const batchNonce = await oldBridgeContract.connect(provider).state_lastBatchNonces(token.addr)
            lastBatchNonce.push(batchNonce);
            const status = await oldBridgeContract.connect(provider).tokenStatus(token.addr)
            tokenStatus.push(status);
        }

        const bridgeMigrateLogicFactory: any = await hre.ethers.getContractFactory("FxBridgeMigrateLogic")
        let proxyFactory: any = await hre.ethers.getContractFactory("TransparentUpgradeableProxy")
        const initData = bridgeMigrateLogicFactory.interface.encodeFunctionData('migrateInit', [
            fxBridgeId,
            powerThreshold,
            lastEventNonce,
            lastOracleSetCheckpoint,
            lastOracleSetNonce,
            bridgeTokenAddress,
            lastBatchNonce,
            tokenStatus
        ])
        let {answer} = await hre.run(SUB_CONFIRM_TRANSACTION, {
            message: `\n${BridgeStateInfoToJson(
                fxBridgeId,
                powerThreshold,
                lastEventNonce,
                lastOracleSetNonce,
                lastOracleSetCheckpoint,
                bridgeTokenAddress,
                lastBatchNonce,
                tokenStatus
            )}\n`,
            disableConfirm: taskArgs.disableConfirm,
        });
        if (!answer) return;

        const {wallet} = await hre.run(SUB_CHECK_PRIVATE_KEY, taskArgs);
        if (!newBridgeLogic) {
            const paramData = bridgeMigrateLogicFactory.interface.encodeDeploy([]);
            const data = bridgeMigrateLogicFactory.bytecode + paramData.slice(2);
            const tx = await hre.run(SUB_CREATE_TRANSACTION, {
                from: await wallet.getAddress(), data: data,
                gasPrice: taskArgs.gasPrice,
                maxFeePerGas: taskArgs.maxFeePerGas,
                maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
                nonce: taskArgs.nonce,
                gasLimit: taskArgs.gasLimit,
            });

            answer = await hre.run(SUB_CONFIRM_TRANSACTION, {
                message: `\n${TransactionToJson(tx)}\n`,
                disableConfirm: taskArgs.disableConfirm,
            });
            if (!answer) return;

            try {
                const deployTx = await wallet.sendTransaction(tx);
                const receipt = await deployTx.wait();
                newBridgeLogic = receipt.contractAddress;
                console.log(`deploy bridge logic, ${newBridgeLogic}`)
            } catch (e) {
                console.log(`Deploy failed, ${e}`)
                return;
            }
        }
        admin = admin || await wallet.getAddress()
        if (!newBridgeProxy) {
            const paramData = proxyFactory.interface.encodeDeploy([newBridgeLogic, admin, initData]);
            const data = proxyFactory.bytecode + paramData.slice(2);
            const tx = await hre.run(SUB_CREATE_TRANSACTION, {
                from: await wallet.getAddress(), data: data,
                gasPrice: taskArgs.gasPrice,
                maxFeePerGas: taskArgs.maxFeePerGas,
                maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
                nonce: taskArgs.nonce,
                gasLimit: taskArgs.gasLimit,
            });

            answer = await hre.run(SUB_CONFIRM_TRANSACTION, {
                message: `\n${TransactionToJson(tx)}\n`,
                disableConfirm: taskArgs.disableConfirm,
            });
            if (!answer) return;

            try {
                const deployTx = await wallet.sendTransaction(tx);
                const receipt = await deployTx.wait();
                newBridgeProxy = receipt.contractAddress;
                console.log(`deploy proxy, ${newBridgeProxy}`)
            } catch (e) {
                console.log(`Deploy failed, ${e}`)
                return;
            }
        } else {
            proxyFactory = await hre.ethers.getContractAt("ITransparentUpgradeableProxy", newBridgeProxy, wallet)
            const data = proxyFactory.interface.encodeFunctionData('upgradeToAndCall', [newBridgeLogic, initData])
            const tx = await hre.run(SUB_CREATE_TRANSACTION, {
                from: await wallet.getAddress(), to: newBridgeProxy, data: data,
                gasPrice: taskArgs.gasPrice,
                maxFeePerGas: taskArgs.maxFeePerGas,
                maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
                nonce: taskArgs.nonce,
                gasLimit: taskArgs.gasLimit,
            });

            answer = await hre.run(SUB_CONFIRM_TRANSACTION, {
                message: `\n${TransactionToJson(tx)}\n`,
                disableConfirm: taskArgs.disableConfirm,
            });
            if (!answer) return;

            try {
                const migrateTx = await wallet.sendTransaction(tx);
                await migrateTx.wait();
                console.log(`migrate success, ${migrateTx.hash}`)
            } catch (e) {
                console.log(`migrate failed, ${e}`)
            }
        }
    })
AddTxParam([deploy, migrateBridge])