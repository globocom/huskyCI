var db = connect("mongodb://localhost/huskyCIDB");

db.createUser(
    {
        user: "huskyCIUser",
        pwd: "huskyCIPassword",
        roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
    }
);
