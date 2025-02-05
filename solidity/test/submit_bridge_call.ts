import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import {
  BridgeCallContextTest,
  ERC20TokenTest,
  FxBridgeLogic,
} from "../typechain-types";
import { encodeBytes32String } from "ethers";
import { getSignerAddresses, submitBridgeCall } from "./common";

describe("submit bridge call tests", function () {
  let deploy: HardhatEthersSigner;
  let admin: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let receiver: HardhatEthersSigner;
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
    receiver = signers[3];

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
      await fxBridge.getAddress(),
      receiver.address
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
      false
    );
    await fxBridge.addBridgeToken(
      await token2.getAddress(),
      encodeBytes32String(""),
      false
    );

    await token1.transferOwnership(await fxBridge.getAddress());
    await token2.transferOwnership(await fxBridge.getAddress());
  });

  it("submit bridge call", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;

    await erc20Token.transfer(await fxBridge.getAddress(), amount);

    const initialBalance = await erc20Token.balanceOf(user1.address);

    await submitBridgeCall(
      gravityId,
      1,
      user1.address,
      user1.address,
      user1.address,
      "0x",
      "0x",
      [erc20TokenAddress],
      [amount],
      timeout,
      0,
      0,
      validators,
      powers,
      fxBridge
    );

    const finalBalance = await erc20Token.balanceOf(user1.address);
    expect(finalBalance.toString()).to.equal(initialBalance + BigInt(amount));
  });

  it("bridge call on bridge call", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const token1Address = await token1.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;
    const bridgeCallContextAddress = await bridgeCallContextTest.getAddress();

    await erc20Token.transfer(await fxBridge.getAddress(), amount);
    await token1.transfer(await fxBridge.getAddress(), amount);

    await submitBridgeCall(
      gravityId,
      1,
      user1.address,
      bridgeCallContextAddress,
      bridgeCallContextAddress,
      "0x",
      "0x",
      [erc20TokenAddress, token1Address],
      [amount, amount],
      timeout,
      0,
      0,
      validators,
      powers,
      fxBridge
    );

    const receiverBalance1 = await erc20Token.balanceOf(receiver.address);
    const receiverBalance2 = await token1.balanceOf(receiver.address);
    expect(receiverBalance1.toString()).to.equal(amount);
    expect(receiverBalance2.toString()).to.equal(amount);
  });

  it("bridge call onRevert", async function () {
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;
    const bridgeCallContextAddress = await bridgeCallContextTest.getAddress();
    const bridgeContractAddress = await fxBridge.getAddress();
    const token1Address = await token1.getAddress();
    const token2Address = await token2.getAddress();

    await token1.transfer(bridgeContractAddress, amount);
    await token2.transfer(bridgeContractAddress, amount);

    await bridgeCallContextTest.setBridgeCallParams(
      "",
      user1.address,
      [token1Address, token2Address],
      [amount, amount],
      bridgeCallContextAddress,
      "0x",
      0,
      300000,
      "0x"
    );

    await submitBridgeCall(
      gravityId,
      1,
      user1.address,
      bridgeCallContextAddress,
      bridgeCallContextAddress,
      "0x",
      "0x",
      [token1Address, token2Address],
      [amount, amount],
      timeout,
      0,
      1,
      validators,
      powers,
      fxBridge
    );

    const allowance1 = await token1.allowance(
      bridgeCallContextAddress,
      bridgeContractAddress
    );
    const allowance2 = await token2.allowance(
      bridgeCallContextAddress,
      bridgeContractAddress
    );
    expect(allowance1.toString()).to.equal(amount);
    expect(allowance2.toString()).to.equal(amount);

    const finalBalance1 = await token1.balanceOf(bridgeContractAddress);
    const finalBalance2 = await token2.balanceOf(bridgeContractAddress);
    expect(finalBalance1.toString()).to.equal("0");
    expect(finalBalance2.toString()).to.equal("0");
  });

  it("bridge call onRevert retry bridge call", async function () {
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;
    const bridgeCallContextAddress = await bridgeCallContextTest.getAddress();
    const bridgeContractAddress = await fxBridge.getAddress();
    const token1Address = await token1.getAddress();
    const token2Address = await token2.getAddress();

    await token1.transfer(bridgeContractAddress, amount);
    await token2.transfer(bridgeContractAddress, amount);

    await bridgeCallContextTest.setBridgeCallParams(
      "",
      user1.address,
      [token1Address, token2Address],
      [amount, amount],
      bridgeCallContextAddress,
      "0x",
      0,
      300000,
      "0x"
    );

    await bridgeCallContextTest.setRetryBridgeCall(true);

    await submitBridgeCall(
      gravityId,
      1,
      user1.address,
      bridgeCallContextAddress,
      bridgeCallContextAddress,
      "0x",
      "0x",
      [token1Address, token2Address],
      [amount, amount],
      timeout,
      0,
      1,
      validators,
      powers,
      fxBridge
    );

    const allowance1 = await token1.allowance(
      bridgeCallContextAddress,
      bridgeContractAddress
    );
    const allowance2 = await token2.allowance(
      bridgeCallContextAddress,
      bridgeContractAddress
    );
    expect(allowance1.toString()).to.equal("0");
    expect(allowance2.toString()).to.equal("0");

    const finalBalance1 = await token1.balanceOf(bridgeContractAddress);
    const finalBalance2 = await token2.balanceOf(bridgeContractAddress);
    expect(finalBalance1.toString()).to.equal(amount);
    expect(finalBalance2.toString()).to.equal(amount);
  });
});
