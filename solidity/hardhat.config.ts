import {HardhatUserConfig} from "hardhat/config"
import "hardhat-dependency-compiler"
import "@nomicfoundation/hardhat-ethers"
import '@typechain/hardhat'
import "hardhat-gas-reporter"

import './tasks/task'

const config: HardhatUserConfig = {
    defaultNetwork: "hardhat",
    networks: {
        hardhat: {},
        localhost: {
            url: `${process.env.LOCAL_URL || "http://127.0.0.1:8545"}`,
        },
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
    gasReporter: {
        enabled: false,
        currency: 'USD',
        gasPrice: 30
    },
};

export default config;
