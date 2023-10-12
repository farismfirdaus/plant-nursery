-- Revert plant-nursery:2_plant from pg

BEGIN;

DROP TABLE IF EXISTS plants;

COMMIT;
