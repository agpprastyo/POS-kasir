CREATE TABLE IF NOT EXISTS settings (
    key VARCHAR(50) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed default values
INSERT INTO settings (key, value, description) VALUES
('app_name', 'POS Kasir', 'Application Name displayed in sidebar and titles'),
('app_logo', '', 'URL or path to the application logo image'),
('footer_text', '© 2024 POS Kasir. All rights reserved.', 'Text displayed in the footer'),
('theme_color', '#000000', 'Primary theme color (hex code)');
