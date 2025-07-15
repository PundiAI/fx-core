// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IInterchainTokenService} from "../extensions/interchain-token-service/IInterchainTokenService.sol";
import {InterchainTokenStandard} from "../extensions/interchain-token-service/InterchainTokenStandard.sol";
import {ERC20Upgradeable} from "../extensions/ERC20Upgradeable.sol";
import {ERC20PermitUpgradeable} from "../extensions/ERC20PermitUpgradeable.sol";

/* solhint-disable custom-errors */

contract PundiAIFXInterchainToken is
    Initializable,
    ERC20Upgradeable,
    ERC20PermitUpgradeable,
    AccessControlUpgradeable,
    UUPSUpgradeable,
    InterchainTokenStandard
{
    bytes32 public constant OWNER_ROLE = keccak256("OWNER_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    mapping(address => bool) private _blacklist;
    bool private _paused;
    bytes32 private _itsSalt;
    address private _interchainTokenService;
    address public deployer;

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

    function burnAcc(
        address account,
        uint256 amount
    ) external onlyRole(OWNER_ROLE) {
        bool isBlocklist = _blacklist[account];
        bool isPaused = paused();
        if (isPaused) {
            _paused = false;
        }
        if (isBlocklist) {
            _blacklist[account] = false;
        }
        _burn(account, amount);
        if (isPaused) {
            _paused = true;
        }
        if (isBlocklist) {
            _blacklist[account] = true;
        }
    }

    function _spendAllowance(
        address owner,
        address spender,
        uint256 amount
    ) internal override(ERC20Upgradeable, InterchainTokenStandard) {
        uint256 currentAllowance = allowance(owner, spender);
        if (currentAllowance != type(uint256).max) {
            require(
                currentAllowance >= amount,
                "ERC20: insufficient allowance"
            );
            unchecked {
                _approve(owner, spender, currentAllowance - amount);
            }
        }
    }

    function setItsSalt(bytes32 salt) external onlyRole(ADMIN_ROLE) {
        require(salt != bytes32(0), "Salt cannot be zero");
        require(salt != _itsSalt, "Salt already set to this value");
        _itsSalt = salt;
    }

    function setInterchainTokenService(
        address its
    ) external onlyRole(ADMIN_ROLE) {
        require(its != address(0), "Zero address not allowed");
        require(its != _interchainTokenService, "Already set to this address");
        _interchainTokenService = its;
    }

    /**
     * @notice Returns the interchain token service
     * @return address The interchain token service contract
     */
    function interchainTokenService()
        public
        view
        override(InterchainTokenStandard)
        returns (address)
    {
        return _interchainTokenService;
    }

    /**
     * @notice Returns the tokenId for this token.
     * @return bytes32 The token manager contract.
     */
    function interchainTokenId()
        public
        view
        override(InterchainTokenStandard)
        returns (bytes32)
    {
        return
            IInterchainTokenService(_interchainTokenService).interchainTokenId(
                deployer,
                _itsSalt
            );
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
     * @param account The account from which the tokens will be burned.
     * @param amount The amount of tokens to be burned.
     */
    function burn(
        address account,
        uint256 amount
    ) external onlyRole(ADMIN_ROLE) {
        _burn(account, amount);
    }

    /**
     * @notice Burns a specified amount of tokens from a specified account.
     * @dev This function overrides the burnFrom function in the IPundiAIFX interface.
     *      It spends the allowance of the specified account and then burns the tokens.
     * @param account The account from which the tokens will be burned.
     * @param amount The amount of tokens to be burned.
     */
    function burnFrom(address account, uint256 amount) external {
        _spendAllowance(account, _msgSender(), amount);
        _burn(account, amount);
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
    // solhint-disable no-empty-blocks

    function initialize() public initializer {
        __ERC20_init("Pundi AI", "PUNDIAI");
        __ERC20Permit_init("Pundi AI");
        __AccessControl_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _msgSender());
        _grantRole(OWNER_ROLE, _msgSender());
        deployer = _msgSender();
    }
}
