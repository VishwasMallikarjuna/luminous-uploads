CREATE TABLE IF NOT EXISTS images
(
    id integer NOT NULL DEFAULT nextval('images_id_seq'::regclass),
    image_data bytea NOT NULL,
    image_hash text COLLATE pg_catalog."default",
    CONSTRAINT images_pkey PRIMARY KEY (id)
)

CREATE TABLE IF NOT EXISTS image_detail
(
    id integer NOT NULL,
    width integer NOT NULL,
    height integer NOT NULL,
    camera_model character varying(255) COLLATE pg_catalog."default",
    location character varying(255) COLLATE pg_catalog."default",
    format character varying(50) COLLATE pg_catalog."default" NOT NULL,
    upload_timestamp timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT image_detail_pkey PRIMARY KEY (id)
)
