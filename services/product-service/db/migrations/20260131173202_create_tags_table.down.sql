DROP INDEX IF EXISTS idx_product_tags_tag_id;
DROP INDEX IF EXISTS idx_product_tags_product_id;
DROP INDEX IF EXISTS idx_tags_shop_id;
DROP TABLE IF EXISTS product_tags CASCADE;
DROP TABLE IF EXISTS tags CASCADE;