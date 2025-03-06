// SPDX-License-Identifier: Apache-2.0

pragma solidity ^0.8.0;

interface IStaking {
    enum ValidatorSortBy {
        Power,
        Missed
    }

    function delegateV2(
        string memory _val,
        uint256 _amount
    ) external returns (bool _result);

    function undelegateV2(
        string memory _val,
        uint256 _amount
    ) external returns (bool _result);

    function redelegateV2(
        string memory _valSrc,
        string memory _valDst,
        uint256 _amount
    ) external returns (bool _result);

    function withdraw(string memory _val) external returns (uint256 _reward);

    function approveShares(
        string memory _val,
        address _spender,
        uint256 _shares
    ) external returns (bool _result);

    function transferShares(
        string memory _val,
        address _to,
        uint256 _shares
    ) external returns (uint256 _token, uint256 _reward);

    function transferFromShares(
        string memory _val,
        address _from,
        address _to,
        uint256 _shares
    ) external returns (uint256 _token, uint256 _reward);

    function delegation(
        string memory _val,
        address _del
    ) external view returns (uint256 _shares, uint256 _delegateAmount);

    function delegationRewards(
        string memory _val,
        address _del
    ) external view returns (uint256 _reward);

    function allowanceShares(
        string memory _val,
        address _owner,
        address _spender
    ) external view returns (uint256 _shares);

    function slashingInfo(
        string memory _val
    ) external view returns (bool _jailed, uint256 _missed);

    function validatorList(
        ValidatorSortBy _sortBy
    ) external view returns (string[] memory);

    event DelegateV2(
        address indexed delegator,
        string validator,
        uint256 amount
    );

    event UndelegateV2(
        address indexed sender,
        string validator,
        uint256 amount,
        uint256 completionTime
    );

    event RedelegateV2(
        address indexed sender,
        string valSrc,
        string valDst,
        uint256 amount,
        uint256 completionTime
    );

    event Withdraw(address indexed sender, string validator, uint256 reward);

    event ApproveShares(
        address indexed owner,
        address indexed spender,
        string validator,
        uint256 shares
    );

    event TransferShares(
        address indexed from,
        address indexed to,
        string validator,
        uint256 shares,
        uint256 token
    );
}
