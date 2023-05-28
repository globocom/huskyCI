set -e

mongo <<EOF
use $MONGO_INITDB_DATABASE

db.createUser({
  user: '$MONGO_INITDB_ROOT_USERNAME',
  pwd: '$MONGO_INITDB_ROOT_PASSWORD',
  roles: [{
    role: 'root',
    db: 'admin'
  }]
})

db.createCollection('dockerAPIAddresses')

db.getCollection('dockerAPIAddresses').insertOne({
    currentHostIndex: 0,
    hostList: ["dockerapi"]
})

EOF
