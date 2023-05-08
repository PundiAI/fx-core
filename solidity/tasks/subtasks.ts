import {subtask} from "hardhat/config";
import {LedgerSigner} from "@ethersproject/hardware-wallets";
import {BigNumber} from "ethers";

const inquirer = require('inquirer')

// sub task name
export const SUB_CHECK_PRIVATE_KEY: string = "sub:check-private-key";
export const SUB_PRIVATE_KEY_WALLET: string = "sub:generate-wallet";
export const SUB_GET_NODE_URL: string = "sub:get-eth-node-url";
export const SUB_CREATE_LEDGER_WALLET: string = "sub:create-ledger-wallet";
export const SUB_CREATE_TRANSACTION: string = "sub:create-transaction";
export const SUB_CONFIRM_TRANSACTION: string = "sub:confirm-transaction";
export const SUB_MNEMONIC_WALLET: string = "sub:mnemonic-wallet";
// public flag
export const DISABLE_CONFIRM_FLAG: string = "disableConfirm";
export const PRIVATE_KEY_FLAG = "privateKey";
export const MNEMONIC_FLAG = "mnemonic";
export const IS_LEDGER_FLAG = "isLedger";
export const DRIVER_PATH_FLAG = "driverPath";
export const NONCE_FLAG = "nonce";
export const GAS_PRICE_FLAG = "gasPrice";
export const MAX_FEE_PER_GAS_FLAG = "maxFeePerGas";
export const MAX_PRIORITY_FEE_PER_GAS_FLAG = "maxPriorityFeePerGas";
export const GAS_LIMIT_FLAG = "gasLimit";

export const DEFAULT_DRIVE_PATH = "m/44'/60'/0'/0/0";
export const DEFAULT_PRIORITY_FEE = "1500000000";
export const PROMPT_CHECK_TRANSACTION_DATA = "Do you want continue?";

type Transaction = {
    from: string,
    to?: string,
    data?: string,
    gasPrice?: BigNumber,
    maxFeePerGas?: BigNumber,
    maxPriorityFeePerGas?: BigNumber,
    nonce: number,
    gasLimit?: number,
    chainId: number
}

subtask(SUB_CREATE_TRANSACTION, "create transaction").setAction(
    async (taskArgs, hre) => {
        let {from, to, data, gasPrice, maxFeePerGas, maxPriorityFeePerGas, nonce, gasLimit, chainId} = taskArgs;
        if (gasPrice && maxFeePerGas) {
            throw new Error("Please provide only one of gasPrice or maxFeePerGas and maxPriorityFeePerGas");
        }
        if (!gasPrice && !maxFeePerGas) {
            await hre.ethers.provider.getBlock("latest").then(
                async (block) => {
                    if (block.baseFeePerGas) {
                        maxFeePerGas = block.baseFeePerGas.add(maxPriorityFeePerGas);
                    } else {
                        gasPrice = await hre.ethers.provider.getGasPrice()
                    }
                }
            );
        }
        if (maxFeePerGas) {
            maxPriorityFeePerGas = maxPriorityFeePerGas ? maxPriorityFeePerGas : BigNumber.from(DEFAULT_PRIORITY_FEE);
            maxFeePerGas = BigNumber.from(maxFeePerGas).add(maxPriorityFeePerGas);
        }
        const transaction: Transaction = {
            from: from,
            to: to,
            data: data,
            nonce: nonce ? nonce : await hre.ethers.provider.getTransactionCount(from),
            gasLimit: gasLimit ? gasLimit : await hre.ethers.provider.estimateGas({
                from: from,
                to: to,
                data: data
            }),
            chainId: chainId ? chainId : await hre.ethers.provider.getNetwork().then(network => network.chainId)
        }
        if (gasPrice) {
            transaction.gasPrice = gasPrice;
        }
        if (maxFeePerGas) {
            transaction.maxFeePerGas = maxFeePerGas;
            transaction.maxPriorityFeePerGas = maxPriorityFeePerGas;
        }
        console.log(
            "New Transaction:\n",
            `from: ${transaction.from}\n`,
            `to: ${transaction.to}\n`,
            `data: ${transaction.data}\n`,
            `gasPrice: ${transaction.gasPrice ? transaction.gasPrice.toString() : "null"}\n`,
            `maxFeePerGas: ${transaction.maxFeePerGas ? transaction.maxFeePerGas.toString() : "null"}\n`,
            `maxPriorityFeePerGas: ${transaction.maxPriorityFeePerGas ? transaction.maxPriorityFeePerGas.toString() : "null"}\n`,
            `nonce: ${transaction.nonce}\n`,
            `gasLimit: ${transaction.gasLimit ? transaction.gasLimit.toString() : "null"}\n`,
            `chainId: ${transaction.chainId}`
        )
        return transaction;
    }
);

subtask(SUB_CHECK_PRIVATE_KEY, "check the method of getting private key").setAction(
    async (taskArgs, hre) => {
        const {privateKey, isLedger, mnemonic} = taskArgs;
        if (
            privateKey && isLedger || privateKey && mnemonic || isLedger && mnemonic
        ) {
            throw new Error("Please provide only one of private key or ledger or mnemonic");
        }
        if (privateKey) {
            const {wallet} = await hre.run(SUB_PRIVATE_KEY_WALLET, taskArgs);
            return {wallet}
        }
        if (mnemonic) {
            return await hre.run(SUB_MNEMONIC_WALLET, taskArgs);
        }
        if (isLedger) {
            return await hre.run(SUB_CREATE_LEDGER_WALLET, taskArgs);
        }
        return (await hre.ethers.getSigners())[0];
    }
);

subtask(SUB_CREATE_LEDGER_WALLET, "create ledger wallet").setAction(
    async (taskArgs, hre) => {
        const {driverPath} = taskArgs;
        const nodeUrl = await hre.run(SUB_GET_NODE_URL);
        const provider = await new hre.ethers.providers.JsonRpcProvider(nodeUrl);

        const _path = driverPath ? driverPath : DEFAULT_DRIVE_PATH;

        const wallet = new LedgerSigner(provider, "hid", _path);
        return {wallet};
    });

subtask(SUB_PRIVATE_KEY_WALLET, "private key wallet account").setAction(
    async (taskArgs, hre) => {
        const {privateKey} = taskArgs;
        const nodeUrl = await hre.run(SUB_GET_NODE_URL);
        const provider = await new hre.ethers.providers.JsonRpcProvider(nodeUrl);
        const wallet = new hre.ethers.Wallet(privateKey, provider);
        return {provider, wallet};
    });

subtask(SUB_MNEMONIC_WALLET, "mnemonic wallet account").setAction(
    async (taskArgs, hre) => {
        const {mnemonic, driverPath} = taskArgs;

        const nodeUrl = await hre.run(SUB_GET_NODE_URL);
        const provider = await new hre.ethers.providers.JsonRpcProvider(nodeUrl);

        const _path = driverPath ? driverPath : DEFAULT_DRIVE_PATH;

        const wallet = hre.ethers.Wallet.fromMnemonic(mnemonic, _path).connect(provider);
        return {provider, wallet};
    }
);

subtask(SUB_GET_NODE_URL, "get node url form hardhat.network").setAction(
    async (taskArgs, hre) => {
        return "url" in hre.network.config ? hre.network.config.url : "";
    },
);

subtask(SUB_CONFIRM_TRANSACTION, "confirm transaction").setAction(
    async (taskArgs, _) => {
        const {message, disableConfirm} = taskArgs;
        let _answer;
        if (!disableConfirm) {
            const {answer} = await inquirer.createPromptModule()({
                type: "confirm",
                name: "answer",
                message,
            });
            _answer = answer;
        } else {
            _answer = true;
        }
        return {answer: _answer};
    });
