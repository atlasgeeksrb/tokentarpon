db.createCollection('community');
db.createUser(
    {
        user: 'tarponuser',
        pwd: 'strongpassword',
        roles: [
            {
                role: 'read',
                db: 'admin'
            },
            {
                role: 'readWrite',
                db: 'community'
            }
        ]
    }
);