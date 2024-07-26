const express = require("express");
const bodyParser = require("body-parser");
const { v4: uuidv4 } = require("uuid");

const app = express();
const port = 8082;

app.use(bodyParser.json());
const receiptStore = new Map();

function getPoints(receipt) {
  let points = 0;
  let regex = /[a-z0-9]/i; //i is for case insensitive

  for (let i = 0; i < receipt.retailer.length; i++) {
    if (regex.test(receipt.retailer[i])) {
      //This will check if our regex character passes the regex test
      points++;
    }
  }

  const total = parseFloat(receipt.total);
  if (Number.isInteger(total)) {
    points += 50;
  }
  if (total % 0.25 === 0) {
    points += 25;
  }

  points += Math.floor(receipt.items.length / 2) * 5;

  for (let i = 0; i < receipt.items.length; i++) {
    if (receipt.items[i].shortDescription.trim().length % 3 == 0) {
      points += Math.ceil(parseFloat(receipt.items[i].price) * 0.2);
    }
  }

  const purchaseDate = new Date(receipt.purchaseDate);
  if (purchaseDate.getDay() % 2 !== 0) {
    points += 6;
  }

  let [hour, minute] = receipt.purchaseTime.split(":").map(Number);
  if (hour >= 14 && hour < 16) points += 10;

  return points;
}

app.post("/receipts/process", (req, res) => {
  const { retailer, purchaseDate, purchaseTime, items, total } = req.body;

  if (!retailer || !purchaseDate || !purchaseTime || !items || !total) {
    return res.status(400).json({ error: "The receipt data is invalid!" });
  }

  const receipt = {
    retailer,
    purchaseDate,
    purchaseTime,
    items,
    total,
  };

  const uuid = uuidv4();
  id =
    uuid.substring(0, 8) +
    "-" +
    uuid.substring(9, 13) +
    "-" +
    uuid.substring(14, 18) +
    "-" +
    uuid.substring(19, 31);
  receiptStore.set(String(id), receipt);
  return res.status(201).json({ message: "Recipe created!", id });
});

app.get("/receipts/:id/points", (req, res) => {
  const id = req.params.id;
  console.log(receiptStore);
  const receipt = receiptStore.get(id);

  if (!receipt) {
    return res.status(404).json({ error: "Receipt not found" });
  }

  const points = getPoints(receipt);
  res.json({ points });
});

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
