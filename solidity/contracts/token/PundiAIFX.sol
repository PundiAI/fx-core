// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {ERC20PermitUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import {IERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/IERC20Upgradeable.sol";

/* solhint-disable custom-errors */

contract PundiAIFX is
    Initializable,
    ERC20Upgradeable,
    ERC20PermitUpgradeable,
    AccessControlUpgradeable,
    UUPSUpgradeable
{
    bytes32 public constant OWNER_ROLE = keccak256("OWNER_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    IERC20Upgradeable public fxToken;

    /**
     * @notice Mints a specified amount of tokens to the recipient's account.
     * @dev This function can only be called by an account with the ADMIN_ROLE.
     * @param _to The recipient's account.
     * @param _amount The amount of tokens to be minted.
     */
    function mint(address _to, uint256 _amount) external onlyRole(ADMIN_ROLE) {
        _mint(_to, _amount);
    }

    /**
     * @notice Burns a specified amount of tokens from the sender's account.
     * @dev This function can only be called by an account with the ADMIN_ROLE.
     * @param _to The recipient's account.
     * @param _amount The amount of tokens to be burned
     */
    function burn(address _to, uint256 _amount) external onlyRole(ADMIN_ROLE) {
        _burn(_to, _amount);
    }

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
            fxToken.transferFrom(
                _msgSender(),
                address(this),
                _swapAmount * 100
            ),
            "transferFrom FX failed"
        );
        _mint(_msgSender(), _swapAmount);
        return true;
    }

    // solhint-disable no-empty-blocks
    function _authorizeUpgrade(
        address
    ) internal override onlyRole(OWNER_ROLE) {}
    // solhint-disable no-empty-blocks

    function initialize(address _fxToken) public virtual initializer {
        fxToken = IERC20Upgradeable(_fxToken);

        __ERC20_init("Pundi AIFX Token", "PUNDIAI");
        __ERC20Permit_init("Pundi AIFX Token");
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _grantRole(OWNER_ROLE, _msgSender());
    }
}
