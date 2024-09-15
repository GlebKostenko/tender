CREATE TABLE IF NOT EXISTS tenders (
    id SERIAL PRIMARY KEY,
    organization_id INT REFERENCES organization(id),
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50),
    description TEXT,
    status VARCHAR(20),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    tender_id INT
);

CREATE TABLE IF NOT EXISTS bids (
    id SERIAL PRIMARY KEY,
    tender_id INT REFERENCES tenders(tender_id),
    description TEXT,
    status VARCHAR(20),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    bid_id INT,
    username VARCHAR(50),
    name VARCHAR(100) NOT NULL
);