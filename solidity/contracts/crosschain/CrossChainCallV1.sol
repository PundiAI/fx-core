// SPDX-License-Identifier: Apache-2.0

/* solhint-disable one-contract-per-file */
pragma solidity ^0.8.0;

library Decode {
    // Deprecated: for fip20 only
    function fip20CrossChain(bytes memory data) internal pure returns (bool) {
        bool result = abi.decode(data, (bool));
        return result;
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

library Encode {
    // Deprecated: for fip20 only
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal pure returns (bytes memory) {
        return
            abi.encodeWithSignature(
                "fip20CrossChain(address,string,uint256,uint256,bytes32,string)",
                _sender,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            );
    }
}

// Deprecated: Please use CrossChainCall
library CrossChainCallV1 {
    address public constant CROSS_CHAIN_ADDRESS =
        address(0x0000000000000000000000000000000000001004);

    // Deprecated: for fip20 only
    function fip20CrossChain(
        address _sender,
        string memory _receipt,
        uint256 _amount,
        uint256 _fee,
        bytes32 _target,
        string memory _memo
    ) internal returns (bool) {
        // solhint-disable-next-line avoid-low-level-calls
        (bool result, bytes memory data) = CROSS_CHAIN_ADDRESS.call(
            Encode.fip20CrossChain(
                _sender,
                _receipt,
                _amount,
                _fee,
                _target,
                _memo
            )
        );
        Decode.ok(result, data, "fip-cross-chain failed");
        return Decode.fip20CrossChain(data);
    }
}
