-- Revert plant-nursery:4_order from pg

BEGIN;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS order_items;

COMMIT;
