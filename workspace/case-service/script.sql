CREATE TABLE IF NOT EXISTS public.tb_case (
    case_id BIGSERIAL PRIMARY KEY,
    customer_number VARCHAR(100),
    doctor_name VARCHAR(100),
    email VARCHAR(100),
    previous_case VARCHAR(10),
    previous_case_number VARCHAR(100),
    patient_name VARCHAR(100),
    dob DATE,
    height_of_patient VARCHAR(100),
    gender VARCHAR(10),
    race VARCHAR(100),
    package_list VARCHAR(100),
    lateral_xray_date DATE,
    consult_date DATE,
    missing_teeth VARCHAR(100),
    adenoids_removed VARCHAR(10),
    comment TEXT,
    status_case VARCHAR(2),
    lateral_xray_image TEXT,
    frontal_xray_image TEXT,
    lower_arch_image TEXT,
    upper_arch_image TEXT,
    handwrist_xray_image TEXT,
    panoramic_xray_image TEXT,
    additional_record_1 TEXT,
    additional_record_2 TEXT,
    additional_record_3 TEXT,
    additional_record_4 TEXT,
    additional_record_5 TEXT,
    time_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE IF EXISTS public.tb_case
    OWNER TO postgres;


CREATE TABLE IF NOT EXISTS public.tb_user (
    user_id BIGSERIAL PRIMARY KEY,
    email VARCHAR(100),
    active VARCHAR(1),
    time_create TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE IF EXISTS public.tb_user
    OWNER TO postgres;



insert into tb_user (email,active) values ('alwinalwin57@gmail.com','1')