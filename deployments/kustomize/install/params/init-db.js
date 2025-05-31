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

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
      print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })

//insert sample data
let result = db[collection].insertMany([
    {
        id: "vp-001",
        name: "Anna Kováčová",
        recordId: "mongo-record-123",
        difficulty: 2,
        symptoms: ["fever", "cough", "fatigue"],
        anamnesis: "Patient reports feeling unwell for the past 3 days"
    },
    {
        id: "vp-002",
        name: "Filip Mocháč",
        recordId: "mongo-record-456",
        difficulty: 3,
        symptoms: ["headache", "nausea", "dizziness"],
        anamnesis: "Patient reports severe headache and dizziness since morning"
    },
    {
        id: "vp-003",
        name: "Peter Novák",
        recordId: "mongo-record-789",
        difficulty: 4,
        symptoms: ["chest pain", "shortness of breath", "sweating"],
        anamnesis: "Patient reports severe chest pain and difficulty breathing"
    },
    {
        id: "vp-004",
        name: "Mária Horváthová",
        recordId: "mongo-record-101",
        difficulty: 1,
        symptoms: ["sore throat", "runny nose", "mild fever"],
        anamnesis: "Patient reports cold-like symptoms for 2 days"
    },
    {
        id: "vp-005",
        name: "Ján Tóth",
        recordId: "mongo-record-202",
        difficulty: 5,
        symptoms: ["severe abdominal pain", "vomiting", "fever", "dehydration"],
        anamnesis: "Patient reports acute abdominal pain and persistent vomiting for 12 hours"
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);