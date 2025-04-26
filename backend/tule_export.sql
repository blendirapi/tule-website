--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8
-- Dumped by pg_dump version 16.8

-- Started on 2025-04-26 00:08:19

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

DROP DATABASE tuledb;
--
-- TOC entry 4869 (class 1262 OID 24718)
-- Name: tuledb; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE tuledb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'Greek_Greece.1252';


ALTER DATABASE tuledb OWNER TO postgres;

\connect tuledb

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
-- TOC entry 220 (class 1259 OID 24742)
-- Name: bookings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bookings (
    booking_id integer NOT NULL,
    firstname character varying(255) NOT NULL,
    lastname character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    phone character varying(255) NOT NULL,
    date date NOT NULL,
    start_time character varying NOT NULL,
    end_time character varying NOT NULL,
    user_id integer,
    service_id integer,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.bookings OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 24741)
-- Name: bookings_booking_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.bookings_booking_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.bookings_booking_id_seq OWNER TO postgres;

--
-- TOC entry 4870 (class 0 OID 0)
-- Dependencies: 219
-- Name: bookings_booking_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bookings_booking_id_seq OWNED BY public.bookings.booking_id;


--
-- TOC entry 216 (class 1259 OID 24720)
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.services (
    service_id integer NOT NULL,
    name character varying(255) NOT NULL,
    color character varying(255) NOT NULL
);


ALTER TABLE public.services OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 24719)
-- Name: services_service_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.services_service_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.services_service_id_seq OWNER TO postgres;

--
-- TOC entry 4871 (class 0 OID 0)
-- Dependencies: 215
-- Name: services_service_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.services_service_id_seq OWNED BY public.services.service_id;


--
-- TOC entry 218 (class 1259 OID 24729)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    name character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 24728)
-- Name: users_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_user_id_seq OWNER TO postgres;

--
-- TOC entry 4872 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- TOC entry 4702 (class 2604 OID 24745)
-- Name: bookings booking_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings ALTER COLUMN booking_id SET DEFAULT nextval('public.bookings_booking_id_seq'::regclass);


--
-- TOC entry 4698 (class 2604 OID 24723)
-- Name: services service_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services ALTER COLUMN service_id SET DEFAULT nextval('public.services_service_id_seq'::regclass);


--
-- TOC entry 4699 (class 2604 OID 24732)
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- TOC entry 4863 (class 0 OID 24742)
-- Dependencies: 220
-- Data for Name: bookings; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.bookings VALUES (7, 'Jasper', 'Fox', 'fadarotuwo@mailinator.com', '+1 (635) 557-5678', '2025-04-30', '10:00', '15:30', 2, 2, '2025-04-24 23:38:54.700821', '2025-04-24 23:38:54.700821');
INSERT INTO public.bookings VALUES (8, 'Clare', 'Case', 'lapakyn@mailinator.com', '+1 (473) 693-6622', '2025-04-30', '20:30', '21:00', 2, 2, '2025-04-25 00:30:41.623642', '2025-04-25 00:30:41.623642');
INSERT INTO public.bookings VALUES (10, 'Idola', 'Whitfield', 'dehy@mailinator.com', '+1 (189) 133-3544', '2025-04-30', '18:30', '19:00', 2, 2, '2025-04-25 00:30:56.227707', '2025-04-25 00:30:56.227707');
INSERT INTO public.bookings VALUES (13, 'Beverly', 'Bradley', 'sukysimenu@mailinator.com', '+1 (121) 822-2768', '2025-04-30', '20:00', '20:30', 2, 2, '2025-04-25 01:01:55.423339', '2025-04-25 01:01:55.423339');
INSERT INTO public.bookings VALUES (18, 'Calista', 'Roman', 'waboda@mailinator.com', '+1 (981) 927-3283', '2025-04-30', '16:30', '17:00', 2, 2, '2025-04-25 11:02:13.649085', '2025-04-25 11:02:13.649085');
INSERT INTO public.bookings VALUES (19, 'Vera', 'Barr', 'boviro@mailinator.com', '+1 (767) 907-2904', '2025-04-30', '17:00', '17:30', 2, 2, '2025-04-25 11:04:59.514409', '2025-04-25 11:04:59.514409');
INSERT INTO public.bookings VALUES (21, 'aaaaaaa', 'aaaaaaa', 'cyno@mailinator.com', '+1 (363) 889-5319', '2025-04-30', '17:30', '18:00', 2, 2, '2025-04-25 11:08:08.280604', '2025-04-25 11:08:08.280604');
INSERT INTO public.bookings VALUES (22, 'Russell', 'Lynn', 'luladadiky@mailinator.com', '+1 (716) 184-4311', '2025-04-30', '18:00', '18:30', 2, 2, '2025-04-25 11:11:03.465392', '2025-04-25 11:11:03.465392');
INSERT INTO public.bookings VALUES (23, 'Aileen', 'Allen', 'zagivuwi@mailinator.com', '+1 (738) 621-7199', '2025-04-30', '18:00', '18:30', 2, 2, '2025-04-25 11:13:16.730387', '2025-04-25 11:13:16.730387');
INSERT INTO public.bookings VALUES (24, 'Armando', 'Dennis', 'ziqodudivy@mailinator.com', '+1 (602) 376-1414', '2025-04-30', '19:00', '19:30', 2, 2, '2025-04-25 11:14:07.451923', '2025-04-25 11:14:07.451923');
INSERT INTO public.bookings VALUES (25, 'Octavia', 'Delgado', 'gociwojyh@mailinator.com', '+1 (885) 169-6914', '2025-04-30', '19:30', '20:00', 2, 2, '2025-04-25 11:15:55.289179', '2025-04-25 11:15:55.289179');
INSERT INTO public.bookings VALUES (26, 'Wanda', 'Osborne', 'qoqugukun@mailinator.com', '+1 (864) 924-8765', '2025-04-29', '12:30', '13:00', 2, 2, '2025-04-25 11:16:51.123559', '2025-04-25 11:16:51.123559');
INSERT INTO public.bookings VALUES (27, 'Bryar', 'Allen', 'rixud@mailinator.com', '+1 (929) 511-8554', '2025-04-26', '11:00', '11:30', 2, 1, '2025-04-25 11:17:23.806409', '2025-04-25 11:17:23.806409');
INSERT INTO public.bookings VALUES (28, 'Chava', 'Workman', 'qaquwygyr@mailinator.com', '+1 (857) 309-9194', '2025-04-26', '17:00', '17:30', 2, 2, '2025-04-25 11:21:43.071032', '2025-04-25 11:21:43.071032');
INSERT INTO public.bookings VALUES (29, 'Risa', 'Madden', 'zanasuhuby@mailinator.com', '+1 (666) 386-1968', '2025-04-26', '09:30', '10:00', 2, 2, '2025-04-25 11:24:04.87172', '2025-04-25 11:24:04.87172');
INSERT INTO public.bookings VALUES (30, 'Dara', 'Branch', 'xigiqic@mailinator.com', '+1 (585) 707-2415', '2025-04-29', '17:30', '18:00', 2, 2, '2025-04-25 11:53:10.705481', '2025-04-25 11:53:10.705481');
INSERT INTO public.bookings VALUES (14, 'Virginia', 'Rivas', 'vahitaqa@mailinator.com', '+1 (904) 707-8472', '2025-04-30', '09:00', '09:30', 1, 2, '2025-04-25 10:33:17.700683', '2025-04-25 10:33:17.700683');
INSERT INTO public.bookings VALUES (15, 'Connor', 'Slater', 'vybuw@mailinator.com', '+1 (548) 204-6219', '2025-04-30', '09:30', '10:00', 1, 2, '2025-04-25 10:33:35.803281', '2025-04-25 10:33:35.803281');
INSERT INTO public.bookings VALUES (16, 'Quemby', 'Moody', 'sefimek@mailinator.com', '+1 (961) 724-4264', '2025-04-30', '15:30', '16:00', 1, 2, '2025-04-25 10:34:51.233914', '2025-04-25 10:34:51.233914');
INSERT INTO public.bookings VALUES (17, 'Jelani', 'Calhoun', 'nykiri@mailinator.com', '+1 (151) 864-1508', '2025-04-30', '16:00', '16:30', 1, 2, '2025-04-25 10:35:06.04996', '2025-04-25 10:35:06.04996');
INSERT INTO public.bookings VALUES (31, 'Danielle', 'Clay', 'nepiv@mailinator.com', '+1 (574) 165-8637', '2025-04-30', '16:00', '16:30', 2, 2, '2025-04-25 23:50:34.186937', '2025-04-25 23:50:34.186937');


--
-- TOC entry 4859 (class 0 OID 24720)
-- Dependencies: 216
-- Data for Name: services; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.services VALUES (1, 'Κούρεμα', '#FF0000');
INSERT INTO public.services VALUES (2, 'Κούρεμα & Γενιάδα', '#00FF00');


--
-- TOC entry 4861 (class 0 OID 24729)
-- Dependencies: 218
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users VALUES (1, 'Tule', 'tule', '123456', '2025-04-24 16:55:13.497852', '2025-04-24 16:55:13.497852');
INSERT INTO public.users VALUES (2, 'Άγγελος', 'ahmeti', '1111', '2025-04-24 18:10:59.621866', '2025-04-24 18:10:59.621866');


--
-- TOC entry 4873 (class 0 OID 0)
-- Dependencies: 219
-- Name: bookings_booking_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bookings_booking_id_seq', 31, true);


--
-- TOC entry 4874 (class 0 OID 0)
-- Dependencies: 215
-- Name: services_service_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.services_service_id_seq', 2, true);


--
-- TOC entry 4875 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_user_id_seq', 4, true);


--
-- TOC entry 4712 (class 2606 OID 24751)
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (booking_id);


--
-- TOC entry 4706 (class 2606 OID 24727)
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (service_id);


--
-- TOC entry 4708 (class 2606 OID 24738)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- TOC entry 4710 (class 2606 OID 24740)
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- TOC entry 4713 (class 2606 OID 24757)
-- Name: bookings bookings_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(service_id);


--
-- TOC entry 4714 (class 2606 OID 24752)
-- Name: bookings bookings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


-- Completed on 2025-04-26 00:08:20

--
-- PostgreSQL database dump complete
--

