-- Menghapus kolom deleted_at dari tabel users untuk membatalkan migrasi.
ALTER TABLE users
    DROP COLUMN deleted_at;
