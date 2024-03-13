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

const groupBy = (arr, property) => {
    return arr.reduce((acc, obj) => {
        const key = obj[property];
        if (!acc[key]) {
            acc[key] = [];
        }
        acc[key].push(obj);
        return acc;
    }, {});
};

function getStarscanUrl() {
    let starscanUrl = "https://starscan.io";
    const network = getQueryVariable("network")
    if (network === "testnet") {
        starscanUrl = "https://testnet.starscan.io";
    }
    return starscanUrl
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

function getRestUrl() {
    let restUrl = "https://fx-rest.functionx.io";
    const network = getQueryVariable("network")
    if (network === "testnet") {
        restUrl = "https://testnet-fx-rest.functionx.io";
    }
    return restUrl
}

function getJsonrpcUrl() {
    let jsonrpcUrl = "https://fx-json.functionx.io:26657";
    const network = getQueryVariable("network")
    if (network === "testnet") {
        jsonrpcUrl = "https://testnet-fx-json.functionx.io:26657";
    }
    return jsonrpcUrl
}

function getValidators() {
    const result = httpGetJson(getRestUrl() + "/cosmos/staking/v1beta1/validators")
    return result.validators
}

function getSigningInfos() {
    const result = httpGetJson(getRestUrl() + "/cosmos/slashing/v1beta1/signing_infos")
    return result.info
}

function getBalances(address) {
    const result = httpGetJson(getRestUrl() + `/cosmos/bank/v1beta1/balances/${address}`)
    return result.balances
}

function getTotalSupply() {
    const result = httpGetJson(getRestUrl() + `/cosmos/bank/v1beta1/supply`)
    return result.supply
}

function getMetadatas() {
    const result = httpGetJson(getRestUrl() + `/cosmos/bank/v1beta1/denoms_metadata`)
    return result.metadatas
}

function getCurrentPlan() {
    const result = httpGetJson(getRestUrl() + `/cosmos/upgrade/v1beta1/current_plan`)
    return result.plan
}

function getBlock(height) {
    const result = httpGetJson(getJsonrpcUrl() + "/block" + (height ? "?height=" + height : ""))
    return result.block
}

function getBlockSignatures() {
    return getBlock().last_commit.signatures
}

function getPeers() {
    const result = httpGetJson(getJsonrpcUrl() + "/net_info")
    return result.peers
}

function getAverageBlockTime(block) {
    const nowHeight = Number(block.header.height);
    const nowTime = block.header.time;
    const oldBlock = getBlock(nowHeight - 20000);
    const oldTime = oldBlock.header.time;
    const avgBlockTime = Math.round((new Date(nowTime) - new Date(oldTime)) / 20000);
    console.log("Average BlockTime", avgBlockTime);
    return avgBlockTime;
}
