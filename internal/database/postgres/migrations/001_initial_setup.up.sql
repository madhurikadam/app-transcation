CREATE TABLE IF NOT EXISTS accounts (
    id uuid PRIMARY KEY,
    document_number TEXT NOT NULL,
    withdrawal_limit float8 not null default 3000,
    credit_limit float8 not null default 3000,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS operations_types (
    id int PRIMARY KEY,
    description TEXT NOT NULL
);

INSERT INTO operations_types (id,description) VALUES (1, 'Normal Purchase');
INSERT INTO operations_types (id,description) VALUES (2, 'Purchase with installments');
INSERT INTO operations_types (id,description) VALUES (3, 'Withdrawal');
INSERT INTO operations_types (id,description) VALUES (4, 'Credit Voucher');


CREATE TABLE IF NOT EXISTS transcations (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL,
    operation_type_id int NOT NULL,
    amount float8 NOT NULL,
    event_at timestamp NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (operation_type_id) REFERENCES operations_types(id)
);
