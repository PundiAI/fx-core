import {ethers} from "hardhat";
import {expect} from "chai";
import {FxBridgeLogicETH} from "../typechain-types";
import {AbiCoder, TransactionRequest} from "ethers"

describe("fork mainnet fx bridge test", function () {
    it.skip("update bridge contract", async function () {
        const gasAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
        const adminAddress = "0x0F413055AdEF9b61e9507928c6856F438d690882"
        const ownerAddress = "0xE77A7EA2F1DC25968b5941a456d99D37b80E98B5"
        const bridgeContractAddress = "0x6f1D09Fed11115d65E1071CD2109eDb300D80A27"

        const gasSigner = await ethers.getImpersonatedSigner(gasAddress)

        await gasSigner.sendTransaction({
            to: adminAddress,
            value: ethers.parseEther("100")
        })

        await gasSigner.sendTransaction({
            to: ownerAddress,
            value: ethers.parseEther("100")
        })

        const adminSigner = await ethers.getImpersonatedSigner(adminAddress)

        const ownerSigner = await ethers.getImpersonatedSigner(ownerAddress)

        const bridgeFactory = await ethers.getContractFactory("FxBridgeLogicETH")

        const bridgeContractV1 = bridgeFactory.attach(bridgeContractAddress) as FxBridgeLogicETH
        const oldFxBridgeId = await bridgeContractV1.state_fxBridgeId()
        const oldPowerThreshold = await bridgeContractV1.state_powerThreshold()
        const fxAddress = await bridgeContractV1.state_fxOriginatedToken()
        const oldLastEventNonce = await bridgeContractV1.state_lastEventNonce()
        const oldCheckpoint = await bridgeContractV1.state_lastOracleSetCheckpoint()
        const oldOracleSetNonce = await bridgeContractV1.state_lastOracleSetNonce()
        const oldBridgeTokens = await bridgeContractV1.getBridgeTokenList()

        let oldBatchNonce = new Map();

        for (const bridgeToken of oldBridgeTokens) {
            const batchNonce = await bridgeContractV1.state_lastBatchNonces(bridgeToken.addr)
            oldBatchNonce.set(bridgeToken.addr.toString(), batchNonce)
        }

        const bridgeContract = await bridgeFactory.deploy()
        await bridgeContract.waitForDeployment()

        const bridgeLogicContractAddress = await bridgeContract.getAddress()

        // 0x3659cfe6 is the signature of the upgradeTo(address) function
        const abiCode = new AbiCoder;

        const data = ethers.concat([
            '0x3659cfe6',
            abiCode.encode(['address'], [bridgeLogicContractAddress])
        ])

        const transaction: TransactionRequest = {
            to: bridgeContractAddress,
            data: data,
        }

        const upgradeTx = await adminSigner.sendTransaction(transaction)
        await upgradeTx.wait()

        // const migrateTx = await bridgeContractV1.connect(ownerSigner).migrate()
        // await migrateTx.wait()

        const fxBridgeId = await bridgeContractV1.state_fxBridgeId()
        const powerThreshold = await bridgeContractV1.state_powerThreshold()
        const lastEventNonce = await bridgeContractV1.state_lastEventNonce()
        const checkpoint = await bridgeContractV1.state_lastOracleSetCheckpoint()
        const oracleSetNonce = await bridgeContractV1.state_lastOracleSetNonce()
        const bridgeTokens = await bridgeContractV1.getBridgeTokenList()

        expect(fxBridgeId).to.equal(oldFxBridgeId)
        expect(powerThreshold.toString()).to.equal(oldPowerThreshold.toString())
        expect(lastEventNonce.toString()).to.equal(oldLastEventNonce.toString())
        expect(checkpoint).to.equal(oldCheckpoint)
        expect(oracleSetNonce.toString()).to.equal(oldOracleSetNonce.toString())

        for (const bridgeToken of bridgeTokens) {
            const status = await bridgeContractV1.tokenStatus(bridgeToken.addr)
            if (bridgeToken.addr.toString() === fxAddress.toString()) {
                expect(status.isOriginated).to.equal(true)
            } else {
                expect(status.isOriginated).to.equal(false)
            }
            expect(status.isActive).to.equal(true)
            expect(status.isExist).to.equal(true)
            const batchNonce = await bridgeContractV1.state_lastBatchNonces(bridgeToken.addr)
            expect(batchNonce.toString()).to.equal(oldBatchNonce.get(bridgeToken.addr).toString())
        }
    }).timeout(100000);
});
