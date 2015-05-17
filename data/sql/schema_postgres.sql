SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;


CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;
SET default_tablespace = '';
SET default_with_oids = false;


--
-- Charon Group
--
CREATE TABLE charon_group (
    id integer NOT NULL,
    name character varying(80) NOT NULL
);
CREATE SEQUENCE charon_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_group_id_seq OWNED BY charon_group.id;


--
-- Charon Group Permissions
--
CREATE TABLE charon_group_permissions (
    id integer NOT NULL,
    group_id integer NOT NULL,
    permission_id integer NOT NULL
);
CREATE SEQUENCE charon_group_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_group_permissions_id_seq OWNED BY charon_group_permissions.id;


--
-- Charon Permission
--
CREATE TABLE charon_permission (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    content_type_id integer NOT NULL,
    codename character varying(100) NOT NULL
);
CREATE SEQUENCE charon_permission_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_permission_id_seq OWNED BY charon_permission.id;

--
-- Charon UserPermission
--
CREATE TABLE charon_user_permissions (
    id integer NOT NULL,
    user_id integer NOT NULL,
    permission_id integer NOT NULL
);
CREATE SEQUENCE charon_user_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_user_permissions_id_seq OWNED BY charon_user_permissions.id;

--
-- Charon User
--
CREATE TABLE charon_user (
    id integer NOT NULL,
    password character varying(128) NOT NULL,
    username character varying(75) NOT NULL,
    first_name character varying(45) NOT NULL,
    last_name character varying(45) NOT NULL,
    is_superuser boolean NOT NULL,
    is_active boolean NOT NULL,
    is_staff boolean NOT NULL,
    is_confirmed boolean NOT NULL,
    confirmation_token character varying(122) NOT NULL,
    last_login_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
CREATE SEQUENCE charon_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_user_id_seq OWNED BY charon_user.id;


--
-- Charon User Groups
--
CREATE TABLE charon_user_groups (
    id integer NOT NULL,
    user_id integer NOT NULL,
    group_id integer NOT NULL
);
CREATE SEQUENCE charon_user_groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_user_groups_id_seq OWNED BY charon_user_groups.id;


--
-- Charon ContentType
--
CREATE TABLE charon_content_type (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    app_label character varying(100) NOT NULL,
    model character varying(100) NOT NULL
);
CREATE SEQUENCE charon_content_type_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE charon_content_type_id_seq OWNED BY charon_content_type.id;


ALTER TABLE ONLY charon_group ALTER COLUMN id SET DEFAULT nextval('charon_group_id_seq'::regclass);
ALTER TABLE ONLY charon_group_permissions ALTER COLUMN id SET DEFAULT nextval('charon_group_permissions_id_seq'::regclass);
ALTER TABLE ONLY charon_permission ALTER COLUMN id SET DEFAULT nextval('charon_permission_id_seq'::regclass);
ALTER TABLE ONLY charon_user ALTER COLUMN id SET DEFAULT nextval('charon_user_id_seq'::regclass);
ALTER TABLE ONLY charon_user_groups ALTER COLUMN id SET DEFAULT nextval('charon_user_groups_id_seq'::regclass);
ALTER TABLE ONLY charon_user_permissions ALTER COLUMN id SET DEFAULT nextval('charon_user_permissions_id_seq'::regclass);
ALTER TABLE ONLY charon_content_type ALTER COLUMN id SET DEFAULT nextval('charon_content_type_id_seq'::regclass);

ALTER TABLE ONLY charon_group
    ADD CONSTRAINT charon_group_name_key UNIQUE (name);

ALTER TABLE ONLY charon_group_permissions
    ADD CONSTRAINT charon_group_permissions_group_id_permission_id_key UNIQUE (group_id, permission_id);

ALTER TABLE ONLY charon_group_permissions
    ADD CONSTRAINT charon_group_permissions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY charon_group
    ADD CONSTRAINT charon_group_pkey PRIMARY KEY (id);

ALTER TABLE ONLY charon_permission
    ADD CONSTRAINT charon_permission_content_type_id_codename_key UNIQUE (content_type_id, codename);

ALTER TABLE ONLY charon_permission
    ADD CONSTRAINT charon_permission_pkey PRIMARY KEY (id);

ALTER TABLE ONLY charon_user_groups
    ADD CONSTRAINT charon_user_groups_pkey PRIMARY KEY (id);

ALTER TABLE ONLY charon_user_groups
    ADD CONSTRAINT charon_user_groups_user_id_3da41bdcd69daabb_uniq UNIQUE (user_id, group_id);

ALTER TABLE ONLY charon_user
    ADD CONSTRAINT charon_user_pkey PRIMARY KEY (id);

ALTER TABLE ONLY charon_user_permissions
    ADD CONSTRAINT charon_user_permissions_id_seq_pkey PRIMARY KEY (id);

ALTER TABLE ONLY charon_user_permissions
    ADD CONSTRAINT charon_user_permissions_id_seq_user_id_505ff7b6d553b31a_uniq UNIQUE (user_id, permission_id);

ALTER TABLE ONLY charon_user
    ADD CONSTRAINT charon_user_username_key UNIQUE (username);

ALTER TABLE ONLY charon_content_type
    ADD CONSTRAINT charon_content_type_app_label_model_key UNIQUE (app_label, model);

ALTER TABLE ONLY charon_content_type
    ADD CONSTRAINT charon_content_type_pkey PRIMARY KEY (id);

CREATE INDEX charon_group_name_like ON charon_group USING btree (name varchar_pattern_ops);
CREATE INDEX charon_group_permissions_group_id ON charon_group_permissions USING btree (group_id);
CREATE INDEX charon_group_permissions_permission_id ON charon_group_permissions USING btree (permission_id);
CREATE INDEX charon_permission_content_type_id ON charon_permission USING btree (content_type_id);
CREATE INDEX charon_user_groups_group_id ON charon_user_groups USING btree (group_id);
CREATE INDEX charon_user_groups_user_id ON charon_user_groups USING btree (user_id);
CREATE INDEX charon_user_permissions_permission_id ON charon_user_permissions USING btree (permission_id);
CREATE INDEX charon_user_permissions_user_id ON charon_user_permissions USING btree (user_id);
CREATE INDEX charon_user_username_like ON charon_user USING btree (username varchar_pattern_ops);


ALTER TABLE ONLY charon_group_permissions
    ADD CONSTRAINT charon_group_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES charon_permission(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE ONLY charon_permission
    ADD CONSTRAINT content_type_id_refs_id_d043b34a FOREIGN KEY (content_type_id) REFERENCES charon_content_type(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE ONLY charon_user_groups
    ADD CONSTRAINT group_id_refs_id_179f149d FOREIGN KEY (group_id) REFERENCES charon_group(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE ONLY charon_group_permissions
    ADD CONSTRAINT group_id_refs_id_f4b32aac FOREIGN KEY (group_id) REFERENCES charon_group(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE ONLY charon_user_permissions
    ADD CONSTRAINT permission_id_refs_id_36d649dc FOREIGN KEY (permission_id) REFERENCES charon_permission(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE ONLY charon_user_groups
    ADD CONSTRAINT user_id_refs_id_0d277b3b FOREIGN KEY (user_id) REFERENCES charon_user(id) DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE ONLY charon_user_permissions
    ADD CONSTRAINT user_id_refs_id_f1e5c798 FOREIGN KEY (user_id) REFERENCES charon_user(id) DEFERRABLE INITIALLY DEFERRED;
