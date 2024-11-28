-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS timedb;

-- Switch to the timedb database
USE timedb;

-- Create the time_log table if it doesn't exist
CREATE TABLE IF NOT EXISTS time_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    timestamp TIMESTAMP
);
