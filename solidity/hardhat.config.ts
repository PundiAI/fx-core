import {HardhatUserConfig} from "hardhat/config"
import "hardhat-dependency-compiler"
import "@nomicfoundation/hardhat-ethers"
import '@typechain/hardhat'
import "hardhat-gas-reporter"
import "@nomicfoundation/hardhat-verify";
import "@nomicfoundation/hardhat-chai-matchers";

import './tasks/task'

const config: HardhatUserConfig = {
    defaultNetwork: "hardhat",
    networks: {
        hardhat: {
            // forking: {
            //     url: `${process.env.MAINNET_URL || "https://mainnet.infura.io/v3/"+process.env.INFURA_KEY}`,
            // }
            chainId: 1337,
        },
        mainnet: {
            url: `${process.env.MAINNET_URL || "https://mainnet.infura.io/v3/" + process.env.INFURA_KEY}`,
            chainId: 1,
        },
        sepolia: {
            url: `${process.env.SEPOLIA_URL || "https://sepolia.infura.io/v3/" + process.env.INFURA_KEY}`,
            chainId: 11155111,
        },
        arbitrumSepolia: {
            url: `${process.env.ARBITRUM_URL || "https://sepolia-rollup.arbitrum.io/rpc"}`,
            chainId: 421614,
        },
        optimisticSepolia: {
            url: `${process.env.OPTIMISTIC_URL || "https://sepolia.optimism.io"}`,
            chainId: 11155420,
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
            mainnet: `${process.env.ETHERSCAN_API_KEY}`,
            sepolia: `${process.env.ETHERSCAN_API_KEY}`,
            arbitrumSepolia: `${process.env.ETHERSCAN_API_KEY}`,
            optimisticSepolia: `${process.env.ETHERSCAN_API_KEY}`,
        },
        customChains: [
            {
                network: "optimisticSepolia",
                chainId: 11155420,
                urls: {
                    apiURL: "https://api-sepolia-optimistic.etherscan.io/api",
                    browserURL: "https://sepolia-optimism.etherscan.io/"
                }
            }
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
