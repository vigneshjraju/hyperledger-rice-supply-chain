
# 📝 Rice Supply Chain using Hyperledger Fabric

## 📌 Project Overview
This project implements a blockchain-based decentralized application (DApp) for simulating a **Rice Supply Chain** using Hyperledger Fabric.

### 👥 Participants
- 👨‍🌾 **Farmer (Org1)** – creates and registers harvested paddy batches
- 🏭 **Miller (Org2)** – places private processing orders and matches with paddy batches
- 🛒 **Retailer (Org3)** – dispatches and distributes final rice batches to the market

---

## 🏗️ Network Setup Instructions

### ✅ Start the Fabric Network (with CA & CouchDB)
```bash
./network.sh up createChannel -c mychannel -ca -s couchdb
```

### ➕ Add Org3 (Retailer) to Network
```bash
cd addOrg3
./addOrg3.sh up -c mychannel -ca -s couchdb
cd ..
```

### 🚀 Deploy Chaincode
```bash
./network.sh deployCC \
  -ccn rice \
  -ccp ../../RiceSupplyChain/Chaincode \
  -ccl go \
  -c mychannel \
  -ccv 1.0 \
  -ccs 1 \
  -cccg ../../RiceSupplyChain/Chaincode/collection_config.json
```

### ⛔ Shut Down the Network
```bash
./network.sh down
```

---

## 🌐 Environment Variables (per organization)

### 📍 General (for everyone)
```bash
export FABRIC_CFG_PATH=$PWD/../config/

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_TLS_ENABLED=true
```

### 👨‍🌾 Org1 - Farmer
```bash
export CORE_PEER_LOCALMSPID=Org1MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

export CORE_PEER_ADDRESS=localhost:7051
```

### 🏭 Org2 - Miller
```bash
export CORE_PEER_LOCALMSPID=Org2MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

export CORE_PEER_ADDRESS=localhost:9051
```

### 🛒 Org3 - Retailer
```bash
export CORE_PEER_LOCALMSPID=Org3MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp

export CORE_PEER_ADDRESS=localhost:11051
```

---

## ⚙️ Chaincode Functions & CLI Commands

### 🧱 Farmer: Create Paddy Batch (Org1)
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n rice \
--peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
--peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
-c '{"function":"CreateRiceBatch","Args":["PADDY001","Sona Masuri","2025-07-07","1000","FarmerA"]}'
```

### 📦 Query All Paddy Batches (Any Org)
```bash
peer chaincode query -C mychannel -n rice -c '{"Args":["GetAllRiceBatches"]}'
```

### 🔎 Range Query Paddy Batches
```bash
peer chaincode query -C mychannel -n rice -c '{"Args":["GetRiceBatchByRange", "PADDY001", "PADDY003"]}'
```

### 📜 Query Paddy Batch History
```bash
peer chaincode query -C mychannel -n rice -c '{"Args":["GetRiceBatchHistory", "PADDY001"]}'
```

---

### 🧾 Miller: Create Processing Order (Org2)
```bash
export VARIETY=$(echo -n "Sona Masuri" | base64 | tr -d '\n')
export QUANTITY=$(echo -n "1000" | base64 | tr -d '\n')
export MILLER_NAME=$(echo -n "MillerA" | base64 | tr -d '\n')

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n rice \
--peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
--peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
--transient "{"variety":"$VARIETY","quantityInKg":"$QUANTITY","millerName":"$MILLER_NAME"}" \
-c '{"Args":["CreateProcessingOrder", "ORDER001"]}'
```

### 📖 Read All Processing Orders
```bash
peer chaincode query -C mychannel -n rice -c '{"Args":["GetAllProcessingOrders"]}'
```

### 🔍 Get Orders by Range
```bash
peer chaincode query -C mychannel -n rice -c '{"Args":["GetProcessingOrdersByRange", "ORDER001", "ORDER003"]}'
```

### 🔁 Get Matching Orders for Rice Batch
```bash
peer chaincode query -C mychannel -n rice -c '{"Args":["GetMatchingOrders", "PADDY001"]}'
```

### 🔗 Match Rice Batch with Order
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n rice \
--peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
--peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
-c '{"function":"MatchProcessingOrder","Args":["PADDY001","ORDER001"]}'
```

### 🚚 Dispatch Rice Batch to Retailer (Org3)
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C mychannel -n rice \
--peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
--peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
-c '{"function":"DispatchToRetailer","Args":["PADDY001","RetailerX"]}'

### 🚚 Run Frontend
```bash
go run main.go client.go connect.go profile.go



