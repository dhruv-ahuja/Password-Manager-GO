# Introduction
This is my first project in Go, a password manager application. A humble attempt at execution of an idea I've had for some time now. The application can save, view and edit credentials. You can add multiple credentials for a single website. A simplistic textual interface is used to get things done.

# How Does It Work
When you first start the application you are prompted to make a master password. The app has four main functions:

1) Adding credentials to the database 
2) Displaying saved credentials for a particular website 
3) Modifying saved credentials 
4) Deleting saved credentials 

# Security 
The master password the user creates is hashed and kept in a file in the same directory as the app. The user is prompted to enter the it when they open the application each time and their entered text is compared to the saved hash of the password. If the plaintext matches, the user is let through. 

The credentials are kept in a table. The password is encrypted with the AES encryption method using an encryption key of length of 32 characters, these encrypted bytes are then converted to base64 string and saved in the database. The AES encryption key is set by the user in the .env file. It is also used when decrypting the password each time to display to the user. 

# Setup
We will be using postgres for the database so please ensure that you have it installed before proceeding.

Go to database/db.go file and fill up the constants username and password according to your Postgres configuration.  

Go ahead and create a database "passwords" for the application. The table will be created automatically. If you decide to make a database by another name, you will have to set the "dbname" constant to the new name in the "database/db.go" file. 
The table will be created automatically once you make a master password. 

You also need to set an "ENC_KEY" variable inside a .env file that acts as the encryption key and will be used in encryption and decryption of passwords. A sample .env.example file is included. Please remember that the ENC_KEY must be 32 characters long.

Once done, just open your terminal and use the command in the project directory: ```go run main.go```

# TODO
1) Perhaps implement copying the password to the users' clipboard instead of just displaying it.
2) Improve the way authentication is handled.
