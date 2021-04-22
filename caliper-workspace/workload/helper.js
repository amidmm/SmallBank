/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

let names = ["Kadie Millington","Bonnie Glenn","Iris Chapman","Mujtaba Cresswell","Jensen Rennie","Faizaan Bradshaw","Kaif Finnegan","Eamon Holt","Chance West","Archibald Dawe","Arman Olsen","Brenda Combs","Kaleb Hensley","Lex Schneider","Madeleine O'Quinn","Mikolaj Nicholson","Giorgia Goddard","Katie Blackwell","Quentin Fenton","Amelia-Grace Fry"];


let id= 0;
let name;
let balance;
let txIndex = 0;

module.exports.createAccount = async function (bc, workerIndex, args) {
    while (txIndex < args.assets) {
        txIndex++;
        name = names[Math.floor(Math.random() * names.length)];
        balance = Math.floor(Math.random() * 100000000);


        let myArgsA = {
            contractId: "banka",
            channel: "bankachannel",
            contractVersion: 'v1',
            contractFunction: 'CreateAccount',
            contractArguments: [JSON.stringify({"id":id.toString(),"name":name,"balance":balance})],
            timeout: 30
        };

        let myArgsB = {
            contractId: "bankb",
            channel: "bankbchannel",
            contractVersion: 'v1',
            contractFunction: 'CreateAccount',
            contractArguments: [JSON.stringify({"id":id.toString(),"name":name,"balance":balance})],
            timeout: 30
        };
        
        bc.sendRequests(myArgsA);
        await bc.sendRequests(myArgsB);
        console.log(id)
        id++
    }
};
