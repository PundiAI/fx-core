import {task} from "hardhat/config";
import {string} from "hardhat/internal/core/params/argumentTypes";
import {
    AddTxParam,
    SUB_CHECK_PRIVATE_KEY,
    SUB_CONFIRM_TRANSACTION,
    SUB_CREATE_TRANSACTION,
    SUB_GET_CONTRACT_ADDR,
    TransactionToJson
} from "./subtasks";

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

AddTxParam([deploy])