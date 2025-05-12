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
4	Breakfast
\.


--
-- Data for Name: menu_items; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.menu_items (id, name, description, price, category_id, order_count, is_active, image_url, created_at) FROM stdin;
10	Escabeche	Onion Soup that rocks	13.00	2	0	t	./ui/static/img/uploads/20250428_110940_2_escabeche.jpg	2025-04-11 08:22:02.356135
22	Beef Soup	Get that protein up with our beef soup special	13.00	2	0	t	./ui/static/img/uploads/20250428_111129_2_beef-soup.jpg	2025-04-28 11:11:29.116869
23	Beef Tostadas	Tostadas, the beef kind of way! Beef lover's paradise	1.50	1	0	t	./ui/static/img/uploads/20250428_111222_1_beef-tostadas.jpg	2025-04-28 11:12:22.290286
24	Boil up	That goodness that is nutritious and yummy!	15.00	2	0	t	./ui/static/img/uploads/20250508_092912_2_boil-up.jpg	2025-04-28 11:13:09.203901
25	Fried tacos	It's not just fried tacos, it's creaminess added!	12.00	1	0	t	./ui/static/img/uploads/20250508_092926_1_fried-tacos.jpg	2025-04-28 11:14:33.289235
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.order_items (id, order_id, menu_item_id, quantity, item_price) FROM stdin;
1	1	24	0	15.00
2	2	24	0	15.00
3	3	24	2	15.00
4	4	24	7	15.00
5	5	25	2	12.00
6	6	23	5	1.50
7	7	24	1	15.00
8	7	25	5	12.00
9	8	22	2	13.00
10	9	24	5	15.00
11	9	25	2	12.00
12	9	23	3	1.50
13	10	24	3	15.00
14	11	25	12	12.00
15	12	22	2	13.00
16	13	24	3	15.00
17	13	25	2	12.00
18	14	24	2	15.00
19	14	25	3	12.00
20	15	24	1	15.00
21	15	23	1	1.50
22	15	25	1	12.00
23	16	24	1	15.00
24	16	25	1	12.00
25	17	22	1	13.00
26	17	10	1	13.00
27	18	25	1	12.00
28	18	24	1	15.00
29	19	25	1	12.00
30	20	24	1	15.00
31	20	23	1	1.50
32	20	25	1	12.00
33	21	25	7	12.00
34	21	24	1	15.00
35	22	25	1	12.00
36	22	24	1	15.00
37	23	25	1	12.00
38	23	24	1	15.00
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: blessed_bites
--

COPY public.orders (id, user_id, total_cost, created_at, status, payment_method) FROM stdin;
1	5	0.00	2025-05-08 13:22:26.261567	pending	\N
2	5	0.00	2025-05-08 13:29:28.108126	pending	\N
3	5	30.00	2025-05-08 13:33:22.101052	pending	\N
4	5	105.00	2025-05-08 13:35:47.658949	pending	\N
5	5	24.00	2025-05-08 13:48:59.176151	pending	\N
6	5	7.50	2025-05-08 14:20:23.457522	pending	\N
7	5	75.00	2025-05-08 14:24:11.062028	pending	\N
8	5	26.00	2025-05-08 14:28:06.338564	pending	\N
9	5	103.50	2025-05-08 15:17:34.413178	pending	\N
10	5	45.00	2025-05-08 15:38:01.626601	pending	\N
11	5	144.00	2025-05-08 15:38:18.569652	pending	\N
12	5	26.00	2025-05-08 15:38:39.874415	pending	\N
13	5	69.00	2025-05-08 15:44:21.161546	pending	\N
14	5	66.00	2025-05-08 16:19:26.715438	pending	\N
15	5	28.50	2025-05-08 16:25:32.268085	pending	\N
16	5	27.00	2025-05-08 16:26:17.905074	pending	\N
17	5	26.00	2025-05-08 16:30:32.920616	pending	\N
18	5	27.00	2025-05-08 16:35:39.200895	pending	\N
19	5	12.00	2025-05-08 16:37:38.439864	pending	\N
20	5	28.50	2025-05-08 16:45:32.526279	pending	\N
21	5	99.00	2025-05-08 16:46:33.033769	pending	\N
22	5	27.00	2025-05-09 13:26:13.068397	pending	\N
23	5	27.00	2025-05-09 13:31:41.54274	pending	\N
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
5	2022156707@ub.edu.bz	Amilcar	6082424	$2a$10$0O4itWrJ7bPUHthgZg7u7.XShh8MxZTuCl5DCZXQVo/Dc0Kr7ZCCe	admin	2025-05-02 12:44:42.201872
6	ingris@blessedbites.bz	Ingris	6741874	$2a$10$0QlVBfKqr13MKkKuMozPd.bgG7oFlgH//Gz5TyavuNslVavJ23q8i	admin	2025-05-02 17:40:53.461998
7	belizeno@gmail.com	noAdminAmilcar	6082424	$2a$10$o5iZQY7dTS7AuqhLq0R4eOzvrYKVc9FEAUDZ3fN.EhKRbjX3Rb6Fa	user	2025-05-03 23:47:32.728312
9	alessia@vns.edu.bz	Alessia Vasquez	6082424	$2a$10$Y/gFSR5bZtPR3oisFdyH8eFdy/7ek/XsUS8lf5FVtNzbdCilx3g6C	user	2025-05-12 08:40:34.628321
\.


--
-- Name: analytics_events_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.analytics_events_id_seq', 1, false);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.categories_id_seq', 5, true);


--
-- Name: menu_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.menu_items_id_seq', 28, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.order_items_id_seq', 38, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.orders_id_seq', 23, true);


--
-- Name: recommendations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.recommendations_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: blessed_bites
--

SELECT pg_catalog.setval('public.users_id_seq', 9, true);


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

