import {ethers} from "hardhat";
import {TryCatchTest, TryCatchTestV2} from "../typechain-types";
import {expect} from "chai";
import {HardhatEthersSigner} from "@nomicfoundation/hardhat-ethers/signers";
import * as wasi from "wasi";

describe("try catch test", function () {
    let tryCatchContract: TryCatchTest;
    let tryCatchV2Address: string;
    let signers: HardhatEthersSigner[];
    let placeholder = ethers.encodeBytes32String("placeholder")

    beforeEach(async function () {
        signers = await ethers.getSigners()

        const proxyFactory = await ethers.getContractFactory('TransparentUpgradeableProxy')
        const tryCatchFactory = await ethers.getContractFactory('TryCatchTest')
        const tryCatch = <TryCatchTest>await tryCatchFactory.deploy()

        const tryCatchV2Factory = await ethers.getContractFactory("TryCatchTestV2")
        const tryCatchV2 = await tryCatchV2Factory.deploy()
        tryCatchV2Address = await tryCatchV2.getAddress()

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

    it("upgrade to v2", async function () {
        await tryCatchContract.tryCatch(2)
        await tryCatchContract.tryCatch(2)
        expect(await tryCatchContract.test()).to.equal(4)
        expect(await tryCatchContract.initialized()).to.equal(true)
        expect((await tryCatchContract.placeholder())).to.equal(placeholder)

        const proxy = await ethers.getContractAt('ITransparentUpgradeableProxy', await tryCatchContract.getAddress())
        await proxy.connect(signers[1]).upgradeTo(tryCatchV2Address)
        const tryCatchV2Factory = await ethers.getContractFactory("TryCatchTestV2")
        const tryCatchV2 = <TryCatchTestV2>tryCatchV2Factory.attach(await tryCatchContract.getAddress())

        await tryCatchV2.setTokenStatus(signers[2].address, false, false, false, 0)
        expect((await tryCatchV2.tokenStatus(signers[2].address))["0"]).to.equal(false)
        expect((await tryCatchV2.tokenStatus(signers[2].address))["1"]).to.equal(false)
        expect((await tryCatchV2.tokenStatus(signers[2].address))["2"]).to.equal(false)
        expect((await tryCatchV2.tokenStatus(signers[2].address))["3"]).to.equal(0)

        expect(await tryCatchV2.test()).to.equal(4)
        expect(await tryCatchV2.initialized()).to.equal(true)
        expect((await tryCatchContract.placeholder())).to.equal(placeholder)
        expect((await tryCatchV2.tokenStatus(signers[0].address))["0"]).to.equal(true)
        expect((await tryCatchV2.tokenStatus(signers[0].address))["1"]).to.equal(true)
        expect((await tryCatchV2.tokenStatus(signers[0].address))["2"]).to.equal(true)
        expect((await tryCatchV2.tokenStatus(signers[0].address))["3"]).to.equal(0)
        await expect(tryCatchV2.setTokenType(signers[1].address, 4)).to.be.revertedWithoutReason()
        await tryCatchV2.setTokenType(signers[1].address, 2)
        expect((await tryCatchV2.tokenStatus(signers[1].address))["0"]).to.equal(false)
        expect((await tryCatchV2.tokenStatus(signers[1].address))["1"]).to.equal(false)
        expect((await tryCatchV2.tokenStatus(signers[1].address))["2"]).to.equal(true)
        expect((await tryCatchV2.tokenStatus(signers[1].address))["3"]).to.equal(2)

        await expect(tryCatchContract.tryCatch(1)).to.be.emit(tryCatchContract, "Log").withArgs("call failed")
        expect(await tryCatchContract.test()).to.equal(4)

        await expect(tryCatchContract.tryCatch(2)).to.be.emit(tryCatchContract, "Log").withArgs("call success")
        expect(await tryCatchContract.test()).to.equal(6)
    })
})