CREATE TABLE enrollments (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    tahun_akademik VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_enrollments_students
        FOREIGN KEY (student_id)
        REFERENCES students(id),

    CONSTRAINT fk_enrollments_courses
        FOREIGN KEY (course_id)
        REFERENCES courses(id),

    CONSTRAINT uq_student_course_year
        UNIQUE (student_id, course_id, tahun_akademik)
);