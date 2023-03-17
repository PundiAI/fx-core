// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

interface IStaking {
    function delegate(
        string memory _val
    ) external payable returns (uint256, uint256);

    function undelegate(
        string memory _val,
        uint256 _shares
    ) external returns (uint256, uint256, uint256);

    function withdraw(string memory _val) external returns (uint256);

    function delegation(
        string memory _val,
        address _del
    ) external view returns (uint256);
}

contract Staking is IStaking {
    address private constant _stakingAddress =
        address(0x0000000000000000000000000000000000000064);

    function delegate(
        string memory _val
    ) external payable virtual override returns (uint256, uint256) {
        return _delegate(_val);
    }

    function _delegate(string memory _val) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call{
            value: msg.value
        }(Encode.delegate(_val));
        Decode.ok(result, data, "delegate failed");
        return Decode.delegate(data);
    }

    function undelegate(
        string memory _val,
        uint256 _shares
    ) external virtual override returns (uint256, uint256, uint256) {
        return _undelegate(_val, _shares);
    }

    function _undelegate(
        string memory _val,
        uint256 _shares
    ) internal returns (uint256, uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.undelegate(_val, _shares)
        );
        Decode.ok(result, data, "undelegate failed");
        return Decode.undelegate(data);
    }

    function withdraw(
        string memory _val
    ) external virtual override returns (uint256) {
        return _withdraw(_val);
    }

    function _withdraw(string memory _val) internal returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.withdraw(_val)
        );
        Decode.ok(result, data, "withdraw failed");
        return Decode.withdraw(data);
    }

    function delegation(
        string memory _val,
        address _del
    ) public view virtual override returns (uint256) {
        return _delegation(_val, _del);
    }

    function _delegation(
        string memory _val,
        address _del
    ) internal view returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.delegation(_val, _del)
        );
        Decode.ok(result, data, "delegation failed");
        return Decode.delegation(data);
    }
}

library Encode {
    function delegate(
        string memory _validator
    ) internal pure returns (bytes memory) {
        return abi.encodeWithSignature("delegate(string)", _validator);
    }

    function undelegate(
        string memory _validator,
        uint256 _shares
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "undelegate(string,uint256)",
                _validator,
                _shares
            );
    }

    function withdraw(
        string memory _validator
    ) internal pure returns (bytes memory) {
        return abi.encodeWithSignature("withdraw(string)", _validator);
    }

    function delegation(
        string memory _validator,
        address _delegate
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "delegation(string,address)",
                _validator,
                _delegate
            );
    }
}

library Decode {
    function delegate(
        bytes memory data
    ) internal pure returns (uint256, uint256) {
        (uint256 shares, uint256 reward) = abi.decode(data, (uint256, uint256));
        return (shares, reward);
    }

    function undelegate(
        bytes memory data
    ) internal pure returns (uint256, uint256, uint256) {
        (uint256 amount, uint256 reward, uint256 endTime) = abi.decode(
            data,
            (uint256, uint256, uint256)
        );
        return (amount, reward, endTime);
    }

    function withdraw(bytes memory data) internal pure returns (uint256) {
        uint256 reward = abi.decode(data, (uint256));
        return reward;
    }

    function delegation(bytes memory data) internal pure returns (uint256) {
        uint256 delegateAmount = abi.decode(data, (uint256));
        return delegateAmount;
    }

    function ok(
        bool _result,
        bytes memory _data,
        string memory _msg
    ) internal pure {
        if (!_result) {
            string memory errMsg = abi.decode(_data, (string));
            if (bytes(_msg).length < 1) {
                revert(errMsg);
            }
            revert(string(abi.encodePacked(_msg, ": ", errMsg)));
        }
    }
}
