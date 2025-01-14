import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import { ERC20TokenTest, PundiAIFX } from "../typechain-types";
import { it } from "mocha";

describe("pundiaifx tests", function () {
  let deploy: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let fxToken: ERC20TokenTest;
  let pundiAIFX: PundiAIFX;
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

    const erc1967ProxyFactory = await ethers.getContractFactory("ERC1967Proxy");
    const erc1967Proxy = await erc1967ProxyFactory
      .connect(deploy)
      .deploy(await pundiAIFXDeploy.getAddress(), "0x");

    pundiAIFX = await ethers.getContractAt(
      "PundiAIFX",
      await erc1967Proxy.getAddress()
    );
    await pundiAIFX.connect(deploy).initialize(await fxToken.getAddress());

    await fxToken
      .connect(deploy)
      .transfer(user1.address, ethers.parseEther(totalSupply));
    await fxToken
      .connect(user1)
      .approve(await pundiAIFX.getAddress(), ethers.parseEther(totalSupply));
  });

  it("Pundi AIFX", async function () {
    expect(await pundiAIFX.name()).to.equal("Pundi AIFX Token");
    expect(await pundiAIFX.symbol()).to.equal("PUNDIAI");
    expect(await pundiAIFX.decimals()).to.equal(18);
    expect(await pundiAIFX.totalSupply()).to.equal(0);
  });

  it("mint, burn and  transfer AIFX", async function () {
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(0);
    expect(await pundiAIFX.totalSupply()).to.equal(0);

    await pundiAIFX.grantRole(await pundiAIFX.ADMIN_ROLE(), deploy.address);
    await pundiAIFX.mint(deploy.address, "1");

    await pundiAIFX.approve(user1.address, "1");
    await pundiAIFX
      .connect(user1)
      .transferFrom(deploy.address, user1.address, "1");

    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(1);
    expect(await pundiAIFX.totalSupply()).to.equal(1);

    await pundiAIFX.connect(user1).transfer(deploy.address, "1");

    expect(await pundiAIFX.balanceOf(deploy.address)).to.equal(1);

    await pundiAIFX.burn("1");

    expect(await pundiAIFX.totalSupply()).to.equal(0);
  });

  it("check role", async function () {
    expect(await pundiAIFX.getRoleAdmin(await pundiAIFX.OWNER_ROLE())).to.equal(
      await pundiAIFX.DEFAULT_ADMIN_ROLE()
    );
    expect(await pundiAIFX.getRoleAdmin(await pundiAIFX.ADMIN_ROLE())).to.equal(
      await pundiAIFX.DEFAULT_ADMIN_ROLE()
    );

    expect(
      await pundiAIFX.hasRole(
        await pundiAIFX.DEFAULT_ADMIN_ROLE(),
        deploy.address
      )
    ).to.equal(true);
    expect(
      await pundiAIFX.hasRole(await pundiAIFX.OWNER_ROLE(), deploy.address)
    ).to.equal(true);

    expect(
      await pundiAIFX.hasRole(await pundiAIFX.ADMIN_ROLE(), user1.address)
    ).to.equal(false);

    await pundiAIFX.grantRole(await pundiAIFX.ADMIN_ROLE(), user1.address);

    expect(
      await pundiAIFX.hasRole(await pundiAIFX.ADMIN_ROLE(), user1.address)
    ).to.equal(true);

    await pundiAIFX.revokeRole(await pundiAIFX.ADMIN_ROLE(), user1.address);

    expect(
      await pundiAIFX.hasRole(await pundiAIFX.ADMIN_ROLE(), user1.address)
    ).to.equal(false);
  });

  it("swap FX", async function () {
    expect(await fxToken.balanceOf(user1.address)).to.equal(
      ethers.parseEther(totalSupply)
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(0);
    expect(await fxToken.balanceOf(await pundiAIFX.getAddress())).to.equal(0);
    expect(await pundiAIFX.totalSupply()).to.equal(0);

    await pundiAIFX.connect(user1).swap(ethers.parseEther("100"));

    expect(await fxToken.balanceOf(user1.address)).to.equal(
      ethers.parseEther((Number(totalSupply) - 100).toString())
    );
    expect(await pundiAIFX.balanceOf(user1.address)).to.equal(
      ethers.parseEther("1")
    );
    expect(await fxToken.balanceOf(await pundiAIFX.getAddress())).to.equal(
      ethers.parseEther("100")
    );
    expect(await pundiAIFX.totalSupply()).to.equal(ethers.parseEther("1"));
  });
});
