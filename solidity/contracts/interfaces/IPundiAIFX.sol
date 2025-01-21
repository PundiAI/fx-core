// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {IERC20MetadataUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20MetadataUpgradeable.sol";
import {IERC20PermitUpgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/extensions/IERC20PermitUpgradeable.sol";
import {IAccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {IERC20Burn} from "./IERC20Burn.sol";

interface IPundiAIFX is
    IERC20MetadataUpgradeable,
    IERC20PermitUpgradeable,
    IAccessControlUpgradeable,
    IERC20Burn
{
    /// @notice Mints new tokens
    /// @param _to The address that will receive the minted tokens
    /// @param _amount The amount of tokens to mint
    /// @dev Caller must have ADMIN_ROLE
    function mint(address _to, uint256 _amount) external;
}
