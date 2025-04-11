--
-- PostgreSQL database dump
--

-- Dumped from database version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.17 (Ubuntu 14.17-0ubuntu0.22.04.1)

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
-- Name: analytics_events; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.analytics_events (
    id integer NOT NULL,
    user_id integer,
    action text NOT NULL,
    menu_item_id integer,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.analytics_events OWNER TO blessed_bites;

--
-- Name: analytics_events_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.analytics_events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.analytics_events_id_seq OWNER TO blessed_bites;

--
-- Name: analytics_events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.analytics_events_id_seq OWNED BY public.analytics_events.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.categories OWNER TO blessed_bites;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.categories_id_seq OWNER TO blessed_bites;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: menu_items; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.menu_items (
    id integer NOT NULL,
    name text NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    category_id integer,
    order_count integer DEFAULT 0,
    is_active boolean DEFAULT true,
    image_url text,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.menu_items OWNER TO blessed_bites;

--
-- Name: menu_items_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.menu_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.menu_items_id_seq OWNER TO blessed_bites;

--
-- Name: menu_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.menu_items_id_seq OWNED BY public.menu_items.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.order_items (
    id integer NOT NULL,
    order_id integer,
    menu_item_id integer,
    quantity integer NOT NULL,
    item_price numeric(10,2) NOT NULL,
    subtotal numeric(10,2) GENERATED ALWAYS AS (((quantity)::numeric * item_price)) STORED
);


ALTER TABLE public.order_items OWNER TO blessed_bites;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.order_items_id_seq OWNER TO blessed_bites;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    user_id integer,
    total_cost numeric(10,2) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    status text DEFAULT 'pending'::text,
    payment_method text
);


ALTER TABLE public.orders OWNER TO blessed_bites;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orders_id_seq OWNER TO blessed_bites;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: recommendations; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.recommendations (
    id integer NOT NULL,
    user_id integer,
    menu_item_id integer,
    reason text,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.recommendations OWNER TO blessed_bites;

--
-- Name: recommendations_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.recommendations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.recommendations_id_seq OWNER TO blessed_bites;

--
-- Name: recommendations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.recommendations_id_seq OWNED BY public.recommendations.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO blessed_bites;

--
-- Name: users; Type: TABLE; Schema: public; Owner: blessed_bites
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL,
    full_name text NOT NULL,
    phone_no text,
    password_hash text NOT NULL,
    role text NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO blessed_bites;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: blessed_bites
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO blessed_bites;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: blessed_bites
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: analytics_events id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.analytics_events ALTER COLUMN id SET DEFAULT nextval('public.analytics_events_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: menu_items id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.menu_items ALTER COLUMN id SET DEFAULT nextval('public.menu_items_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: recommendations id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.recommendations ALTER COLUMN id SET DEFAULT nextval('public.recommendations_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: analytics_events; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.analytics_events (id, user_id, action, menu_item_id, created_at) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.categories (id, name) FROM stdin;
1	Fast Food
2	Lunch Special
3	Drinks
\.


--
-- Data for Name: menu_items; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.menu_items (id, name, description, price, category_id, order_count, is_active, image_url, created_at) FROM stdin;
3	Tacos	test	2.00	2	0	t		2025-04-10 16:31:51.546828
1	Burritos	Dutch cheese and Chicken Burritos	3.50	1	0	t		2025-04-10 13:15:39.070104
4	Salbutes	Fried corn delicacy topped with chicken and pico de gallo	0.33	1	0	t		2025-04-10 17:25:40.632362
9	Breakfast Burritos	testing again	4.50	1	0	t		2025-04-11 08:20:26.837369
10	Escabeche	Onion Soup that rocks	12.00	2	0	t		2025-04-11 08:22:02.356135
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.order_items (id, order_id, menu_item_id, quantity, item_price) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.orders (id, user_id, total_cost, created_at, status, payment_method) FROM stdin;
\.


--
-- Data for Name: recommendations; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.recommendations (id, user_id, menu_item_id, reason, created_at) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.schema_migrations (version, dirty) FROM stdin;
7	f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.users (id, email, full_name, phone_no, password_hash, role, created_at) FROM stdin;
1	2022156707@ub.edu.bz	Amilcar Vasquez	6082424	Firme.2424	user	2025-04-08 23:22:07.710769
\.


--
-- Name: analytics_events_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.analytics_events_id_seq', 1, false);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.categories_id_seq', 3, true);


--
-- Name: menu_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.menu_items_id_seq', 10, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.order_items_id_seq', 1, false);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.orders_id_seq', 1, false);


--
-- Name: recommendations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.recommendations_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.users_id_seq', 2, true);


--
-- Name: analytics_events analytics_events_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.analytics_events
    ADD CONSTRAINT analytics_events_pkey PRIMARY KEY (id);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: menu_items menu_items_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: recommendations recommendations_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.recommendations
    ADD CONSTRAINT recommendations_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: analytics_events analytics_events_menu_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.analytics_events
    ADD CONSTRAINT analytics_events_menu_item_id_fkey FOREIGN KEY (menu_item_id) REFERENCES public.menu_items(id);


--
-- Name: analytics_events analytics_events_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.analytics_events
    ADD CONSTRAINT analytics_events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: menu_items menu_items_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.menu_items
    ADD CONSTRAINT menu_items_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: order_items order_items_menu_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_menu_item_id_fkey FOREIGN KEY (menu_item_id) REFERENCES public.menu_items(id);


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: recommendations recommendations_menu_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.recommendations
    ADD CONSTRAINT recommendations_menu_item_id_fkey FOREIGN KEY (menu_item_id) REFERENCES public.menu_items(id);


--
-- Name: recommendations recommendations_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: blessed_bites
--

ALTER TABLE ONLY public.recommendations
    ADD CONSTRAINT recommendations_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

