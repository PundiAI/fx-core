function getQueryVariable(variable) {
    const query = window.location.search.substring(1);
    const vars = query.split("&");
    for (let i = 0; i < vars.length; i++) {
        let pair = vars[i].split("=");
        if (pair[0] === variable) {
            return pair[1];
        }
    }
    return undefined
}

function toHexString(byteArray) {
    return Array.from(byteArray, function (byte) {
        return ('0' + (byte & 0xFF).toString(16).toUpperCase()).slice(-2);
    }).join('')
}

function httpGetJson(theUrl) {
    let xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", theUrl, false); // false for synchronous request
    xmlHttp.send(null);
    const response = JSON.parse(xmlHttp.responseText)
    if (response.result) {
        return response.result;
    }
    return response;
}

function getFxRestUrl(){
    let fxRestUrl = "https://fx-rest.functionx.io";
    const network = getQueryVariable("network")
    if (network === "testnet") {
        fxRestUrl = "https://testnet-fx-rest.functionx.io";
    }
    return fxRestUrl
}

function getFxJsonrpcUrl(){
    let fxJsonrpcUrl = "https://fx-json.functionx.io:26657";
    const network = getQueryVariable("network")
    if (network === "testnet") {
        fxJsonrpcUrl = "https://testnet-fx-json.functionx.io:26657";
    }
    return fxJsonrpcUrl
}

function getAverageBlockTime(block) {
    const nowHeight = Number(block.block.header.height);
    const nowTime = block.block.header.time;
    const oldBlock = httpGetJson(getFxJsonrpcUrl() + "/block?height=" + (nowHeight - 100).toString());
    const oldTime = oldBlock.block.header.time;
    const avgBlockTime = Math.round((new Date(nowTime) - new Date(oldTime)) / 100);
    console.log("Average BlockTime", avgBlockTime);
    return avgBlockTime;
}

function getValidators() {
    const result = httpGetJson(getFxRestUrl() + "/cosmos/staking/v1beta1/validators")
    return result.validators
}

function getSigningInfos() {
    const result = httpGetJson(getFxRestUrl() + "/cosmos/slashing/v1beta1/signing_infos")
    return result.info
}

function getBlockSignatures() {
    const result = httpGetJson(getFxJsonrpcUrl() + "/block")
    return result.block.last_commit.signatures
}

function getPeers() {
    const result = httpGetJson(getFxJsonrpcUrl() + "/net_info")
    return result.peers
}

function getBalances(address) {
    const result = httpGetJson(getFxRestUrl() + `/cosmos/bank/v1beta1/balances/${address}`)
    return result.balances
}

function getTotalSupply() {
    const result = httpGetJson(getFxRestUrl() + `/cosmos/bank/v1beta1/supply`)
    return result.supply
}

function getMetadatas() {
    const result = httpGetJson(getFxRestUrl() + `/cosmos/bank/v1beta1/denoms_metadata`)
    return result.metadatas
}

function getCurrentPlan() {
    const result = httpGetJson(getFxRestUrl() + `/cosmos/upgrade/v1beta1/current_plan`)
    return result.plan
}