--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Debian 15.4-1.pgdg120+1)
-- Dumped by pg_dump version 15.4 (Debian 15.4-1.pgdg120+1)

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
-- Name: insta_accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.insta_accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.insta_accounts_id_seq OWNER TO postgres;

--
-- Name: insta_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.insta_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.insta_users_id_seq OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: levels; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.levels (
    id bigint NOT NULL,
    level text
);


ALTER TABLE public.levels OWNER TO postgres;

--
-- Name: levels_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.levels_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.levels_id_seq OWNER TO postgres;

--
-- Name: levels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.levels_id_seq OWNED BY public.levels.id;


--
-- Name: positions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.positions (
    id bigint NOT NULL,
    type text NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.positions OWNER TO postgres;

--
-- Name: positions_aliases; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.positions_aliases (
    position_id bigint,
    language text,
    alias text
);


ALTER TABLE public.positions_aliases OWNER TO postgres;

--
-- Name: positions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.positions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.positions_id_seq OWNER TO postgres;

--
-- Name: positions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.positions_id_seq OWNED BY public.positions.id;


--
-- Name: positions_skills; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.positions_skills (
    position_id integer,
    skill_id integer
);


ALTER TABLE public.positions_skills OWNER TO postgres;

--
-- Name: skills; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.skills (
    id bigint NOT NULL,
    type text NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.skills OWNER TO postgres;

--
-- Name: skills_aliases; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.skills_aliases (
    skill_id bigint,
    language text,
    alias text
);


ALTER TABLE public.skills_aliases OWNER TO postgres;

--
-- Name: skills_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.skills_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.skills_id_seq OWNER TO postgres;

--
-- Name: skills_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.skills_id_seq OWNED BY public.skills.id;


--
-- Name: levels id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.levels ALTER COLUMN id SET DEFAULT nextval('public.levels_id_seq'::regclass);


--
-- Name: positions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions ALTER COLUMN id SET DEFAULT nextval('public.positions_id_seq'::regclass);


--
-- Name: skills id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills ALTER COLUMN id SET DEFAULT nextval('public.skills_id_seq'::regclass);


--
-- Data for Name: levels; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.levels (id, level) FROM stdin;
1	middle
2	mid
3	гуру
4	мастер
5	специалист
6	specialist
7	master
8	expert
9	эксперт
10	младший
11	углублённый
12	medium
13	5
14	4
15	3
16	2
17	1
18	5-й
19	4-й
20	3-й
21	2-й
22	1-й
23	первый
24	второй
25	третий
26	четвертый
27	пятый
28	deep
29	trainee
30	graduate
31	выпускник
32	ас
33	начинающий
34	главный
35	cтарший
36	студент
\.


--
-- Data for Name: positions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.positions (id, type, description) FROM stdin;
1	editor	
2	manager	
3	writer	
4	designer	
5	controller	
6	helper	
7	presenter	
8	programmer	
9	operator	
10	artist	
11	engineer	
12	marketer	
13	analyst	
14	officer	
15	entertainer	
16	clerk	
17	teacher	
18	researcher	
19	courier	
20	mathematician	
21	scientist	
22	reviewer	
23	partner	
24	intern	
25	technician	
26	tester	
27	recruiter	
28	advisor	
29	seller	
30	entrepreneur	
31	content creator	
32	linguist	
33	specialist	
34	staff	
35	medic	
36	coach	
37	inspector	
38	reporter	
39	mentor	
40	curator	
41	driver	
42	supplier	
43	lawyer	
44	chemist	
\.


--
-- Data for Name: positions_aliases; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.positions_aliases (position_id, language, alias) FROM stdin;
1	en	editor
2	en	chief
3	en	writer
4	en	designer
3	en	author
2	en	manager
5	en	regulator
6	en	assistant
2	en	administrator
7	en	leading
8	en	developer
9	en	operator
10	en	artist
11	en	electrical engineer
12	en	telemarketer
13	en	economist
2	en	director
14	en	officer
15	en	animator
16	en	secretary
15	en	actor
17	en	art historian
18	en	researcher
19	en	courier
20	en	mathematician
13	en	analyst
21	en	scientist
8	en	programmer
22	en	reviewer
6	en	support
23	en	partner
2	en	coordinator
24	en	pupil
2	en	deputy
25	en	adjuster
2	en	head
11	en	engineer
2	en	devops
11	en	qa-engineer
26	en	autotester
27	en	recruiter
17	en	teacher
3	en	scriptwriter
11	en	software engineer
28	en	consultant
29	en	seller
25	en	fitter
11	en	technologist
2	en	teamlead
28	en	adviser
26	en	tester
30	en	owner
12	en	marketer
25	en	mechanic
4	en	architect
2	en	hr
8	en	coder
3	en	journalist
31	en	blogger
25	en	turner
2	en	leader
2	en	team lead
29	en	promoter
25	en	collector
13	en	methodologist
16	en	receiver
2	en	technical director
25	en	electrician
32	en	translator
2	en	boss
12	en	media buyer
13	en	modeler
12	en	targetologist
3	en	rater
9	en	sentry
33	en	methodist
29	en	retailer
25	en	miller
2	en	integrator
11	en	constructor
25	en	auto mechanic
34	en	collaborator
4	en	visualizer
35	en	doctor
1	en	video editor
2	en	producer
8	en	algorithmist
25	en	welder
2	en	supervisor
12	en	seo
13	en	optimizer
3	en	copywriter
36	en	instructor
37	en	auditor
32	en	linguist
36	en	trainer
4	en	game designer
38	en	correspondent
39	en	mentor
2	en	techlead
40	en	curator
41	en	driver
42	en	supplier
43	en	lawyer
2	en	moderator
30	en	host
34	en	worker
25	en	grinder
11	en	radio engineer
13	en	strategist
25	en	roboticist
10	en	illustrator
25	en	automation
25	en	circuit technician
1	ru	редактор
2	ru	начальник
3	ru	писатель
17	ru	искусствовед
16	ru	секретарь
4	ru	дизайнер
3	ru	автор
2	ru	менеджер
5	ru	регулировщик
6	ru	помощник
2	ru	администратор
7	ru	ведущий
8	ru	разработчик
9	ru	оператор
10	ru	художник
15	ru	артист
11	ru	электромеханик
2	ru	директор
14	ru	офицер
15	ru	аниматор
13	ru	экономист
12	ru	телемаркетолог
18	ru	исследователь
2	ru	управляющий
19	ru	курьер
20	ru	математик
13	ru	аналитик
21	ru	учёный
8	ru	программист
22	ru	ревьюер
6	ru	поддержка
2	ru	координатор
24	ru	ученик
2	ru	заместитель
25	ru	наладчик
2	ru	руководитель
2	ru	глава
11	ru	инженер
11	ru	qa-инженер
26	ru	автотестировщик
27	ru	рекрутер
17	ru	преподаватель
3	ru	сценарист
24	ru	стажёр
11	ru	разработчик программного обеспечения
28	ru	консультант
29	ru	продавец
25	ru	монтажник
11	ru	технолог
2	ru	тимлид
4	ru	проектировщик
28	ru	советник
26	ru	тестировщик
2	ru	режиссёр
30	ru	владелец
12	ru	маркетолог
25	ru	механик
4	ru	архитектор
2	ru	согласующий
6	ru	техподдержка
8	ru	верстальщик
6	ru	ассистент
3	ru	журналист
31	ru	блогер
25	ru	токарь
2	ru	лидер
29	ru	промоутер
25	ru	сборщик
13	ru	методолог
16	ru	приёмщик
2	ru	технический директор
25	ru	электрик
32	ru	переводчик
2	ru	шеф
12	ru	медиабайер
13	ru	моделист
12	ru	таргетолог
3	ru	райтер
9	ru	дежурный
33	ru	методист
29	ru	ритейлер
25	ru	фрезеровщик
2	ru	интегратор
11	ru	конструктор
25	ru	автомеханик
34	ru	сотрудник
4	ru	визуализатор
35	ru	врач
1	ru	видеомонтажёр
2	ru	продюсер
8	ru	алгоритмист
25	ru	сварщик
2	ru	супервайзер
12	ru	сео
13	ru	оптимизатор
3	ru	копирайтер
36	ru	инструктор
37	ru	аудитор
25	ru	установщик
32	ru	лингвист
36	ru	тренер
4	ru	геймдизайнер
38	ru	корреспондент
24	ru	практикант
39	ru	наставник
2	ru	техлид
40	ru	куратор
41	ru	водитель
42	ru	поставщик
43	ru	юрист
2	ru	модератор
30	ru	хозяин
34	ru	рабочий
27	ru	вербовщик
25	ru	электромонтёр
25	ru	шлифовщик
11	ru	радиоинженер
13	ru	стратег
25	ru	робототехник
10	ru	иллюстратор
25	ru	автоматизатор
25	ru	схемотехник
44	ru	фармацевт
44	ru	аптекарь
\.


--
-- Data for Name: positions_skills; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.positions_skills (position_id, skill_id) FROM stdin;
2	561
3	549
3	32
4	579
4	33
4	232
4	171
4	545
4	471
4	472
4	35
4	34
4	594
4	572
4	470
17	481
4	398
4	38
4	43
6	598
4	264
6	583
4	37
4	57
6	594
4	503
4	36
4	40
24	584
4	41
8	251
24	321
3	158
3	299
2	70
2	527
2	197
2	481
2	75
2	545
2	1
2	83
2	50
2	273
2	44
2	539
2	395
2	414
2	46
2	544
24	585
2	454
2	67
2	404
2	299
2	191
2	558
6	445
2	69
2	195
2	72
2	398
2	201
2	54
2	593
2	49
2	193
2	71
2	261
2	456
2	283
2	3
2	56
2	485
2	61
2	65
2	583
2	474
2	82
2	177
2	45
2	24
2	495
2	202
2	51
2	176
2	535
2	560
2	582
2	549
2	501
2	563
6	353
2	76
2	74
2	58
24	503
2	592
2	538
2	68
2	584
2	57
2	536
2	376
2	85
2	350
2	576
2	478
2	348
2	178
2	64
2	53
2	59
2	62
2	188
2	559
2	467
2	282
2	581
2	264
2	580
2	368
2	209
2	48
2	530
2	296
5	397
6	551
6	218
6	523
6	582
6	211
6	414
6	539
2	465
2	393
2	523
6	399
24	456
2	480
2	520
2	301
2	218
2	80
28	488
2	255
2	268
2	553
2	79
2	276
2	213
2	575
28	218
2	256
2	486
2	552
2	591
2	81
2	142
2	244
25	110
2	517
2	367
2	540
2	405
28	534
7	597
7	365
7	311
7	426
7	523
7	85
7	393
7	468
7	527
7	535
7	560
7	226
7	582
7	346
7	488
7	340
7	585
7	350
7	414
7	318
7	599
7	549
7	97
7	481
7	187
7	87
7	478
7	521
7	218
7	573
7	398
7	84
7	545
7	441
7	424
7	321
7	452
7	255
7	88
7	91
7	329
7	238
7	553
7	563
7	442
7	515
7	534
7	378
7	561
7	494
7	504
7	244
7	581
7	103
28	575
7	519
7	51
7	511
7	525
7	83
17	456
7	212
7	538
7	386
7	224
7	486
7	503
7	552
7	456
6	405
7	584
7	171
7	505
7	383
7	86
7	485
7	403
7	524
7	1
7	436
7	517
7	407
7	540
7	225
6	275
8	245
7	404
7	385
7	522
7	373
7	537
7	451
7	299
7	240
7	476
8	127
8	426
8	144
8	468
8	100
8	129
8	226
8	585
8	318
8	250
8	143
8	545
8	334
8	151
8	324
8	441
8	500
8	442
8	515
8	148
8	494
8	326
8	308
8	317
8	122
8	575
8	146
8	525
8	119
8	398
8	109
8	386
8	503
8	149
8	90
8	112
8	472
8	116
8	404
8	445
8	299
8	374
8	311
8	451
8	523
8	488
8	147
8	97
8	429
8	121
8	301
8	98
8	508
8	433
8	96
22	549
8	124
8	91
8	329
6	259
8	107
8	115
8	470
8	103
8	133
8	106
8	261
8	578
8	134
8	102
8	171
8	524
8	336
8	407
6	398
6	201
2	171
8	537
8	597
8	125
8	140
8	94
8	549
8	136
8	131
8	321
8	573
8	501
8	110
8	120
2	180
6	359
6	256
8	111
8	477
8	432
6	487
8	339
8	118
8	51
8	538
8	95
8	150
8	331
8	142
8	244
8	584
8	145
8	307
8	517
8	367
8	594
8	240
8	205
8	340
8	350
8	315
8	520
8	130
8	576
8	99
8	108
8	320
8	218
8	114
8	328
8	211
8	126
6	591
8	232
2	468
8	519
8	138
8	510
8	552
8	92
8	341
24	468
8	139
24	227
8	135
8	444
8	383
8	101
8	113
9	158
9	155
9	170
9	414
9	167
9	165
9	218
9	398
9	110
9	348
9	153
9	162
9	277
9	551
9	83
24	350
9	368
9	156
9	486
9	157
9	471
9	152
9	154
9	160
9	564
9	169
9	544
9	164
9	159
9	540
9	161
9	168
9	166
10	171
10	594
10	172
10	174
10	173
10	232
10	235
10	472
10	329
2	579
24	311
2	488
3	470
8	564
6	545
6	599
8	323
22	114
24	597
6	424
6	400
6	262
6	261
6	540
24	527
2	266
24	6
24	599
24	534
2	304
24	494
2	269
2	294
28	585
2	6
2	91
2	298
2	494
2	292
2	504
28	501
2	281
28	497
2	210
28	538
28	510
28	502
2	272
2	537
28	583
8	325
8	312
28	517
8	335
8	572
28	503
25	566
8	327
8	330
13	554
8	556
13	507
8	460
8	511
8	332
13	598
13	509
8	583
13	510
13	591
8	319
8	338
2	516
2	103
4	522
4	581
2	438
11	585
11	568
11	378
11	317
11	386
11	414
11	1
11	451
11	523
11	377
26	551
11	349
11	551
11	382
11	397
11	583
11	373
11	596
11	549
11	110
11	255
11	432
11	388
11	368
11	357
11	403
11	385
11	365
11	45
26	545
11	348
11	553
11	438
11	389
11	519
26	524
11	375
11	394
2	437
26	552
2	382
30	527
2	406
2	429
12	528
12	577
2	457
26	573
11	450
11	597
11	426
11	420
12	578
11	446
11	509
11	433
11	425
11	573
4	542
11	441
4	280
11	442
11	556
11	447
11	440
11	405
11	54
11	275
11	538
4	535
11	412
11	457
11	552
4	561
11	423
11	102
11	444
11	350
4	386
4	91
11	436
11	409
11	540
11	91
4	534
11	338
4	543
25	573
25	455
25	350
25	454
25	453
25	457
17	597
4	87
4	299
31	559
2	423
2	554
6	556
32	557
32	342
12	559
9	432
34	398
4	583
4	574
12	576
3	580
37	561
36	583
28	539
28	589
39	398
42	592
34	593
25	595
26	598
26	597
8	246
6	462
6	511
8	241
22	525
6	549
6	268
6	432
6	592
24	218
6	299
24	466
24	485
24	414
2	267
24	467
24	299
28	491
28	229
28	494
2	344
28	484
2	2
2	284
28	386
2	279
28	492
28	584
28	536
2	291
2	295
28	540
28	299
2	293
25	541
2	577
13	296
2	508
2	476
2	179
13	506
13	511
2	175
14	561
15	171
15	232
8	421
13	599
1	579
1	398
1	183
1	184
1	503
8	333
1	71
18	579
18	185
18	561
2	203
2	205
2	420
13	503
8	310
2	599
2	509
2	186
2	519
2	572
2	192
4	521
2	198
4	523
8	419
2	77
8	249
8	591
8	174
26	218
2	204
8	536
2	511
26	255
26	599
26	456
2	156
2	503
30	299
12	273
2	187
2	194
2	189
12	594
2	598
12	530
2	594
19	206
20	585
20	503
13	220
13	579
13	535
13	527
13	217
13	226
13	582
13	222
13	488
13	414
13	549
13	481
13	301
13	87
13	545
13	214
13	478
13	208
13	218
13	501
13	398
13	151
13	216
13	432
13	442
13	188
13	467
13	179
13	494
13	561
13	185
13	279
13	213
13	593
13	326
13	244
13	368
13	581
13	103
13	207
13	281
4	218
13	519
13	51
13	592
11	545
13	212
13	538
13	386
13	210
13	224
13	486
13	261
13	552
13	479
13	578
13	539
13	395
13	584
13	221
13	485
13	82
13	540
13	225
13	215
13	227
13	577
13	594
13	373
13	537
13	288
13	299
13	476
21	594
21	536
21	509
10	234
10	230
11	363
11	390
11	525
11	355
10	563
10	40
10	471
10	231
11	505
11	445
8	365
11	299
11	393
11	508
11	452
11	354
8	243
11	359
11	372
11	366
11	82
11	566
11	356
11	23
11	392
11	561
8	255
11	364
11	598
8	238
8	239
11	399
8	354
8	561
11	367
11	522
11	371
11	599
8	242
11	358
8	267
11	362
6	463
8	540
6	503
8	254
8	372
8	248
8	253
6	202
8	486
6	187
6	255
6	401
6	484
6	584
24	83
2	317
2	263
2	386
2	573
25	564
24	91
24	386
2	270
2	297
24	515
2	290
2	477
2	285
2	470
2	339
8	471
2	300
2	278
2	412
28	490
28	500
28	504
2	311
2	274
2	288
28	489
28	486
28	499
28	485
28	487
28	537
28	495
13	512
8	309
13	513
8	598
13	583
8	275
8	554
13	3
2	321
2	515
4	566
28	582
26	398
26	525
2	345
2	343
26	526
2	277
11	381
11	387
11	353
11	396
11	398
11	395
11	472
11	379
11	374
12	398
12	538
11	369
11	593
12	529
11	456
11	524
11	311
11	411
11	582
11	521
11	424
11	88
11	360
11	361
11	511
11	564
11	376
11	503
11	346
11	87
11	218
11	400
11	401
11	581
11	83
11	370
11	402
11	383
11	476
2	447
25	596
2	524
25	533
4	573
4	244
4	582
4	538
4	540
26	519
11	437
11	413
11	428
4	541
11	429
11	576
11	572
11	470
11	478
4	537
11	310
4	539
11	410
11	458
11	463
11	419
11	453
11	415
11	434
11	435
11	554
8	546
11	462
11	416
11	574
11	536
11	454
11	407
11	443
6	481
11	594
11	537
25	554
25	598
25	545
25	599
25	91
25	594
25	572
17	218
2	548
2	550
18	585
25	551
6	555
32	201
2	562
33	75
34	592
34	299
4	575
4	598
4	536
12	545
37	582
37	581
28	588
28	590
23	591
39	321
41	398
42	299
27	594
25	294
25	456
36	414
\.


--
-- Data for Name: skills; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.skills (id, type, description) FROM stdin;
1	communication	
2	projection	
3	supply	
4	self-service	
5	gamer	
6	practice	
7	api	
8	electric transport	
9	hfm	
10	orm	
11	ocr	
12	serm	
13	miner	
14	robotization	
15	clerical work	
16	translated	
17	center for technological assistance	
18	reception	
19	scouting	
20	replacement	
21	support	
22	cyber intelligence	
23	electric power supply	
24	bim	
25	installation	
26	maya	
27	media	
28	warehouse	
29	tuning	
30	claim	
31	center	
32	bpto	
33	economy	
34	sound	
35	narrative	
36	clothing	
37	presentation	
38	motion	
39	landing	
40	concept	
41	photography	
42	fintech	
43	interior	
44	itsm	
45	production	
46	store	
47	product specialist	
48	сommerce	
49	proactivity	
50	scenario	
51	b2b	
52	wholesale	
53	steering	
54	delivery	
55	representative	
56	clinical	
57	advertising	
58	trade	
59	horeca	
60	ppc	
61	ict	
62	french	
63	telemarketing	
64	stm	
65	tagline	
66	bizdev	
67	social	
68	localization	
69	retention	
70	webmaster	
71	social network	
72	media buying	
73	media purchasing	
74	electrical equipment	
75	lesson	
76	order	
77	influence	
78	buyer	
79	soa	
80	credit	
81	lms	
82	storage	
83	service	
84	hire	
85	distribution of goods	
86	youtube	
87	protection	
88	autotest	
89	ethernet	
90	rendering	
91	backend	
92	sharepoint	
93	driver	
94	notification	
95	navigation	
96	php7	
97	andreoid	
98	core	
99	spark	
100	twitter	
101	mssql	
102	sdk	
103	dwh	
104	ar	
105	rn	
106	woocommerce	
107	esb	
108	mysql	
109	agile	
110	cnc	
111	siebel	
112	neural	
113	additional training	
114	code	
115	telegram	
116	magento2	
117	assembler	
118	deep	
119	phyton	
120	goland	
121	playable	
122	unity3d	
123	java rush	
124	vr	
125	seal	
126	waf	
127	bot	
128	microservices	
129	microservice	
130	internet	
131	quant	
132	spring	
133	salute	
134	pega	
135	haxe	
136	moodle	
137	django	
138	blog	
139	apex	
140	billing	
141	parsing	
142	dbms	
143	erlang	
144	nestjs	
145	qml	
146	hyperion	
147	camcorder	
148	olap	
149	ssas	
150	average	
151	sense	
152	plasma	
153	slitting rewinder	
154	sheet bending machine	
155	semiautomatic	
156	video surveillance	
157	contact	
158	automatic	
159	flying	
160	milling	
161	reference	
162	turning	
163	call center	
164	sheet bending	
165	gas	
166	extruder	
167	printer	
168	cutting	
169	bending	
170	saw	
171	3d	
172	cinematic	
173	layout	
174	ue4	
175	executive	
176	petrochemistry	
177	recruiting	
178	metallurgy	
179	retail	
180	bancassurance	
181	cmo	
182	news	
183	sports	
184	musical	
185	governance	
186	spanish	
187	german	
188	monetization	
189	portuguese	
190	swedish	
191	advancement	
192	commercial	
193	community	
194	outsource	
195	conduct	
196	norwegian	
197	project	
198	dutch	
199	retarget	
200	bdm	
201	english	
202	international	
203	tv	
204	programmatic	
205	ad	
206	pedestrian	
207	finance	
208	infographic	
209	segmentation	
210	directorate	
211	analytics	
212	operational	
213	conveyor	
214	financing	
215	investment	
216	evaluation	
217	excel	
218	1c	
219	system	
220	amocrm	
221	market	
222	insurance	
223	css	
224	abs	
225	тv	
226	back	
227	consulting	
228	scrum	
229	record	
230	colorist	
231	props	
232	2d	
233	special effects	
234	special effect	
235	vfx	
236	team leader	
237	leader	
238	bookkeeping	
239	4gl	
240	rail	
241	opencart	
242	1s8	
243	autodesk	
244	bitrix	
245	delphi	
246	buch	
247	yii2	
248	full-stack	
249	kernel	
250	reactjs	
251	poster	
252	webasyst	
253	graphics	
254	opencv	
255	server	
256	fox	
257	team lead	
258	vps	
259	broker	
260	terrasoft creatio	
261	bitrix24	
262	vds	
263	methodology	
264	design	
265	correlation	
266	gas refueling	
267	mechanics	
268	administration	
269	verification	
270	leasing	
271	implementer	
272	realization	
273	commerce	
274	digitalization	
275	trading	
276	assembly	
277	montage	
278	training	
279	simulation	
280	exchequer	
281	sed	
282	pr	
283	brokerage	
284	logistic	
285	administrative	
286	call-center	
287	callcenter	
288	forecasting	
289	stock	
290	share	
291	otm	
292	ifrs	
293	tos	
294	construction	
295	iaas	
296	visualization	
297	mathematical	
298	fit	
299	management	
300	training center	
301	sas	
302	sass	
303	scss	
304	video broadcast	
305	video broadcasts	
306	lecture	
307	abap	
308	angular	
309	сrm	
310	computer	
311	frontend	
312	unreal	
313	unreal engine	
314	ue	
315	vue	
316	crypto	
317	algorithm	
318	js	
319	voip	
320	script	
321	php	
322	.net	
323	net	
324	symfony	
325	rust	
326	rpa	
327	crosschain	
328	wordpress	
329	unity	
330	collaborative	
331	typescript	
332	visual	
333	вackend	
334	blockchain	
335	laravel	
336	nodejs	
337	microsoft azure	
338	azure	
339	go	
340	scala	
341	swift	
342	russian	
343	buying	
344	accounting	
345	merger	
346	video system	
347	video systems	
348	unmanned	
349	stationary	
350	test	
351	tests	
352	telecommunications	
353	telecommunication	
354	plc	
355	ais	
356	content	
357	devops	
358	multimedia	
359	pos	
360	water supply	
361	heating	
362	apple	
363	radio monitoring	
364	maintenance	
365	qt	
366	aviation	
367	database	
368	processing	
369	radio access	
370	dwdm	
371	bs	
372	autonomous	
373	espd	
374	radio-electronic	
375	adjustment	
376	ott	
377	telephony	
378	vmware	
379	audiovisual	
380	audio visualization	
381	astra	
382	pre-processing	
383	fpga	
384	commissioning	
385	emc	
386	automation	
387	epc	
388	pbx	
389	vet	
390	tech lead	
391	tech leads	
392	electromagnetic	
393	virtualization	
394	office equipment	
395	presale	
396	smd	
397	electronics	
398	yandex	
399	terminal	
400	ups	
401	exit	
402	ventilation	
403	szi	
404	saas	
405	deployment	
406	mdm	
407	react	
408	react js	
409	devsecop	
410	nosql	
411	robotics	
412	rnd	
413	implementation	
414	sale	
415	build	
416	smartspeech	
417	smart-speech	
418	smart speech	
419	automate	
420	engineering	
421	kubernetes	
422	streaming	
423	stream	
424	helpdesk	
425	appsec	
426	node	
427	machinization	
428	mlops	
429	hadoop	
430	algorithmization	
431	algorithms	
432	monitoring	
433	flutter	
434	asic	
435	powerscale	
436	autoqa	
437	migration	
438	sre	
439	computer vision	
440	orchestration	
441	ruby	
442	fullstack	
443	sustain	
444	kotlin	
445	front	
446	backup	
447	manage	
448	native js	
449	nativejs	
450	ipp	
451	pcs	
452	auto-test	
453	sdet	
454	e2e	
455	rslogix	
456	testing	
457	observability	
458	kernelcare	
459	forensic	
460	tibco	
461	mining	
462	san	
463	workflow	
464	workflow engine	
465	unix	
466	datascience	
467	wfm	
468	asp	
469	asp .net	
470	video	
471	graphic	
472	ui	
473	certification	
474	optimization	
475	cisco	
476	cybersecurity	
477	architecture	
478	exploration	
479	budgeting	
480	elk	
481	marketing	
482	ims	
483	alm	
484	hcm	
485	real estate	
486	wms	
487	oebs	
488	axapta	
489	cpq	
490	srm	
491	tms	
492	dax	
493	dynamics 365	
494	dynamics	
495	edo	
496	consultation	
497	ngen	
498	edi	
499	midas	
500	mdg	
501	creatio	
502	fo	
503	remote	
504	bpc	
505	water treatment	
506	malware	
507	merchandise	
508	ml	
509	research	
510	tableau	
511	desk	
512	guidewire	
513	soc	
514	business intelligence	
515	golang	
516	recruitment	
517	etl	
518	data science	
519	python	
520	postgresql	
521	hovik	
522	plumbing	
523	acs	
524	ios	
525	android	
526	etrading	
527	analysis	
528	banking	
529	cryptocurrency	
530	b2c	
531	cleaning	
532	special machinery	
533	motor transport	
534	nsi	
535	cloudy	
536	cloud	
537	bi	
538	crm	
539	logistics	
540	erp	
541	lan	
542	сisco	
543	psu	
544	online	
545	web	
546	html	
547	merry-go-round	
548	chapter	
549	development	
550	cluster	
551	pc	
552	sql	
553	linux	
554	network	
555	ceo	
556	ecommerce	
557	turkish	
558	target	
559	smm	
560	advertisement	
561	security	
562	nft	
563	art	
564	machine	
565	full stack	
566	rea	
567	thermal power engineering	
568	cad	
569	freight transport	
570	videos	
571	photo	
572	mobile	
573	java	
574	hardware	
575	oracle	
576	google	
577	wildberry	
578	ozon	
579	ux	
580	instagram	
581	ib	
582	it	
583	sales	
584	sap	
585	c	
586	c#	
587	c++	
588	ibp	
589	4hana	
590	hana	
591	business	
592	document	
593	planning	
594	lead	
595	radio	
596	repair	
597	javascript	
598	software	
599	qa	
\.


--
-- Data for Name: skills_aliases; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.skills_aliases (skill_id, language, alias) FROM stdin;
1	ru	общение
2	ru	проектирование
3	ru	поставка
4	ru	самообслуживание
5	ru	геймер
6	ru	практика
7	en	api
8	ru	электротранспорт
9	en	hfm
10	en	orm
11	en	ocr
12	en	serm
13	ru	майнер
14	ru	роботизация
15	ru	делопроизводство
16	ru	перевод
17	ru	цто
18	ru	приём
19	ru	скаутинг
20	ru	подменный
21	ru	саппорт
22	ru	киберразведка
23	ru	электроснабжение
24	en	bim
25	ru	установка
26	en	maya
27	ru	медиаконтент
28	ru	складской
29	ru	настройка
30	ru	рекламация
31	en	center
32	ru	бпнуть
33	en	economy
34	en	sound
35	ru	нарративный
35	en	narrative
36	ru	одежда
37	ru	презентация
38	ru	моушн
39	ru	лендинг
40	ru	концепт
41	ru	фотография
42	en	fintech
43	ru	интерьер
44	en	itsm
45	ru	продукция
46	ru	лавка
47	ru	продуктолог
48	en	сommerce
49	ru	проактивность
50	ru	сценарий
51	ru	в2в
52	ru	опт
53	ru	пилотирование
54	ru	доставка
55	en	representative
56	ru	клинический
57	ru	рекламный
58	ru	торговля
59	en	horeca
1	ru	коммуникационный
60	en	ppc
61	en	ict
62	ru	французский
63	ru	телемаркетинг
64	ru	стм
65	ru	тэглайна
66	en	bizdev
67	ru	социальный
68	en	localization
69	en	retention
69	ru	удержание
70	ru	вебмастер
71	ru	соцсеть
72	ru	медиазакупка
72	ru	медиа-закупка
73	ru	медиа закупка
74	ru	электрооборудование
75	ru	урок
76	ru	заказ
77	en	influence
78	en	buyer
79	en	soa
80	ru	кредитный
81	ru	сдо
82	ru	схд
83	ru	услуга
84	ru	наём
85	ru	товародвижение
86	en	youtube
87	ru	охрана
87	ru	защита
88	ru	автотест
89	en	ethernet
90	ru	рендеринг
91	ru	бэкэнд
92	en	sharepoint
93	ru	драйвер
94	ru	нотификация
95	ru	навигация
96	en	php7
97	ru	andrоid
98	en	core
99	en	spark
91	ru	бекэнд
100	ru	вконтакте
101	en	mssql
102	en	sdk
103	en	dwh
104	en	ar
105	en	rn
106	en	woocommerce
107	en	esb
108	en	mysql
109	en	agile
110	en	cnc
111	en	siebel
112	ru	нейронный
113	ru	дообучение
114	en	code
115	en	telegram
116	en	magento2
117	en	assembler
118	ru	глубокий
119	en	phyton
120	en	goland
121	en	playable
122	en	unity3d
123	en	java rush
124	en	vr
125	ru	печать
126	en	waf
127	ru	бот
128	ru	микросервисы
129	ru	микросервис
130	en	internet
131	en	quant
132	en	spring
133	ru	салют
134	en	pega
135	en	haxe
136	en	moodle
137	en	django
138	ru	блог
139	en	apex
140	ru	биллинга
140	ru	биллинг
141	ru	парсинга
141	ru	парсинг
142	ru	субд
143	en	erlang
144	en	nestjs
145	en	qml
146	en	hyperion
147	ru	видеокамера
148	en	olap
149	en	ssas
150	ru	среднее
151	en	sense
152	ru	плазменный
152	ru	плазменным
153	ru	бобинорезательный
154	ru	листогиб
155	ru	полуавтоматический
156	ru	видеонаблюдение
157	ru	контактный
158	ru	автоматический
159	ru	летательный
160	ru	фрезерный
161	ru	справочный
162	ru	токарный
163	ru	колл центр
163	ru	колл-центр
164	ru	листогибочный
165	ru	газовый
166	ru	экструдер
167	ru	печатник
168	ru	резка
169	ru	гибочный
170	ru	раскроечный
171	ru	3д
172	en	cinematic
173	en	layout
174	en	ue4
175	ru	исполнительный
176	ru	нефтехимия
177	ru	рекрутинг
178	ru	металлургия
179	ru	ритейл
180	ru	банкострахование
181	en	cmo
182	ru	новости
183	ru	спортивный
184	ru	музыкальный
185	en	governance
186	en	spanish
68	ru	локализация
187	en	german
188	en	monetization
189	en	portuguese
190	en	swedish
191	ru	продвижение
192	en	commercial
193	ru	комьюнити
194	en	outsource
195	ru	ведение
196	en	norwegian
197	ru	проджект
198	en	dutch
62	en	french
199	en	retarget
200	en	bdm
67	en	social
193	en	community
201	en	english
202	en	international
203	en	tv
204	en	programmatic
205	en	ad
206	ru	пеший
207	ru	финансы
208	ru	инфографик
209	ru	сегментация
210	ru	дирекция
211	en	analytics
212	ru	оперативный
213	ru	конвейер
188	ru	монетизация
214	ru	финансирование
215	ru	инвестиция
216	ru	оценка
217	en	excel
218	ru	1с
219	ru	cистемный
220	en	amocrm
215	ru	инвестиции
221	ru	рыночный
222	ru	страховой
222	ru	страховка
223	en	css
224	ru	абс
225	en	тv
226	ru	бэк
227	ru	консалтинг
228	en	scrum
229	ru	запись
230	ru	колорист
231	ru	пропс
232	en	2d
232	ru	2д
233	ru	спецэффекты
234	ru	спецэффект
235	en	vfx
236	ru	тимлидер
238	ru	бухгалтерия
239	en	4gl
240	en	rail
241	en	opencart
242	ru	1с8
243	en	autodesk
244	ru	битрикс
245	en	delphi
246	ru	бух
247	en	yii2
248	ru	фулстек
249	en	kernel
250	en	reactjs
251	en	poster
252	en	webasyst
253	ru	графика
254	en	opencv
255	ru	сервера
114	ru	код
187	ru	немецкий
256	ru	лис
258	en	vps
259	ru	брокер
260	en	terrasoft creatio
261	ru	битрикс24
262	en	vds
263	ru	методология
264	ru	конструирование
265	ru	корреляция
266	ru	газозаправочный
267	ru	механика
268	ru	администрирование
269	ru	верификация
270	ru	лизинг
271	ru	реализатор
272	ru	реализация
273	ru	коммерция
274	ru	цифровизация
58	ru	трейд
275	ru	трейдинг
276	ru	монтажный
277	ru	монтаж
278	ru	тренинг
279	ru	моделирование
280	ru	казначейство
281	ru	сэд
282	en	pr
283	ru	брокерский
284	ru	логистический
285	ru	административный
286	en	call-center
287	en	callcenter
163	en	call center
288	ru	прогнозирование
289	ru	акции
290	ru	акция
291	en	otm
292	ru	мсфо
293	en	tos
294	ru	строительство
295	en	iaas
296	ru	визуализация
297	ru	математический
298	en	fit
299	ru	менеджмент
300	ru	уц
301	en	sas
302	en	sass
303	en	scss
304	ru	видеотрансляция
305	ru	видеотрансляции
306	ru	лекция
306	ru	лекции
307	en	abap
308	en	angular
309	en	сrm
310	en	computer
311	ru	фронтэнд
312	en	unreal
313	en	unreal engine
314	en	ue
315	en	vue
316	en	crypto
317	en	algorithm
318	en	js
319	en	voip
320	en	script
321	ru	рнр
322	en	.net
323	en	net
324	en	symfony
325	en	rust
326	en	rpa
327	en	crosschain
328	en	wordpress
329	en	unity
330	en	collaborative
331	en	typescript
332	en	visual
333	en	вackend
334	en	blockchain
335	en	laravel
336	en	nodejs
171	en	3d
337	en	microsoft azure
338	en	azure
339	en	go
340	en	scala
341	en	swift
342	en	russian
343	en	buying
344	ru	аккаунтинг
345	en	merger
346	ru	видеосистема
347	ru	видеосистемы
348	ru	беспилотный
349	ru	стационарный
350	ru	тест
351	ru	тесты
352	ru	телекоммуникации
353	ru	телекоммуникация
354	en	plc
355	ru	аис
356	ru	содержание
357	ru	devоps
358	ru	мультимедийный
359	en	pos
360	ru	водоснабжение
361	ru	отопление
362	en	apple
363	ru	радиоконтроль
364	ru	техобслуживание
365	en	qt
366	ru	авиационный
367	ru	бд
368	ru	обработка
369	ru	радиодоступ
368	ru	процессинг
370	en	dwdm
371	ru	бс
372	ru	автономный
373	ru	еспд
1	ru	коммуникация
374	ru	радиоэлектронный
375	ru	наладка
310	ru	компьютер
376	ru	отт
377	ru	телефония
378	en	vmware
379	ru	аудиовизуальный
82	ru	хранение
380	ru	аудиовизуализация
381	en	astra
382	ru	предпроцессинг
383	ru	плис
384	ru	пусконаладочный
385	ru	эмс
353	ru	телекоммуникационный
386	ru	автоматизация
387	en	epc
388	ru	атс
389	ru	пто
390	ru	техлида
390	ru	тех-лид
391	ru	тех лид
392	ru	электромагнитный
393	ru	виртуализация
394	ru	оргтехника
395	en	presale
396	en	smd
397	ru	электроника
398	en	yandex
399	ru	терминал
400	ru	ибп
401	ru	выездной
402	ru	вентиляция
403	ru	сзи
404	en	saas
405	ru	развёртывание
406	en	mdm
407	en	react
408	en	react js
409	en	devsecop
410	en	nosql
411	en	robotics
412	en	rnd
413	en	implementation
414	en	sale
415	en	build
416	en	smartspeech
417	en	smart-speech
418	en	smart speech
419	en	automate
420	en	engineering
421	en	kubernetes
422	en	streaming
423	en	stream
424	en	helpdesk
425	en	appsec
426	en	node
393	en	virtualization
427	en	machinization
428	en	mlops
429	en	hadoop
430	ru	алгоритмитизация
431	ru	алгоритмы
317	ru	алгоритм
432	en	monitoring
433	en	flutter
434	en	asic
435	en	powerscale
436	en	autoqa
437	ru	миграция
437	ru	миграции
438	en	sre
439	en	computer vision
440	en	orchestration
441	en	ruby
442	en	fullstack
443	en	sustain
45	en	production
444	en	kotlin
445	en	front
446	en	backup
447	en	manage
448	en	native js
449	en	nativejs
450	en	ipp
451	ru	асутп
452	ru	автотестирование
453	en	sdet
454	en	e2e
455	en	rslogix
350	en	test
351	en	tests
456	en	testing
457	en	observability
458	en	kernelcare
459	en	forensic
460	en	tibco
461	en	mining
275	en	trading
405	en	deployment
462	en	san
463	en	workflow
464	en	workflow engine
465	en	unix
466	en	datascience
467	en	wfm
91	ru	бэкенд
227	ru	консультирование
468	en	asp
469	en	asp .net
470	en	video
471	en	graphic
255	en	server
470	ru	видео
472	en	ui
473	ru	сертификация
474	ru	оптимизация
475	en	cisco
476	ru	кибербезопасность
477	ru	архитектура
478	ru	исследование
479	ru	бюджетирование
480	en	elk
445	ru	фронт
481	ru	сбыт
482	en	ims
483	en	alm
395	ru	пресейл
484	en	hcm
485	ru	недвижимость
486	en	wms
487	en	oebs
488	en	axapta
489	en	cpq
490	en	srm
491	en	tms
492	en	dax
493	en	dynamics 365
494	en	dynamics
495	ru	эдо
496	ru	консультации
496	ru	консультация
497	en	ngen
498	en	edi
499	en	midas
500	en	mdg
501	en	creatio
502	en	fo
503	ru	дистанционный
504	en	bpc
505	ru	водоподготовка
45	ru	производство
506	en	malware
507	en	merchandise
296	en	visualization
508	en	ml
509	en	research
3	en	supply
510	en	tableau
511	en	desk
512	en	guidewire
513	en	soc
514	en	business intelligence
515	en	golang
516	en	recruitment
517	en	etl
518	en	data science
519	en	python
520	en	postgresql
521	ru	овик
522	ru	вк
523	ru	асу
524	en	ios
525	en	android
255	ru	сервер
526	en	etrading
527	ru	анализ
54	en	delivery
528	ru	банкинг
529	ru	криптовалюта
530	en	b2c
531	ru	клининг
532	ru	спецтехника
533	ru	автотранспорт
534	ru	нси
535	ru	облачный
536	ru	облако
537	en	bi
538	en	crm
539	ru	логистика
91	ru	бекенд
244	en	bitrix
540	en	erp
91	en	backend
541	ru	лвс
542	en	сisco
543	ru	бп
544	en	online
202	ru	международный
311	ru	фронтенд
545	en	web
546	en	html
211	ru	аналитика
547	ru	карусель
548	ru	чаптер
423	ru	стрим
549	ru	разработка
550	ru	кластер
551	ru	пк
552	en	sql
553	en	linux
554	en	network
367	en	database
555	en	ceo
556	en	ecommerce
342	ru	русский
201	ru	английский
557	ru	турецкий
503	en	remote
51	en	b2b
558	ru	таргетировать
559	en	smm
560	ru	реклама
561	en	security
481	ru	маркетинг
562	en	nft
563	en	art
563	ru	арт
516	ru	подбор персонал
516	ru	подбор персонала
432	ru	мониторинг
564	ru	станок
110	ru	чпу
565	en	full stack
248	en	full-stack
261	en	bitrix24
566	ru	рэа
567	ru	теплоэнергетика
568	ru	сапр
273	en	commerce
264	en	design
264	ru	дизайн
569	ru	грузовой транспорт
470	ru	видеоролик
570	ru	видеоролики
571	ru	фото
536	en	cloud
572	en	mobile
573	en	java
574	en	hardware
218	en	1c
575	en	oracle
299	en	management
83	ru	обслуживание
576	en	google
577	en	wildberry
578	en	ozon
545	ru	веб
579	en	ux
503	ru	удалённый
580	ru	инстаграм
581	ru	иб
582	ru	ит
561	ru	безопасность
583	en	sales
549	en	development
584	en	sap
539	en	logistics
585	en	c
586	en	c#
587	en	c++
588	en	ibp
589	en	4hana
590	en	hana
591	en	business
321	en	php
311	en	frontend
398	ru	яндекс
299	ru	управление
592	ru	документооборот
593	ru	планирование
594	en	lead
294	ru	сооружение
595	ru	радиофикация
596	ru	ремонт
597	en	javascript
598	en	software
471	ru	графический
456	ru	тестирование
599	en	qa
414	ru	продажа
583	ru	продажи
\.


--
-- Name: insta_accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.insta_accounts_id_seq', 1, false);


--
-- Name: insta_users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.insta_users_id_seq', 1, false);


--
-- Name: levels_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.levels_id_seq', 36, true);


--
-- Name: positions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.positions_id_seq', 44, true);


--
-- Name: skills_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.skills_id_seq', 603, true);


--
-- Name: levels levels_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.levels
    ADD CONSTRAINT levels_pkey PRIMARY KEY (id);


--
-- Name: positions positions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions
    ADD CONSTRAINT positions_pkey PRIMARY KEY (id);


--
-- Name: positions_skills positions_skills_position_id_skill_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions_skills
    ADD CONSTRAINT positions_skills_position_id_skill_id_key UNIQUE (position_id, skill_id);


--
-- Name: skills skills_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills
    ADD CONSTRAINT skills_pkey PRIMARY KEY (id);


--
-- Name: positions_aliases unique_positions_aliases; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions_aliases
    ADD CONSTRAINT unique_positions_aliases UNIQUE (position_id, alias);


--
-- Name: skills_aliases unique_skills_aliases; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills_aliases
    ADD CONSTRAINT unique_skills_aliases UNIQUE (skill_id, alias);


--
-- Name: skills unique_type; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills
    ADD CONSTRAINT unique_type UNIQUE (type);


--
-- Name: positions unique_type_positions; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions
    ADD CONSTRAINT unique_type_positions UNIQUE (type);


--
-- Name: positions_aliases positions_aliases_position_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions_aliases
    ADD CONSTRAINT positions_aliases_position_id_fkey FOREIGN KEY (position_id) REFERENCES public.positions(id);


--
-- Name: positions_skills positions_skills_position_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions_skills
    ADD CONSTRAINT positions_skills_position_id_fkey FOREIGN KEY (position_id) REFERENCES public.positions(id);


--
-- Name: positions_skills positions_skills_skill_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.positions_skills
    ADD CONSTRAINT positions_skills_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES public.skills(id);


--
-- Name: skills_aliases skills_aliases_skill_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.skills_aliases
    ADD CONSTRAINT skills_aliases_skill_id_fkey FOREIGN KEY (skill_id) REFERENCES public.skills(id);


--
-- PostgreSQL database dump complete
--

