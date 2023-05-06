import {ethers} from "hardhat";
import {BridgeInfo} from "./deploy_bridge";
import fs from "fs";
import axios from "axios";

async function main() {
    const signers = await ethers.getSigners()
    const vote_power = 2834678415

    const out_file = process.env.OUT_PATH || "./out.json"
    const rest_rpc = process.env.REST_RPC || "http://127.0.0.1:1317"

    if (!out_file) {
        console.error("OUT_PATH is not set")
        return
    }
    const bridge_info: BridgeInfo[] = JSON.parse(fs.readFileSync(out_file, 'utf8'))

    for (let i = 0; i < bridge_info.length; i++) {
        const bridge_logic_factory = await ethers.getContractFactory(bridge_info[i].bridge_contract)

        const proxy = await ethers.getContractAt("ITransparentUpgradeableProxy", bridge_info[i].bridge_contract_address as string)

        const oracle_set = await GetOracleSet(rest_rpc, bridge_info[i].chain_name)
        const gravity_id_str = await GetGravityId(rest_rpc, bridge_info[i].chain_name)
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

        await proxy.upgradeToAndCall(bridge_info[i].bridge_logic_address, init_data)
        await proxy.changeAdmin(signers[signers.length - 1].address)

        const bridge_contract = await ethers.getContractAt(bridge_info[i].bridge_contract, bridge_info[i].bridge_contract_address as string)

        for (let j = 0; j < bridge_info[i].bridge_token.length; j++) {
            const ibc = ethers.utils.formatBytes32String(bridge_info[i].bridge_token[j].target_ibc || "")
            bridge_contract.addBridgeToken(bridge_info[i].bridge_token[j].address, ibc, bridge_info[i].bridge_token[j].is_original)
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