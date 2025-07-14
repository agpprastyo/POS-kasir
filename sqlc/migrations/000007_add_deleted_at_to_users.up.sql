-- Menambahkan kolom deleted_at ke tabel users untuk implementasi soft delete.
-- Kolom ini akan NULL secara default, yang berarti pengguna tersebut aktif.
ALTER TABLE users
    ADD COLUMN deleted_at TIMESTAMPTZ NULL;
