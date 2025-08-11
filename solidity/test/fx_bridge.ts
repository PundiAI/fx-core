import { ethers } from "hardhat";
import { expect } from "chai";
import { FxBridgeLogic } from "../typechain-types";
import { AbiCoder, TransactionRequest } from "ethers";

describe("fork network and fx bridge test", function () {
  let gasAddress: string;
  let bridgeAddress: string;
  let adminAddress: string;
  let bridgeContractName: string;
  let expectBridgeCallStatus: boolean;
  const abiCode = new AbiCoder();

  beforeEach(async function () {
    if (!process.env.FORK_ENABLE) {
      return;
    }
    const network = await ethers.provider.getNetwork();
    switch (network.chainId.toString()) {
      case "1":
        gasAddress = "0x00000000219ab540356cBB839Cbe05303d7705Fa";
        bridgeAddress = "0x6f1D09Fed11115d65E1071CD2109eDb300D80A27";
        adminAddress = "0x0F413055AdEF9b61e9507928c6856F438d690882";
        bridgeContractName = "FxBridgeLogicETH";
        expectBridgeCallStatus = false;
        break;
      case "11155111":
        gasAddress = "0x6Cc9397c3B38739daCbfaA68EaD5F5D77Ba5F455";
        bridgeAddress = "0xd384a8e8822Ea845e83eb5AA2877239150615C18";
        adminAddress = "0xcF8049f0B918650614D5bf18CF15af080eFdDEe1";
        bridgeContractName = "FxBridgeLogicETH";
        expectBridgeCallStatus = true;
        break;
      case "84532":
        gasAddress = "0x4200000000000000000000000000000000000016";
        bridgeAddress = "0x9164D153b8Af6D94d41E7876E814DD3Db1AEC320";
        adminAddress = "0xcF8049f0B918650614D5bf18CF15af080eFdDEe1";
        bridgeContractName = "FxBridgeLogic";
        expectBridgeCallStatus = true;
        break;
      case "8453":
        gasAddress = "0x4200000000000000000000000000000000000016";
        bridgeAddress = "0x7a986bA67227acfab86385FF33436a80E2BB4CC5";
        adminAddress = "0xE77A7EA2F1DC25968b5941a456d99D37b80E98B5";
        bridgeContractName = "FxBridgeLogic";
        expectBridgeCallStatus = false;
        break;
      default:
        throw new Error("Unsupported network");
    }
  });

  it("upgrade bridge contract", async function () {
    if (!process.env.FORK_ENABLE) {
      return;
    }
    const gasSigner = await ethers.getImpersonatedSigner(gasAddress);
    await gasSigner.sendTransaction({
      to: adminAddress,
      value: ethers.parseEther("100"),
    });
    const adminSigner = await ethers.getImpersonatedSigner(adminAddress);

    const bridgeFactory = await ethers.getContractFactory(bridgeContractName);

    const bridgeContract = bridgeFactory.attach(bridgeAddress) as FxBridgeLogic;

    const oldFxBridgeId = await bridgeContract.state_fxBridgeId();
    const oldPowerThreshold = await bridgeContract.state_powerThreshold();
    const oldLastEventNonce = await bridgeContract.state_lastEventNonce();
    const oldCheckpoint = await bridgeContract.state_lastOracleSetCheckpoint();
    const oldOracleSetNonce = await bridgeContract.state_lastOracleSetNonce();
    const oldBridgeTokens = await bridgeContract.getBridgeTokenList();

    let oldTokenStatus = new Map();
    for (const bridgeToken of oldBridgeTokens) {
      const status = await bridgeContract.tokenStatus(bridgeToken.addr);
      const batchNonce = await bridgeContract.state_lastBatchNonces(
        bridgeToken.addr
      );
      oldTokenStatus.set(bridgeToken.addr.toString(), {
        batchNonce: batchNonce,
        status: status,
      });
    }

    const bridgeLogicContract = await bridgeFactory.deploy();
    await bridgeLogicContract.waitForDeployment();
    const bridgeLogicContractAddress = await bridgeLogicContract.getAddress();

    // 0x3659cfe6 is the signature of the upgradeTo(address) function
    const data = ethers.concat([
      "0x3659cfe6",
      abiCode.encode(["address"], [bridgeLogicContractAddress]),
    ]);

    const transaction: TransactionRequest = {
      to: bridgeAddress,
      data: data,
    };

    const upgradeTx = await adminSigner.sendTransaction(transaction);
    await upgradeTx.wait();

    const fxBridgeId = await bridgeContract.state_fxBridgeId();
    const powerThreshold = await bridgeContract.state_powerThreshold();
    const lastEventNonce = await bridgeContract.state_lastEventNonce();
    const checkpoint = await bridgeContract.state_lastOracleSetCheckpoint();
    const oracleSetNonce = await bridgeContract.state_lastOracleSetNonce();
    const bridgeTokens = await bridgeContract.getBridgeTokenList();

    expect(fxBridgeId).to.equal(oldFxBridgeId);
    expect(powerThreshold.toString()).to.equal(oldPowerThreshold.toString());
    expect(lastEventNonce.toString()).to.equal(oldLastEventNonce.toString());
    expect(checkpoint).to.equal(oldCheckpoint);
    expect(oracleSetNonce.toString()).to.equal(oldOracleSetNonce.toString());

    for (const bridgeToken of bridgeTokens) {
      const status = await bridgeContract.tokenStatus(bridgeToken.addr);
      expect(status.isOriginated).to.equal(
        oldTokenStatus.get(bridgeToken.addr).status.isOriginated
      );
      expect(status.isActive).to.equal(
        oldTokenStatus.get(bridgeToken.addr).status.isActive
      );
      expect(status.isExist).to.equal(
        oldTokenStatus.get(bridgeToken.addr).status.isExist
      );
      const batchNonce = await bridgeContract.state_lastBatchNonces(
        bridgeToken.addr
      );
      expect(batchNonce.toString()).to.equal(
        oldTokenStatus.get(bridgeToken.addr).batchNonce.toString()
      );
      expect(await bridgeContract.state_lastBridgeCallNonces(1)).to.equal(
        expectBridgeCallStatus
      );
    }
  }).timeout(100000);
});

describe("fx bridge submitBatch", function () {
  let hacker = "0x26bC046BFA81ff9F38d0c701D456BfDf34b7F69c";
  let bridgeWrapperContract: any;
  let erc20Address: string;
  let erc20Contract: any;
  let signerAddress: string;
  beforeEach(async function () {
    const signer = (await ethers.getSigners())[0];
    signerAddress = await signer.getAddress();

    const bridgeTest = await ethers.getContractFactory("FxBridgeTest");
    const bridgeTestContract = await bridgeTest.deploy();
    const bridgeAddress = await bridgeTestContract.getAddress();

    const erc20 = await ethers.getContractFactory("ERC20TokenTest");
    erc20Contract = await erc20.deploy("test", "TEST", 18, 0);
    erc20Address = await erc20Contract.getAddress();
    await erc20Contract.mint(bridgeAddress, ethers.parseEther("100"));

    const bridgeWrapper = await ethers.getContractFactory("FxBridgeWrapper");
    bridgeWrapperContract = await bridgeWrapper.deploy(bridgeAddress);
  });
  it("should user success", async () => {
    const receiveAddress = ethers.Wallet.createRandom().address;
    await bridgeWrapperContract.submitBatch(
      [],
      [],
      [],
      [],
      [],
      [ethers.parseEther("1")],
      [receiveAddress],
      [ethers.parseEther("1")],
      [0, 0],
      erc20Address,
      0,
      signerAddress
    );
    expect(await erc20Contract.balanceOf(receiveAddress)).to.equal(
      ethers.parseEther("1")
    );
    expect(await erc20Contract.balanceOf(signerAddress)).to.equal(
      ethers.parseEther("1")
    );
  });

  it("should hacker failed", async () => {
    const zeroAddr = "0x0000000000000000000000000000000000000000";
    await expect(
      bridgeWrapperContract.submitBatch(
        [],
        [],
        [],
        [],
        [],
        [ethers.parseEther("1")],
        [hacker],
        [0],
        [0, 0],
        erc20Address,
        0,
        zeroAddr
      )
    ).to.be.rejectedWith("Balance mismatch after batch submission");
    expect(await erc20Contract.balanceOf(hacker)).to.equal(
      ethers.parseEther("0")
    );
    expect(await erc20Contract.balanceOf(signerAddress)).to.equal(
      ethers.parseEther("0")
    );
    expect(await erc20Contract.balanceOf(zeroAddr)).to.equal(
      ethers.parseEther("0")
    );
  });

  it("should hacker success", async () => {
    await bridgeWrapperContract.submitBatch(
      [],
      [],
      [],
      [],
      [],
      [ethers.parseEther("1")],
      [hacker],
      [ethers.parseEther("1")],
      [0, 0],
      erc20Address,
      0,
      signerAddress
    );
    expect(await erc20Contract.balanceOf(signerAddress)).to.equal(
      ethers.parseEther("2")
    );
    expect(await erc20Contract.balanceOf(hacker)).to.equal(
      ethers.parseEther("0")
    );
  });
});
