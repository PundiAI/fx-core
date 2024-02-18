import {HardhatUserConfig} from "hardhat/config"
import "hardhat-dependency-compiler"
import "@nomicfoundation/hardhat-ethers"
import '@typechain/hardhat'
import "hardhat-gas-reporter"
import "@nomicfoundation/hardhat-verify";

import './tasks/task'

const config: HardhatUserConfig = {
    defaultNetwork: "hardhat",
    networks: {
        hardhat: {
            // forking: {
            //     url: `${process.env.MAINNET_URL || "https://mainnet.infura.io/v3/infura-key"}`,
            // }
            chainId: 1337
        },
        mainnet: {
            url: `${process.env.MAINNET_URL || "https://mainnet.infura.io/v3/infura-key"}`,
            chainId: 1,
        },
        goerli: {
            url: `${process.env.GOERLI_URL || "https://goerli.infura.io/v3/infura-key"}`,
            chainId: 5,
        },
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
    etherscan: {
        apiKey: {
            mainnet: `${process.env.ETHERSCAN_API_KEY || "scan-key"}`,
            goerli: `${process.env.ETHERSCAN_API_KEY || "scan-key"}`,
        }
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
