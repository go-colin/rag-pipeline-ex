CREATE TABLE wallet (
    address VARCHAR(255) PRIMARY KEY,
    balance DECIMAL(18, 8) DEFAULT 0
);

CREATE TABLE transaction (
    txid VARCHAR(255) PRIMARY KEY,
    wallet_id VARCHAR(255) REFERENCES wallet(address),
    amount DECIMAL(18, 8) DEFAULT 0,
    token_id VARCHAR(255) REFERENCES token(address),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE token (
    address VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);