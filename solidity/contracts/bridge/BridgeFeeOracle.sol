// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {EnumerableSetUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/structs/EnumerableSetUpgradeable.sol";
import {IBridgeFeeOracle} from "./IBridgeFee.sol";
import {ICrosschain} from "./ICrosschain.sol";

contract BridgeFeeOracle is
    IBridgeFeeOracle,
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable,
    ReentrancyGuardUpgradeable
{
    using EnumerableSetUpgradeable for EnumerableSetUpgradeable.AddressSet;

    bytes32 public constant QUOTE_ROLE = keccak256("QUOTE_ROLE");

    address public crossChainContract;
    address public defaultOracle;

    struct State {
        bool isBlacklisted;
        bool isActive;
    }

    EnumerableSetUpgradeable.AddressSet private oracles;
    mapping(address => State) public oracleStatus;

    function initialize(address _crossChain) public initializer {
        __AccessControl_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();

        crossChainContract = _crossChain;

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
    }

    /**
     * @notice Checks if an oracle is online for a given chain.
     * @param _chainName The name of the chain.
     * @param _oracle The address of the oracle.
     * @return bool indicating if the oracle is online.
     */
    function isOnline(
        string memory _chainName,
        address _oracle
    ) external onlyRole(QUOTE_ROLE) nonReentrant returns (bool) {
        if (oracleStatus[_oracle].isActive) return true;
        if (oracleStatus[_oracle].isBlacklisted) return false;
        if (!ICrosschain(crossChainContract).hasOracle(_chainName, _oracle)) {
            return false;
        }
        if (
            !ICrosschain(crossChainContract).isOracleOnline(_chainName, _oracle)
        ) {
            return false;
        }
        oracleStatus[_oracle] = State(false, true);
        oracles.add(_oracle);
        return true;
    }

    function blackOracle(
        address _oracle
    ) external onlyRole(DEFAULT_ADMIN_ROLE) nonReentrant {
        if (oracleStatus[_oracle].isBlacklisted) return;
        if (oracleStatus[_oracle].isActive) {
            oracleStatus[_oracle].isActive = false;
            oracles.remove(_oracle);
        }
        oracleStatus[_oracle].isBlacklisted = true;
    }

    function setDefaultOracle(
        address _defaultOracle
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (!oracles.contains(_defaultOracle)) {
            oracleStatus[_defaultOracle] = State(false, true);
            oracles.add(_defaultOracle);
        }
        defaultOracle = _defaultOracle;
    }

    function getOracleList() external view returns (address[] memory) {
        return oracles.values();
    }

    function _authorizeUpgrade(
        address
    ) internal override onlyRole(DEFAULT_ADMIN_ROLE) {} // solhint-disable-line no-empty-blocks
}
