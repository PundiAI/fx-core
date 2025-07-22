import { ethers } from "hardhat";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";
import { expect } from "chai";
import { PundiAIFX } from "../typechain-types";
import { it } from "mocha";

describe("pundiaifx tests", function () {
  let deploy: HardhatEthersSigner;
  let user1: HardhatEthersSigner;
  let pundiAIFX: PundiAIFX;

  beforeEach(async function () {
    const signers = await ethers.getSigners();
    deploy = signers[0];
    user1 = signers[1];

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
    await pundiAIFX.connect(deploy).initialize();
  });

  it("burn acc", async function () {
    const addr = ethers.Wallet.createRandom().address;
    await pundiAIFX.grantRole(await pundiAIFX.ADMIN_ROLE(), deploy.address);
    await pundiAIFX.mint(deploy.address, "20");
    await pundiAIFX.connect(deploy).transfer(addr, "20");
    await pundiAIFX.connect(deploy).addToBlacklist(addr);
    await pundiAIFX.connect(deploy).burnAcc(addr, "10");
    expect(await pundiAIFX.balanceOf(addr)).to.equal("10");

    await pundiAIFX.connect(deploy).pause();
    await pundiAIFX.connect(deploy).burnAcc(addr, "10");
    expect(await pundiAIFX.balanceOf(addr)).to.equal("0");
  });

  it("token info", async function () {
    expect(await pundiAIFX.name()).to.equal("Pundi AI");
    expect(await pundiAIFX.symbol()).to.equal("PUNDIAI");
    expect(await pundiAIFX.decimals()).to.equal(18);
    expect(await pundiAIFX.totalSupply()).to.equal(0);
    expect((await pundiAIFX.eip712Domain()).name).to.equal("Pundi AI");
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

  it("set name", async function () {
    expect(await pundiAIFX.name()).to.equal("Pundi AI");
    expect((await pundiAIFX.eip712Domain()).name).to.equal("Pundi AI");

    const newName = "xxxx";
    await pundiAIFX.connect(deploy).setName(newName);

    expect(await pundiAIFX.name()).to.equal(newName);
    expect((await pundiAIFX.eip712Domain()).name).to.equal(newName);
  });

  it("upgrade contract", async function () {
    const newPundiAIFXFactory = await ethers.getContractFactory("PundiAIFX");
    const newPundiAIFXDeploy = await newPundiAIFXFactory.deploy();

    await pundiAIFX
      .connect(deploy)
      .upgradeTo(await newPundiAIFXDeploy.getAddress());

    // Verify state preservation
    expect(await pundiAIFX.name()).to.equal("Pundi AI");
    expect(await pundiAIFX.symbol()).to.equal("PUNDIAI");
  });
});
