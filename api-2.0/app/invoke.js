const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const EventStrategies = require('fabric-network/lib/impl/event/defaulteventhandlerstrategies');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')

const helper = require('./helper');
const query = require('./query');

const { blockListener, contractListener } = require('./Listeners');
const channelName = "mychannel"
const chaincodeName = "fabcar"
var org_name 


const invokeTransaction = async (fcn,username,args) => {
    try {
        org_name = "Org2";
        const ccp = await helper.getCCP(org_name);

        const walletPath = await helper.getWalletPath(org_name);
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            return;
        }

        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: true }
            // eventHandlerOptions: EventStrategies.NONE
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);

        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        // await contract.addContractListener(contractListener);
        // await network.addBlockListener(blockListener);


        // Multiple smartcontract in one chaincode
        let result;
        let message;

        switch (fcn) {
            case "CreateData":
            case "ChangeData":
                console.log(`User name is ${username}`)
                var new_args = {};
                new_args["Name"] = args["Name"];
                new_args["AadharNumber"] = args["AadharNumber"];
                new_args["PhoneNumber"] = args["PhoneNumber"];
                new_args["Status"] = "inactive";
                new_args["Money"] = 0;
                new_args["Transaction_type"] = "info";
                console.log(JSON.stringify(new_args));
                result = await contract.submitTransaction('SmartContract:'+fcn, JSON.stringify(new_args));
                result = {txid: result.toString()}
                break;

            case "BuyService":
                result = await contract.submitTransaction('SmartContract:'+fcn, args["Service_name"],args["Price"]);
                result = {txid: result.toString()}
                break;
            // case "CreateAadharData":
            // case "CreateDrivingLicenceData":
            //     result = await contract.submitTransaction('SmartContract:'+fcn, JSON.stringify(args));
            //     result = {txid: result.toString()}
            //     break;
            default:
                break;
        }

        await gateway.disconnect();

        // result = JSON.parse(result.toString());

        let response = {
            message: message,
            result
        }

        return response;


    } catch (error) {

        console.log(`Getting error: ${error}`)
        return error.message

    }
}

exports.invokeTransaction = invokeTransaction;