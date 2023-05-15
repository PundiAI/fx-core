import {ethers} from "hardhat";

const fs = require('fs')

export type BridgeTokenList = {
    one_to_one: OneToOne[]
    one_to_many: OneToMany[]
}

export type OneToOne = {
    chain_name: string
    name: string
    symbol: string
    decimals: number
    total_supply: string
    is_original: boolean
    target_ibc: string
    address?: string
}

export type OneToManyChain = {
    chain_name: string
    total_supply: string
    is_original: boolean
    target_ibc?: string
    address?: string
}

export type OneToMany = {
    name: string
    symbol: string
    decimals: number
    base_denom: string
    chain_list: OneToManyChain[]
}

export type ExternalChain = {
    chain_name: string
    bridge_contract: string
    bridge_logic_address?: string
    bridge_contract_address?: string
}

export type BridgeInfo = {
    external_chain_list: ExternalChain[]
    bridge_token_list: BridgeTokenList
}

async function main() {
    const signers = await ethers.getSigners()

    const config_file = process.env.BRIDGE_CONFIG_FILE || "./bridge.json"
    const out_file = process.env.CONFIG_OUT_FILE || "./bridge.json"

    if (!config_file) {
        console.error("BRIDGE_CONFIG_FILE || CONFIG_OUT_FILE is not set")
        return
    }

    const bridge_info: BridgeInfo = JSON.parse(fs.readFileSync(config_file, 'utf8'))

    for (let i = 0; i < bridge_info.external_chain_list.length; i++) {
        const bridge_logic_factory = await ethers.getContractFactory(bridge_info.external_chain_list[i].bridge_contract)
        const bridge_logic = await bridge_logic_factory.deploy();
        await bridge_logic.deployed();
        bridge_info.external_chain_list[i].bridge_logic_address = bridge_logic.address

        const proxy_factory = await ethers.getContractFactory("TransparentUpgradeableProxy");
        const proxy = await proxy_factory.deploy(bridge_logic.address, signers[0].address, "0x");
        await proxy.deployed();
        bridge_info.external_chain_list[i].bridge_contract_address = proxy.address
    }

    const erc20_factory = await ethers.getContractFactory("ERC20TokenTest");

    for (let i = 0; i < bridge_info.bridge_token_list.one_to_one.length; i++) {
        const erc20 = await erc20_factory.deploy(
            bridge_info.bridge_token_list.one_to_one[i].name,
            bridge_info.bridge_token_list.one_to_one[i].symbol,
            bridge_info.bridge_token_list.one_to_one[i].decimals,
            bridge_info.bridge_token_list.one_to_one[i].total_supply
        );
        await erc20.deployed();
        bridge_info.bridge_token_list.one_to_one[i].address = erc20.address
    }

    for (let i = 0; i < bridge_info.bridge_token_list.one_to_many.length; i++) {
        for (let j = 0; j < bridge_info.bridge_token_list.one_to_many[i].chain_list.length; j++) {
            const erc20 = await erc20_factory.deploy(
                bridge_info.bridge_token_list.one_to_many[i].name,
                bridge_info.bridge_token_list.one_to_many[i].symbol,
                bridge_info.bridge_token_list.one_to_many[i].decimals,
                bridge_info.bridge_token_list.one_to_many[i].chain_list[j].total_supply
            );
            await erc20.deployed();
            bridge_info.bridge_token_list.one_to_many[i].chain_list[j].address = erc20.address
        }
    }
    console.log(JSON.stringify(bridge_info, null, 2))
    fs.writeFile(out_file, JSON.stringify(bridge_info, null, 2), function (err: any) {
        if (err) return console.error(err);
    });
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});