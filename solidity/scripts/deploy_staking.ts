import {ethers} from "hardhat";

const fs = require('fs')

async function main() {
    const staking_factory = await ethers.getContractFactory("StakingTest");
    const staking = await staking_factory.deploy()
    await staking.deployed();
    console.log(staking.address);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});