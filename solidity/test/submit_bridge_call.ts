import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import {
  ERC20TokenTest,
  FxBridgeLogic,
  BridgeCallContextTest,
} from "../typechain-types";
import { encodeBytes32String } from "ethers";
import { getSignerAddresses, submitBridgeCall } from "./common";

describe("submit bridge call tests", function () {
  let deploy: HardhatEthersSigner;
  let admin: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let erc20Token: ERC20TokenTest;
  let fxBridge: FxBridgeLogic;
  let bridgeCallContextTest: BridgeCallContextTest;

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
  });

  it("submit bridge call", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;

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
      1,
      0,
      validators,
      powers,
      fxBridge
    );
  });

  it("submit bridge call with bridge context on bridge call", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;
    const bridgeCallContextAddress = await bridgeCallContextTest.getAddress();

    await erc20Token.transfer(
      await fxBridge.getAddress(),
      ethers.parseEther("1")
    );

    const ownerBal1 = await erc20Token.balanceOf(bridgeCallContextAddress);
    await submitBridgeCall(
      gravityId,
      1,
      user1.address,
      user1.address,
      bridgeCallContextAddress,
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
    const ownerBal2 = await erc20Token.balanceOf(bridgeCallContextAddress);
    expect(ownerBal2).to.equal(ownerBal1 + BigInt(amount));
    expect(await bridgeCallContextTest.callFlag()).to.equal(true);
    expect(await bridgeCallContextTest.revertFlag()).to.equal(false);
  });

  it("submit bridge call with bridge context on revert", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;
    const bridgeCallContextAddress = await bridgeCallContextTest.getAddress();

    await erc20Token.transfer(
      await fxBridge.getAddress(),
      ethers.parseEther("1")
    );

    const ownerBal1 = await erc20Token.balanceOf(user1.address);
    await submitBridgeCall(
      gravityId,
      1,
      user1.address,
      user1.address,
      bridgeCallContextAddress,
      "0x",
      "0x",
      [erc20TokenAddress],
      [amount],
      timeout,
      0,
      1,
      validators,
      powers,
      fxBridge
    );
    const ownerBal2 = await erc20Token.balanceOf(user1.address);
    expect(ownerBal2).to.equal(ownerBal1);
    expect(await bridgeCallContextTest.revertFlag()).to.equal(true);
    expect(await bridgeCallContextTest.callFlag()).to.equal(false);
  });

  it("submit bridge call with refund", async function () {
    const erc20TokenAddress = await erc20Token.getAddress();
    const amount = "1000";
    const timeout = (await ethers.provider.getBlockNumber()) + 1000;

    await erc20Token.transfer(
      await fxBridge.getAddress(),
      ethers.parseEther("1")
    );

    const ownerBal1 = await erc20Token.balanceOf(user1.address);
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
      1,
      validators,
      powers,
      fxBridge
    );
    const ownerBal2 = await erc20Token.balanceOf(user1.address);
    expect(ownerBal2).to.equal(ownerBal1 + amount);
  });

  describe("submit bridge call batch test", function () {
    let token1: ERC20TokenTest;
    let token2: ERC20TokenTest;
    let token3: ERC20TokenTest;
    let token4: ERC20TokenTest;

    beforeEach(async function () {
      const erc2TokenFactory = await ethers.getContractFactory(
        "ERC20TokenTest"
      );
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
      token3 = await erc2TokenFactory.deploy(
        "Token3",
        "TTT",
        "18",
        ethers.parseEther(totalSupply)
      );
      token4 = await erc2TokenFactory.deploy(
        "Token4",
        "TTTT",
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
      await fxBridge.addBridgeToken(
        await token3.getAddress(),
        encodeBytes32String(""),
        true
      );
      await fxBridge.addBridgeToken(
        await token4.getAddress(),
        encodeBytes32String(""),
        true
      );

      await token1.transferOwnership(await fxBridge.getAddress());
      await token2.transferOwnership(await fxBridge.getAddress());
      await token3.transferOwnership(await fxBridge.getAddress());
      await token4.transferOwnership(await fxBridge.getAddress());
    });

    it("submit bridge call 2 token", async function () {
      const tokens = [await token1.getAddress(), await token2.getAddress()];
      const amounts = ["1", "2"];
      const timeout = (await ethers.provider.getBlockNumber()) + 1000;

      await submitBridgeCall(
        gravityId,
        1,
        user1.address,
        user1.address,
        user1.address,
        "0x",
        "0x",
        tokens,
        amounts,
        timeout,
        0,
        0,
        validators,
        powers,
        fxBridge
      );
    });

    it("submit bridge call 3 token", async function () {
      const tokens = [
        await token1.getAddress(),
        await token2.getAddress(),
        await token3.getAddress(),
      ];
      const amounts = ["1", "2", "3"];
      const timeout = (await ethers.provider.getBlockNumber()) + 1000;

      await submitBridgeCall(
        gravityId,
        1,
        user1.address,
        user1.address,
        user1.address,
        "0x",
        "0x",
        tokens,
        amounts,
        timeout,
        0,
        0,
        validators,
        powers,
        fxBridge
      );
    });

    it("submit bridge call 4 token", async function () {
      const tokens = [
        await token1.getAddress(),
        await token2.getAddress(),
        await token3.getAddress(),
        await token4.getAddress(),
      ];
      const amounts = ["1", "2", "3", "4"];
      const timeout = (await ethers.provider.getBlockNumber()) + 1000;

      await submitBridgeCall(
        gravityId,
        1,
        user1.address,
        user1.address,
        user1.address,
        "0x",
        "0x",
        tokens,
        amounts,
        timeout,
        0,
        0,
        validators,
        powers,
        fxBridge
      );
    });
  });
});
