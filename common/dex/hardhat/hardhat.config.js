require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.19",
  networks: {
    gethDev: {
      url: "http://127.0.0.1:8547",
      accounts: ["0x140068712d8b1e77451ec8962055ea6e14408cbf47eb04c468cb8886513f1f34"]
    }
  }
};
