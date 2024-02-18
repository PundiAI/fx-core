// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {ERC721Burnable} from "@openzeppelin/contracts/token/ERC721/extensions/ERC721Burnable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract ERC721TokenTest is ERC721, ERC721Burnable, Ownable {
    constructor(
        string memory name,
        string memory symbol
    ) ERC721(name, symbol) {}

    function mint(address _to, uint256 _id) external onlyOwner {
        _safeMint(_to, _id);
    }

    function _baseURI() internal pure override returns (string memory) {
        return "ipfs://test-url";
    }
}
