package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// RiceContract contract for managing CRUD for Rice Batches and Processing Orders
type RiceContract struct {
	contractapi.Contract
}

type RiceBatch struct {
	AssetType    string `json:"assetType"`
	BatchID      string `json:"batchID"`
	Variety      string `json:"variety"`
	HarvestDate  string `json:"harvestDate"`
	QuantityInKg int    `json:"quantityInKg"`
	ProducedBy   string `json:"producedBy"`
	Status       string `json:"status"`
}

type ProcessingOrder struct {
	AssetType    string `json:"assetType"`
	OrderID      string `json:"orderID"`
	Variety      string `json:"variety"`
	MillerName   string `json:"millerName"`
	QuantityInKg int    `json:"quantityInKg"`
}

type HistoryQueryResult struct {
	Record    *RiceBatch `json:"record"`
	TxId      string     `json:"txId"`
	Timestamp string     `json:"timestamp"`
	IsDelete  bool       `json:"isDelete"`
}

func getCollectionName() string {
	return "ProcessingOrderCollection"
}

// RiceBatchExists returns true when rice batch with given ID exists in world state
func (c *RiceContract) RiceBatchExists(ctx contractapi.TransactionContextInterface, batchID string) (bool, error) {
	data, err := ctx.GetStub().GetState(batchID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

// CreateRiceBatch creates a new rice batch (only by farmer)
func (c *RiceContract) CreateRiceBatch(ctx contractapi.TransactionContextInterface, batchID string, variety string, harvestDate string, quantityInKg int, farmerName string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "Org1MSP" {
		exists, err := c.RiceBatchExists(ctx, batchID)
		if err != nil {
			return "", err
		} else if exists {
			return "", fmt.Errorf("the batch %s already exists", batchID)
		}

		rice := RiceBatch{
			AssetType:    "riceBatch",
			BatchID:      batchID,
			Variety:      variety,
			HarvestDate:  harvestDate,
			QuantityInKg: quantityInKg,
			ProducedBy:   farmerName,
			Status:       "Harvested",
		}
		bytes, _ := json.Marshal(rice)

		return fmt.Sprintf("successfully added rice batch %v", batchID), ctx.GetStub().PutState(batchID, bytes)
	} else {
		return "", fmt.Errorf("user under following MSPID: %v can't perform this action", clientOrgID)
	}
}

// ReadRiceBatch retrieves an instance of RiceBatch
func (c *RiceContract) ReadRiceBatch(ctx contractapi.TransactionContextInterface, batchID string) (*RiceBatch, error) {
	bytes, err := ctx.GetStub().GetState(batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("the rice batch %s does not exist", batchID)
	}

	var batch RiceBatch
	err = json.Unmarshal(bytes, &batch)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type RiceBatch")
	}
	return &batch, nil
}

// DeleteRiceBatch removes batch from world state
func (c *RiceContract) DeleteRiceBatch(ctx contractapi.TransactionContextInterface, batchID string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}
	if clientOrgID == "Org1MSP" {
		exists, err := c.RiceBatchExists(ctx, batchID)
		if err != nil {
			return "", fmt.Errorf("Could not read from world state. %s", err)
		} else if !exists {
			return "", fmt.Errorf("The asset %s does not exist", batchID)
		}
		err = ctx.GetStub().DelState(batchID)
		return fmt.Sprintf("Rice batch %v deleted", batchID), err
	} else {
		return "", fmt.Errorf("User under following MSP:%v cannot perform deletion", clientOrgID)
	}
}

// GetAllRiceBatches retrieves all rice batches
func (c *RiceContract) GetAllRiceBatches(ctx contractapi.TransactionContextInterface) ([]*RiceBatch, error) {
	queryString := `{"selector":{"assetType":"riceBatch"}, "sort":[{ "batchID": "desc"}]}`
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return riceBatchIterator(resultsIterator)
}

// Iterator function for rice batch
func riceBatchIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*RiceBatch, error) {
	var batches []*RiceBatch
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var batch RiceBatch
		err = json.Unmarshal(queryResult.Value, &batch)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}
	return batches, nil
}

// GetRiceBatchHistory returns the history
func (c *RiceContract) GetRiceBatchHistory(ctx contractapi.TransactionContextInterface, batchID string) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(batchID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var batch RiceBatch
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &batch)
			if err != nil {
				return nil, err
			}
		} else {
			batch = RiceBatch{BatchID: batchID}
		}
		timestamp := response.Timestamp.AsTime()
		formattedTime := timestamp.Format(time.RFC1123)
		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &batch,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}
	return records, nil
}

// CreateProcessingOrder creates new private order
func (c *RiceContract) CreateProcessingOrder(ctx contractapi.TransactionContextInterface, orderID string) (string, error) {
	clientOrgID, _ := ctx.GetClientIdentity().GetMSPID()
	if clientOrgID != "Org2MSP" {
		return "", fmt.Errorf("Only Org2MSP (Miller) can create order")
	}

	exists, _ := c.ProcessingOrderExists(ctx, orderID)
	if exists {
		return "", fmt.Errorf("Order ID already exists")
	}

	transientData, _ := ctx.GetStub().GetTransient()
	if len(transientData) == 0 {
		return "", fmt.Errorf("Provide transient fields: variety, millerName, quantityInKg")
	}

	order := &ProcessingOrder{
		AssetType:    "processingOrder",
		OrderID:      orderID,
		Variety:      string(transientData["variety"]),
		MillerName:   string(transientData["millerName"]),
		QuantityInKg: parseInt(string(transientData["quantityInKg"])),
	}

	bytes, _ := json.Marshal(order)
	return fmt.Sprintf("Processing order %v created", orderID), ctx.GetStub().PutPrivateData(getCollectionName(), orderID, bytes)
}

// ProcessingOrderExists checks in private data
func (c *RiceContract) ProcessingOrderExists(ctx contractapi.TransactionContextInterface, orderID string) (bool, error) {
	data, err := ctx.GetStub().GetPrivateDataHash(getCollectionName(), orderID)
	return data != nil, err
}

// ReadProcessingOrder from private collection
func (c *RiceContract) ReadProcessingOrder(ctx contractapi.TransactionContextInterface, orderID string) (*ProcessingOrder, error) {
	bytes, err := ctx.GetStub().GetPrivateData(getCollectionName(), orderID)
	if err != nil || bytes == nil {
		return nil, fmt.Errorf("Order does not exist or cannot be read")
	}
	var order ProcessingOrder
	err = json.Unmarshal(bytes, &order)
	return &order, err
}

// MatchProcessingOrder with rice batch
func (c *RiceContract) MatchProcessingOrder(ctx contractapi.TransactionContextInterface, batchID string, orderID string) (string, error) {
	batch, err := c.ReadRiceBatch(ctx, batchID)
	if err != nil {
		return "", err
	}

	order, err := c.ReadProcessingOrder(ctx, orderID)
	if err != nil {
		return "", err
	}

	if batch.Variety == order.Variety && batch.QuantityInKg >= order.QuantityInKg {
		batch.Status = fmt.Sprintf("Assigned to Miller %v", order.MillerName)
		bytes, _ := json.Marshal(batch)

		ctx.GetStub().DelPrivateData(getCollectionName(), orderID)
		err = ctx.GetStub().PutState(batchID, bytes)
		return fmt.Sprintf("Order %v fulfilled by batch %v", orderID, batchID), err
	} else {
		return "", fmt.Errorf("Variety or quantity mismatch for order pairing")
	}
}

// Final registration: Retailer dispatch batch
func (c *RiceContract) DispatchToRetailer(ctx contractapi.TransactionContextInterface, batchID string, retailer string) (string, error) {
	clientOrgID, _ := ctx.GetClientIdentity().GetMSPID()
	if clientOrgID != "Org3MSP" {
		return "", fmt.Errorf("Only Org3MSP (Retailer) can dispatch batch")
	}

	batch, err := c.ReadRiceBatch(ctx, batchID)
	if err != nil {
		return "", err
	}
	batch.Status = fmt.Sprintf("Dispatched to %v", retailer)
	batch.ProducedBy = retailer

	bytes, _ := json.Marshal(batch)
	err = ctx.GetStub().PutState(batchID, bytes)
	return fmt.Sprintf("Batch %v dispatched to %v", batchID, retailer), err
}

// GetRiceBatchByRange retrieves rice batches within a key range
func (c *RiceContract) GetRiceBatchByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*RiceBatch, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	return riceBatchIterator(resultsIterator)
}

// processingOrderIterator is used to iterate over private data query results
func processingOrderIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*ProcessingOrder, error) {
	var orders []*ProcessingOrder

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var order ProcessingOrder
		err = json.Unmarshal(queryResult.Value, &order)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

// GetMatchingOrders finds orders that match the rice batch variety and quantity
func (c *RiceContract) GetMatchingOrders(ctx contractapi.TransactionContextInterface, batchID string) ([]*ProcessingOrder, error) {
	exists, err := c.RiceBatchExists(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("Could not check batch existence: %s", err)
	} else if !exists {
		return nil, fmt.Errorf("Batch %s does not exist", batchID)
	}

	batch, err := c.ReadRiceBatch(ctx, batchID)
	if err != nil {
		return nil, fmt.Errorf("Error reading batch: %v", err)
	}

	queryString := fmt.Sprintf(`{
		"selector": {
			"assetType": "processingOrder",
			"variety": "%s"
		}
	}`, batch.Variety)

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(getCollectionName(), queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return processingOrderIterator(resultsIterator)
}

// GetAllProcessingOrders returns all processing orders from private data collection
func (c *RiceContract) GetAllProcessingOrders(ctx contractapi.TransactionContextInterface) ([]*ProcessingOrder, error) {
	queryString := `{"selector":{"assetType":"processingOrder"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(getCollectionName(), queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return processingOrderIterator(resultsIterator)
}

// GetProcessingOrdersByRange returns all processing orders by key range
func (c *RiceContract) GetProcessingOrdersByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*ProcessingOrder, error) {
	collectionName := getCollectionName()
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return processingOrderIterator(resultsIterator)
}


// Helper: Parse integer
func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}