// SPDX-License-Identifier: UNLICENSED
/* solhint-disable one-contract-per-file */
pragma solidity ^0.8.0;

interface IBridgeFeeQuote {
    struct QuoteInput {
        string chainName;
        address token;
        address oracle;
        uint256 fee;
        uint256 gasLimit;
        uint256 expiry;
        bytes signature;
    }

    struct QuoteInfo {
        uint256 id;
        string chainName;
        address token;
        address oracle;
        uint256 fee;
        uint256 gasLimit;
        uint256 expiry;
    }

    function quote(QuoteInput[] memory _inputs) external returns (bool);

    function getQuoteList(
        string memory _chainName
    ) external view returns (QuoteInfo[] memory);

    function getQuoteById(uint256 _id) external view returns (QuoteInfo memory);

    function getQuoteByToken(
        string memory _chainName,
        address _token,
        uint256 _amount
    ) external view returns (QuoteInfo memory, bool);
}

interface IBridgeFeeOracle {
    function defaultOracle() external view returns (address);

    function isOnline(
        string memory _chainName,
        address _oracle
    ) external returns (bool);

    function getOracleList() external view returns (address[] memory);
}
