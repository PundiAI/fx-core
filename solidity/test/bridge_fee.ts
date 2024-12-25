import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import {
  ERC20TokenTest,
  FxBridgeLogic,
  BridgeCallContextTest,
} from "../typechain-types";
import { AbiCoder, encodeBytes32String } from "ethers";
import { getSignerAddresses, submitBridgeCall } from "./common";

describe("submit bridge call tests", function () {
  let deploy: HardhatEthersSigner;
  let admin: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let erc20Token: ERC20TokenTest;
  let fxBridge: FxBridgeLogic;
  let bridgeCallContextTest: BridgeCallContextTest;

  let token1: ERC20TokenTest;
  let token2: ERC20TokenTest;

  let totalSupply = "10000";
  const gravityId: string = encodeBytes32String("eth-fxcore");
  const powerThreshold = 6666;
  const powers: number[] = [
    1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000,
  ];

  let validators: HardhatEthersSigner[];
  let valAddresses: string[];

  beforeEach(async function () {
    const signers = await ethers.getSigners();
    deploy = signers[0];
    admin = signers[1];
    user1 = signers[2];

    validators = [
      signers[3],
      signers[4],
      signers[5],
      signers[6],
      signers[7],
      signers[8],
      signers[9],
      signers[10],
      signers[11],
      signers[12],
    ];
    valAddresses = await getSignerAddresses(validators);

    const erc20TokenFactory = await ethers.getContractFactory("ERC20TokenTest");
    erc20Token = await erc20TokenFactory.deploy(
      "ERC20 Token",
      "ERC20",
      "18",
      ethers.parseEther(totalSupply)
    );
    const erc20TokenAddress = await erc20Token.getAddress();
    expect(await erc20Token.balanceOf(deploy.address)).to.equal(
      ethers.parseEther("10000")
    );

    const fxBridgeLogicFactory = await ethers.getContractFactory(
      "FxBridgeLogic"
    );
    const fxBridgeLogic = await fxBridgeLogicFactory.deploy();
    const fxBridgeLogicAddress = await fxBridgeLogic.getAddress();

    const transparentUpgradeableProxyFactory = await ethers.getContractFactory(
      "TransparentUpgradeableProxy"
    );
    const fxBridgeLogicProxy = await transparentUpgradeableProxyFactory.deploy(
      fxBridgeLogicAddress,
      admin.address,
      "0x"
    );
    const fxBridgeLogicProxyAddress = await fxBridgeLogicProxy.getAddress();

    fxBridge = <FxBridgeLogic>(
      fxBridgeLogicFactory.attach(fxBridgeLogicProxyAddress)
    );

    await fxBridge.init(gravityId, powerThreshold, valAddresses, powers);
    await fxBridge.addBridgeToken(
      erc20TokenAddress,
      encodeBytes32String(""),
      true
    );

    await erc20Token.transferOwnership(await fxBridge.getAddress());

    const bridgeCallContextTestFactory = await ethers.getContractFactory(
      "BridgeCallContextTest"
    );
    bridgeCallContextTest = await bridgeCallContextTestFactory.deploy(
      await fxBridge.getAddress()
    );

    const erc2TokenFactory = await ethers.getContractFactory("ERC20TokenTest");
    token1 = await erc2TokenFactory.deploy(
      "Token1",
      "T",
      "18",
      ethers.parseEther(totalSupply)
    );
    token2 = await erc2TokenFactory.deploy(
      "Token2",
      "TT",
      "18",
      ethers.parseEther(totalSupply)
    );

    await fxBridge.addBridgeToken(
      await token1.getAddress(),
      encodeBytes32String(""),
      true
    );
    await fxBridge.addBridgeToken(
      await token2.getAddress(),
      encodeBytes32String(""),
      true
    );

    await token1.transferOwnership(await fxBridge.getAddress());
    await token2.transferOwnership(await fxBridge.getAddress());
  });

  it("bridge call fee", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;
    const bridgeCallContextAddress = await bridgeCallContextTest.getAddress();

    await erc20Token.transfer(
      await fxBridge.getAddress(),
      ethers.parseEther("1")
    );

    await token1.transfer(await fxBridge.getAddress(), ethers.parseEther("1"));

    const memo = new AbiCoder().encode(
      ["address", "bytes"],
      [await bridgeCallContextTest.getAddress(), "0x"]
    );

    const deployBal1 = await token1.balanceOf(bridgeCallContextAddress);
    const callFlag1 = await bridgeCallContextTest.callFlag();
    await submitBridgeCall(
      gravityId,
      user1.address,
      bridgeCallContextAddress,
      bridgeCallContextAddress,
      "0x",
      memo,
      [erc20TokenAddress, await token1.getAddress()],
      [amount, amount],
      1,
      timeout,
      0,
      validators,
      powers,
      fxBridge
    );

    const deployBal2 = await token1.balanceOf(bridgeCallContextAddress);
    expect(deployBal2).to.be.equal(deployBal1 + BigInt(amount));
    const callFlag2 = await bridgeCallContextTest.callFlag();
    expect(callFlag2).to.be.equal(!callFlag1);
  });
});
