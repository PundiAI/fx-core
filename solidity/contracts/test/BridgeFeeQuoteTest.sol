// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

/* solhint-disable custom-errors */

contract BridgeFeeQuoteTest {
    struct OracleState {
        bool registered;
        bool online;
    }

    mapping(address => OracleState) public oracleStatus;

    function setOracle(address _oracle, OracleState memory _state) external {
        oracleStatus[_oracle] = _state;
    }

    function hasOracle(
        string memory,
        address _externalAddress
    ) external view returns (bool _result) {
        return oracleStatus[_externalAddress].registered;
    }

    function isOracleOnline(
        string memory,
        address _externalAddress
    ) external view returns (bool _result) {
        return oracleStatus[_externalAddress].online;
    }
}
