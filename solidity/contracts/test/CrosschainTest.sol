// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

/* solhint-disable no-global-import */
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "../interfaces/ICrosschain.sol";
import {IBridgeFeeQuote} from "../interfaces/IBridgeFee.sol";
import {IBridgeCallContext} from "../interfaces/IBridgeCallContext.sol";

/* solhint-enable no-global-import */
/* solhint-disable custom-errors */

contract CrosschainTest is IBridgeCallContext {
    address public constant CROSS_CHAIN_ADDRESS =
        address(0x0000000000000000000000000000000000001004);
    address public constant BRIDGE_FEE_ADDRESS =
        address(0x0000000000000000000000000000000000001005);

    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable returns (bool) {
        if (_token != address(0)) {
            IERC20(_token).transferFrom(
                msg.sender,
                address(this),
                _amount + _fee
            );
            IERC20(_token).approve(CROSS_CHAIN_ADDRESS, _amount + _fee);
        } else {
            require(
                msg.value == _amount + _fee,
                "msg.value not equal amount + fee"
            );
        }

        return
            ICrosschain(CROSS_CHAIN_ADDRESS).crossChain{value: msg.value}(
                _token,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            );
    }

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) external view returns (uint256) {
        return
            ICrosschain(CROSS_CHAIN_ADDRESS).bridgeCoinAmount(_token, _target);
    }

    function bridgeCall(
        string memory _dstChain,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _quoteId,
        uint256 _gasLimit,
        bytes memory _memo
    ) external payable returns (uint256) {
        require(
            _tokens.length == _amounts.length,
            "token and amount length not equal"
        );
        for (uint256 i = 0; i < _tokens.length; i++) {
            IERC20(_tokens[i]).transferFrom(
                msg.sender,
                address(this),
                _amounts[i]
            );
            IERC20(_tokens[i]).approve(CROSS_CHAIN_ADDRESS, _amounts[i]);
        }

        IBridgeFeeQuote.QuoteInfo memory info = IBridgeFeeQuote(
            BRIDGE_FEE_ADDRESS
        ).getQuoteById(_quoteId);

        uint256 msgValue = msg.value;
        // check if the fee token is origin token
        if (
            info.tokenName ==
            bytes32(
                0x6170756e64696169000000000000000000000000000000000000000000000000
            )
        ) {
            require(msg.value >= info.amount, "msg.value not enough");
            msgValue = msgValue - info.amount;
        } else {
            (address tokenAddress, ) = ICrosschain(CROSS_CHAIN_ADDRESS)
                .getERC20Token(info.tokenName);
            IERC20(tokenAddress).transferFrom(
                msg.sender,
                address(this),
                info.amount
            );
        }

        return
            ICrosschain(CROSS_CHAIN_ADDRESS).bridgeCall{value: msgValue}(
                _dstChain,
                _receiver,
                _tokens,
                _amounts,
                _to,
                _data,
                _quoteId,
                _gasLimit,
                _memo
            );
    }

    function executeClaim(
        string memory _chain,
        uint256 _eventNonce
    ) external returns (bool _result) {
        return
            ICrosschain(CROSS_CHAIN_ADDRESS).executeClaim(_chain, _eventNonce);
    }

    function onBridgeCall(
        address,
        address,
        address[] memory,
        uint256[] memory,
        bytes memory,
        bytes memory
    ) external view override {
        require(msg.sender == CROSS_CHAIN_ADDRESS, "only cross-chain address");
    }

    function onRevert(uint256, bytes memory) external view override {
        require(msg.sender == CROSS_CHAIN_ADDRESS, "only cross-chain address");
    }

    function hasOracle(
        bytes32 _chain,
        address _externalAddress
    ) external view returns (bool _result) {
        return
            ICrosschain(CROSS_CHAIN_ADDRESS).hasOracle(
                _chain,
                _externalAddress
            );
    }

    function isOracleOnline(
        bytes32 _chain,
        address _externalAddress
    ) external view returns (bool _result) {
        return
            ICrosschain(CROSS_CHAIN_ADDRESS).isOracleOnline(
                _chain,
                _externalAddress
            );
    }

    function getERC20Token(
        bytes32 _denom
    ) external view returns (address _token, bool _enable) {
        return ICrosschain(CROSS_CHAIN_ADDRESS).getERC20Token(_denom);
    }
}
