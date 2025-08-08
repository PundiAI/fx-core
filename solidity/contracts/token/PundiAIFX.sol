// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.10;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ERC20Upgradeable} from "../extensions/ERC20Upgradeable.sol";
import {ERC20PermitUpgradeable} from "../extensions/ERC20PermitUpgradeable.sol";
import {IPundiAIFX} from "../interfaces/IPundiAIFX.sol";

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

    mapping(address => bool) private _blacklist;
    bool private _paused;

    constructor() initializer {}

    function paused() public view virtual returns (bool) {
        return _paused;
    }

    function pause() external onlyRole(OWNER_ROLE) {
        require(!_paused, "Contract is already paused");
        _paused = true;
    }

    function unpause() external onlyRole(OWNER_ROLE) {
        require(_paused, "Contract is not paused");
        _paused = false;
    }

    function isBlacklisted(address account) public view returns (bool) {
        return _blacklist[account];
    }

    function addToBlacklist(address account) external onlyRole(OWNER_ROLE) {
        require(account != address(0), "Invalid address");
        require(!_blacklist[account], "Account already blacklisted");
        _blacklist[account] = true;
    }

    function removeFromBlacklist(
        address account
    ) external onlyRole(OWNER_ROLE) {
        require(_blacklist[account], "Account not blacklisted");
        _blacklist[account] = false;
    }

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

    /**
     * @notice setName set a new name for PundiAIFX.
     * @dev This function can only be called by an account with the OWNER_ROLE.
     * @param newName The new name to set.
     */
    function setName(string memory newName) external onlyRole(OWNER_ROLE) {
        _setName(newName);
        _setEIP712Name(newName);
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256
    ) internal view override {
        require(!paused(), "Pausable: paused");
        require(!_blacklist[from], "Sender is blacklisted");
        require(!_blacklist[to], "Recipient is blacklisted");
    }

    // solhint-disable no-empty-blocks
    function _authorizeUpgrade(
        address
    ) internal override onlyRole(OWNER_ROLE) {}

    function initialize() public initializer {
        __ERC20_init("Pundi AI", "PUNDIAI");
        __ERC20Permit_init("Pundi AI");
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _grantRole(OWNER_ROLE, _msgSender());
    }
}
