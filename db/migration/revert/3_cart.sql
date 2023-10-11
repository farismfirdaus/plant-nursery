-- Revert plant-nursery:3_cart from pg

BEGIN;

DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS cart_items;

COMMIT;
