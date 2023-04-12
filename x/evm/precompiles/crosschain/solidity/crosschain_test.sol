// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./CrossChain.sol";

interface IERC20 {
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    function approve(address spender, uint256 amount) external returns (bool);

    function allowance(address owner, address spender) external view returns (uint256);
}

contract crosschain_test is CrossChain {
    address private constant _crossChainAddress =
    address(0x0000000000000000000000000000000000001004);

    function crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) external payable override returns (bool) {
        if (_token != address(0)) {
            IERC20(_token).transferFrom(
                msg.sender,
                address(this),
                _amount + _fee
            );
            IERC20(_token).approve(_crossChainAddress, _amount + _fee);
        }
        return _crossChain(_token, _receipt, _amount, _fee, _target, _memo);
    }

    function cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) external override returns (bool) {
        return _cancelSendToExternal(_chain, _txID);
    }

    function increaseBridgeFee(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _fee
    ) external payable override returns (bool) {
        return _increaseBridgeFee(_chain, _txID, _token, _fee);
    }

    function _crossChain(
        address _token,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal returns (bool) {
        if (_token != address(0)) {
            uint256 allowance = IERC20(_token).allowance(
                address(this),
                _crossChainAddress
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

        (bool result, bytes memory data) = _crossChainAddress.call{
        value : msg.value
        }(Encode.crossChain(_token, _receipt, _amount, _fee, _target, _memo));
        Decode.ok(result, data, "cross-chain failed");
        return Decode.crossChain(data);
    }

    function _cancelSendToExternal(
        string memory _chain,
        uint256 _txID
    ) internal returns (bool) {
        (bool result, bytes memory data) = _crossChainAddress.call(
            Encode.cancelSendToExternal(_chain, _txID)
        );
        Decode.ok(result, data, "cancel send to external failed");
        return Decode.cancelSendToExternal(data);
    }

    function _increaseBridgeFee(
        string memory _chain,
        uint256 _txID,
        address _token,
        uint256 _fee
    ) internal returns (bool) {
        (bool result, bytes memory data) = _crossChainAddress.call(
            Encode.increaseBridgeFee(_chain, _txID, _token, _fee)
        );
        Decode.ok(result, data, "increase bridge fee failed");
        return Decode.increaseBridgeFee(data);
    }
}
