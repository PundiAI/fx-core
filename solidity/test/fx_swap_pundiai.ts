import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import { ERC20TokenTest, FXSwapPundiAI, PundiAIFX } from "../typechain-types";
import { it } from "mocha";

describe("pundiaifx tests", function () {
  let deploy: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let fxToken: ERC20TokenTest;
  let pundiAIFX: PundiAIFX;
  let fxSwapPundiAI: FXSwapPundiAI;
  let totalSupply = "100000000";

  beforeEach(async function () {
    const signers = await ethers.getSigners();
    deploy = signers[0];
    user1 = signers[1];

    const erc20TokenFactory = await ethers.getContractFactory("ERC20TokenTest");
    fxToken = await erc20TokenFactory
      .connect(deploy)
      .deploy("FX Token", "FX", "18", ethers.parseEther(totalSupply));

    const pundiAIFXFactory = await ethers.getContractFactory("PundiAIFX");
    const pundiAIFXDeploy = await pundiAIFXFactory.deploy();

    const pundiAIFXProxyFactory = await ethers.getContractFactory(
      "ERC1967Proxy"
    );
    const pundiAIFXProxy = await pundiAIFXProxyFactory
      .connect(deploy)
      .deploy(await pundiAIFXDeploy.getAddress(), "0x");

    pundiAIFX = await ethers.getContractAt(
      "PundiAIFX",
      await pundiAIFXProxy.getAddress()
    );
    await pundiAIFX.connect(deploy).initialize();

    const fxSwapPundiAIFactory = await ethers.getContractFactory(
      "FXSwapPundiAI"
    );
    const fxSwapPundiAIDeploy = await fxSwapPundiAIFactory.deploy();

    const fxSwapPundiAIProxyFactory = await ethers.getContractFactory(
      "ERC1967Proxy"
    );
    const fxSwapPundiAIProxy = await fxSwapPundiAIProxyFactory
      .connect(deploy)
      .deploy(await fxSwapPundiAIDeploy.getAddress(), "0x");

    fxSwapPundiAI = await ethers.getContractAt(
      "FXSwapPundiAI",
      await fxSwapPundiAIProxy.getAddress()
    );
    await fxSwapPundiAI
      .connect(deploy)
      .initialize(await fxToken.getAddress(), await pundiAIFX.getAddress());

    await fxToken
      .connect(deploy)
      .transfer(user1.address, ethers.parseEther(totalSupply));
    await fxToken
      .connect(user1)
      .approve(
        await fxSwapPundiAI.getAddress(),
        ethers.parseEther(totalSupply)
      );
    await pundiAIFX
      .connect(deploy)
      .grantRole(
        await pundiAIFX.ADMIN_ROLE(),
        await fxSwapPundiAI.getAddress()
      );
  });

  it("swap", async function () {
    expect(await fxToken.balanceOf(user1.address)).to.equal(
      ethers.parseEther(totalSupply)
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(0);
    expect(await fxToken.balanceOf(await fxSwapPundiAI.getAddress())).to.equal(
      0
    );
    expect(await pundiAIFX.totalSupply()).to.equal(0);

    await fxSwapPundiAI.connect(user1).swap(ethers.parseEther("100"));

    expect(await fxToken.balanceOf(user1.address)).to.equal(
      ethers.parseEther((Number(totalSupply) - 100).toString())
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(
      ethers.parseEther("1")
    );
    expect(await fxToken.balanceOf(await fxSwapPundiAI.getAddress())).to.equal(
      ethers.parseEther("100")
    );
    expect(await pundiAIFX.totalSupply()).to.equal(ethers.parseEther("1"));
  });

  it("burn FX", async function () {
    await fxToken
      .connect(user1)
      .transfer(await fxSwapPundiAI.getAddress(), ethers.parseEther("100"));
    expect(await fxToken.balanceOf(await fxSwapPundiAI.getAddress())).to.equal(
      ethers.parseEther("100")
    );

    await fxSwapPundiAI.burnFXToken(ethers.parseEther("100"));
    expect(await fxToken.balanceOf(await fxSwapPundiAI.getAddress())).to.equal(
      ethers.parseEther("0")
    );
    expect(await fxToken.totalSupply()).to.equal(
      ethers.parseEther((Number(totalSupply) - 100).toString())
    );
  });

  it("upgrade contract", async function () {
    const newFxSwapPundiAIFactory = await ethers.getContractFactory(
      "FXSwapPundiAI"
    );
    const newFxSwapPundiAIDeploy = await newFxSwapPundiAIFactory.deploy();

    await fxSwapPundiAI
      .connect(deploy)
      .upgradeTo(await newFxSwapPundiAIDeploy.getAddress());
  });
});
