import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import { ERC20TokenTest, FxBridgeLogic } from "../typechain-types";
import { BigNumberish, encodeBytes32String } from "ethers";
import { it } from "mocha";
import {
  encodeFunctionData,
  examplePowers,
  getSignerAddresses,
} from "./common";

describe("bridge call tests", function () {
  let deploy: HardhatEthersSigner;
  let admin: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let erc20Token: ERC20TokenTest;
  let fxBridge: FxBridgeLogic;

  let totalSupply = "10000";
  const gravityId: string = encodeBytes32String("eth-fxcore");
  const powerThreshold = 6666;

  beforeEach(async function () {
    const signers = await ethers.getSigners();
    deploy = signers[0];
    admin = signers[1];
    user1 = signers[2];

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

    const powers: number[] = examplePowers();
    const validators = signers.slice(0, powers.length);
    const valAddresses = await getSignerAddresses(validators);

    await fxBridge.init(gravityId, powerThreshold, valAddresses, powers);
    await fxBridge.addBridgeToken(
      erc20TokenAddress,
      encodeBytes32String(""),
      true
    );
  });

  const fxcoreChainId = "";

  it("bridge call erc20", async function () {
    const tokens = [await erc20Token.getAddress()];
    const amount = ethers.parseEther("1");
    const amounts: BigNumberish[] = [amount];

    await erc20Token.approve(await fxBridge.getAddress(), amount);
    const lastEventNonce = await fxBridge.state_lastEventNonce();

    await fxBridge.bridgeCall(
      fxcoreChainId,
      user1.address,
      tokens,
      amounts,
      user1.address,
      "0x",
      0,
      "0x"
    );

    const lastEventNonceAfter = await fxBridge.state_lastEventNonce();
    expect(lastEventNonceAfter).to.equal(BigInt(Number(lastEventNonce) + 1));

    const balance1 = await erc20Token.balanceOf(deploy.address);
    expect(balance1).to.equal(
      ethers.parseEther((Number(totalSupply) - 1).toString())
    );
    const balance3 = await erc20Token.balanceOf(await fxBridge.getAddress());
    expect(balance3).to.equal(BigInt(0));
    const newTotalSupply = await erc20Token.totalSupply();
    expect(newTotalSupply).to.equal(
      ethers.parseEther((Number(totalSupply) - 1).toString())
    );
  });

  it("bridge call erc20 with value", async function () {
    const tokens = [await erc20Token.getAddress()];
    const amount = ethers.parseEther("1");
    const amounts: BigNumberish[] = [amount];
    const quoteId = 1;
    const lastEventNonce = await fxBridge.state_lastEventNonce();

    await erc20Token.approve(await fxBridge.getAddress(), amount);

    await expect(
      fxBridge.bridgeCall(
        fxcoreChainId,
        user1.address,
        tokens,
        amounts,
        user1.address,
        "0x",
        quoteId,
        "0x"
      )
    )
      .to.be.emit(fxBridge, "BridgeCallEvent")
      .withArgs(
        deploy.address,
        user1.address,
        user1.address,
        deploy.address,
        BigInt(Number(lastEventNonce) + 1),
        fxcoreChainId,
        tokens,
        amounts,
        "0x",
        quoteId,
        "0x"
      );
  });

  it("bridge call erc20 with data", async function () {
    const tokens = [await erc20Token.getAddress()];
    const amount = ethers.parseEther("1");
    const amounts: BigNumberish[] = [amount];
    const quoteId = 1;
    const lastEventNonce = await fxBridge.state_lastEventNonce();
    const data = encodeFunctionData(
      erc20Token.interface.formatJson(),
      "transfer",
      [deploy.address, amount]
    );
    await erc20Token.approve(await fxBridge.getAddress(), amount);

    await expect(
      fxBridge.bridgeCall(
        fxcoreChainId,
        user1.address,
        tokens,
        amounts,
        user1.address,
        data,
        quoteId,
        "0x"
      )
    )
      .to.be.emit(fxBridge, "BridgeCallEvent")
      .withArgs(
        deploy.address,
        user1.address,
        user1.address,
        deploy.address,
        BigInt(Number(lastEventNonce) + 1),
        fxcoreChainId,
        tokens,
        amounts,
        data,
        quoteId,
        "0x"
      );
  });

  it("bridge call erc20 with memo", async function () {
    const tokens = [await erc20Token.getAddress()];
    const amount = ethers.parseEther("1");
    const amounts: BigNumberish[] = [amount];
    const quoteId = 1;
    const lastEventNonce = await fxBridge.state_lastEventNonce();
    const data = encodeFunctionData(
      erc20Token.interface.formatJson(),
      "transfer",
      [deploy.address, amount]
    );
    const demo = encodeFunctionData(
      '[{"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"delegate","outputs":[],"stateMutability":"nonpayable","type":"function"}]',
      "delegate",
      [user1.address]
    );

    await erc20Token.approve(await fxBridge.getAddress(), amount);

    await expect(
      fxBridge.bridgeCall(
        fxcoreChainId,
        user1.address,
        tokens,
        amounts,
        user1.address,
        data,
        quoteId,
        demo
      )
    )
      .to.be.emit(fxBridge, "BridgeCallEvent")
      .withArgs(
        deploy.address,
        user1.address,
        user1.address,
        deploy.address,
        BigInt(Number(lastEventNonce) + 1),
        fxcoreChainId,
        tokens,
        amounts,
        data,
        quoteId,
        demo
      );
  });

  describe("bridge call batch transfer test", function () {
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
        false
      );
      await fxBridge.addBridgeToken(
        await token2.getAddress(),
        encodeBytes32String(""),
        false
      );
      await fxBridge.addBridgeToken(
        await token3.getAddress(),
        encodeBytes32String(""),
        false
      );
      await fxBridge.addBridgeToken(
        await token4.getAddress(),
        encodeBytes32String(""),
        false
      );

      await token1.approve(await fxBridge.getAddress(), totalSupply);
      await token2.approve(await fxBridge.getAddress(), totalSupply);
      await token3.approve(await fxBridge.getAddress(), totalSupply);
      await token4.approve(await fxBridge.getAddress(), totalSupply);
    });

    it("bridge call transfer 2 token", async function () {
      const tokens = [await token1.getAddress(), await token2.getAddress()];
      const amounts = [BigInt(1), BigInt(2)];

      await fxBridge.bridgeCall(
        fxcoreChainId,
        user1.address,
        tokens,
        amounts,
        user1.address,
        "0x",
        0,
        "0x"
      );

      const balance1 = await token1.balanceOf(deploy.address);
      const balance2 = await token2.balanceOf(deploy.address);
      expect(balance1).to.equal(ethers.parseEther(totalSupply) - BigInt(1));
      expect(balance2).to.equal(ethers.parseEther(totalSupply) - BigInt(2));

      const balance3 = await token1.balanceOf(await fxBridge.getAddress());
      const balance4 = await token2.balanceOf(await fxBridge.getAddress());
      expect(balance3).to.equal(BigInt(1));
      expect(balance4).to.equal(BigInt(2));
    });

    it("bridge call transfer 3 token", async function () {
      const tokens = [
        await token1.getAddress(),
        await token2.getAddress(),
        await token3.getAddress(),
      ];
      const amounts = [BigInt(1), BigInt(2), BigInt(3)];

      await fxBridge.bridgeCall(
        fxcoreChainId,
        user1.address,
        tokens,
        amounts,
        user1.address,
        "0x",
        0,
        "0x"
      );
    });

    it("bridge call transfer 4 token", async function () {
      const tokens = [
        await token1.getAddress(),
        await token2.getAddress(),
        await token3.getAddress(),
        await token4.getAddress(),
      ];
      const amounts = [BigInt(1), BigInt(2), BigInt(3), BigInt(4)];

      await fxBridge.bridgeCall(
        fxcoreChainId,
        user1.address,
        tokens,
        amounts,
        user1.address,
        "0x",
        0,
        "0x"
      );
    });
  });
});
