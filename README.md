# user_management_system

This GO web application performs User Authentication and Management.

To run this application you should have already installed GO and setup Cassandra on docker.
# Setting Cassandra DB on Docker

#Below command pulls the latest Docker Cassandra Image
$ docker pull cassandra

#To check the Docker Image run the below command
$ docker images

#Below command creates a Cassandra container from the Cassandra Image
$ docker run --name cassandraDB -d cassandra:latest

#Below command will open interactive termial for Cassandra cql to create Keyspace, table and sample data.
$ docker exec -it cassandraDB bash 
cassandraDB@<continerID>:~$ cqlsh

#Create the below keyspace and emps tables on Cassandra db.
Create keyspace in cassandra DB

cqlsh> CREATE KEYSPACE userdb WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
Create table in the keyspace created above

cqlsh> CREATE TABLE userdb.user_role_details (
user_id text PRIMARY KEY,
first_name text,
last_name text,
email_id text,
password text,
role_name text,
manager_id text,
date_created date,
date_modified timestamp
);

CREATE TABLE userdb.messages (
    msg_id text PRIMARY KEY,
    date_created timestamp,
    msg_from text,
    msg_header text,
    msg_text text,
    msg_to text
)



# To create docker image for the web application run the below commands from project root directory
$ go get github.com/gocql/gocql

$ go get github.com/gorilla/mux

# Command to build GO apllication
$ go build
 
# Get the executable generated from the above command and run below command to create docer image
docker build -t executable .

# Please update IP address in the database connection section accordingly.
cluster := gocql.NewCluster("172.17.0.2")


