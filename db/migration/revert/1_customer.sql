-- Revert plant-nursery:1_customer from pg

BEGIN;

DROP TABLE IF EXISTS customers;

COMMIT;
