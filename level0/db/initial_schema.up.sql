SET session_replication_role = 'replica';

CREATE TABLE orders (
                        order_uid VARCHAR(255) PRIMARY KEY,
                        track_number VARCHAR(255) NOT NULL UNIQUE,
                        entry VARCHAR(255),
                        locale VARCHAR(10),
                        customer_id VARCHAR(255),
                        delivery_service VARCHAR(255),
                        shardkey VARCHAR(10),
                        sm_id INT,
                        date_created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        oof_shard VARCHAR(10)
);


CREATE TABLE delivery (
                          order_uid VARCHAR(255) PRIMARY KEY,
                          name VARCHAR(255),
                          phone VARCHAR(255),
                          zip VARCHAR(255),
                          city VARCHAR(255),
                          address VARCHAR(255),
                          region VARCHAR(255),
                          email VARCHAR(255),

                          CONSTRAINT fk_delivery_orders FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE payment (

                         transaction_uid VARCHAR(255) PRIMARY KEY,
                         order_uid VARCHAR(255) NOT NULL UNIQUE,
                         request_id VARCHAR(255),
                         currency VARCHAR(10),
                         provider VARCHAR(255),
                         amount NUMERIC,
                         payment_dt BIGINT,
                         bank VARCHAR(255),
                         delivery_cost NUMERIC,
                         goods_total INT,
                         custom_fee NUMERIC,
                         CONSTRAINT fk_payment_orders FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE items (
                       chrt_id BIGINT PRIMARY KEY,
                       order_uid VARCHAR(255) NOT NULL,
                       track_number VARCHAR(255),
                       price NUMERIC,
                       rid VARCHAR(255),
                       name VARCHAR(255),
                       sale INT,
                       size VARCHAR(255),
                       total_price NUMERIC,
                       nm_id BIGINT,
                       brand VARCHAR(255),
                       status INT,
                       CONSTRAINT fk_items_orders FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

SET session_replication_role = 'origin';