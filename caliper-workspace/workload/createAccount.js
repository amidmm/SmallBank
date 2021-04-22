'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const helper = require('./helper');


let names = ["Kadie Millington","Bonnie Glenn","Iris Chapman","Mujtaba Cresswell","Jensen Rennie","Faizaan Bradshaw","Kaif Finnegan","Eamon Holt","Chance West","Archibald Dawe","Arman Olsen","Brenda Combs","Kaleb Hensley","Lex Schneider","Madeleine O'Quinn","Mikolaj Nicholson","Giorgia Goddard","Katie Blackwell","Quentin Fenton","Amelia-Grace Fry"];
let name;
let balance;

class queryAccountWorkLoad extends WorkloadModuleBase {

    constructor() {
        super();
        this.txIndex = 0;
        this.limitIndex = 0;
    }
    
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.limitIndex = this.roundArguments.assets;
        await helper.createAccount(this.sutAdapter, this.workerIndex, this.roundArguments);
    }

    async submitTransaction() {
        this.txIndex++;
        let id = this.txIndex
        name = names[Math.floor(Math.random() * names.length)];
        balance = Math.floor(Math.random() * 100000000);
        // console.log(id)

        let myArgsA = {
            contractId: "banka",
            channel: "bankachannel",
            contractVersion: 'v1',
            contractFunction: 'CreateAccount',
            contractArguments: [JSON.stringify({"id":id.toString(),"name":name,"balance":balance})],
            timeout: 30
        };

        let result = await this.sutAdapter.sendRequests(myArgsA);

        let myArgsB = {
            contractId: "bankb",
            channel: "bankbchannel",
            contractVersion: 'v1',
            contractFunction: 'CreateAccount',
            contractArguments: [JSON.stringify({"id":id.toString(),"name":name,"balance":balance})],
            timeout: 30
        };
        result = await this.sutAdapter.sendRequests(myArgsB);
        id++
    }
}


function createWorkloadModule() {
    return new queryAccountWorkLoad();
}

module.exports.createWorkloadModule = createWorkloadModule;
