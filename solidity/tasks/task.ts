import {task} from "hardhat/config";
import {string} from "hardhat/internal/core/params/argumentTypes";
import {
    AddTxParam,
    SUB_CHECK_PRIVATE_KEY,
    SUB_CONFIRM_TRANSACTION,
    SUB_CREATE_TRANSACTION,
    SUB_GET_NODE_URL,
    SUB_SEND_ETH,
    TransactionToJson
} from "./subtasks";

import "./bridge_tasks"
import "./contract_task"

const send = task("send", "send tx, Example: npx hardhat send 0x... transfer(address,uint256) 0x... 1000000000000000000 --privateKey ...")
    .addVariadicPositionalParam("params", "send tx params", undefined, string, true)
    .setAction(async (taskArgs, hre) => {
        const {wallet} = await hre.run(SUB_CHECK_PRIVATE_KEY, taskArgs);
        const from = await wallet.getAddress();

        const {params} = taskArgs;
        const to = params[0];
        const func = params[1];
        params.splice(0, 2);

        if (!to) {
            throw new Error("Please provide to address");
        }

        if (!func) {
            await hre.run(SUB_SEND_ETH, {
                to: to, value: taskArgs.value, wallet: wallet,
                gasPrice: taskArgs.gasPrice,
                maxFeePerGas: taskArgs.maxFeePerGas,
                maxPriorityFeePerGas: taskArgs.maxPriorityFeePerGas,
                nonce: taskArgs.nonce,
                gasLimit: taskArgs.gasLimit,
            });
            return
        }
        const abi = parseAbiItemFromSignature(func)
        const abiInterface = new hre.ethers.Interface([abi])

        if (abi.inputs && (abi.inputs.length !== params.length)) {
            throw new Error(`Please provide ${abi.inputs.length} params`)
        }

        const data = abiInterface.encodeFunctionData(abi.name as string, params);

        const tx = await hre.run(SUB_CREATE_TRANSACTION, {
            from: from, to: to, data: data, value: taskArgs.value,
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
            const txRes = await wallet.sendTransaction(tx)
            console.log(txRes.hash)
            await txRes.wait()
        } catch (e) {
            console.error(`send tx failed, ${e}`)
            return;
        }
    });

task("call", "call contract, Example: npx hardhat call 0x... balanceOf(address)(uint256) 0x...")
    .addParam("from", "", undefined, string, true)
    .addVariadicPositionalParam("params", "call contract params", undefined, string, true)
    .setAction(async (taskArgs, hre) => {
            const nodeUrl = await hre.run(SUB_GET_NODE_URL, taskArgs)
            const provider = new hre.ethers.JsonRpcProvider(nodeUrl);

            const {from, params} = taskArgs;
            const to = params[0];
            const func = params[1];
            params.splice(0, 2);

            if (!to || !func) {
                throw new Error("Please provide to address and func");
            }

            const abi = parseAbiItemFromSignature(func)
            const abiInterface = new hre.ethers.Interface([abi])

            if (abi.inputs && (abi.inputs.length !== params.length)) {
                throw new Error(`Please provide ${abi.inputs.length} params`)
            }
            const data = abiInterface.encodeFunctionData(abi.name as string, params)
            const result = await provider.call({
                    to: to,
                    from: from,
                    data: data
                }
            )
            console.log(splitByComma(abiInterface.decodeFunctionResult(abi.name as string, result).toString()))
        }
    )

interface AbiItem {
    name?: string;
    inputs?: { name: string; type: string; }[];
    outputs?: { name: string; type: string; }[];
    constant?: boolean;
    payable?: boolean;
    type: string;
}

function parseAbiItemFromSignature(signature: string): AbiItem {
    const data = {
        name: '',
        inputs: [] as { name: string, type: string }[],
        outputs: [] as { name: string, type: string }[],
        constant: false,
        payable: false,
        type: 'function',
    };

    const funcRegex = /^(\w+)\((.*)\)\((.*)\)$/s;
    let match = signature.match(funcRegex);

    if (!match) {
        const funcRegex = /^(\w+)\((.*)\)$/s;
        match = signature.match(funcRegex);
        if (!match) {
            throw new Error('Invalid signature');
        }
    }

    data.name = match[1];
    if (match[2]) {
        data.inputs = match[2].split(',').map((input) => {
            let [type, name] = input.trim().split(' ');
            return {name: name, type: type};
        });
    }

    if (match[3]) {
        data.outputs = match[3].split(',').map((output) => {
            let [type, name] = output.trim().split(' ');
            return {name: name, type: type};
        });
    }

    data.type = 'function';
    data.payable = false;
    data.constant = false;

    return data;
}

function splitByComma(str: string): string {
    const arr = str.split(",");
    if (arr[arr.length - 1] === "") arr.pop();
    return arr.join("\n");
}

AddTxParam([send])