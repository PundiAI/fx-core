// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";

error TransferFailed();
error InvalidInput();
error InsufficientBalance();

contract DelegateRefund is
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable,
    ReentrancyGuardUpgradeable
{
    bytes32 public constant REFUNDER_ROLE = keccak256("REFUNDER_ROLE");
    bytes32 public constant UPGRADE_ROLE = keccak256("UPGRADE_ROLE");

    event RefundExecuted(address indexed to, uint256 amount);

    function initialize() public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(REFUNDER_ROLE, msg.sender);
    }

    function batchRefund(
        address[] calldata recipients,
        uint256[] calldata amounts
    ) external onlyRole(REFUNDER_ROLE) nonReentrant {
        if (recipients.length == 0 || recipients.length != amounts.length) {
            revert InvalidInput();
        }

        uint256 totalAmount;
        for (uint256 i = 0; i < amounts.length; i++) {
            totalAmount += amounts[i];
        }

        if (address(this).balance < totalAmount) {
            revert InsufficientBalance();
        }

        for (uint256 i = 0; i < recipients.length; i++) {
            if (recipients[i] == address(0) || amounts[i] == 0) {
                continue;
            }

            (bool success, ) = recipients[i].call{value: amounts[i]}("");
            if (!success) {
                revert TransferFailed();
            }

            emit RefundExecuted(recipients[i], amounts[i]);
        }
    }

    receive() external payable {}

    function _authorizeUpgrade(
        address
    ) internal override onlyRole(UPGRADE_ROLE) {} // solhint-disable-line no-empty-blocks
}
