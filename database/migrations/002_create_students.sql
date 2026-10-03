
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    nim VARCHAR(12) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    prodi VARCHAR(100) NOT NULL,
    angkatan INTEGER NOT NULL,
    ipk_terakhir NUMERIC(3,2) DEFAULT 0.00,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_students_users
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT chk_ipk
        CHECK (ipk_terakhir >= 0 AND ipk_terakhir <= 4.00)
);