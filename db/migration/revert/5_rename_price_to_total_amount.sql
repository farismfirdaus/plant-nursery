-- Revert plant-nursery:5_rename_price_to_total_amount from pg

BEGIN;

ALTER TABLE carts RENAME COLUMN total_amount TO price;

COMMIT;
