'use strict';

const {
    WorkloadModuleBase
} = require('@hyperledger/caliper-core');


class queryAccountWorkLoad extends WorkloadModuleBase {

    constructor() {
        super();
        this.txIndex = 0;
        this.limitIndex = 0;
    }


    async submitTransaction() {
        this.txIndex++;
        let senderId, recipientId
        let isB 
        if (Math.floor(Math.random() * 100)<= this.roundArguments.selfShard) {
            senderId = this.roundArguments.prefix + Math.floor(Math.random() * this.roundArguments.assets)
            recipientId = this.roundArguments.prefix + Math.floor(Math.random() * this.roundArguments.assets)
            isB = false
        } else {
            senderId = this.roundArguments.otherPrefix + Math.floor(Math.random() * this.roundArguments.assets)
            recipientId = this.roundArguments.otherPrefix + this.txIndex % this.roundArguments.assets
            isB = true
        }
        let contractId,channelId
        if (!isB) {
            contractId = "banka"
            channelId = "bankachannel"
        }else{
            contractId = "bankb"
            channelId = "bankbchannel"
        }

        let args = {
            contractId: contractId,
            channel: channelId,
            contractVersion: 'v1',
            contractFunction: 'LockSenderMoney',
            contractArguments: [JSON.stringify({
                "senderID": senderId,
                "recipientID": recipientId,
                "amount": 10,
                "tx": ""
            })],
            timeout: 30,
        };

        if (this.txIndex === this.limitIndex) {
            this.txIndex = 0;
        }
        let result = await this.sutAdapter.sendRequests(args);
        if (result.GetStatus() == "success") {
            let txId = Buffer.from(Buffer.from(result.GetResult()).toJSON().data).toString()
                args = {
                    contractId: contractId,
                    channel: channelId,
                    contractVersion: 'v1',
                    contractFunction: 'ReceiveMoney',
                    contractArguments: [JSON.stringify({
                        "senderID": senderId,
                        "recipientID": recipientId,
                        "amount": 10,
                        "tx": txId
                    })],
                    timeout: 30,
                };
            result = await this.sutAdapter.sendRequests(args);
            if (result.GetStatus()!="success") {
                args = {
                    contractId: contractId,
                    channel: channelId,
                    contractVersion: 'v1',
                    contractFunction: 'UnlockSenderMoney',
                    contractArguments: [JSON.stringify({"senderID":senderId,"recipientID":recipientId,"amount":10,"tx":txId})],
                    timeout: 30,
                };
                 this.sutAdapter.sendRequests(args);
            }else{
                args = {
                    contractId: contractId,
                    channel: channelId,
                    contractVersion: 'v1',
                    contractFunction: 'FinializeTransfer',
                    contractArguments: [JSON.stringify({
                        "senderID": senderId,
                        "recipientID": recipientId,
                        "amount": 10,
                        "tx": txId
                    })],
                    timeout: 30,
                };
                this.sutAdapter.sendRequests(args);
            }
        }


    }
}


function createWorkloadModule() {
    return new queryAccountWorkLoad();
}

module.exports.createWorkloadModule = createWorkloadModule;