// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {IERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IPundiAIFX} from "../interfaces/IPundiAIFX.sol";
import {IERC20Burn} from "../interfaces/IERC20Burn.sol";

/* solhint-disable custom-errors */

contract FXSwapPundiAI is Initializable, OwnableUpgradeable, UUPSUpgradeable {
    IERC20Upgradeable public fxToken;
    IPundiAIFX public pundiAIToken;

    /**
     * @notice Swaps a specified amount of tokens for a smaller amount of new tokens.
     * @dev The function calculates the swap amount as 1% of the input amount.
     *      It then transfers the original tokens from the sender to the contract
     *      and mints new tokens to the sender.
     * @param _amount The amount of tokens to be swapped.
     * @return bool Returns true if the swap is successful.
     */
    function swap(uint256 _amount) external returns (bool) {
        uint256 _swapAmount = _amount / 100;
        require(_swapAmount > 0, "swap amount is too small");

        require(
            fxToken.transferFrom(_msgSender(), address(this), _amount),
            "transferFrom FX failed"
        );
        pundiAIToken.mint(_msgSender(), _swapAmount);

        emit Swap(_msgSender(), _amount, _swapAmount);
        return true;
    }

    function burnFXToken(uint256 _amount) external onlyOwner {
        IERC20Burn(address(fxToken)).burn(_amount);
    }

    event Swap(address indexed from, uint256 fxAmt, uint256 pundiaiAmt);

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function initialize(
        address _fxToken,
        address _pundiAIToken
    ) public virtual initializer {
        fxToken = IERC20Upgradeable(_fxToken);
        pundiAIToken = IPundiAIFX(_pundiAIToken);

        __Ownable_init();
        __UUPSUpgradeable_init();
    }
}
