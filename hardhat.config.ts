import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

const config: HardhatUserConfig = {
  solidity: "0.8.19",
  defaultNetwork: "sandbox",
  networks: {
    sandbox: {
      url: "http://localhost:8080/wasp/api/v1/chains/snd1prdgxuwvy8legrgxvdusfyu77kcftf6yvx85wmz7ar8swg3rh538sqq22ut/evm",
      accounts: {
        mnemonic:
          "echo mass dentist hood neutral claw fragile scare magnet float citizen pink order hundred village edge quality excuse donor spawn heart happy extend resist",
      },
    },
  },
};

export default config;
