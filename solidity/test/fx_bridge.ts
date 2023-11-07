import {ethers, network} from "hardhat";
import {TransactionRequest} from "@ethersproject/abstract-provider/src.ts";
import {expect} from "chai";

describe("fork ethereum", function () {
    it("update-bridge-contract", async function () {
        const adminAddress = "0x0F413055AdEF9b61e9507928c6856F438d690882"
        const ownerAddress = "0xE77A7EA2F1DC25968b5941a456d99D37b80E98B5"
        const bridgeContractAddress = "0x6f1D09Fed11115d65E1071CD2109eDb300D80A27"

        await network.provider.request({
            method: "hardhat_impersonateAccount",
            params: [adminAddress],
        });
        const adminSigner = ethers.provider.getSigner(adminAddress)

        await network.provider.request({
            method: "hardhat_impersonateAccount",
            params: [ownerAddress],
        });

        const ownerSigner = ethers.provider.getSigner(ownerAddress)
        const bridgeFactory = await ethers.getContractFactory("FxBridgeLogicETH")

        const bridgeContractV1 = bridgeFactory.attach(bridgeContractAddress)
        const oldLastEventNonce = await bridgeContractV1.state_lastEventNonce()
        const oldCheckpoint = await bridgeContractV1.state_lastOracleSetCheckpoint()
        const oldOracleSetNonce = await bridgeContractV1.state_lastOracleSetNonce()
        const fxAddress = await bridgeContractV1.state_fxOriginatedToken()

        const bridgeContract = await bridgeFactory.deploy()
        await bridgeContract.deployed()

        const data = ethers.utils.hexConcat([
            '0x3659cfe6',
            ethers.utils.defaultAbiCoder.encode(['address'], [bridgeContract.address])
        ])

        const transaction: TransactionRequest = {
            to: bridgeContractAddress,
            data: data,
        }

        const upgradeTx = await adminSigner.sendTransaction(transaction)
        await upgradeTx.wait()

        const migrateTx = await bridgeContractV1.connect(ownerSigner).migrate()
        await migrateTx.wait()

        const lastEventNonce = await bridgeContractV1.state_lastEventNonce()
        const checkpoint = await bridgeContractV1.state_lastOracleSetCheckpoint()
        const oracleSetNonce = await bridgeContractV1.state_lastOracleSetNonce()
        const bridgeTokens = await bridgeContractV1.getBridgeTokenList()

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
        }

    });
});