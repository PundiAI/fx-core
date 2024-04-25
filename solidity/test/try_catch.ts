import {ethers} from "hardhat";
import {TryCatchTest} from "../typechain-types";
import {expect} from "chai";
import {HardhatEthersSigner} from "@nomicfoundation/hardhat-ethers/signers";

describe("try catch test", function () {
    let tryCatchContract: TryCatchTest;
    let signers: HardhatEthersSigner[];
    let placeholder = ethers.encodeBytes32String("placeholder")

    beforeEach(async function () {
        signers = await ethers.getSigners()

        const proxyFactory = await ethers.getContractFactory('TransparentUpgradeableProxy')
        const tryCatchFactory = await ethers.getContractFactory('TryCatchTest')
        const tryCatch = <TryCatchTest>await tryCatchFactory.deploy()

        const proxy = await proxyFactory.deploy(await tryCatch.getAddress(), signers[1].address, "0x")
        tryCatchContract = <TryCatchTest>tryCatchFactory.attach(await proxy.getAddress())
        await tryCatchContract.initialize(placeholder)

        await tryCatchContract.setTokenStatus(signers[0].address, true, true, true)
        await tryCatchContract.setTokenStatus(signers[1].address, false, false, true)
    })

    it("should call failed", async function () {
        await expect(tryCatchContract.tryCatch(0)).to.be.emit(tryCatchContract, "Log").withArgs("call failed")
        expect(await tryCatchContract.test()).to.equal(0)
    })

    it("should call failed", async function () {
        await expect(tryCatchContract.tryCatch(1)).to.be.emit(tryCatchContract, "Log").withArgs("call failed")
        expect(await tryCatchContract.test()).to.equal(0)
    })

    it("should call success", async function () {
        await expect(tryCatchContract.tryCatch(2)).to.be.emit(tryCatchContract, "Log").withArgs("call success")
        expect(await tryCatchContract.test()).to.equal(2)
    })
})
