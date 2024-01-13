import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-chai-matchers";
import "@nomicfoundation/hardhat-ethers";

const config: HardhatUserConfig = {
  solidity: "0.8.19",
  defaultNetwork: "sandbox",
  networks: {
    sandbox: {
      url: "http://localhost:8080/wasp/api/v1/chains/snd1pr6fw846t2ugw74urqjlxzl4jrkayu80w8p79e3vpd3s4psus9p8c0hyxsu/evm",
      accounts: {
        mnemonic:
          "echo mass dentist hood neutral claw fragile scare magnet float citizen pink order hundred village edge quality excuse donor spawn heart happy extend resist",
      },
    },
  },
};

export default config;
