import {ethers} from "hardhat";
import {HardhatEthersSigner} from "@nomicfoundation/hardhat-ethers/signers";
import {expect} from "chai";
import {ERC20TokenTest, FxBridgeLogic} from "../typechain-types"
import {ZeroAddress, encodeBytes32String, keccak256, AbiCoder, Signature} from "ethers"
import {arrayify} from "@ethersproject/bytes";

export async function getSignerAddresses(signers: HardhatEthersSigner[]) {
    return await Promise.all(signers.map(signer => signer.getAddress()));
}

export function makeSubmitBridgeCallHash(
    gravityId: string, sender: string,receiver: string, tokens: string[],amounts: string[],
    to: string, data: string, memo: string, nonce: number | string, timeout: number | string
){
    let methodName = encodeBytes32String("bridgeCall");
    let abiCoder = new AbiCoder()
    return keccak256(
        abiCoder.encode(
            ["bytes32", "bytes32", "address", "address", "address[]", "uint256[]", "address", "bytes", "bytes", "uint256", "uint256"],
            [gravityId, methodName, sender, receiver, tokens, amounts, to, data, memo, nonce, timeout]
        )
    );
}

export async function signHash(signers: HardhatEthersSigner[], hash: string) {
    let v: number[] = [];
    let r: string[] = [];
    let s: string[] = [];

    const signMessage = arrayify(hash)
    for (let i = 0; i < signers.length; i = i + 1) {
        const sig = await signers[i].signMessage(signMessage);
        const signature = Signature.from(sig);

        v.push(signature.v);
        r.push(signature.r);
        s.push(signature.s);
    }
    return {v, r, s};
}

describe("submit bridge call tests", function () {
    let deploy: HardhatEthersSigner;
    let admin: HardhatEthersSigner;
    let user1: HardhatEthersSigner;
    let oracle: HardhatEthersSigner;
    let erc20Token: ERC20TokenTest;
    let fxBridge: FxBridgeLogic;


    let totalSupply = "10000"
    const gravityId: string = encodeBytes32String("eth-fxcore");
    const powerThreshold = 6666;
    const powers: number[] = [10000];

    let validators: any;
    let valAddresses: any;

    beforeEach(async function () {
        const signers = await ethers.getSigners()
        deploy = signers[0]
        admin = signers[1]
        oracle = signers[2]
        user1 = signers[3]

        validators = [oracle];
        valAddresses = await getSignerAddresses(validators);

        const erc20TokenFactory = await ethers.getContractFactory('ERC20TokenTest')
        erc20Token = await erc20TokenFactory.deploy("ERC20 Token", "ERC20", "18", ethers.parseEther(totalSupply))
        const erc20TokenAddress = await erc20Token.getAddress()
        expect(await erc20Token.balanceOf(deploy.address)).to.equal(ethers.parseEther("10000"))


        const fxBridgeLogicFactory = await ethers.getContractFactory('FxBridgeLogic')
        const fxBridgeLogic = await fxBridgeLogicFactory.deploy()
        const fxBridgeLogicAddress = await fxBridgeLogic.getAddress()

        const transparentUpgradeableProxyFactory = await ethers.getContractFactory("TransparentUpgradeableProxy");
        const fxBridgeLogicProxy = await transparentUpgradeableProxyFactory.deploy(fxBridgeLogicAddress, admin.address, "0x")
        const fxBridgeLogicProxyAddress = await fxBridgeLogicProxy.getAddress()

        fxBridge = <FxBridgeLogic>fxBridgeLogicFactory.attach(fxBridgeLogicProxyAddress)

        await fxBridge.init(gravityId, powerThreshold, valAddresses, powers)
        await fxBridge.addBridgeToken(erc20TokenAddress, encodeBytes32String(""), true)
    })

    it("should submit bridge call", async function () {
        const erc20TokenAddress = await erc20Token.getAddress()
        const amount = "1000";
        const timeout = await ethers.provider.getBlockNumber() + 1000;

        const digest = makeSubmitBridgeCallHash(
            gravityId,user1.address,user1.address,[erc20TokenAddress],
            [amount], ZeroAddress,"0x","0x",1,timeout)

        const {v, r, s} = await signHash(validators, digest)

        const bridgeCallData: FxBridgeLogic.BridgeCallDataStruct = {
            sender: user1.address,
            receiver: user1.address,
            tokens: [erc20TokenAddress],
            amounts: [amount],
            to: ZeroAddress,
            data: "0x",
            memo: "0x",
            timeout: timeout
        };

        await fxBridge.submitBridgeCall(
            valAddresses,
            powers,
            v,
            r,
            s,
            [0, 1],
            bridgeCallData
        );
    })
})