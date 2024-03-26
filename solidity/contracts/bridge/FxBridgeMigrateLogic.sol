// SPDX-License-Identifier: MIT

pragma experimental ABIEncoderV2;
pragma solidity ^0.8.0;

import {FxBridgeLogic} from "./FxBridgeLogic.sol";

/* solhint-disable custom-errors */
contract FxBridgeMigrateLogic is FxBridgeLogic {
    /* =============== MIGRATE =============== */
    function migrateInit(
        bytes32 _fxBridgeId,
        uint256 _powerThreshold,
        uint256 _lastEventNonce,
        bytes32 _lastOracleSetCheckpoint,
        uint256 _lastOracleSetNonce,
        address[] memory _bridgeTokens,
        uint256[] memory _lastBatchNonces,
        TokenStatus[] memory _tokenStatuses
    ) public initializer {
        __Pausable_init();
        __Ownable_init();
        __ReentrancyGuard_init();
        require(
            _lastBatchNonces.length == _bridgeTokens.length &&
                _bridgeTokens.length == _tokenStatuses.length,
            "Malformed last batch token information."
        );
        state_fxBridgeId = _fxBridgeId;
        state_powerThreshold = _powerThreshold;
        state_lastEventNonce = _lastEventNonce;
        state_lastOracleSetCheckpoint = _lastOracleSetCheckpoint;
        state_lastOracleSetNonce = _lastOracleSetNonce;
        for (uint256 i = 0; i < _bridgeTokens.length; i++) {
            bridgeTokens.push(_bridgeTokens[i]);
            state_lastBatchNonces[_bridgeTokens[i]] = _lastBatchNonces[i];
            tokenStatus[_bridgeTokens[i]] = _tokenStatuses[i];
        }
        version = "1.0.0";
    }
}
