CREATE TABLE "users_balance"
(
    user_id		    INTEGER PRIMARY KEY,
    balance         DECIMAL(19,4) NOT NULL DEFAULT 0.0 CONSTRAINT positive_balance CHECK (balance >= 0.0)
);

CREATE TABLE "transfers"
(
    sender_id		INTEGER,
    recipient_id    INTEGER,
    amount          DECIMAL(19,4) NOT NULL,
    date            TIMESTAMP with time zone,
    message         TEXT,
    CHECK (sender_id != recipient_id)
);
CREATE INDEX ind_recipient ON transfers (recipient_id);
CREATE INDEX ind_sender ON transfers (sender_id);