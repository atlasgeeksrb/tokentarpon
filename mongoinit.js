db.createCollection('mydb');
db.createUser(
    {
        user: 'myuser',
        pwd: 'mypassword',
        roles: [
            {
                role: 'readWrite',
                db: 'mydb'
            }
        ]
    }
);