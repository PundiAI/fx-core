// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title Airdrop
 * @dev A contract for distributing tokens to multiple addresses.
 */
contract Airdrop is Ownable {
    /**
     * @dev Distributes tokens to multiple addresses.
     * @param token The address of the token contract.
     * @param recipients An array of addresses to receive the tokens.
     * @param amounts An array of amounts corresponding to each recipient.
     */
    function distributeTokens(
        IERC20 token,
        address[] calldata recipients,
        uint256[] calldata amounts
    ) external onlyOwner {
        require(
            recipients.length == amounts.length,
            "Mismatched input lengths"
        );

        for (uint256 i = 0; i < recipients.length; i++) {
            require(
                token.balanceOf(recipients[i]) == 0,
                "Recipient already has tokens"
            );
            require(recipients[i] != address(0), "Invalid recipient address");
            require(amounts[i] > 0, "Amount must be greater than zero");
            require(
                token.transfer(recipients[i], amounts[i]),
                "Transfer failed"
            );
        }
        require(
            token.balanceOf(address(this)) == 0,
            "Contract should not hold any tokens after distribution"
        );
    }

    /**
     * @dev Withdraws any remaining tokens in the contract to the owner.
     * @param token The address of the token contract.
     * @param amount The amount of tokens to withdraw.
     */
    function withdrawTokens(IERC20 token, uint256 amount) external onlyOwner {
        require(
            token.balanceOf(address(this)) >= amount,
            "Insufficient balance"
        );
        require(token.transfer(owner(), amount), "Transfer failed");
    }

    /**
     * @dev Returns the balance of the contract in tokens.
     * @param token The address of the token contract.
     * @return The balance of the contract in tokens.
     */
    function getTokenBalance(IERC20 token) external view returns (uint256) {
        return token.balanceOf(address(this));
    }
}
