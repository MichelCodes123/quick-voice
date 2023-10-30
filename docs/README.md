# quick-voice
Sleek and powerful invoicing application that lets users create and store business invoices. 

## Demo

## Core Features
Quick voice supports the following :
- Quick and easy invoice creation tools
- Support for up to 5 different businesses/ sender profiles
- Invoice summary per preset

Soon to be Added: 
- Useful data analytics such as..
    - Average invoice total 
    - Average number of items invoiced
    - Monlthy and Yearly income trends

## Motivations
I had been looking to learn SQL and how to use a DBMS, so I thought a project would be the best place to start.

## Database Schema
My database is setup as followed: 
![Alt text](image.png)
## How to run the project.
1. Install Go
> https://go.dev/dl/

2. Install PostgresSQL and setup a local database
> https://www.postgresql.org/download/
> Run the sql file to setup the database schema

3. Clone this repository.

4. In the project directory, run `go run .` from the terminal
 

 ## Why isnt this app online?
The focus of this project was learning golang and how to utilize databases. Putting the project online would require extra work setting up how to handle multiple/concurrent users in the same database. This was out of scope of the project. For now, it serves as a useful tool that anybody can run for themselves to store their own invoices.

