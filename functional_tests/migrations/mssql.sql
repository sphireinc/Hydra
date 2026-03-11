-- Create the Person table
CREATE TABLE Person (
                        id INT IDENTITY(1,1) PRIMARY KEY,
                        first_name NVARCHAR(50) NOT NULL,
                        last_name NVARCHAR(50) NOT NULL,
                        sex CHAR(1) NOT NULL CHECK (sex IN ('M', 'F', 'O')),
                        date_of_birth DATE NOT NULL,
                        created DATETIME DEFAULT GETDATE(),
                        updated DATETIME DEFAULT GETDATE(),
                        deleted DATETIME NULL
);

-- Create the Addresses table
CREATE TABLE Addresses (
                           id INT IDENTITY(1,1) PRIMARY KEY,
                           user_id INT NOT NULL,
                           address_1 NVARCHAR(100) NOT NULL,
                           address_2 NVARCHAR(100),
                           city NVARCHAR(50) NOT NULL,
                           state NVARCHAR(50),
                           province NVARCHAR(50),
                           postal_code NVARCHAR(20),
                           country NVARCHAR(50) NOT NULL,
                           created DATETIME DEFAULT GETDATE(),
                           updated DATETIME DEFAULT GETDATE(),
                           deleted DATETIME NULL,
                           CONSTRAINT FK_UserAddresses FOREIGN KEY (user_id) REFERENCES Person(id)
);

-- Insert sample data into Person
INSERT INTO Person (first_name, last_name, sex, date_of_birth)
VALUES
    ('John', 'Doe', 'M', '1985-06-12'),
    ('Jane', 'Smith', 'F', '1990-11-24'),
    ('Alice', 'Johnson', 'F', '1978-09-05'),
    ('Bob', 'Brown', 'M', '1988-02-14'),
    ('Charlie', 'Williams', 'M', '1995-04-18'),
    ('Emily', 'Davis', 'F', '2000-12-12'),
    ('George', 'White', 'M', '1973-08-08'),
    ('Hannah', 'Moore', 'F', '1983-10-20'),
    ('Isaac', 'Taylor', 'M', '1999-01-01'),
    ('Zara', 'Clark', 'F', '2001-07-15');

-- Insert sample data into Addresses
INSERT INTO Addresses (user_id, address_1, address_2, city, state, postal_code, country)
VALUES
    (1, '123 Main St', 'Apt 4', 'New York', 'NY', '10001', 'USA'),
    (2, '456 Elm St', NULL, 'Los Angeles', 'CA', '90001', 'USA'),
    (3, '789 Oak St', 'Unit 10', 'Chicago', 'IL', '60601', 'USA'),
    (4, '321 Pine St', 'Suite 500', 'San Francisco', 'CA', '94101', 'USA'),
    (5, '654 Cedar St', NULL, 'Houston', 'TX', '77001', 'USA'),
    (6, '987 Maple St', 'Apt 12B', 'Miami', 'FL', '33101', 'USA'),
    (7, '555 Spruce St', NULL, 'Seattle', 'WA', '98101', 'USA'),
    (8, '444 Birch St', NULL, 'Denver', 'CO', '80201', 'USA'),
    (9, '333 Walnut St', NULL, 'Boston', 'MA', '02101', 'USA'),
    (10, '222 Ash St', 'Unit 7', 'Atlanta', 'GA', '30301', 'USA');