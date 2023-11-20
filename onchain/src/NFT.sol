// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract NFT is ERC721, Ownable {

    string private uri;

    constructor(string memory name_, string memory symbol_, string memory uri_) ERC721(name_, symbol_) Ownable(_msgSender()) {
       uri = uri_;
    }

    function _baseURI() override internal view returns(string memory) {
        return uri;
    }

    function mint(address to, uint256 tokenId) external onlyOwner {
        _safeMint(to, tokenId);
    }
}