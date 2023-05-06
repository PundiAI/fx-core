import {HardhatUserConfig} from "hardhat/config";
import "hardhat-dependency-compiler"
import "@nomiclabs/hardhat-ethers";

import './tasks/bridge_task';

const port = process.env.LOCAL_PORT || 8535

const config: HardhatUserConfig = {
    defaultNetwork: "hardhat",
    networks: {
        hardhat: {
            mining: {
                interval: 1000
            },
            accounts: {
                mnemonic: "test test test test test test test test test test test junk",
                initialIndex: 0,
                count: 10,
                accountsBalance: "1000000000000000000000000",
            },
        },
        localhost: {
            url: `http://localhost:${port}`,
        }
    },
    solidity: {
        compilers: [
            {
                version: "0.8.0",
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200
                    }
                }
            },
            {
                version: "0.8.1",
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200
                    }
                }
            },
            {
                version: "0.8.2",
                settings: {
                    optimizer: {
                        enabled: true,
                        runs: 200
                    }
                }
            },
        ]
    },
    dependencyCompiler: {
        paths: [
            "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol",
            "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol",
        ],
    },
};

export default config;