import {ethers} from "hardhat";
import {expect} from "chai";
import {FxBridgeLogicETH} from "../typechain-types";
import {AbiCoder, TransactionRequest} from "ethers"

describe("fork network and fx bridge test", function () {
    let gasAddress: string
    let bridgeAddress: string
    let adminAddress: string
    const abiCode = new AbiCoder;

    beforeEach(async function () {
        if (!process.env.FORK_ENABLE) {
            this.skip()
        }
        const network = await ethers.provider.getNetwork()
        switch (network.chainId.toString()) {
            case "1":
                gasAddress = "0x00000000219ab540356cBB839Cbe05303d7705Fa"
                bridgeAddress = "0x6f1D09Fed11115d65E1071CD2109eDb300D80A27"
                adminAddress = "0x0F413055AdEF9b61e9507928c6856F438d690882"
                break
            case "11155111":
                gasAddress = "0x6Cc9397c3B38739daCbfaA68EaD5F5D77Ba5F455"
                bridgeAddress = "0xd384a8e8822Ea845e83eb5AA2877239150615C18"
                adminAddress = "0xcF8049f0B918650614D5bf18CF15af080eFdDEe1"
                break
            default:
                throw new Error("Unsupported network")
        }
    })

    it("upgrade bridge contract", async function () {
        const gasSigner = await ethers.getImpersonatedSigner(gasAddress)
        await gasSigner.sendTransaction({
            to: adminAddress,
            value: ethers.parseEther("100")
        })
        const adminSigner = await ethers.getImpersonatedSigner(adminAddress)

        const bridgeFactory = await ethers.getContractFactory("FxBridgeLogicETH")

        const bridgeContractV1 = bridgeFactory.attach(bridgeAddress) as FxBridgeLogicETH

        const oldFxBridgeId = await bridgeContractV1.state_fxBridgeId()
        const oldPowerThreshold = await bridgeContractV1.state_powerThreshold()
        const oldLastEventNonce = await bridgeContractV1.state_lastEventNonce()
        const oldCheckpoint = await bridgeContractV1.state_lastOracleSetCheckpoint()
        const oldOracleSetNonce = await bridgeContractV1.state_lastOracleSetNonce()
        const oldBridgeTokens = await bridgeContractV1.getBridgeTokenList()

        let oldTokenStatus = new Map();
        for (const bridgeToken of oldBridgeTokens) {
            const status = await bridgeContractV1.tokenStatus(bridgeToken.addr)
            const batchNonce = await bridgeContractV1.state_lastBatchNonces(bridgeToken.addr)
            oldTokenStatus.set(bridgeToken.addr.toString(), {batchNonce: batchNonce, status: status})
        }

        const bridgeLogicContract = await bridgeFactory.deploy()
        await bridgeLogicContract.waitForDeployment()
        const bridgeLogicContractAddress = await bridgeLogicContract.getAddress()

        // 0x3659cfe6 is the signature of the upgradeTo(address) function
        const data = ethers.concat([
            '0x3659cfe6',
            abiCode.encode(['address'], [bridgeLogicContractAddress])
        ])

        const transaction: TransactionRequest = {
            to: bridgeAddress,
            data: data,
        }

        const upgradeTx = await adminSigner.sendTransaction(transaction)
        await upgradeTx.wait()

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
            expect(status.isOriginated).to.equal(oldTokenStatus.get(bridgeToken.addr).status.isOriginated)
            expect(status.isActive).to.equal(oldTokenStatus.get(bridgeToken.addr).status.isActive)
            expect(status.isExist).to.equal(oldTokenStatus.get(bridgeToken.addr).status.isExist)
            const batchNonce = await bridgeContractV1.state_lastBatchNonces(bridgeToken.addr)
            expect(batchNonce.toString()).to.equal(oldTokenStatus.get(bridgeToken.addr).batchNonce.toString())
            expect(await bridgeContractV1.state_lastBridgeCallNonces(1)).to.equal(false)
        }
    }).timeout(100000);
});
