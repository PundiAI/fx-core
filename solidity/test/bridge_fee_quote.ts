import { ethers } from "hardhat";
import { expect } from "chai";
import { AddressLike, HDNodeWallet, Wallet } from "ethers";
import {
  BridgeFeeOracle,
  BridgeFeeQuote,
  BridgeFeeQuoteTest,
  IBridgeFeeQuote,
} from "../typechain-types";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";

const messagePrefix = "\x19Ethereum Signed Message:\n32";

describe("BridgeFeeQuoteUpgradeable", function () {
  let bridgeFeeQuote: BridgeFeeQuote;
  let bridgeFeeQuoteTest: BridgeFeeQuoteTest;
  let bridgeFeeOracle: BridgeFeeOracle;
  let chainName = "TestChain";
  let oracle: HDNodeWallet;
  let owner: HardhatEthersSigner;
  let token1: any;
  let token2: any;
  let tokens: AddressLike[];

  beforeEach(async function () {
    oracle = Wallet.createRandom();

    [owner, token1, token2] = await ethers.getSigners();

    const BridgeFeeQuoteTest = await ethers.getContractFactory(
      "BridgeFeeQuoteTest"
    );
    bridgeFeeQuoteTest = await BridgeFeeQuoteTest.deploy();
    const quoteTest = await bridgeFeeQuoteTest.getAddress();

    const state: BridgeFeeQuoteTest.OracleStateStruct = {
      registered: true,
      online: true,
    };

    await bridgeFeeQuoteTest.setOracle(oracle.getAddress(), state);

    const BridgeFeeOracle = await ethers.getContractFactory("BridgeFeeOracle");

    bridgeFeeOracle = await BridgeFeeOracle.deploy();
    await bridgeFeeOracle.initialize(bridgeFeeQuoteTest.getAddress());

    const BridgeFeeQuote = await ethers.getContractFactory("BridgeFeeQuote");
    bridgeFeeQuote = await BridgeFeeQuote.deploy();

    await bridgeFeeQuote.initialize(bridgeFeeOracle.getAddress());

    const role = await bridgeFeeOracle.QUOTE_ROLE();
    await bridgeFeeOracle.grantRole(role, bridgeFeeQuote.getAddress());

    tokens = [token1.address, token2.address];

    await bridgeFeeQuote.registerChain(chainName, []);
    await bridgeFeeQuote.registerToken(chainName, tokens);
  });

  describe("Oracle Management", function () {
    it("should block an oracle correctly", async function () {
      await bridgeFeeOracle.blackOracle(oracle.address);
      const oracleStatus = await bridgeFeeOracle.oracleStatus(oracle.address);
      expect(oracleStatus.isBlackListed).to.be.true;
    });
  });

  describe("Quote Management", function () {
    it("should create a new quote", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input = await newBridgeFeeQuote(
        chainName,
        token1.address,
        fee,
        gasLimit,
        expiry,
        oracle
      );

      await expect(bridgeFeeQuote.quote([input]))
        .to.be.emit(bridgeFeeQuote, "NewQuote")
        .withArgs(
          0,
          input.oracle,
          chainName,
          input.token,
          input.fee,
          input.gasLimit,
          input.expiry
        );

      const quoteList = await bridgeFeeQuote.getQuoteList(chainName);
      expect(quoteList.length).to.be.equal(1);
    });

    it("should revert when trying to get quotes for an inactive chain", async function () {
      const chainName = ethers.encodeBytes32String("InactiveChain");
      await expect(
        bridgeFeeQuote.getQuoteList(chainName)
      ).to.be.revertedWithCustomError(bridgeFeeQuote, "ChainNameInvalid");
    });

    it("should revert when trying to create a quote with an expired expiry", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) - 3600;

      const signature = await generateSignature(
        chainName,
        token1.address,
        fee,
        gasLimit,
        expiry,
        oracle
      );

      const quoteInput: IBridgeFeeQuote.QuoteInputStruct = {
        chainName: chainName,
        token: token1.address,
        oracle: oracle.address,
        fee: fee,
        gasLimit: gasLimit,
        expiry: expiry,
        signature: signature,
      };
      await expect(
        bridgeFeeQuote.quote([quoteInput])
      ).to.be.revertedWithCustomError(bridgeFeeQuote, "QuoteExpired");
    });

    it("should revert when trying to create a quote without new oracle", async function () {
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      const input = await newBridgeFeeQuote(
        chainName,
        token1.address,
        fee,
        gasLimit,
        expiry,
        oracle
      );
      const input2 = await newBridgeFeeQuote(
        chainName,
        token2.address,
        fee,
        gasLimit,
        expiry,
        oracle
      );
      await bridgeFeeQuote.quote([input, input2]);

      const quoteList = await bridgeFeeQuote.getQuoteList(chainName);
      expect(quoteList.length).to.be.equal(2);

      const oracles = await bridgeFeeOracle.getOracleList();
      expect(oracles.length).to.be.equal(1);
    });

    it("test 1 ~ 5 quote gas limit", async function () {
      const number = 5;
      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;
      const singers = await ethers.getSigners();
      let tokens: AddressLike[] = [];
      for (let i = 0; i < number; i++) {
        tokens.push(singers[i + 10].address);
      }
      await bridgeFeeQuote.registerToken(chainName, tokens);
      let quoteList: IBridgeFeeQuote.QuoteInputStruct[] = [];
      for (let i = 0; i < number; i++) {
        const input = await newBridgeFeeQuote(
          chainName,
          singers[i + 10].address,
          fee,
          gasLimit,
          expiry,
          oracle
        );
        quoteList.push(input);
        await bridgeFeeQuote.quote(quoteList);

        const quoteL = await bridgeFeeQuote.getQuoteList(chainName);
        expect(quoteL.length).to.be.equal(i + 1);
      }
    });
    it("first oracle quote", async function () {
      await bridgeFeeOracle.setDefaultOracle(oracle);

      const fee = 1;
      const gasLimit = 0;
      const expiry = (await currentTime()) + 3600;

      await bridgeFeeQuote.quote([
        {
          chainName: chainName,
          token: token1.address,
          oracle: oracle.address,
          fee: fee,
          gasLimit: gasLimit,
          expiry: expiry,
          signature: await generateSignature(
            chainName,
            token1.address,
            fee,
            gasLimit,
            expiry,
            oracle
          ),
        },
      ]);

      const [quote, expire] = await bridgeFeeQuote.getQuoteByToken(
        chainName,
        token1.address,
        0
      );
      expect(expire).to.be.true;
      expect(quote.fee).to.be.equal(fee);
      expect(quote.gasLimit).to.be.equal(gasLimit);
    });
  });
});

async function generateSignature(
  chainName: string,
  token: string,
  fee: number,
  gasLimit: number,
  expiry: number,
  wallet: HDNodeWallet
): Promise<string> {
  const hash = ethers.solidityPackedKeccak256(
    ["string", "address", "uint256", "uint256", "uint256"],
    [chainName, token, fee, gasLimit, expiry]
  );

  const messageHash = ethers.solidityPackedKeccak256(
    ["string", "bytes32"],
    [messagePrefix, hash]
  );

  const signatureW = wallet.signingKey.sign(messageHash);
  let v = "0x1b";
  if (signatureW.v === 28) {
    v = "0x1c";
  }
  return ethers.concat([signatureW.r, signatureW.s, v]);
}

async function newBridgeFeeQuote(
  chainName: string,
  token: string,
  fee: number,
  gasLimit: number,
  expiry: number,
  oracle: HDNodeWallet
): Promise<IBridgeFeeQuote.QuoteInputStruct> {
  const signature = await generateSignature(
    chainName,
    token,
    fee,
    gasLimit,
    expiry,
    oracle
  );
  return {
    chainName: chainName,
    token: token,
    oracle: oracle.address,
    fee: fee,
    gasLimit: gasLimit,
    expiry: expiry,
    signature: signature,
  };
}

async function currentTime(): Promise<number> {
  const blockNumber = await ethers.provider.getBlockNumber();
  const block = await ethers.provider.getBlock(blockNumber);
  return block ? block.timestamp : Math.floor(Date.now() / 1000);
}
