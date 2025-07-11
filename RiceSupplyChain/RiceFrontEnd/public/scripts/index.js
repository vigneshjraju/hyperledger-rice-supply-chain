// Add Rice Batch
const addData = async (event) => {
    event.preventDefault();

    const batchID = document.getElementById("batchID").value;
    const variety = document.getElementById("variety").value;
    const harvestDate = document.getElementById("harvestDate").value;
    const quantity = document.getElementById("quantity").value;
    const farmerName = document.getElementById("farmerName").value;

    const riceData = {
        batchID,
        variety,
        harvestDate,
        quantity,
        farmerName
    };

    if (!batchID || !variety || !harvestDate || !quantity || !farmerName) {
        alert("Please enter all fields properly.");
    } else {
        try {
            const response = await fetch("/api/rice", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(riceData),
            });

            const result = await response.json();
            console.log("RESULT: ", result);
            alert("Rice Batch Created Successfully");
        } catch (err) {
            alert("Error while creating RiceBatch");
            console.log(err);
        }
    }
};

// Read Rice Batch by ID
const readData = async (event) => {
    event.preventDefault();

    const batchID = document.getElementById("batchIDInput").value;

    if (!batchID) {
        alert("Enter a valid rice batch ID");
    } else {
        try {
            const response = await fetch(`/api/rice/${batchID}`);
            const result = await response.json();
            alert(JSON.stringify(result));
        } catch (err) {
            alert("Error while reading Rice Batch");
            console.log(err);
        }
    }
};

// ========== ORG1: Farmer - Create Rice Batch ==========
const createRiceBatch = async (event) => {
  event.preventDefault();

  const batchID = document.getElementById("batchID").value;
  const variety = document.getElementById("variety").value;
  const harvestDate = document.getElementById("harvestDate").value;
  const quantity = document.getElementById("quantity").value;
  const farmerName = document.getElementById("farmerName").value;

  if (!batchID || !variety || !harvestDate || !quantity || !farmerName) {
    alert("All fields are required.");
    return;
  }

  const payload = {
    batchID,
    variety,
    harvestDate,
    quantity,
    farmerName
  };

  const res = await fetch("/api/rice", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  const result = await res.json();
  alert(result.message || "Rice Batch Created");
};

// ========== ORG1: Query Batches ==========
const getAllRiceBatches = async () => {
  const res = await fetch("/api/rice/all");
  const result = await res.json();
  alert("Rice Batches:\n" + result.data);
};

const getRiceByRange = async () => {
  const start = document.getElementById("rangeStart").value;
  const end = document.getElementById("rangeEnd").value;

  const res = await fetch(`/api/rice/range?start=${start}&end=${end}`);
  const result = await res.json();
  alert("Rice Batches by Range:\n" + result.data);
};

const getRiceHistory = async () => {
  const batchID = document.getElementById("historyBatchID").value;
  const res = await fetch(`/api/rice/history/${batchID}`);
  const result = await res.json();
  alert("Batch History:\n" + result.data);
};

// ========== ORG2: Create Processing Order ==========
const createOrder = async () => {
  const orderID = document.getElementById("orderID").value;
  const variety = document.getElementById("orderVariety").value;
  const quantity = document.getElementById("orderQuantity").value;
  const miller = document.getElementById("millerName").value;

  if (!orderID || !variety || !quantity || !miller) {
    alert("All fields are required.");
    return;
  }

  const payload = {
    orderID,
    variety,
    quantityInKg: quantity,
    millerName: miller
  };

  const res = await fetch("/api/orders", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  const result = await res.json();
  alert("Create Order: " + result.message);
};

// ========== ORG2: Match Order with Rice ==========
const matchOrder = async () => {
  const batchID = document.getElementById("matchBatchID").value;
  const orderID = document.getElementById("matchOrderID").value;

  if (!batchID || !orderID) {
    alert("All fields are required!");
    return;
  }

  const payload = {
    batchID,
    orderID
  };

  const res = await fetch("/api/orders/match", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  const result = await res.json();
  alert("Match Result: " + result.result);
};

// ========== ORG3: Dispatch Rice to Retailer ==========
const dispatchToRetailer = async () => {
  const batchID = document.getElementById("dispatchBatchID").value;
  const retailer = document.getElementById("retailerName").value;

  const payload = {
    batchID,
    retailerName: retailer
  };

  const res = await fetch("/api/rice/dispatch", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  const result = await res.json();
  alert(result.message);
};