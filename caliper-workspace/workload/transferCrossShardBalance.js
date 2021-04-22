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
        senderId = this.roundArguments.prefix + Math.floor(Math.random() * this.roundArguments.assets)
        if (Math.floor(Math.random() * 100)<= this.roundArguments.selfShard) {
            recipientId = this.roundArguments.prefix + Math.floor(Math.random() * this.roundArguments.assets)
            isB = false
        } else {
            recipientId = this.roundArguments.otherPrefix + this.txIndex % this.roundArguments.assets
            isB = true
        }

        let args = {
            contractId: this.roundArguments.contractId,
            channel: this.roundArguments.channel,
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
            if (isB) {
                args = {
                    contractId: this.roundArguments.otherContractId,
                    channel: this.roundArguments.otherChannel,
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
            } else {
                args = {
                    contractId: this.roundArguments.contractId,
                    channel: this.roundArguments.channel,
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
            }
            result = await this.sutAdapter.sendRequests(args);
            if (result.GetStatus()!="success") {
                args = {
                    contractId: this.roundArguments.contractId,
                    channel: this.roundArguments.channel,
                    contractVersion: 'v1',
                    contractFunction: 'UnlockSenderMoney',
                    contractArguments: [JSON.stringify({"senderID":senderId,"recipientID":recipientId,"amount":10,"tx":txId})],
                    timeout: 30,
                };
                await this.sutAdapter.sendRequests(args);
            }else{
                args = {
                    contractId: this.roundArguments.contractId,
                    channel: this.roundArguments.channel,
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
                result = await this.sutAdapter.sendRequests(args);
            }
        }


    }
}


function createWorkloadModule() {
    return new queryAccountWorkLoad();
}

module.exports.createWorkloadModule = createWorkloadModule;