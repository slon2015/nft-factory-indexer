// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";
import {NFTFactory} from '../src/NFTFactory.sol';

contract NFTFactoryTest {
    function test_DeployCollection() public {
        NFTFactory factory = new NFTFactory();

        factory.createCollection("Collection", "CLTN", "http://base.link/");
    }

    function test_MintToken() public {
        NFTFactory factory = new NFTFactory();

        address collection = factory.createCollection("Collection", "CLTN", "http://base.link/");

        factory.mint(collection, msg.sender, 10);
    }
}