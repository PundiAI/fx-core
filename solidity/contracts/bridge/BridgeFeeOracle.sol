// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {IBridgeFeeOracle} from "./IBridgeFee.sol";
import {IBridgeOracle} from "./IBridgeOracle.sol";
import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {EnumerableSetUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/structs/EnumerableSetUpgradeable.sol";

contract BridgeFeeOracle is
    IBridgeFeeOracle,
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable,
    ReentrancyGuardUpgradeable
{
    using EnumerableSetUpgradeable for EnumerableSetUpgradeable.AddressSet;

    bytes32 public constant QUOTE_ROLE = keccak256("QUOTE_ROLE");
    bytes32 public constant OWNER_ROLE = keccak256("OWNER_ROLE");
    bytes32 public constant UPGRADE_ROLE = keccak256("UPGRADE_ROLE");

    address public crosschainContract;
    address public defaultOracle;

    struct State {
        bool isBlack;
        bool isActive;
    }

    mapping(bytes32 => EnumerableSetUpgradeable.AddressSet) private oracles;
    mapping(bytes32 => mapping(address => State)) public oracleStatus;

    function initialize(address _crosschain) public initializer {
        crosschainContract = _crosschain;

        __AccessControl_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(OWNER_ROLE, msg.sender);
    }

    /**
     * @notice Checks if an oracle is online for a given chain.
     * @param _chainName The name of the chain.
     * @param _oracle The address of the oracle.
     * @return bool indicating if the oracle is online.
     */
    function isOnline(
        bytes32 _chainName,
        address _oracle
    ) external onlyRole(QUOTE_ROLE) nonReentrant returns (bool) {
        if (oracleStatus[_chainName][_oracle].isActive) {
            return true;
        }
        if (oracleStatus[_chainName][_oracle].isBlack) {
            return false;
        }
        if (_oracle == defaultOracle) {
            oracleStatus[_chainName][_oracle].isActive = true;
            oracles[_chainName].add(_oracle);
            return true;
        }
        if (!IBridgeOracle(crosschainContract).hasOracle(_chainName, _oracle)) {
            return false;
        }
        if (
            !IBridgeOracle(crosschainContract).isOracleOnline(
                _chainName,
                _oracle
            )
        ) {
            return false;
        }
        oracleStatus[_chainName][_oracle].isActive = true;
        oracles[_chainName].add(_oracle);
        return true;
    }

    function getOracleList(
        bytes32 _chainName
    ) external view returns (address[] memory) {
        return oracles[_chainName].values();
    }

    function blackOracle(
        bytes32 _chainName,
        address _oracle
    ) external onlyRole(OWNER_ROLE) {
        if (oracleStatus[_chainName][_oracle].isBlack) {
            return;
        }
        if (oracleStatus[_chainName][_oracle].isActive) {
            oracleStatus[_chainName][_oracle].isActive = false;
            oracles[_chainName].remove(_oracle);
        }
        oracleStatus[_chainName][_oracle].isBlack = true;
    }

    function activeOracle(
        bytes32 _chainName,
        address _oracle
    ) external onlyRole(OWNER_ROLE) {
        if (oracleStatus[_chainName][_oracle].isActive) {
            return;
        }
        oracleStatus[_chainName][_oracle] = State(false, true);
        oracles[_chainName].add(_oracle);
    }

    function setDefaultOracle(
        address _defaultOracle
    ) external onlyRole(OWNER_ROLE) {
        defaultOracle = _defaultOracle;
    }

    function _authorizeUpgrade(
        address
    ) internal override onlyRole(UPGRADE_ROLE) {} // solhint-disable-line no-empty-blocks
}
