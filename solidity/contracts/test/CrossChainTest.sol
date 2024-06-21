// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

/* solhint-disable no-global-import */
import "../bridge/ICrossChain.sol";
import "../fip20/IFIP20Upgradable.sol";

/* solhint-enable no-global-import */
/* solhint-disable custom-errors */

contract CrossChainTest {
    address public constant CROSS_CHAIN_ADDRESS =
        address(0x0000000000000000000000000000000000001004);

    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable returns (bool) {
        if (_token != address(0)) {
            IFIP20Upgradable(_token).transferFrom(
                msg.sender,
                address(this),
                _amount + _fee
            );
            IFIP20Upgradable(_token).approve(
                CROSS_CHAIN_ADDRESS,
                _amount + _fee
            );
        }

        if (_token != address(0)) {
            uint256 allowance = IFIP20Upgradable(_token).allowance(
                address(this),
                CROSS_CHAIN_ADDRESS
            );
            require(
                allowance == _amount + _fee,
                "allowance not equal amount + fee"
            );
        } else {
            require(
                msg.value == _amount + _fee,
                "msg.value not equal amount + fee"
            );
        }

        return
            ICrossChain(CROSS_CHAIN_ADDRESS).crossChain{value: msg.value}(
                _token,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            );
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) external returns (bool) {
        return
            ICrossChain(CROSS_CHAIN_ADDRESS).cancelSendToExternal(
                _chain,
                _txID
            );
    }

    function increaseBridgeFee(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _fee
    ) external payable returns (bool) {
        return
            ICrossChain(CROSS_CHAIN_ADDRESS).increaseBridgeFee(
                _chain,
                _txID,
                _token,
                _fee
            );
    }

    function bridgeCoinAmount(
        address _token,
        bytes32 _target
    ) external view returns (uint256) {
        return
            ICrossChain(CROSS_CHAIN_ADDRESS).bridgeCoinAmount(_token, _target);
    }

    function bridgeCall(
        string memory _dstChain,
        address _receiver,
        address[] memory _tokens,
        uint256[] memory _amounts,
        address _to,
        bytes memory _data,
        uint256 _value,
        bytes memory _memo
    ) internal returns (uint256) {
        return
            ICrossChain(CROSS_CHAIN_ADDRESS).bridgeCall(
                _dstChain,
                _receiver,
                _tokens,
                _amounts,
                _to,
                _data,
                _value,
                _memo
            );
    }
}
