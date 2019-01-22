## Example of a blockchain smart-contract for Digital Rights Management

> NOTE: The materials refererened in this README will be made available when they can be distributed under the [open source license](https://github.com/mitmedialab/2019-MIT-Computational-Law-Course/blob/master/LICENSE) designated for files in this repository. 

The smart-contract is designed for Hyperledger Fabric platform:
https://www.hyperledger.org/projects/fabric

The code is implemented in Go language and has the following functions:

- createWallet – creates a new wallet for a participant
- createLicense – creates a new license
- searchForLicense – searches for all licenses based on criteria
- transferRights – transfers rights to another participant, calculates fees and updates balances in wallets

The smart-contract code is in the drm.go file.
The drm_test.go file contains a suite of tests of the smart-contract. These can be executed by the following command: “go test” on the folder.

The smart-contract can be deployed on both, a stand-alone or cloud-based instance of Fabric.

IBM operates Hyperledger Fabric based IBM Blockchain Platform within Bluemix cloud and offers a free Starter Plan to run blockchain networks:
- High-level overview: https://www.ibm.com/blockchain/platform
- Blockchain services within Bluemix: https://console.bluemix.net/catalog/services/blockchain
- Documentation: https://console.bluemix.net/docs/services/blockchain/index.html#ibm-blockchain-platform

To connect, execute and integrate the smart-contract into an application, follow this page:
https://console.bluemix.net/docs/services/blockchain/v10_application.html#dev_app

In order to learn about the design of smart-contracts and Fabric enabled solutions, refer to our book: **Hands-On Blockchain with Hyperledger: Building decentralized applications with Hyperledger Fabric and Composer**
https://rebrand.ly/hobdwh 
