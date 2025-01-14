// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

/* solhint-disable custom-errors */

contract BridgeFeeQuoteTest {
    struct OracleState {
        bool registered;
        bool online;
    }

    mapping(string => mapping(address => OracleState)) public oracleStatus;

    function setOracle(
        string memory _chainName,
        address _oracle,
        OracleState memory _state
    ) external {
        oracleStatus[_chainName][_oracle] = _state;
    }

    function hasOracle(
        string memory _chainName,
        address _externalAddress
    ) external view returns (bool _result) {
        return oracleStatus[_chainName][_externalAddress].registered;
    }

    function isOracleOnline(
        string memory _chainName,
        address _externalAddress
    ) external view returns (bool _result) {
        return oracleStatus[_chainName][_externalAddress].online;
    }
}
