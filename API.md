# NKN.OVH API


# Basic information about structures

API requests have the following structures:

Golang structure:
```go
type WSQuery struct {
	Method string `json:"Method"`
	Value map[string]interface{} `json:"Value,omitempty"`
}

```

JSON structure:
```json
{
"Method":"",
"Value": {}
}

```
______

The server responds to incoming API requests in the format of JSON and has the following structures:

Golang structure:
```go
type WSReply struct {
	Method string `json:"Method"`
	Code int `json:"Code"`
	Error bool `json:"Error,omitempty`
	ErrMessage string `json:"ErrMessage,omitempty"`
	Value interface{} `json:"Value, omitempty"`
}
```

JSON sturcture:
```json
{
"Method": "(string)",
"Code": (integer),
"Error": (boolean, can be absent/empty),
"ErrMessage": "(string, can be absent/empty)",
"Value": (any types but usually it contains JSON objects, can be absent/empty)
}
```
________

The server can receive only POST requests in the following formats:
 - application/json
 - application/x-www-form-urlencoded


**Note!**  
If you use "application/x-www-form-urlencoded" requests, all the values of your request's body will be automatically imported as WsQuery["Value"]["KEY"] = %KEY'S_VALUE% (except for the WsQuery["Method"] value)

________

# Response codes

All response codes (except for the codes: 0, 2, and 3) inform of occurred errors.

- If a request fails, the Error key will be "true" and the ErrMessage key will contain a short description.
- The Code 0 means that your request has completed correctly.
- The Code 2 refers to "addnodes" method and the code means that passed nodes added partially (Error: false; because the request has made some changes)
- The Code 3 means that the server did not find requested information (but Error: false; because it is normal behavior)

```go

o.Web.Response[1] = WSReply{Code: 1, Error: true, ErrMessage: "Cannot execute SQL query"}
o.Web.Response[2] = WSReply{Code: 2, Error: false, ErrMessage: "Nodes added partially"}
o.Web.Response[3] = WSReply{Code: 3, Error: false, ErrMessage: "No info/entries in a database"}
o.Web.Response[4] = WSReply{Code: 4, Error: true, ErrMessage: "An error occured"}
o.Web.Response[5] = WSReply{Code: 5, Error: true, ErrMessage: "Incorrect query"}
o.Web.Response[6] = WSReply{Code: 6, Error: true, ErrMessage: "Name of node(-s) too long or incorrect format."}
o.Web.Response[7] = WSReply{Code: 7, Error: true, ErrMessage: "Multiple variable must be boolean"}
o.Web.Response[8] = WSReply{Code: 8, Error: true, ErrMessage: "Ip variable must be valid string"}
o.Web.Response[9] = WSReply{Code: 9, Error: true, ErrMessage: "Incorrect ip address(-es)"}
o.Web.Response[10] = WSReply{Code: 10, Error: true, ErrMessage: "Passed ip address(-es) was not IPv4"}
o.Web.Response[11] = WSReply{Code: 11, Error: true, ErrMessage: "Passed ip address(-es) was not in public network"}
o.Web.Response[12] = WSReply{Code: 12, Error: true, ErrMessage: "Nodes limit exceeded"}
o.Web.Response[13] = WSReply{Code: 13, Error: true, ErrMessage: "Wrong delimiter passed"}
o.Web.Response[14] = WSReply{Code: 14, Error: true, ErrMessage: "The nodes weren't added since they had been already created."}
o.Web.Response[15] = WSReply{Code: 15, Error: true, ErrMessage: "Wrong data of NodesId passed"}
o.Web.Response[16] = WSReply{Code: 16, Error: true, ErrMessage: "Wallets overflow"}
o.Web.Response[17] = WSReply{Code: 17, Error: true, ErrMessage: "One or more of the passed wallets are not in the correct format"}
o.Web.Response[18] = WSReply{Code: 18, Error: true, ErrMessage: "One or more IDs of the passed nodes are not found. No changes were made."}

//Link to apiGetNodeDetails
o.Web.Response[19] = WSReply{Code: 19, Error: true, ErrMessage: "Wrong data of NodeId passed"}
o.Web.Response[20] = WSReply{Code: 20, Error: true, ErrMessage: "The node is offline / No reply recieved from the node within the timeout period"}
o.Web.Response[21] = WSReply{Code: 21, Error: true, ErrMessage: "Cannot decode json of the node response (getnodestate)"}
o.Web.Response[22] = WSReply{Code: 22, Error: true, ErrMessage: "The node is online, but no information about neighbors has been received within the timeout period"}
o.Web.Response[23] = WSReply{Code: 23, Error: true, ErrMessage: "Cannot decode json of the node response (getneighbor)"}
o.Web.Response[24] = WSReply{Code: 24, Error: true, ErrMessage: "Query returned an error (getneighbor)"}

//Link to apiGetNodeIpByPublicKey
o.Web.Response[25] = WSReply{Code: 25, Error: true, ErrMessage: "PublicKey is not set"}
o.Web.Response[26] = WSReply{Code: 26, Error: true, ErrMessage: "Wrong PublicKey passed"}

//Link to apiRmNodesByIp
o.Web.Response[27] = WSReply{Code: 27, Error: true, ErrMessage: "Wrong data of NodesIp passed"}
o.Web.Response[28] = WSReply{Code: 28, Error: true, ErrMessage: "One or more IPs of the passed nodes weren't found. No changes were made."}

o.Web.Response[230] = WSReply{Code: 230, Error: true, ErrMessage: "No view variable passed, the variable must be string"}
o.Web.Response[231] = WSReply{Code: 231, Error: true, ErrMessage: "No Locale variable passed, the variable must be string"}
o.Web.Response[232] = WSReply{Code: 232, Error: true, ErrMessage: "Locale or View passed variables were overflowed"}
o.Web.Response[233] = WSReply{Code: 233, Error: true, ErrMessage: "Requested language pack was not found in the package slice"}
o.Web.Response[234] = WSReply{Code: 234, Error: true, ErrMessage: "Passed language pack was not found as JSON file"}

o.Web.Response[240] = WSReply{Code: 240, Error: true, ErrMessage: "GenRandomSHA256 returned error"}
o.Web.Response[252] = WSReply{Code: 252, Error: true, ErrMessage: "You have already created 3 IDs for the latest 30 minutes"}
o.Web.Response[253] = WSReply{Code: 253, Error: true, ErrMessage: "You have no authorization"}
o.Web.Response[254] = WSReply{Code: 254, Error: true, ErrMessage: "Incorrect ID's length"}
o.Web.Response[255] = WSReply{Code: 255, Error: true, ErrMessage: "Passed ID was not found"}
o.Web.Response[500] = WSReply{Code: 500, Error: true, ErrMessage: "Internal server error"}

//Main errors
o.Web.Response[1000] = WSReply{Code: 1000, Error: true, ErrMessage: "The method variable didn't pass or was in a wrong format"}
o.Web.Response[1001] = WSReply{Code: 1001, Error: true, ErrMessage: "The passed Method was not found"}
o.Web.Response[1002] = WSReply{Code: 1002, Error: true, ErrMessage: "Connections limit has been reached"}
o.Web.Response[1003] = WSReply{Code: 1003, Error: true, ErrMessage: "Passed JSON is incorrect"}
```
________

# API methods

The methods are divided into methods that require authorization and methods that don't require it.  
All the methods are case-sensitive and must be typed in lower case.

The following methods require authorization:

```
addnodes
rmnodes
rmnodesbyip
getfullstack
getmynodes
getnetstatus
getmywallets
getprices
getnodedetails
savemysettings
```
  
The following methods don't require authorization:

```
genid
getnodeipbypublickey
getdaemon
```
  
The following methods are useless if called by POST requests and work well only with WebSocket.  
So you should not use them with POST requests:

```
auth
getlanguage
logout
```

# Examples using API

## Method `genid`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to generate a client's ID (Hash).  


#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{"Method":"genid"}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

  
#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=genid' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```

  
#### The server returns:
```json
{
"Method":"genid",
"Code":0,
"Error":false,
"Value":
	{
	"Hash":"3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}
```

</details>

## Method `addnodes`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to add nodes into your account.  

- The key "Multiple" (boolean) must be set to a valid boolean (true/false for json requests or t/true/1 and f/false/0 for urlencoded requests)  
- The key "Name" (string) must be set but can be an empty string.   
- The key "Ip" (string) must be set:  
If the **Multiple** key is set to false, the **Ip** must contain a single IP address.  
If the **Multiple** key is set to true, the **Ip** must contain at least two IP addresses which are separated by either commas, spaces or line breaks.  
- The key "Hash" (string) must be set to valid Hash.  



#### cURL example #1 (application/json), adding a single node:  

```bash
curl -X POST  \
-d '{
"Method":"addnodes",
"Value": {
	"Multiple": false,
	"Name": "MySuperNode",
	"Ip": "1.1.1.1",
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded), adding a single node:  

```bash
curl -X POST  \
-d 'Method=addnodes' \
-d 'Multiple=false' \
-d 'Name=MySuperNode' \
-d 'Ip=1.1.1.1' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```

#### The server returns:

```json

{
"Method": "addnodes",
"Code": 0,
"Error": false,
"Value":
	{
	"Info": "Your node added"
	}
}

```
</details>

## Method `rmnodesbyip`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to delete nodes by IP addresses.

- The key "NodesIp" must be set:  
If you use urlencoded request, the key must contain **string** with at least one IP address. Multiple IP addresses must be separated by comma.  
If you use json request, the key must contain **array of strings** with at least one IP address. Multiple IP addresses must be separate elements of the array.  
- The key "Hash" (string) must be set to valid Hash.

#### cURL example #1 (application/json), deleting multiple nodes:  

```bash
curl -X POST  \
-d '{
"Method":"rmnodesbyip",
"Value": {
	"NodesIp": ["1.1.1.1", "1.2.3.4", "8.8.8.8"],
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded), deleting multiple nodes:  

```bash
curl -X POST  \
-d 'Method=rmnodesbyip' \
-d 'NodesIp=1.1.1.1, 1.2.3.4, 8.8.8.8' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```

#### The server returns:

```json
{
"Method":"rmnodesbyip",
"Code":0,
"Error":false,
"Value":
	{
	"Data":"Nodes removed successfully",
	"NodesId":[98390,98392,98393]
	}
}
```

The **NodesId** key in the returned result contains an array of node IDs which have been removed by your request.
</details>

## Method `rmnodes`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to delete nodes by their ID.

- The key "NodesId" must be set:  
If you use urlencoded request, the key must contain **string** with at least one node id. Multiple node IDs must be separated by commas.  
If you use json request, the key must contain **array of integers** with at least one node id. Multiple node IDs must be separate elements of the array.  
- The key "Hash" (string) must be set to valid Hash.

	
#### cURL example #1 (application/json), deleting multiple nodes:  

```bash
curl -X POST  \
-d '{
"Method":"rmnodes",
"Value": {
	"NodesId": [98390,98392,98393],
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded), deleting multiple nodes:  

```bash
curl -X POST  \
-d 'Method=rmnodes' \
-d 'NodesId=98390,98392,98393' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```

#### The server returns:

```json
{
"Method":"rmnodesbyip",
"Code":0,
"Error":false,
"Value":
	{
	"Data":"Nodes removed successfully",
	"NodesId":[98390,98392,98393]
	}
}
```

The **NodesId** key in the returned result contains an array of node IDs which have been removed by your request.  

</details>

## Method `getmynodes`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get your nodes list.

- The key "Hash" (string) must be set to valid Hash.

	
#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"getmynodes",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=getmynodes' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```

#### The server returns:

```json
"Method":"getmynodes",
"Code":0,
"Error":false,
"Value":
	{
		"List": [{
			"Currtimestamp": 0,
			"Err": 1,
			"Height": 0,
			"Ip": "1.1.1.1",
			"LatestUpdate": "2021-09-16 11:46:00",
			"Name": "NodeName",
			"NodeId": 10792,
			"ProposalSubmitted": -1,
			"RelayMessageCount": 0,
			"RelaysPerHour": 0,
			"RelaysPerHour10": 0,
			"RelaysPerHour60": 0,
			"SyncState": "_OFFLINE_",
			"Uptime": 0,
			"Version": ""
			},
			{
			"Currtimestamp": 1631781964,
			"Height": 3097336,
			"Ip": "2.2.2.2",
			"LatestUpdate": "2021-09-16 11:46:05",
			"Name": "HomeNode",
			"NodeId": 36241,
			"ProposalSubmitted": 0,
			"RelayMessageCount": 10137516,
			"RelaysPerHour": 36704,
			"RelaysPerHour10": 44904,
			"RelaysPerHour60": 40034,
			"SyncState": "PERSIST_FINISHED",
			"Uptime": 994301,
			"Version": "v2.1.6"
			}
		//...
		]
	}
}
```

The **List** key in the returned result is an array of node objects.

Note!  
- If you have no nodes added in your account, the server returns the code 3.
- The **RelayMessageCount** key can contain a big unsigned integer value, so if you use the API to program in a language with strict types, use uint64 type while decoding the value.  
- NKNOVH uses transparent replacement of the **SyncState** key. The key's value may not match an actual value returned by a node. 
	
- The best way to detect nodes which are not mining is to check the **Err** key of a node object.  
If the **Err** key is not found or the key equals 0, the node is online and has a status you can see in the "SyncState" key. 

- Other possible values for ["Value"]["List"][n]["Err"]:

```
Err equals 1: the node is offline (SyncState == "_OFFLINE_") 
Err equals 2: the node is waiting for the first update from nknovh programm. (SyncState == "Waiting for first update")
Err equals 3: the node is online but has the status "Out of Network" (SyncState == "_OUT_")

If a node object has no Err key or Err key equals 0, the node has the SyncState:
SYNC_STARTED
WAIT_FOR_SYNCING
SYNC_FINISHED
PERSIST_FINISHED
PRUNING DB
GENERATING ID
```

</details>

## Method `getnetstatus`
<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get statistics of the NKN Network.

- The key "Hash" (string) must be set to valid Hash.


#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"getnetstatus",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=getnetstatus' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```
	
#### The server returns

```json
{
  "Method": "getnetstatus",
  "Code": 0,
  "Error": false,
  "Value": {
    "relays": 8566525058388,
    "average_uptime": 960430,
    "average_relays": 296567,
    "relays_per_hour": 31332092001,
    "proposalSubmitted": 44227,
    "persist_nodes_count": 100732,
    "nodes_count": 105649,
    "last_height": 3097516,
    "last_timestamp": 1631786147,
    "average_blockTime": 22.589506,
    "average_blocksPerDay": 3824.785,
    "latest_update": "2021-09-16 12:55:47"
  }
}
```
	
- The keys **average_blockTime** and **average_blocksPerDay** are calculated since the NKN Mainnet launch.
- The **last_height** key indicates the highest height of the NKN nodes.
- The **relays** key indicates the summary of NKN nodes' relays. The key's value may be high.


</details>

## Method `getmywallets`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get your wallets and balances.

- The key "Hash" (string) must be set to valid Hash.


#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"getmywallets",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=getmywallets' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```
	
#### The server returns

```json
{
  "Method": "getmywallets",
  "Code": 0,
  "Error": false,
  "Value": {
    "Wallets": [
      {
        "Balance": 343.49090321,
        "Id": 124,
        "NknWallet": "NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE"
      }
    ]
  }
}
```

- The **Wallets** key contains **an array of objects**.
- If you have no wallets in your account, the server returns the code 3.

</details>

## Method `getprices`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get a price of the NKN coin.

- The key "Hash" (string) must be set to valid Hash.


#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"getprices",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=getprices' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```
	
#### The server returns

```json
{
  "Method": "getprices",
  "Code": 0,
  "Error": false,
  "Value": {
    "usd": 0.396526
  }
}
```

</details>

## Method `getdaemon`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get information about the NKNOVH programm.

#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"getdaemon"
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api
```

#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=getdaemon' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```
	
#### The server returns

```json
{
  "Method": "getdaemon",
  "Code": 0,
  "Error": false,
  "Value": {
    "Timezone": "+0300",
    "Version": "1.1.0-dirty-6"
  }
}
```

</details>

## Method `getfullstack`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get information about your nodes, your wallets, prices of the NKN coin, NKN Network's statistics and about the NKNOVH programm  
The method calls methods: getmynodes, getmywallets, getprices, getnetstatus, getdaemon; and returns a single JSON. 

- The key "Hash" (string) must be set to valid Hash.
	

#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"getfullstack",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api

```

#### cURL example #2 (application/x-www-form-urlencoded):  

```bash
curl -X POST  \
-d 'Method=getfullstack' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```
	

#### The server returns

```json
{
  "Method": "getfullstack",
  "Code": 0,
  "Error": false,
  "Value": {
    "Daemon": {
      "Method": "getfullstack",
      "Code": 0,
      "Error": false,
      "Value": {
        "Timezone": "+0300",
        "Version": "1.1.0-dirty-6"
      }
    },
    "Netstatus": {
      "Method": "getfullstack",
      "Code": 0,
      "Error": false,
      "Value": {
        "relays": 9090394530670,
        "average_uptime": 999636,
        "average_relays": 312446,
        "relays_per_hour": 33423030261,
        "proposalSubmitted": 47102,
        "persist_nodes_count": 105555,
        "nodes_count": 106972,
        "last_height": 3098117,
        "last_timestamp": 1631800061,
        "average_blockTime": 22.589615,
        "average_blocksPerDay": 3824.7664,
        "latest_update": "2021-09-16 16:47:41"
      }
    },
    "Nodes": {
      "Method": "getfullstack",
      "Code": 0,
      "Error": false,
      "Value": {
        "List": [
          {
            "Currtimestamp": 1631800205,
            "Height": 3098125,
            "Ip": "1.1.1.1",
            "LatestUpdate": "2021-09-16 16:50:05",
            "Name": "Home",
            "NodeId": 36241,
            "ProposalSubmitted": 0,
            "RelayMessageCount": 10331962,
            "RelaysPerHour": 36734,
            "RelaysPerHour10": 42660,
            "RelaysPerHour60": 38364,
            "SyncState": "PERSIST_FINISHED",
            "Uptime": 1012541,
            "Version": "v2.1.6"
          }
        ]
      }
    },
    "Prices": {
      "Method": "getfullstack",
      "Code": 0,
      "Error": false,
      "Value": {
        "usd": 0.39364
      }
    },
    "Wallets": {
      "Method": "getfullstack",
      "Code": 0,
      "Error": false,
      "Value": {
        "Wallets": [
          {
            "Balance": 343.49090321,
            "Id": 124,
            "NknWallet": "NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE"
          }
        ]
      }
    }
  }
}
```


</details>

## Method `savemysettings`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to change your account settings.

- The key "Hash" (string) must be set to valid Hash.
- If you want to change your wallets, add the "Wallets" key in your request. The key must contain **array of strings**.
- If you want to delete all your wallets, pass an empty array in the **Wallets** key.

#### cURL example #1 (application/json):  

```bash
curl -X POST  \
-d '{
"Method":"savemysettings",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39",
	"Wallets": ["NKNZKKF9u1MUQWnK272YoFiMTn5tjZh7uRQE"]
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api
```

**The method is not implemented in application/x-www-form-urlencoded.**

	
#### Server returns:

```json
{
  "Method": "savemysettings",
  "Code": 0,
  "Error": false,
  "Value": {
    "Data": "All settings saved"
  }
}
```

</details>

## Method `getnodedetails`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get node details information by node id.  
The method returns online data.

- The key "Hash" (string) must be set to valid Hash.
- The key "NodeId" (integer) must be set to node id.


#### cURL example #1 (application/json):  	
	
```bash
curl -X POST  \
-d '{
"Method":"getnodedetails",
"Value": {
	"Hash": "3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39",
	"NodeId": 36241
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api
```


#### cURL example #2 (application/x-www-form-urlencoded):

```bash
curl -X POST  \
-d 'Method=getnodedetails' \
-d 'Hash=3397f7beaec0c6921d6b1888e2f66d5559e81e4c8ccad3b149ab04dd3a8baf39' \
-d 'NodeId=36241' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```
	
#### The server returns:

```json
{
	"Method":"getnodedetails",
	"Code":0,
	"Error":false,
	"Value":{
		"DebugInfo":{
				"GetneighborTime":"680.055301ms",
				"GetnodestateTime":"144.97108ms",
				"HandlingTime":"825.278162ms"
		},
		"NodeStats":{
				"MinPing":60,
				"AvgPing":221,
				"MaxPing":3155,
				"NeighborCount":237,
				"NeighborPersist":62,
				"RelaysPerHour":36731,
				"NodeState": {
						"id":"1",
						"jsonrpc":"2.0",
						"result":{
							"addr":"tcp://1.1.1.1:30001",
							"currTimeStamp":1631803058,
							"height":3098247,
							"id":"%nkn_node_id%",
							"jsonRpcPort":30003,
							"proposalSubmitted":0,
							"protocolVersion":40,
							"publicKey":"%nkn_node_pubkey%",
							"relayMessageCount":10360133,
							"syncState":"PERSIST_FINISHED",
							"tlsJsonRpcDomain":"1-1-1-1.ipv4.nknlabs.io",
							"tlsJsonRpcPort":30005,
							"tlsWebsocketDomain":"1-1-1-1.ipv4.nknlabs.io",
							"tlsWebsocketPort":30004,
							"uptime":1015395,
							"version":"v2.1.6",
							"websocketPort":30002
						}
				}
		}
	}
}
```


</details>

## Method `getnodeipbypublickey`

<details>
  <summary>Show description and examples for the method</summary>

______

The method serves to get the node IP address by passed PublicKey.  
The method works with the NKN Network only.

- The key "PublicKey" (string) must be set and contain a node public key.

#### cURL example #1 (application/json):  	
	
```bash
curl -X POST  \
-d '{
"Method":"getnodeipbypublickey",
"Value": {
	"PublicKey": "ab8ecc50adab32f9090ac9afa88b21889a32b1a01c729334100d56d777a2b60e"
	}
}' \
-H "Content-Type: application/json" \
https://nkn.ovh/api
```


#### cURL example #2 (application/x-www-form-urlencoded):

```bash
curl -X POST  \
-d 'Method=getnodeipbypublickey' \
-d 'PublicKey=ab8ecc50adab32f9090ac9afa88b21889a32b1a01c729334100d56d777a2b60e' \
-H "Content-Type: application/x-www-form-urlencoded" \
https://nkn.ovh/api
```

#### The server returns:

```json
{
	"Method":"getnodeipbypublickey",
	"Code":0,
	"Error":false,
	"Value":
		{
		"IpList": ["1.1.1.1"]
		}
}

```

- If the public key is not found, the server returns the code 3.
</details>
