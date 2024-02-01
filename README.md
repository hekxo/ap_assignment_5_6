# BarberShop Web Application

This is a web application for a BarberShop implemented using Go, PostgreSQL, and HTML. It allows users to browse and filter barbers based on their categories, experience, and sort them by name or price. Pagination is also implemented to display a limited number of barbers per page.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Setup](#setup)
- [Usage](#usage)
- [Screenshots](#screenshots)

## Introduction

The BarberShop web application is designed to provide users with a convenient way to find suitable barbers based on their preferences such as category, experience, and pricing. It utilizes Go for backend logic, PostgreSQL for data storage, and HTML for front-end rendering.

## Features

- **Filtering**: Users can filter barbers based on their category and experience.
- **Sorting**: Barbers can be sorted by name or price.
- **Pagination**: Displaying a limited number of barbers per page for better user experience.

## Setup

To run the BarberShop web application locally, follow these steps:

1. Install Go and PostgreSQL if you haven't already.
2. Clone this repository to your local machine.
3. Set up a PostgreSQL database and adjust the connection string in the Go code accordingly.
4. Run the SQL script to create the necessary tables and insert sample data.
5. Install the required dependencies by running `go mod tidy`.
6. Run the Go server using `go run main.go`.
7. Access the application through your web browser at `http://localhost:8080`.

## Usage

Once the application is set up and running, users can visit the main page to browse barbers or navigate to `/barbers` to view all available barbers. The filtering, sorting, and pagination functionality is available on the `/filtered-barbers` page.

### Filtering

Users can filter barbers by selecting the desired category and experience level from the dropdown menus.

### Sorting

Barbers can be sorted by name or price by selecting the preferred sorting option from the dropdown menu.

### Pagination

Users can navigate through multiple pages of barbers by entering the desired page number in the input field.

## Screenshots
#### Create table "barbers" in PostgreSql
![Create table "barbers" in PostgreSql](postgresql.jpg)
#### Aplly filter by status
![Aplly filter by status](filter_status.jpg)
#### Aplly filter by experience
![Apply filter by experience](filter_exp.jpg)
#### Aplly sorting by name
![Apply sort by name](page1.jpg)
#### Aplly sorting by price
![Apply sort by price](sort_price.jpg)
#### Pagination
Page - 1
![Change the page](page1.jpg)
 Page - 2
![Change the page](page2.jpg)







