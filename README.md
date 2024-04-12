# Days 7-8: Basic authentication

Set up custom user authentication. Where a user in the bookstore, enters his name and password before he can get access to any book.

Hint: You will need to create a User Model. Implementing password hashing attracts a higher score

## Setup
- Navigate to the root of this repo.
- Run the command ```go run ./main.go``` to start the server.
- Visit the following url endpoints:
    |METHOD|DESCRIPTION|ENDPOINT|SAMPLE BODY|
    |------|-----------|--------|----|
    |POST  |Register user|http://127.0.0.1:3000/register|{"username":"username","password":"password"}|
    |POST  |Login user   |http://127.0.0.1:3000/login|{"username":"username","password":"password"}|
    |GET   |Get all books added by logged-in user|http://127.0.0.1:3000/books|-|
    |GET   |Get a book added by logged-in user|http://127.0.0.1:3000/books/{id}|-|
    |POST  |Create a book for logged-in user|http://127.0.0.1:3000/books|{"author":"book_author","title":"book_title"}|
