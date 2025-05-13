--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8
-- Dumped by pg_dump version 16.8

-- Started on 2025-05-14 00:30:53

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
-- TOC entry 4889 (class 1262 OID 24718)
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
    firstname character varying(255),
    lastname character varying(255),
    email character varying(255),
    phone character varying(255),
    date date NOT NULL,
    start_time character varying NOT NULL,
    end_time character varying NOT NULL,
    user_id integer,
    service_id integer,
    bath boolean DEFAULT false NOT NULL
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
-- TOC entry 4890 (class 0 OID 0)
-- Dependencies: 219
-- Name: bookings_booking_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bookings_booking_id_seq OWNED BY public.bookings.booking_id;


--
-- TOC entry 224 (class 1259 OID 24799)
-- Name: clients; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clients (
    client_id integer NOT NULL,
    firstname character varying(50),
    lastname character varying(50),
    email character varying(100),
    phone character varying(20)
);


ALTER TABLE public.clients OWNER TO postgres;

--
-- TOC entry 223 (class 1259 OID 24798)
-- Name: clients_client_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.clients_client_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.clients_client_id_seq OWNER TO postgres;

--
-- TOC entry 4891 (class 0 OID 0)
-- Dependencies: 223
-- Name: clients_client_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.clients_client_id_seq OWNED BY public.clients.client_id;


--
-- TOC entry 216 (class 1259 OID 24720)
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.services (
    service_id integer NOT NULL,
    name character varying(255) NOT NULL,
    color character varying(255) NOT NULL,
    "time" character varying(15)
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
-- TOC entry 4892 (class 0 OID 0)
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
    is_visible boolean,
    email character varying(255)
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
-- TOC entry 4893 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_user_id_seq OWNED BY public.users.user_id;


--
-- TOC entry 222 (class 1259 OID 24777)
-- Name: working_hours; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.working_hours (
    working_hour_id integer NOT NULL,
    user_id integer,
    day_of_week character varying(50),
    start_time character varying(10),
    end_time character varying(10)
);


ALTER TABLE public.working_hours OWNER TO postgres;

--
-- TOC entry 221 (class 1259 OID 24776)
-- Name: working_hours_working_hour_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.working_hours_working_hour_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.working_hours_working_hour_id_seq OWNER TO postgres;

--
-- TOC entry 4894 (class 0 OID 0)
-- Dependencies: 221
-- Name: working_hours_working_hour_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.working_hours_working_hour_id_seq OWNED BY public.working_hours.working_hour_id;


--
-- TOC entry 4710 (class 2604 OID 24745)
-- Name: bookings booking_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings ALTER COLUMN booking_id SET DEFAULT nextval('public.bookings_booking_id_seq'::regclass);


--
-- TOC entry 4713 (class 2604 OID 24802)
-- Name: clients client_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients ALTER COLUMN client_id SET DEFAULT nextval('public.clients_client_id_seq'::regclass);


--
-- TOC entry 4708 (class 2604 OID 24723)
-- Name: services service_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services ALTER COLUMN service_id SET DEFAULT nextval('public.services_service_id_seq'::regclass);


--
-- TOC entry 4709 (class 2604 OID 24732)
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_user_id_seq'::regclass);


--
-- TOC entry 4712 (class 2604 OID 24780)
-- Name: working_hours working_hour_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.working_hours ALTER COLUMN working_hour_id SET DEFAULT nextval('public.working_hours_working_hour_id_seq'::regclass);


--
-- TOC entry 4879 (class 0 OID 24742)
-- Dependencies: 220
-- Data for Name: bookings; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.bookings VALUES (79, '', '', '', '', '2025-05-14', '09:00', '21:00', 1, 19, false);
INSERT INTO public.bookings VALUES (81, 'Zephania', 'Robinson', 'qalici@mailinator.com', '+1 (473) 981-3421', '2025-05-21', '13:30', '14:00', 1, 2, true);
INSERT INTO public.bookings VALUES (82, 'Nichole', 'Sullivan', 'kepejo@mailinator.com', '+1 (187) 853-4104', '2025-05-21', '13:00', '13:30', 1, 2, false);
INSERT INTO public.bookings VALUES (83, 'Ivory', 'Riggs', 'xuvoro@mailinator.com', '1', '2025-05-21', '12:30', '13:00', 1, 1, false);
INSERT INTO public.bookings VALUES (84, 'Ulric', 'Kim', 'zynute@mailinator.com', '2', '2025-05-22', '12:30', '13:00', 1, 1, false);
INSERT INTO public.bookings VALUES (85, 'Kirestin', 'Schneider', 'teceko@mailinator.com', '1', '2025-05-24', '12:30', '13:00', 1, 1, false);
INSERT INTO public.bookings VALUES (86, 'Paula', 'Mcclain', 'nady@mailinator.com', '+1 (294) 875-9516', '1990-06-30', '05:57', '10:57', 1, 15, true);
INSERT INTO public.bookings VALUES (74, 'Ulric', 'Kim', 'zynute@mailinator.com', '1111111112', '2025-05-28', '10:00', '10:30', 1, 1, false);
INSERT INTO public.bookings VALUES (88, 'Keely', 'Wheeler', 'huwade@mailinator.com', '+1 (236) 695-6962', '2009-01-28', '23:43', '03:43', 1, 14, false);
INSERT INTO public.bookings VALUES (89, 'Lillian', 'Skinner', 'juvi@mailinator.com', '+1 (524) 973-3368', '2025-05-21', '11:00', '11:30', 1, 1, false);
INSERT INTO public.bookings VALUES (90, 'Reed', 'Wilkinson', 'jocesale@mailinator.com', '+1 (495) 141-8955', '2025-05-15', '14:44', '15:14', 1, 2, true);
INSERT INTO public.bookings VALUES (91, 'Quail', 'Kinney', 'bitoxe@mailinator.com', '+1 (169) 461-7176', '2024-05-14', '12:51', '16:51', 2, 14, true);
INSERT INTO public.bookings VALUES (87, 'aaaaaa', 'aaaaaaaa', 'aaaaaaaaaaa', '1111111112', '2025-05-13', '14:00', '14:30', 1, 1, false);
INSERT INTO public.bookings VALUES (94, 'Grant', 'Stark', 'tawilexery@mailinator.com', '+1 (117) 682-1586', '2025-05-13', '14:41', '15:41', 1, 11, true);
INSERT INTO public.bookings VALUES (95, 'nhbgmjhngb', 'njhbg', 'nyhbgfnhgbfvd', 'nbgvfmnhgbf', '2025-05-13', '10:00', '14:00', 2, 14, false);


--
-- TOC entry 4883 (class 0 OID 24799)
-- Dependencies: 224
-- Data for Name: clients; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.clients VALUES (1, 'Ivoryhfgtdtbgrvfd', 'Riggs', 'xuvoro@mailinator.com', '11111111111111');
INSERT INTO public.clients VALUES (3, 'nhgbf', 'njhgbf', 'njhgbfv', 'jnhgb');
INSERT INTO public.clients VALUES (5, 'Keely', 'Wheeler', 'huwade@mailinator.com', '+1 (236) 695-6962');
INSERT INTO public.clients VALUES (6, 'Lillian', 'Skinner', 'juvi@mailinator.com', '+1 (524) 973-3368');
INSERT INTO public.clients VALUES (7, 'Reed', 'Wilkinson', 'jocesale@mailinator.com', '+1 (495) 141-8955');
INSERT INTO public.clients VALUES (8, 'Quail', 'Kinney', 'bitoxe@mailinator.com', '+1 (169) 461-7176');
INSERT INTO public.clients VALUES (2, 'aaaaaa', 'aaaaaaaa', 'aaaaaaaaaaa', '1111111112');
INSERT INTO public.clients VALUES (10, 'Jayme', 'Wall', 'socosih@mailinator.com', '+1 (617) 803-2026');
INSERT INTO public.clients VALUES (11, '', '', '', '');
INSERT INTO public.clients VALUES (12, 'Grant', 'Stark', 'tawilexery@mailinator.com', '+1 (117) 682-1586');
INSERT INTO public.clients VALUES (13, 'nhbg', 'njhbg', 'nyhbgf', 'nbgvf');
INSERT INTO public.clients VALUES (14, 'nhbgmjhngb', 'njhbg', 'nyhbgfnhgbfvd', 'nbgvfmnhgbf');


--
-- TOC entry 4875 (class 0 OID 24720)
-- Dependencies: 216
-- Data for Name: services; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.services VALUES (10, 'Tattoo', '#000000', '00:30');
INSERT INTO public.services VALUES (11, 'Tattoo', '#000000', '01:00');
INSERT INTO public.services VALUES (12, 'Tattoo', '#000000', '02:00');
INSERT INTO public.services VALUES (13, 'Tattoo', '#000000', '03:00');
INSERT INTO public.services VALUES (14, 'Tattoo', '#000000', '04:00');
INSERT INTO public.services VALUES (15, 'Tattoo', '#000000', '05:00');
INSERT INTO public.services VALUES (16, 'Tattoo', '#000000', '06:00');
INSERT INTO public.services VALUES (2, 'Κούρεμα & Γενιάδα', '#0c16e2', '00:30');
INSERT INTO public.services VALUES (1, 'Κούρεμα', '#1cbf20', '00:30');
INSERT INTO public.services VALUES (19, 'Κενο', '#dd1313', '00:00');


--
-- TOC entry 4877 (class 0 OID 24729)
-- Dependencies: 218
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users VALUES (2, 'Άγγελος', 'aggelos', '$argon2id$v=19$m=65536,t=1,p=16$v//E9Dz9W82cdGxGIpAgKg$ZLTqhhetTjVXXMyCzPedtyzR/lFEABymnmcyo73uyt0', true, 'aggelos@tule.gr');
INSERT INTO public.users VALUES (8, 'admin', 'admin', '$argon2id$v=19$m=65536,t=1,p=16$AnV5GTtxnkzypKu+zclMQw$E2TGfHWN5E0e2mHmhQ3P36a4DPIDpYOf/+2gLQ+Tf1o', false, 'blendirapi@gmail.com');
INSERT INTO public.users VALUES (1, 'Tule', 'tule', '$argon2id$v=19$m=65536,t=1,p=16$GV6xr5FM4Y9tDh+iRZcKRw$BLZkRCA7O+j/XE7JGgmjQjtWJFFUCpAmzSbD8ajeRmE', true, 'tule@tule.gr');


--
-- TOC entry 4881 (class 0 OID 24777)
-- Dependencies: 222
-- Data for Name: working_hours; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.working_hours VALUES (211, 1, 'monday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (212, 1, 'monday', '17:00', '21:00');
INSERT INTO public.working_hours VALUES (213, 1, 'tuesday', '17:00', '21:00');
INSERT INTO public.working_hours VALUES (214, 1, 'tuesday', '10:00', '14:00');
INSERT INTO public.working_hours VALUES (215, 1, 'wednesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (216, 1, 'wednesday', '10:00', '14:00');
INSERT INTO public.working_hours VALUES (217, 1, 'thursday', '17:00', '21:00');
INSERT INTO public.working_hours VALUES (218, 1, 'thursday', '10:00', '14:00');
INSERT INTO public.working_hours VALUES (219, 1, 'friday', '17:00', '21:00');
INSERT INTO public.working_hours VALUES (220, 1, 'friday', '10:00', '14:00');
INSERT INTO public.working_hours VALUES (221, 1, 'saturday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (222, 1, 'saturday', '10:00', '18:00');
INSERT INTO public.working_hours VALUES (223, 1, 'sunday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (224, 1, 'sunday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (183, 2, 'sunday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (184, 2, 'sunday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (185, 2, 'monday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (186, 2, 'monday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (187, 2, 'tuesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (188, 2, 'tuesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (189, 2, 'wednesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (190, 2, 'wednesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (191, 2, 'thursday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (192, 2, 'thursday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (193, 2, 'friday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (194, 2, 'friday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (195, 2, 'saturday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (196, 2, 'saturday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (197, 8, 'saturday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (198, 8, 'saturday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (199, 8, 'sunday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (200, 8, 'sunday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (201, 8, 'monday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (202, 8, 'monday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (203, 8, 'tuesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (204, 8, 'tuesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (205, 8, 'wednesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (206, 8, 'wednesday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (207, 8, 'thursday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (208, 8, 'thursday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (209, 8, 'friday', '00:00', '00:00');
INSERT INTO public.working_hours VALUES (210, 8, 'friday', '00:00', '00:00');


--
-- TOC entry 4895 (class 0 OID 0)
-- Dependencies: 219
-- Name: bookings_booking_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bookings_booking_id_seq', 95, true);


--
-- TOC entry 4896 (class 0 OID 0)
-- Dependencies: 223
-- Name: clients_client_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.clients_client_id_seq', 15, true);


--
-- TOC entry 4897 (class 0 OID 0)
-- Dependencies: 215
-- Name: services_service_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.services_service_id_seq', 21, true);


--
-- TOC entry 4898 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_user_id_seq', 33, true);


--
-- TOC entry 4899 (class 0 OID 0)
-- Dependencies: 221
-- Name: working_hours_working_hour_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.working_hours_working_hour_id_seq', 252, true);


--
-- TOC entry 4721 (class 2606 OID 24751)
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (booking_id);


--
-- TOC entry 4725 (class 2606 OID 24806)
-- Name: clients clients_phone_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_phone_key UNIQUE (phone);


--
-- TOC entry 4727 (class 2606 OID 24804)
-- Name: clients clients_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (client_id);


--
-- TOC entry 4715 (class 2606 OID 24727)
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (service_id);


--
-- TOC entry 4717 (class 2606 OID 24738)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- TOC entry 4719 (class 2606 OID 24740)
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- TOC entry 4723 (class 2606 OID 24782)
-- Name: working_hours working_hours_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.working_hours
    ADD CONSTRAINT working_hours_pkey PRIMARY KEY (working_hour_id);


--
-- TOC entry 4728 (class 2606 OID 24757)
-- Name: bookings bookings_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(service_id);


--
-- TOC entry 4729 (class 2606 OID 24793)
-- Name: bookings bookings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


--
-- TOC entry 4730 (class 2606 OID 24788)
-- Name: working_hours working_hours_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.working_hours
    ADD CONSTRAINT working_hours_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) ON DELETE CASCADE;


-- Completed on 2025-05-14 00:30:54

--
-- PostgreSQL database dump complete
--

