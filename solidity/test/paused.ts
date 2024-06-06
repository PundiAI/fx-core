import { beforeEach } from "mocha";
import { ethers } from "hardhat";
import { ERC20TokenTest, FxBridgeLogic } from "../typechain-types";
import { examplePowers, getSignerAddresses } from "./common";
import { encodeBytes32String } from "ethers";
import { expect } from "chai";

describe("bridge pause test", function () {
  let bridgeToken: ERC20TokenTest;
  let fxBridgeContract: FxBridgeLogic;

  beforeEach(async function () {
    const signers = await ethers.getSigners();
    const erc20Factory = await ethers.getContractFactory("ERC20TokenTest");
    bridgeToken = await erc20Factory.deploy(
      "ERC20 Token",
      "ERC20",
      "18",
      ethers.parseEther("100000")
    );
    await bridgeToken.waitForDeployment();

    const fxBridgeLogicFactory = await ethers.getContractFactory(
      "FxBridgeLogic"
    );
    const fxBridgeLogic = await fxBridgeLogicFactory.deploy();
    await fxBridgeLogic.waitForDeployment();
    const fxBridgeLogicAddress = await fxBridgeLogic.getAddress();

    const transparentUpgradeableProxyFactory = await ethers.getContractFactory(
      "TransparentUpgradeableProxy"
    );
    const proxyContract = await transparentUpgradeableProxyFactory.deploy(
      fxBridgeLogicAddress,
      signers[0].address,
      "0x"
    );
    await proxyContract.waitForDeployment();

    fxBridgeContract = <FxBridgeLogic>(
      fxBridgeLogicFactory.attach(await proxyContract.getAddress())
    );

    const powers: number[] = examplePowers();
    const validators = signers.slice(0, powers.length);
    const valAddresses = await getSignerAddresses(validators);

    const proxy = await ethers.getContractAt(
      "ITransparentUpgradeableProxy",
      await proxyContract.getAddress()
    );
    await proxy.connect(signers[0]).changeAdmin(signers[1].address);

    await fxBridgeContract.init(
      ethers.encodeBytes32String("fx-eth-bridge"),
      1000,
      valAddresses,
      powers
    );
    await fxBridgeContract.addBridgeToken(
      await bridgeToken.getAddress(),
      encodeBytes32String(""),
      true
    );
  });

  it("should pause and unpause", async function () {
    let paused = await fxBridgeContract.paused();
    expect(paused).to.be.false;

    const signers = await ethers.getSigners();
    await bridgeToken
      .connect(signers[0])
      .approve(await fxBridgeContract.getAddress(), ethers.parseEther("100"));
    await expect(
      fxBridgeContract
        .connect(signers[0])
        .sendToFx(
          await bridgeToken.getAddress(),
          ethers.encodeBytes32String("destination"),
          ethers.encodeBytes32String(""),
          ethers.parseEther("10")
        )
    ).to.be.emit(fxBridgeContract, "SendToFxEvent");

    await fxBridgeContract.pause();
    paused = await fxBridgeContract.paused();
    expect(paused).to.be.true;

    await expect(
      fxBridgeContract
        .connect(signers[0])
        .sendToFx(
          await bridgeToken.getAddress(),
          ethers.encodeBytes32String("destination"),
          ethers.encodeBytes32String(""),
          ethers.parseEther("10")
        )
    ).to.be.revertedWith("Pausable: paused");

    await fxBridgeContract.unpause();
    paused = await fxBridgeContract.paused();
    expect(paused).to.be.false;

    await expect(
      fxBridgeContract
        .connect(signers[0])
        .bridgeCall(
          "",
          signers[0].address,
          [await bridgeToken.getAddress()],
          [ethers.parseEther("10")],
          signers[0].address,
          "0x",
          0,
          "0x"
        )
    ).to.be.emit(fxBridgeContract, "BridgeCallEvent");

    await fxBridgeContract.pause();
    paused = await fxBridgeContract.paused();
    expect(paused).to.be.true;

    await expect(
      fxBridgeContract
        .connect(signers[0])
        .bridgeCall(
          "",
          signers[0].address,
          [await bridgeToken.getAddress()],
          [ethers.parseEther("10")],
          signers[0].address,
          "0x",
          0,
          "0x"
        )
    ).to.be.revertedWith("Pausable: paused");
  });
});
