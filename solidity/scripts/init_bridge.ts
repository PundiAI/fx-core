import {ethers} from "hardhat";
import {BridgeInfo} from "./deploy_bridge";
import fs from "fs";
import axios from "axios";

async function main() {
    const signers = await ethers.getSigners()
    const vote_power = 2834678415

    const config_file = process.env.CONFIG_FILE || "./bridge.json"
    const rest_rpc = process.env.REST_RPC || "http://127.0.0.1:1317"

    if (!config_file || !rest_rpc) {
        console.error("BRIDGE_CONFIG_FILE or REST_RPC is not set")
        return
    }
    const bridge_info: BridgeInfo = JSON.parse(fs.readFileSync(config_file, 'utf8'))


    for (const externalChain of bridge_info.external_chain_list) {
        const bridge_logic_factory = await ethers.getContractFactory(externalChain.bridge_contract)
        const proxy = await ethers.getContractAt("ITransparentUpgradeableProxy", externalChain.bridge_contract_address as string)


        const oracle_set = await GetOracleSet(rest_rpc, externalChain.chain_name)
        const gravity_id_str = await GetGravityId(rest_rpc, externalChain.chain_name)
        const gravity_id = ethers.utils.formatBytes32String(gravity_id_str);

        const external_addresses = [];
        const powers = [];
        let powers_sum = 0;

        for (let i = 0; i < oracle_set.members.length; i++) {
            external_addresses.push(oracle_set.members[i].external_address);
            powers.push(oracle_set.members[i].power);
            powers_sum += oracle_set.members[i].power;
        }


        if (powers_sum < vote_power) {
            console.error("Incorrect power! Please inspect the oracle set")
            console.log(`Current oracle set:\n${oracle_set}`)
            return
        }

        const init_data = bridge_logic_factory.interface.encodeFunctionData('init', [
            gravity_id, vote_power, external_addresses, powers
        ])
        
        await proxy.upgradeToAndCall(externalChain.bridge_logic_address, init_data)
        await proxy.changeAdmin(signers[signers.length - 1].address)

        const bridge_contract = await ethers.getContractAt(externalChain.bridge_contract, externalChain.bridge_contract_address as string)

        for (const oneToOne of bridge_info.bridge_token_list.one_to_one) {
            if (oneToOne.chain_name == externalChain.chain_name) {
                const ibc = ethers.utils.formatBytes32String(oneToOne.target_ibc || "")
                bridge_contract.addBridgeToken(oneToOne.address, ibc, oneToOne.is_original)
            }
        }
        for (const oneToMany of bridge_info.bridge_token_list.one_to_many) {
            for (const chain of oneToMany.chain_list) {
                if (chain.chain_name == externalChain.chain_name) {
                    const ibc = ethers.utils.formatBytes32String(chain.target_ibc || "")
                    bridge_contract.addBridgeToken(chain.address, ibc, chain.is_original)
                }
            }
        }
    }
}

type Oracle = {
    power: number;
    external_address: string
};
type OracleSet = {
    members: Oracle[];
    nonce: number;
};

export async function GetOracleSet(restRpc: string, chainName: string): Promise<OracleSet> {
    const request_string = restRpc + `/fx/crosschain/v1/oracle_set/current?chain_name=${chainName}`
    const response = await axios.get(request_string);
    return response.data.oracle_set;
}

export async function GetGravityId(restRpc: string, chainName: string): Promise<string> {
    const request_string = restRpc + `/fx/crosschain/v1/params?chain_name=${chainName}`
    const response = await axios.get(request_string);
    return response.data.params.gravity_id;
}

main()
    .catch((error) => {
        console.error(error);
        process.exitCode = 1;
    });