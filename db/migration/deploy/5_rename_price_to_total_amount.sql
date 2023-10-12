-- Deploy plant-nursery:5_rename_price_to_total_amount to pg

BEGIN;

ALTER TABLE carts RENAME COLUMN price TO total_amount;

COMMIT;
