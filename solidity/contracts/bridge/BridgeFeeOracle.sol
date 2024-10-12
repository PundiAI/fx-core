// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/security/ReentrancyGuardUpgradeable.sol";
import {IBridgeFeeOracle} from "./IBridgeFee.sol";
import {ICrossChain} from "./ICrossChain.sol";

contract BridgeFeeOracle is
    IBridgeFeeOracle,
    Initializable,
    UUPSUpgradeable,
    AccessControlUpgradeable,
    ReentrancyGuardUpgradeable
{
    bytes32 public constant QUOTE_ROLE = keccak256("QUOTE_ROLE");

    address public crossChainContract;
    address public defaultOracle;

    struct State {
        bool isBlackListed;
        bool isActive;
    }

    address[] public oracles;
    mapping(address => State) public oracleStatus;

    function initialize(address _crossChain) public initializer {
        crossChainContract = _crossChain;

        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);

        __AccessControl_init();
        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    function isOnline(
        string memory _chainName,
        address _oracle
    ) external onlyRole(QUOTE_ROLE) returns (bool) {
        if (oracleStatus[_oracle].isActive) return true;
        if (oracleStatus[_oracle].isBlackListed) return false;
        if (!ICrossChain(crossChainContract).hasOracle(_chainName, _oracle)) {
            return false;
        }
        if (
            !ICrossChain(crossChainContract).isOracleOnline(_chainName, _oracle)
        ) {
            return false;
        }
        oracleStatus[_oracle] = State(false, true);
        oracles.push(_oracle);
        return true;
    }

    function blackOracle(
        address _oracle
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        if (oracleStatus[_oracle].isBlackListed) return;
        if (oracleStatus[_oracle].isActive) {
            oracleStatus[_oracle].isActive = false;
            removeOracle(_oracle);
        }
        oracleStatus[_oracle].isBlackListed = true;
    }

    function removeOracle(address _oracle) internal {
        for (uint256 i = 0; i < oracles.length; i++) {
            if (oracles[i] == _oracle) {
                oracles[i] = oracles[oracles.length - 1];
                oracles.pop();
                break;
            }
        }
    }

    function setDefaultOracle(
        address _defaultOracle
    ) external onlyRole(DEFAULT_ADMIN_ROLE) {
        defaultOracle = _defaultOracle;
    }

    function getOracleList() external view returns (address[] memory) {
        return oracles;
    }

    function _authorizeUpgrade(
        address
    ) internal override onlyRole(DEFAULT_ADMIN_ROLE) {} // solhint-disable-line no-empty-blocks
}
