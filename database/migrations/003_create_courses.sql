CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    kode_mk VARCHAR(20) UNIQUE NOT NULL,
    nama_mk VARCHAR(255) NOT NULL,
    sks INTEGER NOT NULL,
    semester INTEGER NOT NULL,
    kuota INTEGER NOT NULL,

    CONSTRAINT chk_courses_sks
        CHECK (sks > 0),

    CONSTRAINT chk_courses_semester
        CHECK (semester > 0),

    CONSTRAINT chk_courses_kuota
        CHECK (kuota >= 0)
);