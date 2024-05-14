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
    to: string, data: string, memo: string, nonce: number | string, timeout: number | string, eventNonce: number | string
){
    let methodName = encodeBytes32String("bridgeCall");
    let abiCoder = new AbiCoder()
    return keccak256(
        abiCoder.encode(
            ["bytes32", "bytes32", "address", "address", "address[]", "uint256[]", "address", "bytes", "bytes", "uint256", "uint256", "uint256"],
            [gravityId, methodName, sender, receiver, tokens, amounts, to, data, memo, nonce, timeout, eventNonce]
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
    let erc20Token: ERC20TokenTest;
    let fxBridge: FxBridgeLogic;


    let totalSupply = "10000"
    const gravityId: string = encodeBytes32String("eth-fxcore");
    const powerThreshold = 6666;
    const powers: number[] = [1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000];

    let validators: any;
    let valAddresses: any;

    beforeEach(async function () {
        const signers = await ethers.getSigners()
        deploy = signers[0]
        admin = signers[1]
        user1 = signers[2]


        validators = [signers[3], signers[4], signers[5], signers[6], signers[7], signers[8], signers[9], signers[10], signers[11], signers[12]];
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

        await erc20Token.transferOwnership(await fxBridge.getAddress())
    })

    async function submitBridgeCall(tokens: string[], amounts: string[], timeout: number, eventNonce: number) {
        const digest = makeSubmitBridgeCallHash(gravityId,user1.address,user1.address,tokens,
            amounts, ZeroAddress,"0x","0x",1,timeout, eventNonce)

        const {v, r, s} = await signHash(validators, digest)

        const bridgeCallData: FxBridgeLogic.BridgeCallDataStruct = {
            sender: user1.address,
            receiver: user1.address,
            tokens: tokens,
            amounts: amounts,
            to: ZeroAddress,
            data: "0x",
            memo: "0x",
            timeout: timeout,
            eventNonce: eventNonce
        };
        await fxBridge.submitBridgeCall(valAddresses, powers, v, r, s, [0, 1], bridgeCallData);
    }

    it("should submit bridge call", async function () {
        const erc20TokenAddress = await erc20Token.getAddress()
        const amount = "1000";
        const timeout = await ethers.provider.getBlockNumber() + 1000;

        await submitBridgeCall([erc20TokenAddress], [amount], timeout,0)
    })

    describe("submit bridge call batch test", function () {
        let token1: ERC20TokenTest;
        let token2: ERC20TokenTest;
        let token3: ERC20TokenTest;
        let token4: ERC20TokenTest;

        beforeEach(async function () {
            const erc2TokenFactory = await ethers.getContractFactory('ERC20TokenTest')
            token1 = await erc2TokenFactory.deploy("Token1", "T", "18", ethers.parseEther(totalSupply))
            token2 = await erc2TokenFactory.deploy("Token2", "TT", "18", ethers.parseEther(totalSupply))
            token3 = await erc2TokenFactory.deploy("Token3", "TTT", "18", ethers.parseEther(totalSupply))
            token4 = await erc2TokenFactory.deploy("Token4", "TTTT", "18", ethers.parseEther(totalSupply))

            await fxBridge.addBridgeToken(await token1.getAddress(), encodeBytes32String(""), true)
            await fxBridge.addBridgeToken(await token2.getAddress(), encodeBytes32String(""), true)
            await fxBridge.addBridgeToken(await token3.getAddress(), encodeBytes32String(""), true)
            await fxBridge.addBridgeToken(await token4.getAddress(), encodeBytes32String(""), true)

            await token1.transferOwnership(await fxBridge.getAddress())
            await token2.transferOwnership(await fxBridge.getAddress())
            await token3.transferOwnership(await fxBridge.getAddress())
            await token4.transferOwnership(await fxBridge.getAddress())
        })

        it("submit bridge call 2 token", async function () {
            const tokens = [await token1.getAddress(), await token2.getAddress()]
            const amounts = ["1", "2"]
            const timeout = await ethers.provider.getBlockNumber() + 1000;

            await submitBridgeCall(tokens, amounts, timeout,0)
        })

        it("submit bridge call 3 token", async function () {
            const tokens = [await token1.getAddress(), await token2.getAddress(), await token3.getAddress()]
            const amounts = ["1", "2", "3"]
            const timeout = await ethers.provider.getBlockNumber() + 1000;

            await submitBridgeCall(tokens, amounts, timeout,0)
        })

        it("submit bridge call 4 token", async function () {
            const tokens = [await token1.getAddress(), await token2.getAddress(), await token3.getAddress(), await token4.getAddress()]
            const amounts = ["1", "2", "3", "4"]
            const timeout = await ethers.provider.getBlockNumber() + 1000;

            await submitBridgeCall(tokens, amounts, timeout,0)
        })
    })
})