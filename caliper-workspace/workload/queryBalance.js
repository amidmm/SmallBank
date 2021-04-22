'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const helper = require('./helper');


class queryAccountWorkLoad extends WorkloadModuleBase {

    constructor() {
        super();
        this.txIndex = 0;
        this.limitIndex = 0;
    }
    
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.limitIndex = this.roundArguments.assets;
        // await helper.createAccount(this.sutAdapter, this.workerIndex, this.roundArguments);
    }

    async submitTransaction() {
        this.txIndex++;
        let id 
        if (Math.floor(Math.random() * 100)<= this.roundArguments.selfShard) {
            id = this.roundArguments.prefix + Math.floor(Math.random() * this.roundArguments.assets)
        } else {
            id = this.roundArguments.otherPrefix + this.txIndex % this.roundArguments.assets
        }
        let args = {
            contractId: this.roundArguments.contractId,
            channel: this.roundArguments.channel,
            contractVersion: 'v1',
            contractFunction: 'GetBalance',
            contractArguments: [id],
            timeout: 30,
            readOnly: true
        };

        if (this.txIndex === this.limitIndex) {
            this.txIndex = 0;
        }
        await this.sutAdapter.sendRequests(args);
    }
}


function createWorkloadModule() {
    return new queryAccountWorkLoad();
}

module.exports.createWorkloadModule = createWorkloadModule;
