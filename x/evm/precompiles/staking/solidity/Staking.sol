// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

import "./IStaking.sol";
import "./Encode.sol";
import "./Decode.sol";

contract Staking is IStaking {
    address private constant _stakingAddress =
        address(0x0000000000000000000000000000000000001003);

    function delegate(
        string memory _val
    )
        external
        payable
        virtual
        override
        returns (uint256 _shares, uint256 _reward)
    {
        // return _delegate(_val);
        emit Delegate(address(0), _val, 0, 0);
        return (0, 0);
    }

    function undelegate(
        string memory _val,
        uint256 _shares
    )
        external
        virtual
        override
        returns (uint256 _amount, uint256 _reward, uint256 _completionTime)
    {
        // return _undelegate(_val, _shares);
        emit Undelegate(address(0), _val, 0, 0, 0);
        return (0, 0, 0);
    }

    function withdraw(
        string memory _val
    ) external virtual override returns (uint256 _reward) {
        // return _withdraw(_val);
        emit Withdraw(address(0), _val, 0);
        return 0;
    }

    function approveShares(
        string memory _val,
        address _spender,
        uint256 _shares
    ) external virtual override returns (bool _result) {
        // return _approveShares(_val, _spender, _shares);
        emit ApproveShares(address(0), _spender, _val, 0);
        return false;
    }

    function transferShares(
        string memory _val,
        address _to,
        uint256 _shares
    ) external virtual override returns (uint256 _token, uint256 _reward) {
        // return _transferShares(_val, _to, _shares);
        emit TransferShares(address(0), _to, _val, 0, 0);
        return (0, 0);
    }

    function transferFromShares(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) external virtual override returns (uint256 _token, uint256 _reward) {
        // return _transferFromShares(_val, _from, _to, _shares);
        emit TransferShares(_from, _to, _val, 0, 0);
        return (0, 0);
    }

    function delegation(
        string memory _val,
        address _del
    )
        public
        view
        virtual
        override
        returns (uint256 _shares, uint256 _delegateAmount)
    {
        // return _delegation(_val, _del);
        return (0, 0);
    }

    function delegationRewards(
        string memory _val,
        address _del
    ) public view virtual override returns (uint256 _reward) {
        // return _delegationRewards(_val, _del);
        return 0;
    }

    function allowanceShares(
        string memory _val,
        address _owner,
        address _spender
    ) public view virtual override returns (uint256 _shares) {
        // return _allowanceShares(_val, _owner, _spender);
        return 0;
    }

    function _delegate(string memory _val) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call{
            value: msg.value
        }(Encode.delegate(_val));
        Decode.ok(result, data, "delegate failed");
        return Decode.delegate(data);
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

    function _withdraw(string memory _val) internal returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.withdraw(_val)
        );
        Decode.ok(result, data, "withdraw failed");
        return Decode.withdraw(data);
    }

    function _approveShares(
        string memory _val,
        address _spender,
        uint256 _shares
    ) internal returns (bool) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.approveShares(_val, _spender, _shares)
        );
        Decode.ok(result, data, "approve shares failed");
        return Decode.approveShares(data);
    }

    function _transferShares(
        string memory _val,
        address _to,
        uint256 _shares
    ) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.transferShares(_val, _to, _shares)
        );
        Decode.ok(result, data, "transfer shares failed");
        return Decode.transferShares(data);
    }

    function _transferFromShares(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) internal returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.call(
            Encode.transferFromShares(_val, _from, _to, _shares)
        );
        Decode.ok(result, data, "transferFrom shares failed");
        return Decode.transferFromShares(data);
    }

    function _delegation(
        string memory _val,
        address _del
    ) internal view returns (uint256, uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.delegation(_val, _del)
        );
        Decode.ok(result, data, "delegation failed");
        return Decode.delegation(data);
    }

    function _delegationRewards(
        string memory _val,
        address _del
    ) internal view returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.delegationRewards(_val, _del)
        );
        Decode.ok(result, data, "delegationRewards failed");
        return Decode.delegationRewards(data);
    }

    function _allowanceShares(
        string memory _val,
        address _owner,
        address _spender
    ) internal view returns (uint256) {
        (bool result, bytes memory data) = _stakingAddress.staticcall(
            Encode.allowanceShares(_val, _owner, _spender)
        );
        Decode.ok(result, data, "allowance shares failed");
        return Decode.allowanceShares(data);
    }
}
