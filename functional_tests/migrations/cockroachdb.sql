DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS person;

CREATE TABLE person (
                        id INT PRIMARY KEY,
                        first_name STRING NOT NULL,
                        last_name STRING NOT NULL,
                        sex STRING NOT NULL
);

CREATE TABLE addresses (
                           id INT PRIMARY KEY,
                           user_id INT NOT NULL,
                           address_1 STRING NOT NULL,
                           address_2 STRING NOT NULL,
                           city STRING NOT NULL,
                           state STRING NOT NULL,
                           postal_code STRING NOT NULL,
                           country STRING NOT NULL,
                           CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES person (id)
);

INSERT INTO person (id, first_name, last_name, sex) VALUES
                                                        (1, 'John', 'Doe', 'M'),
                                                        (2, 'Jane', 'Doe', 'F'),
                                                        (3, 'Alice', 'Smith', 'F'),
                                                        (4, 'Bob', 'Smith', 'M'),
                                                        (5, 'Charlie', 'Brown', 'M'),
                                                        (6, 'Diana', 'Prince', 'F'),
                                                        (7, 'Eve', 'Adams', 'F'),
                                                        (8, 'Frank', 'Miller', 'M'),
                                                        (9, 'Grace', 'Lee', 'F'),
                                                        (10, 'Henry', 'Clark', 'M');

INSERT INTO addresses (id, user_id, address_1, address_2, city, state, postal_code, country) VALUES
                                                                                                 (1, 1, '123 Main St', 'Apt 4', 'New York', 'NY', '10001', 'USA'),
                                                                                                 (2, 2, '456 Oak Ave', 'Suite 12', 'Los Angeles', 'CA', '90001', 'USA'),
                                                                                                 (3, 3, '789 Pine Rd', '', 'Chicago', 'IL', '60601', 'USA'),
                                                                                                 (4, 4, '101 Maple Dr', 'Unit B', 'Houston', 'TX', '77001', 'USA'),
                                                                                                 (5, 5, '202 Cedar Ln', '', 'Phoenix', 'AZ', '85001', 'USA'),
                                                                                                 (6, 6, '303 Birch Blvd', 'Floor 3', 'Philadelphia', 'PA', '19019', 'USA'),
                                                                                                 (7, 7, '404 Walnut St', '', 'San Antonio', 'TX', '78201', 'USA'),
                                                                                                 (8, 8, '505 Cherry Ave', 'Apt 8', 'San Diego', 'CA', '92101', 'USA'),
                                                                                                 (9, 9, '606 Aspen Ct', '', 'Dallas', 'TX', '75201', 'USA'),
                                                                                                 (10, 10, '707 Spruce Way', 'Suite 20', 'San Jose', 'CA', '95101', 'USA');