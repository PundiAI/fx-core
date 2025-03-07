import { expect } from "chai";
import { ethers } from "hardhat";
import { DelegateRefund } from "../typechain-types";
import { HardhatEthersSigner } from "@nomicfoundation/hardhat-ethers/signers";

describe("staking refund test", function () {
  let delegateRefund: DelegateRefund;
  let user1: HardhatEthersSigner;
  let user2: HardhatEthersSigner;
  let user3: HardhatEthersSigner;

  beforeEach(async function () {
    const signers = await ethers.getSigners();
    user1 = signers[5];
    user2 = signers[6];
    user3 = signers[7];

    const delegateRefundFactory = await ethers.getContractFactory(
      "DelegateRefund"
    );
    delegateRefund = await delegateRefundFactory.deploy();
    await delegateRefund.initialize();
    const address = await delegateRefund.getAddress();

    await signers[0].sendTransaction({
      to: address,
      value: ethers.parseEther("10.0"),
    });
  });

  describe("delegate refund", function () {
    it("success refund", async function () {
      const recipients = [user1.address, user2.address, user3.address];
      const amounts = [
        ethers.parseEther("1.0"),
        ethers.parseEther("2.0"),
        ethers.parseEther("3.0"),
      ];

      await expect(delegateRefund.batchRefund(recipients, amounts))
        .to.emit(delegateRefund, "RefundExecuted")
        .withArgs(user1.address, amounts[0])
        .to.emit(delegateRefund, "RefundExecuted")
        .withArgs(user2.address, amounts[1])
        .to.emit(delegateRefund, "RefundExecuted")
        .withArgs(user3.address, amounts[2]);
    });

    it("failed refund no role", async function () {
      const recipients = [user1.address];
      const amounts = [ethers.parseEther("1.0")];

      await expect(
        delegateRefund.connect(user1).batchRefund(recipients, amounts)
      ).to.be.reverted;
    });

    it("failed refund balance insufficient", async function () {
      const recipients = [user1.address];
      const amounts = [ethers.parseEther("100.0")];

      await expect(
        delegateRefund.batchRefund(recipients, amounts)
      ).to.be.revertedWithCustomError(delegateRefund, "InsufficientBalance");
    });

    it("failed refund input error", async function () {
      await expect(
        delegateRefund.batchRefund(
          [user1.address],
          [ethers.parseEther("1.0"), ethers.parseEther("1.0")]
        )
      ).to.be.revertedWithCustomError(delegateRefund, "InvalidInput");

      await expect(
        delegateRefund.batchRefund([], [])
      ).to.be.revertedWithCustomError(delegateRefund, "InvalidInput");
    });
  });
});
