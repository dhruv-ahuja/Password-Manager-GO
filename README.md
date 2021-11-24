# Introduction
This is my first project in Go, a password manager application. A humble attempt at execution of an idea I've had for some time now. The application can save, view, edit and delete credentials. You can add multiple credentials for a single website. A simplistic textual interface is used to get things done.

# How Does It Work
When you first start the application you are prompted to make a master password. The app has four main functions:

1) Adding credentials to the database 
2) Displaying saved credentials for a particular website 
3) Modifying saved credentials 
4) Deleting saved credentials 

# Security 
The master password the user creates is hashed and saved to disk. The user is prompted to re-enter the master password each time they open the application and their entered text is compared to the saved hash of the password. The user is let through if the entry matches the hash. 

User credentials are kept in the table "info", with the password being secured by an encryption key. 

The encryption key is 32 bytes long generated using the GO library's "crypto/rand" package. It is then sealed using a 24 byte long nonce and master password hash. The sealed output is finally saved to disk, ready for future use.
The encryption key's purpose is to encrypt and subsequently decrypt passwords being saved to the database.  

# Setup
We will be using PostgreSQL for the database so please ensure that you have it installed before proceeding.

Create a .env file and set the environment variables using the ".env.example" file included in the project. If you decide that you do not want the databases' name to be "passwords", change it to the name of your choice opposite the "DB_NAME" variable. 

Once done with the .env file, go ahead and create a database for the application. The table will be created automatically alongside the master password and encryption key on application first run.

Finally, open your terminal and use the command in the project directory: ```go run main.go```

# TODO
1) Perhaps implement copying the password to the users' clipboard instead of just displaying it.
2) Improve the way authentication is handled.
