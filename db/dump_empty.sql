--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Drop databases (except postgres and template1)
--

DROP DATABASE construction;




--
-- Drop roles
--

DROP ROLE db_admin;


--
-- Roles
--

CREATE ROLE db_admin;
ALTER ROLE db_admin WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:OjEdCs0c5my05J0xtmu90Q==$wjLWz01FphSqAHjtj3TQuO20eOZEXop3JuXEqtjog14=:mIvwCjsEgaOhAthXtWCZJlBYM1cdVlG3g1fbQNNr/WQ=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

UPDATE pg_catalog.pg_database SET datistemplate = false WHERE datname = 'template1';
DROP DATABASE template1;
--
-- Name: template1; Type: DATABASE; Schema: -; Owner: db_admin
--

CREATE DATABASE template1 WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE template1 OWNER TO db_admin;

\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: COMMENT; Schema: -; Owner: db_admin
--

COMMENT ON DATABASE template1 IS 'default template for new databases';


--
-- Name: template1; Type: DATABASE PROPERTIES; Schema: -; Owner: db_admin
--

ALTER DATABASE template1 IS_TEMPLATE = true;


\connect template1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE template1; Type: ACL; Schema: -; Owner: db_admin
--

REVOKE CONNECT,TEMPORARY ON DATABASE template1 FROM PUBLIC;
GRANT CONNECT ON DATABASE template1 TO PUBLIC;


--
-- PostgreSQL database dump complete
--

--
-- Database "construction" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: construction; Type: DATABASE; Schema: -; Owner: db_admin
--

CREATE DATABASE construction WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE construction OWNER TO db_admin;

\connect construction

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: configurations; Type: TABLE; Schema: public; Owner: db_admin
--

CREATE TABLE public.configurations (
    configuration_id integer NOT NULL,
    description character varying(100) NOT NULL,
    engine character varying(50) NOT NULL,
    gearbox_type character varying(50) NOT NULL,
    engine_management_system boolean NOT NULL,
    monitoring_control_system character varying(50) NOT NULL,
    starting_system boolean NOT NULL,
    power_transmission boolean NOT NULL,
    auxiliary_pto character varying(50)[] DEFAULT '{}'::character varying[] NOT NULL,
    oil_system character varying(50)[] DEFAULT '{}'::character varying[] NOT NULL,
    fuel_system character varying(50)[] DEFAULT '{}'::character varying[] NOT NULL,
    cooling_system character varying(50) NOT NULL,
    exhaust_system boolean NOT NULL,
    mounting_system boolean NOT NULL,
    gearbox_option character varying(50)[] DEFAULT '{}'::character varying[] NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.configurations OWNER TO db_admin;

--
-- Name: configurations_configuration_id_seq; Type: SEQUENCE; Schema: public; Owner: db_admin
--

CREATE SEQUENCE public.configurations_configuration_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.configurations_configuration_id_seq OWNER TO db_admin;

--
-- Name: configurations_configuration_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: db_admin
--

ALTER SEQUENCE public.configurations_configuration_id_seq OWNED BY public.configurations.configuration_id;


--
-- Name: cooling_exhaust_analyzer; Type: TABLE; Schema: public; Owner: db_admin
--

CREATE TABLE public.cooling_exhaust_analyzer (
    job_id integer NOT NULL,
    configuration_id integer NOT NULL,
    success boolean DEFAULT false NOT NULL,
    cooling_system character varying(50) DEFAULT ''::character varying NOT NULL,
    exhaust_system boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    order_created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.cooling_exhaust_analyzer OWNER TO db_admin;

--
-- Name: COLUMN cooling_exhaust_analyzer.order_created_at; Type: COMMENT; Schema: public; Owner: db_admin
--

COMMENT ON COLUMN public.cooling_exhaust_analyzer.order_created_at IS 'Timestamp, when the order as a whole was triggered (not only of this specific microservice)';


--
-- Name: cooling_exhaust_analyzer_job_id_seq; Type: SEQUENCE; Schema: public; Owner: db_admin
--

CREATE SEQUENCE public.cooling_exhaust_analyzer_job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cooling_exhaust_analyzer_job_id_seq OWNER TO db_admin;

--
-- Name: cooling_exhaust_analyzer_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: db_admin
--

ALTER SEQUENCE public.cooling_exhaust_analyzer_job_id_seq OWNED BY public.cooling_exhaust_analyzer.job_id;


--
-- Name: engine_control_analyzer; Type: TABLE; Schema: public; Owner: db_admin
--

CREATE TABLE public.engine_control_analyzer (
    job_id integer NOT NULL,
    configuration_id integer NOT NULL,
    success boolean DEFAULT false NOT NULL,
    engine_management_system boolean DEFAULT false NOT NULL,
    monitoring_control_system character varying(50) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    order_created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.engine_control_analyzer OWNER TO db_admin;

--
-- Name: COLUMN engine_control_analyzer.order_created_at; Type: COMMENT; Schema: public; Owner: db_admin
--

COMMENT ON COLUMN public.engine_control_analyzer.order_created_at IS 'Timestamp, when the order as a whole was triggered (not only of this specific microservice)';


--
-- Name: engine_control_analyzer_job_id_seq; Type: SEQUENCE; Schema: public; Owner: db_admin
--

CREATE SEQUENCE public.engine_control_analyzer_job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.engine_control_analyzer_job_id_seq OWNER TO db_admin;

--
-- Name: engine_control_analyzer_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: db_admin
--

ALTER SEQUENCE public.engine_control_analyzer_job_id_seq OWNED BY public.engine_control_analyzer.job_id;


--
-- Name: mounting_analyzer; Type: TABLE; Schema: public; Owner: db_admin
--

CREATE TABLE public.mounting_analyzer (
    job_id integer NOT NULL,
    configuration_id integer NOT NULL,
    success boolean DEFAULT false NOT NULL,
    mounting_system boolean DEFAULT false NOT NULL,
    gearbox_options character varying(50)[] DEFAULT '{}'::character varying[] NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    order_created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.mounting_analyzer OWNER TO db_admin;

--
-- Name: COLUMN mounting_analyzer.order_created_at; Type: COMMENT; Schema: public; Owner: db_admin
--

COMMENT ON COLUMN public.mounting_analyzer.order_created_at IS 'Timestamp, when the order as a whole was triggered (not only of this specific microservice)';


--
-- Name: mounting_analyzer_job_id_seq; Type: SEQUENCE; Schema: public; Owner: db_admin
--

CREATE SEQUENCE public.mounting_analyzer_job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.mounting_analyzer_job_id_seq OWNER TO db_admin;

--
-- Name: mounting_analyzer_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: db_admin
--

ALTER SEQUENCE public.mounting_analyzer_job_id_seq OWNED BY public.mounting_analyzer.job_id;


--
-- Name: propulsion_analyzer; Type: TABLE; Schema: public; Owner: db_admin
--

CREATE TABLE public.propulsion_analyzer (
    job_id integer NOT NULL,
    configuration_id integer NOT NULL,
    success boolean DEFAULT false NOT NULL,
    starting_system boolean DEFAULT false NOT NULL,
    power_transmission boolean DEFAULT false NOT NULL,
    auxiliary_pto character varying[] DEFAULT '{}'::character varying[] NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    order_created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.propulsion_analyzer OWNER TO db_admin;

--
-- Name: COLUMN propulsion_analyzer.order_created_at; Type: COMMENT; Schema: public; Owner: db_admin
--

COMMENT ON COLUMN public.propulsion_analyzer.order_created_at IS 'Timestamp, when the order as a whole was triggered (not only of this specific microservice)';


--
-- Name: propulsion_analyzer_job_id_seq; Type: SEQUENCE; Schema: public; Owner: db_admin
--

CREATE SEQUENCE public.propulsion_analyzer_job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.propulsion_analyzer_job_id_seq OWNER TO db_admin;

--
-- Name: propulsion_analyzer_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: db_admin
--

ALTER SEQUENCE public.propulsion_analyzer_job_id_seq OWNED BY public.propulsion_analyzer.job_id;


--
-- Name: supply_analyzer; Type: TABLE; Schema: public; Owner: db_admin
--

CREATE TABLE public.supply_analyzer (
    job_id integer NOT NULL,
    configuration_id integer NOT NULL,
    success boolean DEFAULT false NOT NULL,
    oil_system character varying[] DEFAULT '{}'::character varying[] NOT NULL,
    fuel_system character varying[] DEFAULT '{}'::character varying[] NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    last_modified timestamp without time zone DEFAULT now() NOT NULL,
    status smallint DEFAULT 0 NOT NULL,
    order_created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.supply_analyzer OWNER TO db_admin;

--
-- Name: COLUMN supply_analyzer.order_created_at; Type: COMMENT; Schema: public; Owner: db_admin
--

COMMENT ON COLUMN public.supply_analyzer.order_created_at IS 'Timestamp, when the order as a whole was triggered (not only of this specific microservice)';


--
-- Name: supply_analyzer_job_id_seq; Type: SEQUENCE; Schema: public; Owner: db_admin
--

CREATE SEQUENCE public.supply_analyzer_job_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.supply_analyzer_job_id_seq OWNER TO db_admin;

--
-- Name: supply_analyzer_job_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: db_admin
--

ALTER SEQUENCE public.supply_analyzer_job_id_seq OWNED BY public.supply_analyzer.job_id;


--
-- Name: configurations configuration_id; Type: DEFAULT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.configurations ALTER COLUMN configuration_id SET DEFAULT nextval('public.configurations_configuration_id_seq'::regclass);


--
-- Name: cooling_exhaust_analyzer job_id; Type: DEFAULT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.cooling_exhaust_analyzer ALTER COLUMN job_id SET DEFAULT nextval('public.cooling_exhaust_analyzer_job_id_seq'::regclass);


--
-- Name: engine_control_analyzer job_id; Type: DEFAULT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.engine_control_analyzer ALTER COLUMN job_id SET DEFAULT nextval('public.engine_control_analyzer_job_id_seq'::regclass);


--
-- Name: mounting_analyzer job_id; Type: DEFAULT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.mounting_analyzer ALTER COLUMN job_id SET DEFAULT nextval('public.mounting_analyzer_job_id_seq'::regclass);


--
-- Name: propulsion_analyzer job_id; Type: DEFAULT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.propulsion_analyzer ALTER COLUMN job_id SET DEFAULT nextval('public.propulsion_analyzer_job_id_seq'::regclass);


--
-- Name: supply_analyzer job_id; Type: DEFAULT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.supply_analyzer ALTER COLUMN job_id SET DEFAULT nextval('public.supply_analyzer_job_id_seq'::regclass);


--
-- Data for Name: configurations; Type: TABLE DATA; Schema: public; Owner: db_admin
--

COPY public.configurations (configuration_id, description, engine, gearbox_type, engine_management_system, monitoring_control_system, starting_system, power_transmission, auxiliary_pto, oil_system, fuel_system, cooling_system, exhaust_system, mounting_system, gearbox_option, created_at, last_modified) FROM stdin;
\.


--
-- Data for Name: cooling_exhaust_analyzer; Type: TABLE DATA; Schema: public; Owner: db_admin
--

COPY public.cooling_exhaust_analyzer (job_id, configuration_id, success, cooling_system, exhaust_system, created_at, last_modified, status, order_created_at) FROM stdin;
\.


--
-- Data for Name: engine_control_analyzer; Type: TABLE DATA; Schema: public; Owner: db_admin
--

COPY public.engine_control_analyzer (job_id, configuration_id, success, engine_management_system, monitoring_control_system, created_at, last_modified, status, order_created_at) FROM stdin;
\.


--
-- Data for Name: mounting_analyzer; Type: TABLE DATA; Schema: public; Owner: db_admin
--

COPY public.mounting_analyzer (job_id, configuration_id, success, mounting_system, gearbox_options, created_at, last_modified, status, order_created_at) FROM stdin;
\.


--
-- Data for Name: propulsion_analyzer; Type: TABLE DATA; Schema: public; Owner: db_admin
--

COPY public.propulsion_analyzer (job_id, configuration_id, success, starting_system, power_transmission, auxiliary_pto, created_at, last_modified, status, order_created_at) FROM stdin;
\.


--
-- Data for Name: supply_analyzer; Type: TABLE DATA; Schema: public; Owner: db_admin
--

COPY public.supply_analyzer (job_id, configuration_id, success, oil_system, fuel_system, created_at, last_modified, status, order_created_at) FROM stdin;
\.


--
-- Name: configurations_configuration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: db_admin
--

SELECT pg_catalog.setval('public.configurations_configuration_id_seq', 1, false);


--
-- Name: cooling_exhaust_analyzer_job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: db_admin
--

SELECT pg_catalog.setval('public.cooling_exhaust_analyzer_job_id_seq', 1, false);


--
-- Name: engine_control_analyzer_job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: db_admin
--

SELECT pg_catalog.setval('public.engine_control_analyzer_job_id_seq', 1, false);


--
-- Name: mounting_analyzer_job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: db_admin
--

SELECT pg_catalog.setval('public.mounting_analyzer_job_id_seq', 1, false);


--
-- Name: propulsion_analyzer_job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: db_admin
--

SELECT pg_catalog.setval('public.propulsion_analyzer_job_id_seq', 1, false);


--
-- Name: supply_analyzer_job_id_seq; Type: SEQUENCE SET; Schema: public; Owner: db_admin
--

SELECT pg_catalog.setval('public.supply_analyzer_job_id_seq', 1, false);


--
-- Name: configurations configurations_pkey; Type: CONSTRAINT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.configurations
    ADD CONSTRAINT configurations_pkey PRIMARY KEY (configuration_id);


--
-- Name: cooling_exhaust_analyzer cooling_exhaust_analyzer_pkey; Type: CONSTRAINT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.cooling_exhaust_analyzer
    ADD CONSTRAINT cooling_exhaust_analyzer_pkey PRIMARY KEY (job_id);


--
-- Name: engine_control_analyzer engine_control_analyzer_pkey; Type: CONSTRAINT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.engine_control_analyzer
    ADD CONSTRAINT engine_control_analyzer_pkey PRIMARY KEY (job_id);


--
-- Name: mounting_analyzer mounting_analyzer_pkey; Type: CONSTRAINT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.mounting_analyzer
    ADD CONSTRAINT mounting_analyzer_pkey PRIMARY KEY (job_id);


--
-- Name: propulsion_analyzer propulsion_analyzer_pkey; Type: CONSTRAINT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.propulsion_analyzer
    ADD CONSTRAINT propulsion_analyzer_pkey PRIMARY KEY (job_id);


--
-- Name: supply_analyzer supply_analyzer_pkey; Type: CONSTRAINT; Schema: public; Owner: db_admin
--

ALTER TABLE ONLY public.supply_analyzer
    ADD CONSTRAINT supply_analyzer_pkey PRIMARY KEY (job_id);


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3 (Debian 16.3-1.pgdg120+1)
-- Dumped by pg_dump version 16.3 (Debian 16.3-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE postgres;
--
-- Name: postgres; Type: DATABASE; Schema: -; Owner: db_admin
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE postgres OWNER TO db_admin;

\connect postgres

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: db_admin
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

