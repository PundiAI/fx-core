import {ethers} from "hardhat";

const fs = require('fs')

export type BridgeToken = {
    name: string
    symbol: string
    decimals: number
    total_supply: string
    is_original: boolean
    target_ibc?: string
    address?: string
}

export type BridgeInfo = {
    chain_name: string
    bridge_contract: string
    bridge_logic_address?: string
    bridge_contract_address?: string
    bridge_token: BridgeToken[]
}

async function main() {
    const signers = await ethers.getSigners()
    
    const config_file = process.env.CONFIG_PATH || "./config.json"
    const out_file = process.env.OUT_PATH || "./out.json"

    if (!config_file || !out_file) {
        console.error("CONFIG_PATH or OUT_PATH is not set")
        return
    }

    const bridge_info: BridgeInfo[] = JSON.parse(fs.readFileSync(config_file, 'utf8'))

    for (let i = 0; i < bridge_info.length; i++) {
        for (let j = 0; j < bridge_info[i].bridge_token.length; j++) {
            const erc20_factory = await ethers.getContractFactory("ERC20TokenTest");

            const erc20 = await erc20_factory.deploy(
                bridge_info[i].bridge_token[j].name,
                bridge_info[i].bridge_token[j].symbol,
                bridge_info[i].bridge_token[j].decimals,
                bridge_info[i].bridge_token[j].total_supply
            );
            await erc20.deployed();

            bridge_info[i].bridge_token[j].address = erc20.address
        }

        const bridge_logic_factory = await ethers.getContractFactory(bridge_info[i].bridge_contract)
        const bridge_logic = await bridge_logic_factory.deploy();
        await bridge_logic.deployed();
        bridge_info[i].bridge_logic_address = bridge_logic.address

        const proxy_factory = await ethers.getContractFactory("TransparentUpgradeableProxy");
        const proxy = await proxy_factory.deploy(bridge_logic.address, signers[0].address, "0x");
        await proxy.deployed();

        bridge_info[i].bridge_contract_address = proxy.address
    }

    fs.writeFile(out_file, JSON.stringify(bridge_info, null, 2), function (err: any) {
        if (err) return console.error(err);
    });
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});