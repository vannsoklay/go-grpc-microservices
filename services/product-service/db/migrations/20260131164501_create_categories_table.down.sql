DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_categories_shop_id;
DROP TABLE IF EXISTS categories CASCADE;