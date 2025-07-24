import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import { ERC20TokenTest, FXtoPUNDIAISwap, PundiAIFX } from "../typechain-types";
import { it } from "mocha";

describe("FXtoPUNDIAISwap tests", function () {
  let deploy: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let fxToken: ERC20TokenTest;
  let pundiAIFX: PundiAIFX;
  let fxToPundiAISwap: FXtoPUNDIAISwap;
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

    const fxToPundiAISwapFactory = await ethers.getContractFactory(
      "FXtoPUNDIAISwap"
    );
    fxToPundiAISwap = await fxToPundiAISwapFactory.deploy(
      await fxToken.getAddress(),
      await pundiAIFX.getAddress()
    );

    await pundiAIFX
      .connect(deploy)
      .grantRole(
        await pundiAIFX.ADMIN_ROLE(),
        await fxToPundiAISwap.getAddress()
      );
  });

  it("swap", async function () {
    await fxToken
      .connect(deploy)
      .transfer(user1.address, ethers.parseEther(totalSupply));
    expect(await fxToken.balanceOf(user1.address)).to.equal(
      ethers.parseEther(totalSupply)
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(0);
    expect(
      await fxToken.balanceOf(await fxToPundiAISwap.getAddress())
    ).to.equal(0);
    expect(await pundiAIFX.totalSupply()).to.equal(0);

    await fxToPundiAISwap.connect(user1).swap();

    expect(await fxToken.balanceOf(user1.address)).to.equal(
      ethers.parseEther(totalSupply)
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
    expect(
      await fxToken.balanceOf(await fxToPundiAISwap.getAddress())
    ).to.equal(ethers.parseEther("0"));
    expect(await pundiAIFX.totalSupply()).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
    expect(await fxToPundiAISwap.totalMinted()).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
    expect(await fxToPundiAISwap.MAX_TOTAL_MINT()).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
  });

  it("swapFor", async function () {
    let contractAddress = await fxToken.getAddress();
    await fxToken
      .connect(deploy)
      .transfer(contractAddress, ethers.parseEther(totalSupply));
    expect(await fxToken.balanceOf(contractAddress)).to.equal(
      ethers.parseEther(totalSupply)
    );
    expect(await pundiAIFX.balanceOf(contractAddress)).to.equal(0);
    expect(
      await fxToken.balanceOf(await fxToPundiAISwap.getAddress())
    ).to.equal(0);
    expect(await pundiAIFX.totalSupply()).to.equal(0);

    await fxToPundiAISwap.connect(deploy).swapFor(contractAddress, user1);

    expect(await fxToken.balanceOf(contractAddress)).to.equal(
      ethers.parseEther(totalSupply)
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
    expect(
      await fxToken.balanceOf(await fxToPundiAISwap.getAddress())
    ).to.equal(ethers.parseEther("0"));
    expect(await pundiAIFX.totalSupply()).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
    expect(await fxToPundiAISwap.totalMinted()).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
    expect(await fxToPundiAISwap.MAX_TOTAL_MINT()).to.equal(
      ethers.parseEther((Number(totalSupply) / 100).toString())
    );
  });
});
