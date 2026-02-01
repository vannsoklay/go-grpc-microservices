DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP INDEX IF EXISTS idx_products_is_featured;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_barcode;
DROP INDEX IF EXISTS idx_products_sku;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_shop_id;
DROP TABLE IF EXISTS products CASCADE;