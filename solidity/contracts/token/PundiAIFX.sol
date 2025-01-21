// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {ERC20PermitUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/ERC20PermitUpgradeable.sol";
import {IPundiAIFX} from "../interfaces/IPundiAIFX.sol";

/* solhint-disable custom-errors */

contract PundiAIFX is
    Initializable,
    ERC20Upgradeable,
    ERC20PermitUpgradeable,
    AccessControlUpgradeable,
    UUPSUpgradeable,
    IPundiAIFX
{
    bytes32 public constant OWNER_ROLE = keccak256("OWNER_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

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
     * @dev This function overrides the burn function in the IPundiAIFX interface.
     * @param amount The amount of tokens to be burned.
     */
    function burn(uint256 amount) external override {
        _burn(_msgSender(), amount);
    }

    /**
     * @notice Burns a specified amount of tokens from a specified account.
     * @dev This function overrides the burnFrom function in the IPundiAIFX interface.
     *      It spends the allowance of the specified account and then burns the tokens.
     * @param account The account from which the tokens will be burned.
     * @param amount The amount of tokens to be burned.
     */
    function burnFrom(address account, uint256 amount) external override {
        _spendAllowance(account, _msgSender(), amount);
        _burn(account, amount);
    }

    // solhint-disable no-empty-blocks
    function _authorizeUpgrade(
        address
    ) internal override onlyRole(OWNER_ROLE) {}
    // solhint-disable no-empty-blocks

    function initialize() public virtual initializer {
        __ERC20_init("Pundi AIFX Token", "PUNDIAI");
        __ERC20Permit_init("Pundi AIFX Token");
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _grantRole(OWNER_ROLE, _msgSender());
    }
}
