var db = connect("mongodb://localhost/huskyDB");

db.createUser(
    {
        user: "husky",
        pwd: "superENVPassword",
        roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
    }
);