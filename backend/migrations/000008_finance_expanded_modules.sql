--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8
-- Dumped by pg_dump version 16.8

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

ALTER TABLE IF EXISTS ONLY public.expenses DROP CONSTRAINT IF EXISTS expenses_branch_id_fkey;
ALTER TABLE IF EXISTS ONLY public.expenses DROP CONSTRAINT IF EXISTS expenses_account_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budgets DROP CONSTRAINT IF EXISTS budgets_branch_id_fkey;
ALTER TABLE IF EXISTS ONLY public.budgets DROP CONSTRAINT IF EXISTS budgets_account_id_fkey;
ALTER TABLE IF EXISTS ONLY public.bank_reconciliations DROP CONSTRAINT IF EXISTS bank_reconciliations_matched_journal_id_fkey;
ALTER TABLE IF EXISTS ONLY public.bank_reconciliations DROP CONSTRAINT IF EXISTS bank_reconciliations_branch_id_fkey;
ALTER TABLE IF EXISTS ONLY public.expenses DROP CONSTRAINT IF EXISTS expenses_pkey;
ALTER TABLE IF EXISTS ONLY public.expenses DROP CONSTRAINT IF EXISTS expenses_expense_number_key;
ALTER TABLE IF EXISTS ONLY public.budgets DROP CONSTRAINT IF EXISTS budgets_pkey;
ALTER TABLE IF EXISTS ONLY public.budgets DROP CONSTRAINT IF EXISTS budgets_period_category_key;
ALTER TABLE IF EXISTS ONLY public.bank_reconciliations DROP CONSTRAINT IF EXISTS bank_reconciliations_pkey;
DROP TABLE IF EXISTS public.expenses;
DROP TABLE IF EXISTS public.budgets;
DROP TABLE IF EXISTS public.bank_reconciliations;
SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: bank_reconciliations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bank_reconciliations (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    branch_id uuid,
    bank_name character varying(60) NOT NULL,
    period character varying(20) NOT NULL,
    statement_date date NOT NULL,
    ref_number character varying(100) NOT NULL,
    description text NOT NULL,
    amount numeric(14,2) NOT NULL,
    type character varying(10) NOT NULL,
    match_status character varying(20) DEFAULT 'unmatched'::character varying,
    matched_journal_id uuid,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: budgets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.budgets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    branch_id uuid,
    period character varying(20) NOT NULL,
    category character varying(60) NOT NULL,
    account_id uuid,
    budget_amount numeric(14,2) NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: expenses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.expenses (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    branch_id uuid,
    category character varying(60) NOT NULL,
    account_id uuid,
    expense_number character varying(50) NOT NULL,
    expense_date date DEFAULT CURRENT_DATE NOT NULL,
    amount numeric(14,2) NOT NULL,
    submitted_by character varying(100) NOT NULL,
    description text NOT NULL,
    receipt_url text,
    status character varying(30) DEFAULT 'submitted'::character varying,
    approved_by character varying(100),
    approval_date timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


--
-- Data for Name: bank_reconciliations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.bank_reconciliations (id, branch_id, bank_name, period, statement_date, ref_number, description, amount, type, match_status, matched_journal_id, notes, created_at) FROM stdin;
222e44dd-5327-4b4e-bb37-566427d8b75a	\N	Bank BCA Operasional (1102)	2026-09	2026-09-02	BCA-TX-00101	Setoran Penjualan Kasir POS 01-Sep	4520000.00	credit	matched	\N	\N	2026-09-27 10:58:45.841481+07
02603371-3553-42c1-be83-90e6dcf07628	\N	Bank BCA Operasional (1102)	2026-09	2026-09-02	BCA-TX-00102	Pembayaran Sewa Outlet Bulan Sep	15000000.00	debit	matched	\N	\N	2026-09-27 10:58:45.841481+07
147629bc-29f5-42ec-9092-ac6dc779661f	\N	Bank BCA Operasional (1102)	2026-09	2026-09-03	BCA-TX-00103	Auto-debet Tagihan Listrik PLN	2500000.00	debit	matched	\N	\N	2026-09-27 10:58:45.841481+07
443dca43-54ea-4ce6-aecc-7c3cebf5f231	\N	Bank BCA Operasional (1102)	2026-09	2026-09-04	BCA-TX-00104	Transfer Supplier Biji Kopi PO-0891	6000000.00	debit	matched	\N	\N	2026-09-27 10:58:45.841481+07
9b6e825e-4a9b-4def-9fad-f4ca0c10c745	\N	Bank BCA Operasional (1102)	2026-09	2026-09-10	BCA-TX-00105	QRIS Settlement BCA Merchant H+1	3840000.00	credit	matched	\N	\N	2026-09-27 10:58:45.841481+07
385204d6-0654-4275-9666-11a2f9a073a6	\N	Bank BCA Operasional (1102)	2026-09	2026-09-15	BCA-TX-00106	Biaya Administrasi Bulanan Rekening	25000.00	debit	unmatched	\N	\N	2026-09-27 10:58:45.841481+07
a82a2a0a-4c48-47ba-b7d3-732e50d3f730	\N	Bank BCA Operasional (1102)	2026-09	2026-09-15	BCA-TX-00107	Pendapatan Bunga Jasa Giro	45600.00	credit	unmatched	\N	\N	2026-09-27 10:58:45.841481+07
32a6bb23-3bba-47ae-a890-fb75fe73db14	\N	Bank BCA Operasional (1102)	2026-09	2026-09-20	BCA-TX-00108	Transfer Payroll Gaji Karyawan	28260969.00	debit	matched	\N	\N	2026-09-27 10:58:45.841481+07
48a17046-bc0a-4c55-b92a-aa998bbdfaec	\N	Bank BCA Operasional (1102)	2026-09	2026-09-25	BCA-TX-00109	EDC BCA Settlement Weekend Sales	5120000.00	credit	unmatched	\N	\N	2026-09-27 10:58:45.841481+07
1686db58-8ee8-47ac-b9c6-6b66217337dc	\N	Bank BCA Operasional (1102)	2026-09	2026-09-02	BCA-TX-00101	Setoran Penjualan Kasir POS 01-Sep	4520000.00	credit	matched	\N	\N	2026-09-27 11:01:33.207042+07
fb9e0626-3cd2-4776-971a-30736f7767fe	\N	Bank BCA Operasional (1102)	2026-09	2026-09-02	BCA-TX-00102	Pembayaran Sewa Outlet Bulan Sep	15000000.00	debit	matched	\N	\N	2026-09-27 11:01:33.207042+07
e89330c1-022f-46a4-b05b-540a0a3aa59a	\N	Bank BCA Operasional (1102)	2026-09	2026-09-03	BCA-TX-00103	Auto-debet Tagihan Listrik PLN	2500000.00	debit	matched	\N	\N	2026-09-27 11:01:33.207042+07
a4d15632-27ae-4cf6-9725-91327b9001cf	\N	Bank BCA Operasional (1102)	2026-09	2026-09-04	BCA-TX-00104	Transfer Supplier Biji Kopi PO-0891	6000000.00	debit	matched	\N	\N	2026-09-27 11:01:33.207042+07
1237b500-3cdd-49b6-8074-2faeeb20b0fc	\N	Bank BCA Operasional (1102)	2026-09	2026-09-10	BCA-TX-00105	QRIS Settlement BCA Merchant H+1	3840000.00	credit	matched	\N	\N	2026-09-27 11:01:33.207042+07
b732aa68-27d0-4778-b309-1bfb2b2cc566	\N	Bank BCA Operasional (1102)	2026-09	2026-09-15	BCA-TX-00106	Biaya Administrasi Bulanan Rekening	25000.00	debit	unmatched	\N	\N	2026-09-27 11:01:33.207042+07
34af777f-aef0-4597-bf08-23497a34ff05	\N	Bank BCA Operasional (1102)	2026-09	2026-09-15	BCA-TX-00107	Pendapatan Bunga Jasa Giro	45600.00	credit	unmatched	\N	\N	2026-09-27 11:01:33.207042+07
9c1fc8b7-2396-44fd-86e6-892ee41eef43	\N	Bank BCA Operasional (1102)	2026-09	2026-09-20	BCA-TX-00108	Transfer Payroll Gaji Karyawan	28260969.00	debit	matched	\N	\N	2026-09-27 11:01:33.207042+07
658adbef-e61c-4f72-9403-0f1e3e8b6e9b	\N	Bank BCA Operasional (1102)	2026-09	2026-09-25	BCA-TX-00109	EDC BCA Settlement Weekend Sales	5120000.00	credit	unmatched	\N	\N	2026-09-27 11:01:33.207042+07
\.


--
-- Data for Name: budgets; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.budgets (id, branch_id, period, category, account_id, budget_amount, notes, created_at, updated_at) FROM stdin;
3860aeca-6d61-47a4-9e97-3a8192ccdcf3	\N	2026-09	Bahan Baku & Minuman	ab111111-0000-0000-0000-000000000008	45000000.00	Alokasi bahan kopi, susu, sirup, dan pastry	2026-09-27 10:58:45.841481+07	2026-09-27 10:58:45.841481+07
97dc99dc-b6c6-4cbf-b29d-1a08d74b875c	\N	2026-09	Utilitas Listrik, Air & Gas	ab111111-0000-0000-0000-000000000010	6000000.00	PLN pascabayar, PDAM & tabung gas LPG 12kg	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07
619593bc-659b-4ebc-ade4-ef65098fef88	\N	2026-09	Sewa Tempat & Gedung	ab111111-0000-0000-0000-000000000009	15000000.00	Biaya sewa outlet ruko bulanan	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07
ffd613ff-8537-4be7-9438-fa72e8e87413	\N	2026-09	Promosi & Pemasaran	ab111111-0000-0000-0000-000000000010	4000000.00	Digital ads, voucher loyalty & event komunitas	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07
65a1a26a-4664-4d12-99ae-e31d9eed0f2e	\N	2026-09	Operasional & Kas Kecil	ab111111-0000-0000-0000-000000000010	3000000.00	Kebersihan, ATK, packaging take-away	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07
\.


--
-- Data for Name: expenses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.expenses (id, branch_id, category, account_id, expense_number, expense_date, amount, submitted_by, description, receipt_url, status, approved_by, approval_date, created_at, updated_at, deleted_at) FROM stdin;
ee111111-0000-0000-0000-000000000002	b1111111-0000-0000-0000-000000000001	Bahan Baku	ab111111-0000-0000-0000-000000000008	EXP-2026-09-002	2026-09-12	750000.00	Chef John (Head Chef)	Pembelian darurat es batu kristal & daun mint segar di pasar lokal	https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400	approved	\N	\N	2026-09-27 10:58:45.841481+07	2026-09-27 10:58:45.841481+07	\N
ee111111-0000-0000-0000-000000000001	b1111111-0000-0000-0000-000000000001	Operasional & Kas Kecil	ab111111-0000-0000-0000-000000000010	EXP-2026-09-001	2026-09-05	350000.00	Jane Doe (Cashier SPV)	Pembelian ATK kasir, kertas thermal struk & spidol marker	https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400	reimbursed	\N	\N	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07	\N
ee111111-0000-0000-0000-000000000003	b1111111-0000-0000-0000-000000000001	Pemeliharaan Peralatan	ab111111-0000-0000-0000-000000000010	EXP-2026-09-003	2026-09-18	1250000.00	Sarah Jenkins (Store Manager)	Servis rutin mesin espresso La Marzocco & penggantian gasket seal	https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400	approved	\N	\N	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07	\N
ee111111-0000-0000-0000-000000000004	b1111111-0000-0000-0000-000000000001	Promosi & Pemasaran	ab111111-0000-0000-0000-000000000010	EXP-2026-09-004	2026-09-22	500000.00	Budi Santoso (Marketing)	Instagram & TikTok Ads promo bundling kopi pastry weekend	https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?w=400	submitted	\N	\N	2026-09-27 11:01:33.207042+07	2026-09-27 11:01:33.207042+07	\N
\.


--
-- Name: bank_reconciliations bank_reconciliations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bank_reconciliations
    ADD CONSTRAINT bank_reconciliations_pkey PRIMARY KEY (id);


--
-- Name: budgets budgets_period_category_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budgets
    ADD CONSTRAINT budgets_period_category_key UNIQUE (period, category);


--
-- Name: budgets budgets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budgets
    ADD CONSTRAINT budgets_pkey PRIMARY KEY (id);


--
-- Name: expenses expenses_expense_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_expense_number_key UNIQUE (expense_number);


--
-- Name: expenses expenses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_pkey PRIMARY KEY (id);


--
-- Name: bank_reconciliations bank_reconciliations_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bank_reconciliations
    ADD CONSTRAINT bank_reconciliations_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id) ON DELETE SET NULL;


--
-- Name: bank_reconciliations bank_reconciliations_matched_journal_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bank_reconciliations
    ADD CONSTRAINT bank_reconciliations_matched_journal_id_fkey FOREIGN KEY (matched_journal_id) REFERENCES public.journal_entries(id) ON DELETE SET NULL;


--
-- Name: budgets budgets_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budgets
    ADD CONSTRAINT budgets_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.chart_of_accounts(id) ON DELETE RESTRICT;


--
-- Name: budgets budgets_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.budgets
    ADD CONSTRAINT budgets_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id) ON DELETE SET NULL;


--
-- Name: expenses expenses_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.chart_of_accounts(id) ON DELETE RESTRICT;


--
-- Name: expenses expenses_branch_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.expenses
    ADD CONSTRAINT expenses_branch_id_fkey FOREIGN KEY (branch_id) REFERENCES public.branches(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

