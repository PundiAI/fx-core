import {task} from "hardhat/config";
import {boolean, string} from "hardhat/internal/core/params/argumentTypes";
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
    SUB_CREATE_TRANSACTION,
    SUB_SEND_ETH,
    VALUE_FLAG
} from "./subtasks";

task("send", "send tx, Example: npx hardhat send 0x... transfer(address,uint256) 0x... 1000000000000000000 --privateKey ...")
    .addVariadicPositionalParam("params", "send tx params", undefined, string, true)
    .addParam(NONCE_FLAG, "nonce", undefined, string, true)
    .addParam(GAS_PRICE_FLAG, "gas price", undefined, string, true)
    .addParam(MAX_FEE_PER_GAS_FLAG, "max fee per gas", undefined, string, true)
    .addParam(MAX_PRIORITY_FEE_PER_GAS_FLAG, "max priority fee per gas", undefined, string, true)
    .addParam(GAS_LIMIT_FLAG, "gas limit", undefined, string, true)
    .addParam(VALUE_FLAG, "value", undefined, string, true)
    .addParam(PRIVATE_KEY_FLAG, "send tx by private key", undefined, string, true)
    .addParam(MNEMONIC_FLAG, "send tx by mnemonic", undefined, string, true)
    .addParam(IS_LEDGER_FLAG, "ledger to send tx", false, boolean, true)
    .addParam(DRIVER_PATH_FLAG, "manual HD Path derivation (overrides BIP44 config)", "m/44'/60'/0'/0/0", string, true)
    .addParam(DISABLE_CONFIRM_FLAG, "disable confirm", false, boolean, true)
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
            await hre.run(SUB_SEND_ETH, taskArgs);
            return
        }

        const abi = parseAbiItemFromSignature(func)
        const abiInterface = new hre.ethers.utils.Interface([abi])

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
            message: `${PROMPT_CHECK_TRANSACTION_DATA}(send tx ${abi.name})`,
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
    const match = signature.match(funcRegex);

    if (!match) {
        throw new Error('Invalid signature');
    }

    data.name = match[1];
    let index = 1

    if (match[2]) {
        data.inputs = match[2].split(',').map((input) => {
            let [type, name] = input.trim().split(' ');
            if (!name) {
                name = `param${index}`;
                index++;
            }
            return {name: name, type: type};
        });
    }

    if (match[3]) {
        data.outputs = match[3].split(',').map((output) => {
            let [type, name] = output.trim().split(' ');
            if (!name) {
                name = `param${index}`;
                index++;
            }
            return {name: name, type: type};
        });
    }

    data.type = 'function';
    data.payable = false;
    data.constant = false;

    return data;
}



