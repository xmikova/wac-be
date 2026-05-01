const mongoHost = process.env.AMBULANCE_API_MONGODB_HOST
const mongoPort = process.env.AMBULANCE_API_MONGODB_PORT

const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD

const database = process.env.AMBULANCE_API_MONGODB_DATABASE
const collection = process.env.AMBULANCE_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

const db = connection.getDB(database)

// initialize ambulance collection if not exists
const databases = connection.getDBNames()
if (!databases.includes(database) || !db.getCollectionNames().includes(collection)) {
    db.createCollection(collection)
    db[collection].createIndex({ "id": 1 })
    db[collection].insertMany([
        {
            "id": "bobulova",
            "name": "Dr.Bobulová",
            "roomNumber": "123",
            "predefinedConditions": [
                { "value": "Nádcha", "code": "rhinitis" },
                { "value": "Kontrola", "code": "checkup" }
            ]
        }
    ]);
    print(`Initialized collection '${collection}' in database '${database}'`)
} else {
    print(`Collection '${collection}' already exists in database '${database}'`)
}

// initialize pharmacy collection if not exists
const pharmacyCollection = 'pharmacy'
if (!db.getCollectionNames().includes(pharmacyCollection)) {
    db.createCollection(pharmacyCollection)
    db[pharmacyCollection].createIndex({ "id": 1 })
    db[pharmacyCollection].insertOne({
        "id": "pmdl-pharmacy",
        "medicines": [
            {
                "id": "med-001",
                "name": "Paracetamol 500mg",
                "activeSubstance": "Paracetamolum",
                "dosage": "500mg",
                "batchNumber": "BT2024001",
                "expiryDate": "2026-12-31",
                "minStock": 100,
                "currentStock": 250
            },
            {
                "id": "med-002",
                "name": "Ibuprofen 400mg",
                "activeSubstance": "Ibuprofenum",
                "dosage": "400mg",
                "batchNumber": "BT2024002",
                "expiryDate": "2025-06-30",
                "minStock": 50,
                "currentStock": 30
            }
        ]
    });
    print(`Initialized collection '${pharmacyCollection}' in database '${database}'`)
} else {
    print(`Collection '${pharmacyCollection}' already exists in database '${database}'`)
}

process.exit(0);
