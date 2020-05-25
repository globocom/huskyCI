--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

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
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: accessToken; Type: TABLE; Schema: public; Owner: huskyCIUser
--

CREATE TABLE IF NOT EXISTS public."accessToken" (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    huskytoken text NOT NULL,
    "repositoryURL" text NOT NULL,
    "isValid" boolean NOT NULL,
    "createdAt" timestamp without time zone NOT NULL,
    salt text NOT NULL,
    uuid text NOT NULL
);


ALTER TABLE public."accessToken" OWNER TO "huskyCIUser";

--
-- Name: analysis; Type: TABLE; Schema: public; Owner: huskyCIUser
--

CREATE TABLE IF NOT EXISTS public.analysis (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    "RID" text NOT NULL,
    "repositoryURL" text NOT NULL,
    "repositoryBranch" text NOT NULL,
    "commitAuthors" text[],
    status text NOT NULL,
    result text,
    "errorFound" text,
    containers jsonb,
    "startedAt" timestamp without time zone,
    "finishedAt" timestamp without time zone,
    codes jsonb,
    huskyciresults jsonb
);


ALTER TABLE public.analysis OWNER TO "huskyCIUser";

--
-- Name: repository; Type: TABLE; Schema: public; Owner: huskyCIUser
--

CREATE TABLE IF NOT EXISTS public.repository (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    "repositoryURL" text NOT NULL,
    "repositoryBranch" text,
    "createdAt" timestamp without time zone NOT NULL
);


ALTER TABLE public.repository OWNER TO "huskyCIUser";

--
-- Name: securityTest; Type: TABLE; Schema: public; Owner: huskyCIUser
--

CREATE TABLE IF NOT EXISTS public."securityTest" (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    image text NOT NULL,
    "imageTag" text,
    cmd text NOT NULL,
    type text NOT NULL,
    language text NOT NULL,
    "default" boolean NOT NULL,
    "timeOutSeconds" integer NOT NULL
);


ALTER TABLE public."securityTest" OWNER TO "huskyCIUser";

--
-- Name: user; Type: TABLE; Schema: public; Owner: huskyCIUser
--

CREATE TABLE IF NOT EXISTS public."user" (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username text NOT NULL,
    password text NOT NULL,
    salt text,
    iterations integer,
    keylen integer,
    hashfunction text,
    "newPassword" text,
    "confirmNewPassword" text
);


ALTER TABLE public."user" OWNER TO "huskyCIUser";

--
-- Data for Name: accessToken; Type: TABLE DATA; Schema: public; Owner: huskyCIUser
--

COPY public."accessToken" (id, huskytoken, "repositoryURL", "isValid", "createdAt", salt, uuid) FROM stdin;
\.


--
-- Data for Name: analysis; Type: TABLE DATA; Schema: public; Owner: huskyCIUser
--

COPY public.analysis (id, "RID", "repositoryURL", "repositoryBranch", "commitAuthors", status, result, "errorFound", containers, "startedAt", "finishedAt", codes, huskyciresults) FROM stdin;
\.


--
-- Data for Name: repository; Type: TABLE DATA; Schema: public; Owner: huskyCIUser
--

COPY public.repository (id, "repositoryURL", "repositoryBranch", "createdAt") FROM stdin;
\.


--
-- Data for Name: securityTest; Type: TABLE DATA; Schema: public; Owner: huskyCIUser
--

COPY public."securityTest" (id, name, image, "imageTag", cmd, type, language, "default", "timeOutSeconds") FROM stdin;
\.


--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: huskyCIUser
--

COPY public."user" (id, username, password, salt, iterations, keylen, hashfunction, "newPassword", "confirmNewPassword") FROM stdin;
\.


--
-- Name: accessToken accessToken_huskytoken_key; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='accessToken' 
             AND constraint_name = 'accessToken_huskytoken_key') THEN
        ALTER TABLE ONLY public."accessToken" 
            ADD CONSTRAINT "accessToken_huskytoken_key" UNIQUE (huskytoken);
    END IF;
END $$;


--
-- Name: accessToken accessToken_pkey; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='accessToken' 
             AND constraint_name = 'accessToken_pkey') THEN
        ALTER TABLE ONLY public."accessToken"
            ADD CONSTRAINT "accessToken_pkey" PRIMARY KEY (id);
    END IF;
END $$;


--
-- Name: analysis analysis_RID_key; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='analysis' 
             AND constraint_name = 'analysis_RID_key') THEN
        ALTER TABLE ONLY public.analysis
            ADD CONSTRAINT "analysis_RID_key" UNIQUE ("RID");
    END IF;
END $$;


--
-- Name: analysis analysis_pkey; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='analysis' 
             AND constraint_name = 'analysis_pkey') THEN
        ALTER TABLE ONLY public.analysis
            ADD CONSTRAINT analysis_pkey PRIMARY KEY (id);
    END IF;
END $$;


--
-- Name: repository repository_pkey; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='repository' 
             AND constraint_name = 'repository_pkey') THEN
        ALTER TABLE ONLY public.repository
            ADD CONSTRAINT repository_pkey PRIMARY KEY (id);
    END IF;
END $$;


--
-- Name: repository repository_repositoryURL_key; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='repository' 
             AND constraint_name = 'repository_repositoryURL_key') THEN
        ALTER TABLE ONLY public.repository
            ADD CONSTRAINT "repository_repositoryURL_key" UNIQUE ("repositoryURL");
    END IF;
END $$;


--
-- Name: securityTest securityTest_name_key; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='securityTest' 
             AND constraint_name = 'securityTest_name_key') THEN
        ALTER TABLE ONLY public."securityTest"
            ADD CONSTRAINT "securityTest_name_key" UNIQUE (name);
    END IF;
END $$;


--
-- Name: securityTest securityTest_pkey; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='securityTest' 
             AND constraint_name = 'securityTest_pkey') THEN
        ALTER TABLE ONLY public."securityTest"
            ADD CONSTRAINT "securityTest_pkey" PRIMARY KEY (id);
    END IF;
END $$;


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='user' 
             AND constraint_name = 'user_pkey') THEN
        ALTER TABLE ONLY public."user"
            ADD CONSTRAINT user_pkey PRIMARY KEY (id);
    END IF;
END $$;


--
-- Name: user user_username_key; Type: CONSTRAINT; Schema: public; Owner: huskyCIUser
--

DO $$
BEGIN
    IF NOT EXISTS ( SELECT NULL
            FROM information_schema.table_constraints
           WHERE table_name='user' 
             AND constraint_name = 'user_username_key') THEN
        ALTER TABLE ONLY public."user"
            ADD CONSTRAINT user_username_key UNIQUE (username);
    END IF;
END $$;


--
-- PostgreSQL database dump complete
--
