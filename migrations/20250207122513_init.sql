-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys=OFF;
CREATE TABLE IF NOT EXISTS "om_props" (
	id        INTEGER not null primary key autoincrement,
	object_id INTEGER not null,
	code      TEXT not null,
	value     TEXT,

	FOREIGN KEY(object_id) REFERENCES "om_objects"(id) on update cascade on delete cascade
);
INSERT INTO om_props VALUES(1,82,'interface','I2C');
INSERT INTO om_props VALUES(2,83,'max_error_value','50');
INSERT INTO om_props VALUES(3,83,'value','24.30');
INSERT INTO om_props VALUES(4,83,'write_graph','true');
INSERT INTO om_props VALUES(5,82,'address','0');
INSERT INTO om_props VALUES(7,83,'max_threshold','27');
INSERT INTO om_props VALUES(8,83,'min_error_value','-50');
INSERT INTO om_props VALUES(9,83,'min_threshold','0');
INSERT INTO om_props VALUES(10,83,'unit','℃');
INSERT INTO om_props VALUES(11,78,'update_interval','60');
INSERT INTO om_props VALUES(14,84,'max_error_value','100');
INSERT INTO om_props VALUES(15,84,'min_error_value','0');
INSERT INTO om_props VALUES(16,84,'unit','%');
INSERT INTO om_props VALUES(17,82,'update_interval','60');
INSERT INTO om_props VALUES(18,84,'value','22.80');
INSERT INTO om_props VALUES(19,84,'write_graph','true');
INSERT INTO om_props VALUES(20,84,'max_threshold','80');
INSERT INTO om_props VALUES(21,84,'min_threshold','0');
INSERT INTO om_props VALUES(22,85,'interface','I2C');
INSERT INTO om_props VALUES(23,86,'unit','℃');
INSERT INTO om_props VALUES(24,85,'update_interval','60');
INSERT INTO om_props VALUES(25,86,'value','24.90');
INSERT INTO om_props VALUES(26,85,'address','0');
INSERT INTO om_props VALUES(27,86,'max_error_value','50');
INSERT INTO om_props VALUES(28,86,'min_threshold','0');
INSERT INTO om_props VALUES(29,86,'write_graph','true');
INSERT INTO om_props VALUES(31,86,'max_threshold','30');
INSERT INTO om_props VALUES(32,86,'min_error_value','-50');
INSERT INTO om_props VALUES(33,87,'unit','%');
INSERT INTO om_props VALUES(34,87,'value','43.90');
INSERT INTO om_props VALUES(36,87,'max_threshold','80');
INSERT INTO om_props VALUES(37,87,'min_error_value','0');
INSERT INTO om_props VALUES(38,89,'update_interval','60');
INSERT INTO om_props VALUES(39,87,'write_graph','true');
INSERT INTO om_props VALUES(41,87,'max_error_value','100');
INSERT INTO om_props VALUES(42,87,'min_threshold','0');
INSERT INTO om_props VALUES(43,91,'update_interval','60');
INSERT INTO om_props VALUES(45,88,'max_error_value','3000');
INSERT INTO om_props VALUES(46,88,'min_error_value','300');
INSERT INTO om_props VALUES(47,88,'unit','ppm');
INSERT INTO om_props VALUES(48,88,'value','683.00');
INSERT INTO om_props VALUES(49,88,'write_graph','true');
INSERT INTO om_props VALUES(51,88,'max_threshold','1800');
INSERT INTO om_props VALUES(52,88,'min_threshold','350');
INSERT INTO om_props VALUES(53,89,'interface','ADC');
INSERT INTO om_props VALUES(54,90,'unit','u');
INSERT INTO om_props VALUES(55,90,'value','1023.00');
INSERT INTO om_props VALUES(56,90,'max_threshold','2700');
INSERT INTO om_props VALUES(57,90,'min_error_value','-50');
INSERT INTO om_props VALUES(58,90,'max_error_value','5000');
INSERT INTO om_props VALUES(59,90,'min_threshold','0');
INSERT INTO om_props VALUES(60,93,'update_interval','60');
INSERT INTO om_props VALUES(61,90,'write_graph','true');
INSERT INTO om_props VALUES(62,89,'address','0');
INSERT INTO om_props VALUES(64,93,'interface','I2C');
INSERT INTO om_props VALUES(65,94,'min_threshold','0');
INSERT INTO om_props VALUES(66,94,'unit','u');
INSERT INTO om_props VALUES(67,94,'value','748.50');
INSERT INTO om_props VALUES(69,94,'max_threshold','2700');
INSERT INTO om_props VALUES(70,94,'min_error_value','-50');
INSERT INTO om_props VALUES(71,97,'update_interval','60');
INSERT INTO om_props VALUES(72,94,'write_graph','true');
INSERT INTO om_props VALUES(73,93,'address','0');
INSERT INTO om_props VALUES(74,94,'max_error_value','5000');
INSERT INTO om_props VALUES(75,95,'max_threshold','30');
INSERT INTO om_props VALUES(76,95,'min_error_value','-50');
INSERT INTO om_props VALUES(77,95,'unit','℃');
INSERT INTO om_props VALUES(78,95,'value','24.40');
INSERT INTO om_props VALUES(79,95,'write_graph','true');
INSERT INTO om_props VALUES(82,95,'max_error_value','50');
INSERT INTO om_props VALUES(83,95,'min_threshold','0');
INSERT INTO om_props VALUES(85,96,'max_threshold','2700');
INSERT INTO om_props VALUES(86,96,'min_threshold','0');
INSERT INTO om_props VALUES(87,96,'unit','u');
INSERT INTO om_props VALUES(88,96,'value','39.60');
INSERT INTO om_props VALUES(89,96,'write_graph','true');
INSERT INTO om_props VALUES(92,96,'max_error_value','5000');
INSERT INTO om_props VALUES(93,96,'min_error_value','-50');
INSERT INTO om_props VALUES(95,32,'address','127.0.0.1');
INSERT INTO om_props VALUES(96,200,'number','0');
INSERT INTO om_props VALUES(98,200,'group','inputs');
INSERT INTO om_props VALUES(100,201,'group','inputs');
INSERT INTO om_props VALUES(101,201,'number','1');
INSERT INTO om_props VALUES(102,202,'number','2');
INSERT INTO om_props VALUES(104,202,'group','inputs');
INSERT INTO om_props VALUES(105,203,'group','inputs');
INSERT INTO om_props VALUES(106,203,'number','3');
INSERT INTO om_props VALUES(108,37,'group','inputs');
INSERT INTO om_props VALUES(109,37,'number','4');
INSERT INTO om_props VALUES(111,38,'group','inputs');
INSERT INTO om_props VALUES(112,38,'number','5');
INSERT INTO om_props VALUES(114,39,'group','inputs');
INSERT INTO om_props VALUES(115,39,'number','6');
INSERT INTO om_props VALUES(117,47,'group','outputs');
INSERT INTO om_props VALUES(118,47,'number','7');
INSERT INTO om_props VALUES(121,204,'group','outputs');
INSERT INTO om_props VALUES(122,204,'number','8');
INSERT INTO om_props VALUES(123,205,'group','outputs');
INSERT INTO om_props VALUES(124,205,'number','9');
INSERT INTO om_props VALUES(126,206,'group','outputs');
INSERT INTO om_props VALUES(127,206,'number','10');
INSERT INTO om_props VALUES(129,207,'group','outputs');
INSERT INTO om_props VALUES(130,207,'number','11');
INSERT INTO om_props VALUES(132,208,'group','outputs');
INSERT INTO om_props VALUES(133,208,'number','12');
INSERT INTO om_props VALUES(136,209,'group','outputs');
INSERT INTO om_props VALUES(137,209,'number','13');
INSERT INTO om_props VALUES(138,210,'group','digital');
INSERT INTO om_props VALUES(139,210,'number','14');
INSERT INTO om_props VALUES(142,40,'group','inputs');
INSERT INTO om_props VALUES(143,40,'number','15');
INSERT INTO om_props VALUES(144,41,'group','inputs');
INSERT INTO om_props VALUES(145,41,'number','16');
INSERT INTO om_props VALUES(147,42,'group','inputs');
INSERT INTO om_props VALUES(148,42,'number','17');
INSERT INTO om_props VALUES(150,43,'number','18');
INSERT INTO om_props VALUES(152,43,'group','inputs');
INSERT INTO om_props VALUES(153,44,'group','inputs');
INSERT INTO om_props VALUES(154,44,'number','19');
INSERT INTO om_props VALUES(156,45,'group','inputs');
INSERT INTO om_props VALUES(157,45,'number','20');
INSERT INTO om_props VALUES(159,46,'group','inputs');
INSERT INTO om_props VALUES(160,46,'number','21');
INSERT INTO om_props VALUES(162,211,'group','outputs');
INSERT INTO om_props VALUES(163,211,'number','22');
INSERT INTO om_props VALUES(165,212,'group','outputs');
INSERT INTO om_props VALUES(166,212,'number','23');
INSERT INTO om_props VALUES(168,213,'group','outputs');
INSERT INTO om_props VALUES(169,213,'number','24');
INSERT INTO om_props VALUES(171,214,'group','outputs');
INSERT INTO om_props VALUES(172,214,'number','25');
INSERT INTO om_props VALUES(174,215,'group','outputs');
INSERT INTO om_props VALUES(175,215,'number','26');
INSERT INTO om_props VALUES(177,216,'group','outputs');
INSERT INTO om_props VALUES(178,216,'number','27');
INSERT INTO om_props VALUES(180,217,'group','outputs');
INSERT INTO om_props VALUES(181,217,'number','28');
INSERT INTO om_props VALUES(183,62,'group','digital');
INSERT INTO om_props VALUES(184,62,'number','29');
INSERT INTO om_props VALUES(186,218,'group','digital');
INSERT INTO om_props VALUES(187,218,'number','30');
INSERT INTO om_props VALUES(189,219,'group','digital');
INSERT INTO om_props VALUES(190,219,'number','31');
INSERT INTO om_props VALUES(193,220,'group','digital');
INSERT INTO om_props VALUES(194,220,'number','32');
INSERT INTO om_props VALUES(196,221,'group','digital');
INSERT INTO om_props VALUES(197,221,'number','33');
INSERT INTO om_props VALUES(198,222,'number','34');
INSERT INTO om_props VALUES(200,222,'group','digital');
INSERT INTO om_props VALUES(202,223,'group','digital');
INSERT INTO om_props VALUES(203,223,'number','35');
INSERT INTO om_props VALUES(204,224,'group','digital');
INSERT INTO om_props VALUES(205,224,'number','36');
INSERT INTO om_props VALUES(207,225,'number','37');
INSERT INTO om_props VALUES(209,225,'group','digital');
INSERT INTO om_props VALUES(210,226,'group','digital');
INSERT INTO om_props VALUES(211,226,'number','38');
INSERT INTO om_props VALUES(214,227,'group','digital');
INSERT INTO om_props VALUES(215,227,'number','39');
INSERT INTO om_props VALUES(216,228,'group','digital');
INSERT INTO om_props VALUES(217,228,'number','40');
INSERT INTO om_props VALUES(219,229,'number','41');
INSERT INTO om_props VALUES(221,229,'group','digital');
INSERT INTO om_props VALUES(222,75,'group','digital');
INSERT INTO om_props VALUES(223,75,'number','42');
INSERT INTO om_props VALUES(226,76,'group','digital');
INSERT INTO om_props VALUES(227,76,'number','43');
INSERT INTO om_props VALUES(228,77,'number','44');
INSERT INTO om_props VALUES(230,77,'group','digital');
INSERT INTO om_props VALUES(231,78,'interface','1W');
INSERT INTO om_props VALUES(233,79,'max_threshold','28');
INSERT INTO om_props VALUES(234,79,'min_error_value','20');
INSERT INTO om_props VALUES(235,79,'min_threshold','10');
INSERT INTO om_props VALUES(237,79,'value','24.10');
INSERT INTO om_props VALUES(238,79,'write_graph','true');
INSERT INTO om_props VALUES(239,78,'address','0');
INSERT INTO om_props VALUES(240,79,'max_error_value','50');
INSERT INTO om_props VALUES(241,79,'unit','℃');
INSERT INTO om_props VALUES(242,97,'interface','I2C');
INSERT INTO om_props VALUES(243,98,'min_error_value','-50');
INSERT INTO om_props VALUES(244,98,'unit','lux');
INSERT INTO om_props VALUES(245,98,'value','219.00');
INSERT INTO om_props VALUES(246,98,'max_error_value','5000');
INSERT INTO om_props VALUES(247,98,'max_threshold','2700');
INSERT INTO om_props VALUES(248,98,'min_threshold','0');
INSERT INTO om_props VALUES(250,98,'write_graph','true');
INSERT INTO om_props VALUES(251,97,'address','0');
INSERT INTO om_props VALUES(253,99,'write_graph','true');
INSERT INTO om_props VALUES(255,99,'max_threshold','2700');
INSERT INTO om_props VALUES(256,99,'unit','u');
INSERT INTO om_props VALUES(257,99,'min_threshold','0');
INSERT INTO om_props VALUES(259,99,'value','748.50');
INSERT INTO om_props VALUES(261,99,'max_error_value','5000');
INSERT INTO om_props VALUES(262,99,'min_error_value','-50');
INSERT INTO om_props VALUES(263,100,'write_graph','true');
INSERT INTO om_props VALUES(264,100,'max_error_value','50');
INSERT INTO om_props VALUES(265,100,'min_threshold','0');
INSERT INTO om_props VALUES(266,100,'unit','℃');
INSERT INTO om_props VALUES(267,100,'value','24.30');
INSERT INTO om_props VALUES(271,100,'max_threshold','30');
INSERT INTO om_props VALUES(272,100,'min_error_value','-50');
INSERT INTO om_props VALUES(275,101,'max_threshold','2700');
INSERT INTO om_props VALUES(276,101,'min_error_value','-50');
INSERT INTO om_props VALUES(277,101,'min_threshold','0');
INSERT INTO om_props VALUES(279,101,'max_error_value','5000');
INSERT INTO om_props VALUES(280,101,'unit','u');
INSERT INTO om_props VALUES(281,101,'value','39.70');
INSERT INTO om_props VALUES(282,101,'write_graph','true');
INSERT INTO om_props VALUES(283,2,'address','10.35.16.11');
INSERT INTO om_props VALUES(285,33,'group','inputs');
INSERT INTO om_props VALUES(286,33,'number','0');
INSERT INTO om_props VALUES(287,34,'group','inputs');
INSERT INTO om_props VALUES(288,34,'number','1');
INSERT INTO om_props VALUES(290,35,'number','2');
INSERT INTO om_props VALUES(292,35,'group','inputs');
INSERT INTO om_props VALUES(293,36,'group','inputs');
INSERT INTO om_props VALUES(294,36,'number','3');
INSERT INTO om_props VALUES(296,230,'group','inputs');
INSERT INTO om_props VALUES(297,230,'number','4');
INSERT INTO om_props VALUES(299,231,'group','inputs');
INSERT INTO om_props VALUES(300,231,'number','5');
INSERT INTO om_props VALUES(303,232,'group','inputs');
INSERT INTO om_props VALUES(304,232,'number','6');
INSERT INTO om_props VALUES(305,233,'group','outputs');
INSERT INTO om_props VALUES(306,233,'number','7');
INSERT INTO om_props VALUES(308,48,'group','outputs');
INSERT INTO om_props VALUES(309,48,'number','8');
INSERT INTO om_props VALUES(311,49,'group','outputs');
INSERT INTO om_props VALUES(312,49,'mode','sw');
INSERT INTO om_props VALUES(313,49,'number','9');
INSERT INTO om_props VALUES(316,50,'group','outputs');
INSERT INTO om_props VALUES(317,50,'mode','sw');
INSERT INTO om_props VALUES(318,50,'number','10');
INSERT INTO om_props VALUES(319,51,'mode','sw');
INSERT INTO om_props VALUES(320,51,'number','11');
INSERT INTO om_props VALUES(322,51,'group','outputs');
INSERT INTO om_props VALUES(323,52,'group','outputs');
INSERT INTO om_props VALUES(324,52,'mode','sw');
INSERT INTO om_props VALUES(325,52,'number','12');
INSERT INTO om_props VALUES(327,53,'group','outputs');
INSERT INTO om_props VALUES(328,53,'mode','pwm');
INSERT INTO om_props VALUES(329,53,'number','13');
INSERT INTO om_props VALUES(331,61,'group','digital');
INSERT INTO om_props VALUES(332,61,'number','14');
INSERT INTO om_props VALUES(334,234,'group','inputs');
INSERT INTO om_props VALUES(335,234,'number','15');
INSERT INTO om_props VALUES(337,235,'group','inputs');
INSERT INTO om_props VALUES(338,235,'number','16');
INSERT INTO om_props VALUES(340,236,'group','inputs');
INSERT INTO om_props VALUES(341,236,'number','17');
INSERT INTO om_props VALUES(343,237,'group','inputs');
INSERT INTO om_props VALUES(344,237,'number','18');
INSERT INTO om_props VALUES(346,238,'group','inputs');
INSERT INTO om_props VALUES(347,238,'number','19');
INSERT INTO om_props VALUES(349,239,'number','20');
INSERT INTO om_props VALUES(351,239,'group','inputs');
INSERT INTO om_props VALUES(352,240,'group','inputs');
INSERT INTO om_props VALUES(353,240,'number','21');
INSERT INTO om_props VALUES(355,54,'group','outputs');
INSERT INTO om_props VALUES(356,54,'mode','sw');
INSERT INTO om_props VALUES(357,54,'number','22');
INSERT INTO om_props VALUES(360,55,'group','outputs');
INSERT INTO om_props VALUES(361,55,'number','23');
INSERT INTO om_props VALUES(362,56,'group','outputs');
INSERT INTO om_props VALUES(363,56,'number','24');
INSERT INTO om_props VALUES(365,57,'group','outputs');
INSERT INTO om_props VALUES(366,57,'number','25');
INSERT INTO om_props VALUES(368,58,'group','outputs');
INSERT INTO om_props VALUES(369,58,'number','26');
INSERT INTO om_props VALUES(371,59,'group','outputs');
INSERT INTO om_props VALUES(372,59,'number','27');
INSERT INTO om_props VALUES(374,60,'group','outputs');
INSERT INTO om_props VALUES(375,60,'number','28');
INSERT INTO om_props VALUES(377,300,'group','digital');
INSERT INTO om_props VALUES(378,300,'number','29');
INSERT INTO om_props VALUES(380,63,'group','digital');
INSERT INTO om_props VALUES(381,63,'type','dsen');
INSERT INTO om_props VALUES(382,63,'number','30');
INSERT INTO om_props VALUES(383,63,'mode','1w');
INSERT INTO om_props VALUES(386,64,'group','digital');
INSERT INTO om_props VALUES(387,64,'number','31');
INSERT INTO om_props VALUES(388,65,'group','digital');
INSERT INTO om_props VALUES(389,65,'type','i2c');
INSERT INTO om_props VALUES(390,65,'number','32');
INSERT INTO om_props VALUES(391,65,'mode','sda');
INSERT INTO om_props VALUES(393,66,'group','digital');
INSERT INTO om_props VALUES(394,66,'number','33');
INSERT INTO om_props VALUES(396,67,'number','34');
INSERT INTO om_props VALUES(397,67,'mode','sda');
INSERT INTO om_props VALUES(399,67,'group','digital');
INSERT INTO om_props VALUES(400,67,'type','i2c');
INSERT INTO om_props VALUES(401,68,'group','digital');
INSERT INTO om_props VALUES(402,68,'number','35');
INSERT INTO om_props VALUES(404,69,'group','digital');
INSERT INTO om_props VALUES(405,69,'number','36');
INSERT INTO om_props VALUES(407,70,'group','digital');
INSERT INTO om_props VALUES(408,70,'type','adc');
INSERT INTO om_props VALUES(409,70,'number','37');
INSERT INTO om_props VALUES(411,71,'group','digital');
INSERT INTO om_props VALUES(412,71,'type','i2c');
INSERT INTO om_props VALUES(413,71,'number','38');
INSERT INTO om_props VALUES(415,72,'group','digital');
INSERT INTO om_props VALUES(416,72,'number','39');
INSERT INTO om_props VALUES(418,73,'group','digital');
INSERT INTO om_props VALUES(419,73,'number','40');
INSERT INTO om_props VALUES(421,74,'group','digital');
INSERT INTO om_props VALUES(422,74,'number','41');
INSERT INTO om_props VALUES(424,301,'group','digital');
INSERT INTO om_props VALUES(425,301,'number','42');
INSERT INTO om_props VALUES(427,302,'group','digital');
INSERT INTO om_props VALUES(428,302,'number','43');
INSERT INTO om_props VALUES(431,303,'group','digital');
INSERT INTO om_props VALUES(432,303,'number','44');
INSERT INTO om_props VALUES(433,102,'name','порт touchon-PWM');
INSERT INTO om_props VALUES(434,102,'number','3');
INSERT INTO om_props VALUES(435,102,'parentPort','64');
INSERT INTO om_props VALUES(436,102,'group','digital');
INSERT INTO om_props VALUES(437,91,'interface','I2C');
INSERT INTO om_props VALUES(439,92,'max_threshold','2700');
INSERT INTO om_props VALUES(440,92,'min_error_value','-50');
INSERT INTO om_props VALUES(442,92,'write_graph','true');
INSERT INTO om_props VALUES(443,91,'address','0');
INSERT INTO om_props VALUES(444,92,'max_error_value','5000');
INSERT INTO om_props VALUES(445,92,'min_threshold','0');
INSERT INTO om_props VALUES(446,92,'unit','lux');
INSERT INTO om_props VALUES(447,92,'value','219.00');
INSERT INTO om_props VALUES(448,83,'value_updated_at','28.12.2024 13:37:20');
INSERT INTO om_props VALUES(451,84,'value_updated_at','28.12.2024 13:37:20');
INSERT INTO om_props VALUES(10464,304,'fallback_sensor_value_id','95');
INSERT INTO om_props VALUES(10465,304,'target_sp','30');
INSERT INTO om_props VALUES(10466,304,'below_tolerance','2.5');
INSERT INTO om_props VALUES(10467,304,'enable','true');
INSERT INTO om_props VALUES(10468,304,'type','simple');
INSERT INTO om_props VALUES(10469,304,'min_sp','23');
INSERT INTO om_props VALUES(10470,304,'sensor_value_ttl','30');
INSERT INTO om_props VALUES(10471,304,'max_sp','66');
INSERT INTO om_props VALUES(10472,304,'above_tolerance','3');
INSERT INTO om_props VALUES(10473,304,'complex_tolerance','5');
INSERT INTO om_props VALUES(17819,79,'value_updated_at','26.09.2024 15:27:58');
INSERT INTO om_props VALUES(17825,86,'value_updated_at','26.09.2024 15:27:58');
INSERT INTO om_props VALUES(17826,87,'value_updated_at','26.09.2024 15:27:58');
INSERT INTO om_props VALUES(17829,88,'value_updated_at','26.09.2024 15:27:58');
INSERT INTO om_props VALUES(17831,90,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17833,92,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17835,94,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17837,95,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17839,96,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17841,99,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17843,100,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17845,101,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(17847,98,'value_updated_at','26.09.2024 15:27:59');
INSERT INTO om_props VALUES(18092,305,'complex_tolerance','');
INSERT INTO om_props VALUES(18093,305,'enable','false');
INSERT INTO om_props VALUES(18094,305,'type','simple');
INSERT INTO om_props VALUES(18095,305,'min_sp','1');
INSERT INTO om_props VALUES(18096,305,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18097,305,'max_sp','1');
INSERT INTO om_props VALUES(18098,305,'below_tolerance','1');
INSERT INTO om_props VALUES(18099,305,'target_sp','1');
INSERT INTO om_props VALUES(18100,305,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18101,305,'above_tolerance','1');
INSERT INTO om_props VALUES(18102,306,'complex_tolerance','');
INSERT INTO om_props VALUES(18103,306,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18104,306,'below_tolerance','1');
INSERT INTO om_props VALUES(18105,306,'type','simple');
INSERT INTO om_props VALUES(18106,306,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18107,306,'above_tolerance','1');
INSERT INTO om_props VALUES(18108,306,'enable','false');
INSERT INTO om_props VALUES(18109,306,'min_sp','1');
INSERT INTO om_props VALUES(18110,306,'target_sp','1');
INSERT INTO om_props VALUES(18111,306,'max_sp','1');
INSERT INTO om_props VALUES(18112,307,'max_sp','1');
INSERT INTO om_props VALUES(18113,307,'min_sp','1');
INSERT INTO om_props VALUES(18114,307,'above_tolerance','1');
INSERT INTO om_props VALUES(18115,307,'type','simple');
INSERT INTO om_props VALUES(18116,307,'complex_tolerance','');
INSERT INTO om_props VALUES(18117,307,'target_sp','1');
INSERT INTO om_props VALUES(18118,307,'enable','false');
INSERT INTO om_props VALUES(18119,307,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18120,307,'below_tolerance','1');
INSERT INTO om_props VALUES(18121,307,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18122,308,'complex_tolerance','');
INSERT INTO om_props VALUES(18123,308,'type','simple');
INSERT INTO om_props VALUES(18124,308,'max_sp','1');
INSERT INTO om_props VALUES(18125,308,'above_tolerance','1');
INSERT INTO om_props VALUES(18126,308,'min_sp','1');
INSERT INTO om_props VALUES(18127,308,'target_sp','1');
INSERT INTO om_props VALUES(18128,308,'enable','false');
INSERT INTO om_props VALUES(18129,308,'below_tolerance','1');
INSERT INTO om_props VALUES(18130,308,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18131,308,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18132,309,'min_sp','1');
INSERT INTO om_props VALUES(18133,309,'enable','false');
INSERT INTO om_props VALUES(18134,309,'type','simple');
INSERT INTO om_props VALUES(18135,309,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18136,309,'target_sp','1');
INSERT INTO om_props VALUES(18137,309,'above_tolerance','1');
INSERT INTO om_props VALUES(18138,309,'complex_tolerance','');
INSERT INTO om_props VALUES(18139,309,'max_sp','1');
INSERT INTO om_props VALUES(18140,309,'below_tolerance','1');
INSERT INTO om_props VALUES(18141,309,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18142,310,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18143,310,'below_tolerance','1');
INSERT INTO om_props VALUES(18144,310,'type','simple');
INSERT INTO om_props VALUES(18145,310,'target_sp','1');
INSERT INTO om_props VALUES(18146,310,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18147,310,'max_sp','1');
INSERT INTO om_props VALUES(18148,310,'above_tolerance','1');
INSERT INTO om_props VALUES(18149,310,'enable','false');
INSERT INTO om_props VALUES(18150,310,'min_sp','1');
INSERT INTO om_props VALUES(18151,310,'complex_tolerance','');
INSERT INTO om_props VALUES(18152,311,'type','simple');
INSERT INTO om_props VALUES(18153,311,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18154,311,'complex_tolerance','');
INSERT INTO om_props VALUES(18155,311,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18156,311,'target_sp','1');
INSERT INTO om_props VALUES(18157,311,'enable','false');
INSERT INTO om_props VALUES(18158,311,'min_sp','1');
INSERT INTO om_props VALUES(18159,311,'max_sp','1');
INSERT INTO om_props VALUES(18160,311,'below_tolerance','1');
INSERT INTO om_props VALUES(18161,311,'above_tolerance','1');
INSERT INTO om_props VALUES(18162,312,'above_tolerance','1');
INSERT INTO om_props VALUES(18163,312,'complex_tolerance','');
INSERT INTO om_props VALUES(18164,312,'type','simple');
INSERT INTO om_props VALUES(18165,312,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18166,312,'enable','false');
INSERT INTO om_props VALUES(18167,312,'min_sp','1');
INSERT INTO om_props VALUES(18168,312,'target_sp','1');
INSERT INTO om_props VALUES(18169,312,'below_tolerance','1');
INSERT INTO om_props VALUES(18170,312,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18171,312,'max_sp','1');
INSERT INTO om_props VALUES(18172,313,'above_tolerance','1');
INSERT INTO om_props VALUES(18173,313,'enable','false');
INSERT INTO om_props VALUES(18174,313,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18175,313,'target_sp','1');
INSERT INTO om_props VALUES(18176,313,'below_tolerance','1');
INSERT INTO om_props VALUES(18177,313,'complex_tolerance','');
INSERT INTO om_props VALUES(18178,313,'type','simple');
INSERT INTO om_props VALUES(18179,313,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18180,313,'min_sp','1');
INSERT INTO om_props VALUES(18181,313,'max_sp','1');
INSERT INTO om_props VALUES(18182,314,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18183,314,'enable','false');
INSERT INTO om_props VALUES(18184,314,'min_sp','1');
INSERT INTO om_props VALUES(18185,314,'complex_tolerance','');
INSERT INTO om_props VALUES(18186,314,'type','simple');
INSERT INTO om_props VALUES(18187,314,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18188,314,'target_sp','1');
INSERT INTO om_props VALUES(18189,314,'above_tolerance','1');
INSERT INTO om_props VALUES(18190,314,'max_sp','1');
INSERT INTO om_props VALUES(18191,314,'below_tolerance','1');
INSERT INTO om_props VALUES(18192,315,'min_sp','1');
INSERT INTO om_props VALUES(18193,315,'below_tolerance','1');
INSERT INTO om_props VALUES(18194,315,'type','simple');
INSERT INTO om_props VALUES(18195,315,'target_sp','1');
INSERT INTO om_props VALUES(18196,315,'above_tolerance','1');
INSERT INTO om_props VALUES(18197,315,'max_sp','1');
INSERT INTO om_props VALUES(18198,315,'complex_tolerance','');
INSERT INTO om_props VALUES(18199,315,'enable','false');
INSERT INTO om_props VALUES(18200,315,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18201,315,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18202,316,'enable','false');
INSERT INTO om_props VALUES(18203,316,'complex_tolerance','');
INSERT INTO om_props VALUES(18204,316,'min_sp','1');
INSERT INTO om_props VALUES(18205,316,'max_sp','1');
INSERT INTO om_props VALUES(18206,316,'below_tolerance','1');
INSERT INTO om_props VALUES(18207,316,'above_tolerance','1');
INSERT INTO om_props VALUES(18208,316,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18209,316,'type','simple');
INSERT INTO om_props VALUES(18210,316,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18211,316,'target_sp','1');
INSERT INTO om_props VALUES(18212,317,'enable','false');
INSERT INTO om_props VALUES(18213,317,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18214,317,'complex_tolerance','');
INSERT INTO om_props VALUES(18215,317,'min_sp','1');
INSERT INTO om_props VALUES(18216,317,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18217,317,'target_sp','1');
INSERT INTO om_props VALUES(18218,317,'below_tolerance','1');
INSERT INTO om_props VALUES(18219,317,'above_tolerance','1');
INSERT INTO om_props VALUES(18220,317,'max_sp','1');
INSERT INTO om_props VALUES(18221,317,'type','simple');
INSERT INTO om_props VALUES(18222,318,'fallback_sensor_value_id','');
INSERT INTO om_props VALUES(18223,318,'target_sp','1');
INSERT INTO om_props VALUES(18224,318,'max_sp','1');
INSERT INTO om_props VALUES(18225,318,'complex_tolerance','');
INSERT INTO om_props VALUES(18226,318,'sensor_value_ttl','120');
INSERT INTO om_props VALUES(18227,318,'min_sp','1');
INSERT INTO om_props VALUES(18228,318,'below_tolerance','1');
INSERT INTO om_props VALUES(18229,318,'enable','false');
INSERT INTO om_props VALUES(18230,318,'type','simple');
INSERT INTO om_props VALUES(18231,318,'above_tolerance','1');
INSERT INTO om_props VALUES(22431,33,'type','i2c');
INSERT INTO om_props VALUES(22432,34,'type','i2c');
INSERT INTO om_props VALUES(22433,35,'type','i2c');
INSERT INTO om_props VALUES(22434,36,'type','i2c');
INSERT INTO om_props VALUES(22435,37,'type','i2c');
INSERT INTO om_props VALUES(22436,38,'type','i2c');
INSERT INTO om_props VALUES(22437,39,'type','i2c');
INSERT INTO om_props VALUES(22438,40,'type','i2c');
INSERT INTO om_props VALUES(22439,41,'type','i2c');
INSERT INTO om_props VALUES(22440,42,'type','i2c');
INSERT INTO om_props VALUES(22441,43,'type','i2c');
INSERT INTO om_props VALUES(22442,44,'type','i2c');
INSERT INTO om_props VALUES(22443,45,'type','i2c');
INSERT INTO om_props VALUES(22444,46,'type','i2c');
INSERT INTO om_props VALUES(22445,47,'type','i2c');
INSERT INTO om_props VALUES(22446,48,'type','i2c');
INSERT INTO om_props VALUES(22447,55,'type','i2c');
INSERT INTO om_props VALUES(22448,56,'type','i2c');
INSERT INTO om_props VALUES(22449,57,'type','i2c');
INSERT INTO om_props VALUES(22450,58,'type','i2c');
INSERT INTO om_props VALUES(22451,59,'type','i2c');
INSERT INTO om_props VALUES(22452,60,'type','i2c');
INSERT INTO om_props VALUES(22453,61,'type','i2c');
INSERT INTO om_props VALUES(22454,62,'type','i2c');
INSERT INTO om_props VALUES(22455,64,'type','i2c');
INSERT INTO om_props VALUES(22456,66,'type','i2c');
INSERT INTO om_props VALUES(22457,68,'type','i2c');
INSERT INTO om_props VALUES(22458,69,'type','i2c');
INSERT INTO om_props VALUES(22459,72,'type','i2c');
INSERT INTO om_props VALUES(22460,73,'type','i2c');
INSERT INTO om_props VALUES(22461,74,'type','i2c');
INSERT INTO om_props VALUES(22462,75,'type','i2c');
INSERT INTO om_props VALUES(22463,76,'type','i2c');
INSERT INTO om_props VALUES(22464,77,'type','i2c');
INSERT INTO om_props VALUES(22465,102,'type','i2c');
INSERT INTO om_props VALUES(22466,200,'type','i2c');
INSERT INTO om_props VALUES(22467,201,'type','i2c');
INSERT INTO om_props VALUES(22468,202,'type','i2c');
INSERT INTO om_props VALUES(22469,203,'type','i2c');
INSERT INTO om_props VALUES(22470,204,'type','i2c');
INSERT INTO om_props VALUES(22471,205,'type','i2c');
INSERT INTO om_props VALUES(22472,206,'type','i2c');
INSERT INTO om_props VALUES(22473,207,'type','i2c');
INSERT INTO om_props VALUES(22474,208,'type','i2c');
INSERT INTO om_props VALUES(22475,209,'type','i2c');
INSERT INTO om_props VALUES(22476,210,'type','i2c');
INSERT INTO om_props VALUES(22477,211,'type','i2c');
INSERT INTO om_props VALUES(22478,212,'type','i2c');
INSERT INTO om_props VALUES(22479,213,'type','i2c');
INSERT INTO om_props VALUES(22480,214,'type','i2c');
INSERT INTO om_props VALUES(22481,215,'type','i2c');
INSERT INTO om_props VALUES(22482,216,'type','i2c');
INSERT INTO om_props VALUES(22483,217,'type','i2c');
INSERT INTO om_props VALUES(22484,218,'type','i2c');
INSERT INTO om_props VALUES(22485,219,'type','i2c');
INSERT INTO om_props VALUES(22486,220,'type','i2c');
INSERT INTO om_props VALUES(22487,221,'type','i2c');
INSERT INTO om_props VALUES(22488,222,'type','i2c');
INSERT INTO om_props VALUES(22489,223,'type','i2c');
INSERT INTO om_props VALUES(22490,224,'type','i2c');
INSERT INTO om_props VALUES(22491,225,'type','i2c');
INSERT INTO om_props VALUES(22492,226,'type','i2c');
INSERT INTO om_props VALUES(22493,227,'type','i2c');
INSERT INTO om_props VALUES(22494,228,'type','i2c');
INSERT INTO om_props VALUES(22495,229,'type','i2c');
INSERT INTO om_props VALUES(22496,230,'type','i2c');
INSERT INTO om_props VALUES(22497,231,'type','i2c');
INSERT INTO om_props VALUES(22498,232,'type','i2c');
INSERT INTO om_props VALUES(22499,233,'type','i2c');
INSERT INTO om_props VALUES(22500,234,'type','i2c');
INSERT INTO om_props VALUES(22501,235,'type','i2c');
INSERT INTO om_props VALUES(22502,236,'type','i2c');
INSERT INTO om_props VALUES(22503,237,'type','i2c');
INSERT INTO om_props VALUES(22504,238,'type','i2c');
INSERT INTO om_props VALUES(22505,239,'type','i2c');
INSERT INTO om_props VALUES(22506,240,'type','i2c');
INSERT INTO om_props VALUES(22507,300,'type','i2c');
INSERT INTO om_props VALUES(22508,301,'type','i2c');
INSERT INTO om_props VALUES(22509,302,'type','i2c');
INSERT INTO om_props VALUES(22510,303,'type','i2c');
INSERT INTO om_props VALUES(23491,380,'interface','DISCRETE');
INSERT INTO om_props VALUES(23492,380,'address','0');
INSERT INTO om_props VALUES(23494,380,'enable','false');
INSERT INTO om_props VALUES(23495,381,'value','0.00');
INSERT INTO om_props VALUES(23497,381,'min_error_value','-1');
INSERT INTO om_props VALUES(23498,381,'min_threshold','-1');
INSERT INTO om_props VALUES(23499,381,'max_threshold','2');
INSERT INTO om_props VALUES(23500,381,'value_updated_at','02.11.2024 08:48:44');
INSERT INTO om_props VALUES(23501,381,'write_graph','true');
INSERT INTO om_props VALUES(23502,381,'unit','');
INSERT INTO om_props VALUES(23503,381,'max_error_value','2');
INSERT INTO om_props VALUES(24401,380,'period','8');
INSERT INTO om_props VALUES(24854,383,'address','33');
INSERT INTO om_props VALUES(24856,49,'type','out');
INSERT INTO om_props VALUES(24857,50,'type','out');
INSERT INTO om_props VALUES(24858,51,'type','out');
INSERT INTO om_props VALUES(24859,52,'type','out');
INSERT INTO om_props VALUES(24860,54,'type','out');
INSERT INTO om_props VALUES(24861,53,'type','out');
INSERT INTO om_props VALUES(25378,385,'tries','3');
INSERT INTO om_props VALUES(25379,385,'connection_string','tcp://10.35.16.7:502');
INSERT INTO om_props VALUES(25380,385,'speed','9600');
INSERT INTO om_props VALUES(25381,385,'data_bits','8');
INSERT INTO om_props VALUES(25382,385,'parity','0');
INSERT INTO om_props VALUES(25383,385,'stop_bits','1');
INSERT INTO om_props VALUES(25384,385,'timeout','3');
INSERT INTO om_props VALUES(25607,388,'address','15');
INSERT INTO om_props VALUES(26558,388,'enable','true');
INSERT INTO om_props VALUES(26613,391,'external_temperature','0');
INSERT INTO om_props VALUES(26614,391,'target_temperature','30');
INSERT INTO om_props VALUES(26615,391,'fan_speed','0');
INSERT INTO om_props VALUES(26616,391,'vertical_slats_mode','1');
INSERT INTO om_props VALUES(26617,391,'eco_mode','false');
INSERT INTO om_props VALUES(26618,391,'ionization','false');
INSERT INTO om_props VALUES(26619,391,'self_cleaning','false');
INSERT INTO om_props VALUES(26620,391,'internal_temperature','639');
INSERT INTO om_props VALUES(26621,391,'update_interval','60');
INSERT INTO om_props VALUES(26622,391,'sounds','false');
INSERT INTO om_props VALUES(26623,391,'operating_mode','3');
INSERT INTO om_props VALUES(26624,391,'address','1');
INSERT INTO om_props VALUES(26625,391,'enable','true');
INSERT INTO om_props VALUES(26626,391,'turbo_mode','false');
INSERT INTO om_props VALUES(26627,391,'sleep_mode','false');
INSERT INTO om_props VALUES(26628,391,'power_status','false');
INSERT INTO om_props VALUES(26629,391,'display_backlight','false');
INSERT INTO om_props VALUES(26630,391,'silent_mode','false');
INSERT INTO om_props VALUES(26631,391,'horizontal_slats_mode','1');
CREATE TABLE IF NOT EXISTS "om_scripts"
(
    id          INTEGER primary key autoincrement,
    code        TEXT not null,
    name        TEXT not null,
    description TEXT not null default '',
    params      TEXT not null default '{}',
    body        TEXT not null,
    name_lowercase text default '' not null
);
INSERT INTO om_scripts VALUES(2,'my_script','Сценарий','Пример сценария',X'7b0a20202274696d656f7574223a207b0a2020202022636f6465223a202274696d656f7574222c0a20202020226e616d65223a2022d0a2d0b0d0b9d0bcd0b0d183d182222c0a20202020226465736372697074696f6e223a2022d0a1d0bad0bed0bbd18cd0bad0be20d0b1d183d0b4d0b5d0bc20d181d0bfd0b0d182d18c222c0a202020202274797065223a2022696e74220a20207d0a7d',replace('package main\nimport s "scripts"\n\n// Аргументы сценария\nvar timeout int // Таймаут\n\nfunc init() {\n	args := s.Args()\n	timeout, _ = args["timeout"].(int)\n}\n\nfunc main() {\n	// Код сценария\n	// ...\n\n	s.Ok("sleep " + s.ToString(timeout))\n	// s.Err("status is err")\n\n	// Вспомогательные функции преобразования типов\n	// --------------------------------------------\n	// s.ToInt("123")\n	// s.ToFloat("1.5")\n	// s.ToBool("true")\n	// s.ToString(123)\n	// s.ToString(false)\n	//\n	// Выполнение другого скрипта\n	// --------------------------\n	// Сигнатура:\n	// func(code string, args map[string]interface{}) (interface{}, string)\n	// Пример:\n	// r, err := s.Exec("my_second_script", map[string]interface{}{"ObjectID":88})\n	// if err != "" {\n	//     s.Err(err)	\n	// }\n	//\n	// Выполнение метода объекта\n	// -------------------------\n	// Выполнение метода объекта с ID=1\n	// r, err := s.ExecObjectMethod(1, "", "", "method", map[string]interface{}{"logicArg":false})\n	// s.Ok(r)\n	//\n	// Выполнение метода объектов категории controller и с типом mega_d\n	// r, err := s.ExecObjectMethod(0, "controller", "mega_d", "method", nil)\n	// s.Err(err)\n}\n','\n',char(10)),'сценарий');
CREATE TABLE IF NOT EXISTS "om_objects"
(
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER references "om_objects" on update cascade on delete set null,
    zone_id   INTEGER,
    category  TEXT    not null,
    type      TEXT    not null,
    internal  bool default false not null,
    name      TEXT    not null,
    status    TEXT default 'N/A' not null,
    tags      JSON default '{}' not null
);
INSERT INTO om_objects VALUES(2,NULL,3,'controller','mega_d',0,'dev4','ON','{"controller":true,"mega_d":true}');
INSERT INTO om_objects VALUES(32,NULL,4,'controller','mega_d',0,'MegaD тестовый объект1','N/A','{"controller":true,"mega_d":true}');
INSERT INTO om_objects VALUES(33,2,0,'port','port_mega_d',1,'Порт 0','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(34,2,0,'port','port_mega_d',1,'Порт 1','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(35,2,0,'port','port_mega_d',1,'Порт 2','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(36,2,0,'port','port_mega_d',1,'Порт 3','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(37,32,0,'port','port_mega_d',1,'Порт 4','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(38,32,0,'port','port_mega_d',1,'Порт 5','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(39,32,0,'port','port_mega_d',1,'Порт 6','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(40,32,0,'port','port_mega_d',1,'Порт 15','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(41,32,0,'port','port_mega_d',1,'Порт 16','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(42,32,0,'port','port_mega_d',1,'Порт 17','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(43,32,0,'port','port_mega_d',1,'Порт 18','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(44,32,0,'port','port_mega_d',1,'Порт 19','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(45,32,0,'port','port_mega_d',1,'Порт 20','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(46,32,0,'port','port_mega_d',1,'Порт 21','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(47,32,0,'port','port_mega_d',1,'Порт 7','OFF','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(48,2,0,'port','port_mega_d',1,'Порт 8','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(49,2,0,'port','port_mega_d',1,'Порт 9','OFF','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(50,2,0,'port','port_mega_d',1,'Порт 10','OFF','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(51,2,0,'port','port_mega_d',1,'Порт 11','OFF','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(52,2,0,'port','port_mega_d',1,'Порт 12','ON','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(53,2,0,'port','port_mega_d',1,'Порт 13','OFF','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(54,2,0,'port','port_mega_d',1,'Порт 22','ON','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(55,2,0,'port','port_mega_d',1,'Порт 23','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(56,2,0,'port','port_mega_d',1,'Порт 24','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(57,2,0,'port','port_mega_d',1,'Порт 25','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(58,2,0,'port','port_mega_d',1,'Порт 26','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(59,2,0,'port','port_mega_d',1,'Порт 27','ON','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(60,2,0,'port','port_mega_d',1,'Порт 28','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(61,2,0,'port','port_mega_d',1,'Порт 14','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(62,32,0,'port','port_mega_d',1,'Порт 29','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(63,2,0,'port','port_mega_d',1,'Порт 30','OFF','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(64,2,0,'port','port_mega_d',1,'Порт 31','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(65,2,0,'port','port_mega_d',1,'Порт 32','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(66,2,0,'port','port_mega_d',1,'Порт 33','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(67,2,0,'port','port_mega_d',1,'Порт 34','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(68,2,0,'port','port_mega_d',1,'Порт 35','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(69,2,0,'port','port_mega_d',1,'Порт 36','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(70,2,0,'port','port_mega_d',1,'Порт 37','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(71,2,0,'port','port_mega_d',1,'Порт 38','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(72,2,0,'port','port_mega_d',1,'Порт 39','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(73,2,0,'port','port_mega_d',1,'Порт 40','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(74,2,0,'port','port_mega_d',1,'Порт 41','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(75,32,0,'port','port_mega_d',1,'Порт 42','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(76,32,0,'port','port_mega_d',1,'Порт 43','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(77,32,0,'port','port_mega_d',1,'Порт 44','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(78,NULL,5,'sensor','ds18b20',0,'test_1wire_sensor','disabled','{"ds18b20":true,"sensor":true,"temperature":true}');
INSERT INTO om_objects VALUES(79,78,0,'sensor_value','temperature',1,'Температура','N/A','{"sensor_value":true,"temperature":true}');
INSERT INTO om_objects VALUES(82,NULL,6,'sensor','htu21d',0,'HTU21D Цифровой датчик температуры и влажности','disabled','{"htu21d":true,"humidity":true,"sensor":true,"temperature":true}');
INSERT INTO om_objects VALUES(83,82,0,'sensor_value','temperature',1,'Температура','N/A','{"sensor_value":true,"temperature":true}');
INSERT INTO om_objects VALUES(84,82,0,'sensor_value','humidity',1,'Влажность','N/A','{"humidity":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(85,NULL,7,'sensor','scd4x',0,'Датчик scd4x','disabled','{"co2":true,"humidity":true,"scd4x":true,"sensor":true,"temperature":true}');
INSERT INTO om_objects VALUES(86,85,0,'sensor_value','temperature',1,'Температура','N/A','{"sensor_value":true,"temperature":true}');
INSERT INTO om_objects VALUES(87,85,0,'sensor_value','humidity',1,'Влажность','N/A','{"humidity":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(88,85,0,'sensor_value','co2',1,'СО2','N/A','{"co2":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(89,NULL,8,'sensor','cs',0,'Датчик тока','disabled','{"cs":true,"current":true,"sensor":true}');
INSERT INTO om_objects VALUES(90,89,0,'sensor_value','current',1,'Ток','N/A','{"current":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(91,NULL,13,'sensor','bh1750',0,'Датчик освещенности bh1750','disabled','{"bh1750":true,"illumination":true,"sensor":true}');
INSERT INTO om_objects VALUES(92,91,0,'sensor_value','illumination',1,'Освещенность','N/A','{"illumination":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(93,NULL,14,'sensor','bme280',0,'Датчик bme280','disabled','{"bme280":true,"humidity":true,"pressure":true,"sensor":true,"temperature":true}');
INSERT INTO om_objects VALUES(94,93,0,'sensor_value','pressure',1,'Давление','N/A','{"pressure":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(95,93,0,'sensor_value','temperature',1,'Температура','N/A','{"sensor_value":true,"temperature":true}');
INSERT INTO om_objects VALUES(96,93,0,'sensor_value','humidity',1,'Влажность','N/A','{"humidity":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(97,NULL,9,'sensor','outdoor',0,'Уличный составной датчик','disabled','{"bh1750":true,"bme280":true,"humidity":true,"illumination":true,"outdoor":true,"pressure":true,"sensor":true,"temperature":true}');
INSERT INTO om_objects VALUES(98,97,0,'sensor_value','illumination',1,'Освещенность','N/A','{"illumination":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(99,97,0,'sensor_value','pressure',1,'Давление','N/A','{"pressure":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(100,97,0,'sensor_value','temperature',1,'Температура','N/A','{"sensor_value":true,"temperature":true}');
INSERT INTO om_objects VALUES(101,97,0,'sensor_value','humidity',1,'Влажность','N/A','{"humidity":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(200,32,0,'port','port_mega_d',1,'Порт 0','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(201,32,0,'port','port_mega_d',1,'Порт 1','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(202,32,0,'port','port_mega_d',1,'Порт 2','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(203,32,0,'port','port_mega_d',1,'Порт 3','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(204,32,0,'port','port_mega_d',1,'Порт 8','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(205,32,0,'port','port_mega_d',1,'Порт 9','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(206,32,0,'port','port_mega_d',1,'Порт 10','ON','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(207,32,0,'port','port_mega_d',1,'Порт 11','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(208,32,0,'port','port_mega_d',1,'Порт 12','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(209,32,0,'port','port_mega_d',1,'Порт 13','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(210,32,0,'port','port_mega_d',1,'Порт 14','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(211,32,0,'port','port_mega_d',1,'Порт 22','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(212,32,0,'port','port_mega_d',1,'Порт 23','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(213,32,0,'port','port_mega_d',1,'Порт 24','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(214,32,0,'port','port_mega_d',1,'Порт 25','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(215,32,0,'port','port_mega_d',1,'Порт 26','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(216,32,0,'port','port_mega_d',1,'Порт 27','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(217,32,0,'port','port_mega_d',1,'Порт 28','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(218,32,0,'port','port_mega_d',1,'Порт 30','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(219,32,0,'port','port_mega_d',1,'Порт 31','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(220,32,0,'port','port_mega_d',1,'Порт 32','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(221,32,0,'port','port_mega_d',1,'Порт 33','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(222,32,0,'port','port_mega_d',1,'Порт 34','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(223,32,0,'port','port_mega_d',1,'Порт 35','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(224,32,0,'port','port_mega_d',1,'Порт 36','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(225,32,0,'port','port_mega_d',1,'Порт 37','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(226,32,0,'port','port_mega_d',1,'Порт 38','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(227,32,0,'port','port_mega_d',1,'Порт 39','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(228,32,0,'port','port_mega_d',1,'Порт 40','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(229,32,0,'port','port_mega_d',1,'Порт 41','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(230,2,0,'port','port_mega_d',1,'Порт 4','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(231,2,0,'port','port_mega_d',1,'Порт 5','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(232,2,0,'port','port_mega_d',1,'Порт 6','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(233,2,0,'port','port_mega_d',1,'Порт 7','off','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(234,2,0,'port','port_mega_d',1,'Порт 15','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(235,2,0,'port','port_mega_d',1,'Порт 16','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(236,2,0,'port','port_mega_d',1,'Порт 17','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(237,2,0,'port','port_mega_d',1,'Порт 18','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(238,2,0,'port','port_mega_d',1,'Порт 19','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(239,2,0,'port','port_mega_d',1,'Порт 20','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(240,2,0,'port','port_mega_d',1,'Порт 21','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(300,2,0,'port','port_mega_d',1,'Порт 29','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(301,2,0,'port','port_mega_d',1,'Порт 42','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(302,2,0,'port','port_mega_d',1,'Порт 43','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(303,2,0,'port','port_mega_d',1,'Порт 44','N/A','{"port":true,"port_mega_d":true}');
INSERT INTO om_objects VALUES(304,83,0,'regulator','regulator',0,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(305,79,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(306,84,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(307,86,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(308,87,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(309,88,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(310,90,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(311,92,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(312,94,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(313,95,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(314,96,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(315,98,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(316,99,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(317,100,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(318,101,0,'regulator','regulator',1,'Регулятор','N/A','{}');
INSERT INTO om_objects VALUES(380,NULL,0,'sensor','motion',0,'Датчик движения','disabled','{"motion":true,"sensor":true}');
INSERT INTO om_objects VALUES(381,380,0,'sensor_value','motion',1,'Движение','N/A','{"motion":true,"sensor_value":true}');
INSERT INTO om_objects VALUES(383,2,0,'generic_input','generic_input',0,'Универсальный вход','N/A','{"generic":true,"input":true}');
INSERT INTO om_objects VALUES(385,NULL,0,'modbus','modbus',0,'Шина Modbus TCP','N/A','{"modbus":true}');
INSERT INTO om_objects VALUES(388,385,0,'modbus','wb_mrm2_mini',0,'WB-MRM2-mini Двухканальный модуль реле','N/A','{"modbus":true,"modbus_device":true,"wb_mrm2_mini":true}');
INSERT INTO om_objects VALUES(391,385,0,'conditioner','onokom/hr_1_mb_b',0,'Кондиционер (Onokom/HR-1-MB-B)','N/A','{"conditioner":true,"gateway":true,"hr_1_mb_b":true,"modbus":true,"modbus_device":true,"onokom":true}');
CREATE TABLE ar_cron_tasks (
    id          integer not null primary key autoincrement,
    name        text not null,
    description text,
    period      text,
    enabled     integer default 1
);
INSERT INTO ar_cron_tasks VALUES(1,'Test htu21d 5s','Проверка датчика htu21d','5s',1);
CREATE TABLE ar_cron_actions (
    id          integer not null primary key autoincrement,
    task_id     integer not null,
    target_type text not null default '',
    target_id   integer not null default 0,
    type        text not null,
    name        text not null,
    args        text not null default '{}',
    qos         integer not null default 0,
    enabled     integer default 1,
    sort        int not null default 0,
    comment     text not null default '',

    FOREIGN KEY(task_id) REFERENCES ar_cron_tasks(id) on update cascade on delete cascade
);
INSERT INTO ar_cron_actions VALUES(1,1,'object',82,'method','check','{}',0,1,0,'');
CREATE TABLE ar_events (
    id          integer not null primary key autoincrement,
    target_type text not null,
    target_id   integer not null,
    event_name  text not null,
    enabled     integer default 1
);
INSERT INTO ar_events VALUES(1,'item',1,'item.on_click',0);
INSERT INTO ar_events VALUES(2,'item',7,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(3,'item',7,'item.on_change_state_off',1);
INSERT INTO ar_events VALUES(4,'item',19,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(5,'item',16,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(6,'item',6,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(7,'item',6,'item.on_change_state_off',1);
INSERT INTO ar_events VALUES(8,'item',230,'item.on_click',0);
INSERT INTO ar_events VALUES(9,'item',1,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(10,'item',1,'item.on_change_state_off',1);
INSERT INTO ar_events VALUES(12,'item',334,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(13,'item',351,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(14,'item',351,'item.on_change_state_off',1);
INSERT INTO ar_events VALUES(15,'item',352,'item.on_change_state_on',1);
INSERT INTO ar_events VALUES(16,'item',352,'item.on_change_state_off',1);
INSERT INTO ar_events VALUES(17,'object',385,'object.sensor.on_check',1);
INSERT INTO ar_events VALUES(18,'object',385,'object.sensor.on_alarm',1);
INSERT INTO ar_events VALUES(19,'object',390,'object.sensor.on_check',1);
INSERT INTO ar_events VALUES(20,'object',390,'object.sensor.on_alarm',1);
INSERT INTO ar_events VALUES(21,'object',395,'object.sensor.on_check',1);
CREATE TABLE ar_event_actions (
    id          integer not null primary key autoincrement,
    event_id    integer not null,
    target_type text not null default '',
    target_id   integer not null default 0,
    type        text not null,
    name        text not null,
    args        text not null default '{}',
    qos         integer not null default 0,
    enabled     integer default 1,
    sort        int not null default 0,
    comment     text not null default '',

    FOREIGN KEY(event_id) REFERENCES ar_events(id) on update cascade on delete cascade
);
INSERT INTO ar_event_actions VALUES(1,1,'object',33,'method','toggle','{"a":123}',0,1,0,'клацаем портом');
INSERT INTO ar_event_actions VALUES(2,1,'object',3,'method','on','{"x":"1", "y":"2"}',0,1,1,'');
INSERT INTO ar_event_actions VALUES(3,1,'object',4,'method','toggle','{}',0,1,2,'');
INSERT INTO ar_event_actions VALUES(4,2,'object',4,'method','on','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(5,3,'object',4,'method','off','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(6,4,'item',19,'method','set_state','{"state":"on"}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(7,5,'object',48,'method','on','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(8,5,'not_matters',0,'delay','','{"duration":"1s"}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(9,5,'object',48,'method','off','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(10,6,'object',233,'method','on','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(11,7,'object',233,'method','off','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(12,8,'object',4,'method','on','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(13,9,'item',4,'method','set_state','{"state":"on"}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(14,10,'item',4,'method','set_state','{"state":"off"}',0,1,0,'my comment');
INSERT INTO ar_event_actions VALUES(18,6,'not_matters',0,'notification','','{"type":"normal", "text": "table light - on"}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(19,7,'not_matters',0,'notification','','{"type":"critical", "text": "table light - off"}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(20,12,'script',2,'method','exec','{"duration":3}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(21,13,'object',82,'method','check','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(22,14,'object',82,'method','check','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(23,15,'object',82,'method','check','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(24,16,'object',82,'method','check','{}',0,1,0,'');
INSERT INTO ar_event_actions VALUES(25,17,'object',33,'method','toggle','null',0,1,0,'клацаем портом');
INSERT INTO ar_event_actions VALUES(26,18,'object',34,'method','on','null',0,1,0,'замыкаем порт по тревоге');
INSERT INTO ar_event_actions VALUES(27,19,'object',33,'method','toggle','null',0,1,0,'клацаем портом');
INSERT INTO ar_event_actions VALUES(28,20,'object',34,'method','on','null',0,1,0,'замыкаем порт по тревоге');
INSERT INTO ar_event_actions VALUES(29,21,'object',33,'method','toggle','null',0,1,0,'клацаем портом');
CREATE TABLE zones (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER not null default 0 references zones on update cascade on delete set default,
    name      TEXT    not null default '',
    style     TEXT    not null default '',
    sort      INTEGER not null default 0
);
INSERT INTO zones VALUES(1,0,'1й этаж','green',1);
INSERT INTO zones VALUES(2,0,'2й этаж','blue',2);
INSERT INTO zones VALUES(3,1,'Кухня','green',1);
INSERT INTO zones VALUES(4,1,'Зал','green',2);
INSERT INTO zones VALUES(5,1,'Туалет','green',3);
INSERT INTO zones VALUES(6,1,'Постирочная','green',4);
INSERT INTO zones VALUES(7,1,'Гостевая спальня','green',5);
INSERT INTO zones VALUES(8,2,'Спальня 1','blue',1);
INSERT INTO zones VALUES(9,2,'Спальня 2','blue',2);
INSERT INTO zones VALUES(10,2,'Детская','blue',3);
INSERT INTO zones VALUES(11,2,'Кабинет','blue',4);
INSERT INTO zones VALUES(12,2,'Туалет','blue',5);
INSERT INTO zones VALUES(13,0,'Мансарда','green',3);
INSERT INTO zones VALUES(14,0,'Беседка','orange',4);
INSERT INTO zones VALUES(15,0,'Гараж','orange',5);
INSERT INTO zones VALUES(16,0,'Двор','orange',6);
CREATE TABLE view_items (
    id            INTEGER not null primary key autoincrement,
    parent_id     INTEGER not null default 0 references view_items on update cascade on delete set default,
    zone_id       INTEGER not null default 0 references zones on update cascade on delete set default,
    type          TEXT    not null default '',
    status        TEXT    not null default '',
    icon          TEXT    not null default '',
    title         TEXT    not null default '',
    sort          INTEGER not null default 0,
    params        JSON    not null default '{}',
    color         TEXT    not null default '',
    auth          TEXT    not null default '',
    description   TEXT    not null default '',
    position_left INTEGER not null default 0,
    scene         INTEGER not null default 0,
    position_top  INTEGER not null default 0,
    enabled       bool    not null default true
);
INSERT INTO view_items VALUES(417,0,3,'group','off','lamp','Освещение',2,'{"is_setted_as_group": true}','','','',0,0,0,1);
INSERT INTO view_items VALUES(418,1,3,'light','off','lamp','Столешница',6,'{}','blue','','',0,0,0,1);
INSERT INTO view_items VALUES(419,1,3,'light','off','lamp','Раковина',5,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(420,1,3,'light','off','intercom','Плита',7,'{}','orange','','',0,0,0,1);
INSERT INTO view_items VALUES(421,1,3,'light','off','lamp','Стол',3,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(422,0,3,'switch','off','intercom','Холодильник',9,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(423,0,3,'switch','on','chaynik','Чайник',6,'{}','red','','',0,0,0,1);
INSERT INTO view_items VALUES(424,0,3,'group','off','shtor','Шторы',7,'{"is_setted_as_group": false}','','','',0,0,0,1);
INSERT INTO view_items VALUES(425,9,3,'curtain','off','shtor-left','Левая штора',8,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(426,9,3,'curtain','off','shtor-right','Правая штора',9,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(427,9,3,'curtain','off','jaluzi','Жалюзи',10,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(428,0,4,'group','on','lamp','Освещение',10,'{"is_setted_as_group": true}','','','',0,0,0,1);
INSERT INTO view_items VALUES(429,13,4,'light','on','lamp','Потолок',11,'{}','orange','','',0,0,0,1);
INSERT INTO view_items VALUES(430,13,4,'light','on','lamp','Люстра',12,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(431,0,4,'conditioner','on','conditioner','Кондиционер',14,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(432,0,4,'switch','off','tv','TV',15,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(433,0,6,'switch','off','washing-machine','Стиралка',16,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(434,0,5,'button','on','mirror','Зеркало',17,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(435,0,7,'light','off','lamp','Свет',18,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(436,0,7,'switch','off','tv','TV',19,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(437,0,8,'light','on','lamp','Свет',1,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(438,0,9,'light','off','lamp','Свет спальня',2,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(439,0,10,'light','on','lamp','Свет детская',3,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(440,0,10,'switch','on','tv','TV',4,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(441,0,11,'group','on','lamp','Освещение',5,'{"is_setted_as_group": true}','','','',0,0,0,1);
INSERT INTO view_items VALUES(442,26,11,'light','on','desk-lamp','Лампа',6,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(443,26,11,'light','on','lamp','Люстра',7,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(444,26,11,'light','on','lamp','Подцветка',8,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(445,0,11,'switch','on','pc','Компьютер',9,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(446,0,11,'conditioner','on','conditioner','Кондиционер',10,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(447,0,11,'group','off','socket','Розетки',11,'{"is_setted_as_group": true}','light_purple','','',0,0,0,1);
INSERT INTO view_items VALUES(448,32,11,'switch','off','socket','Розетка ПК',12,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(449,32,11,'switch','off','socket','Розетка оборудование',13,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(450,0,12,'button','off','mirror','Зеркало',14,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(451,0,13,'conditioner','off','conditioner','Кондиционер',4,'{}','dark_purple','','',0,0,0,1);
INSERT INTO view_items VALUES(452,0,0,'scenario','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(453,0,0,'scenario','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(454,0,0,'scenario','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(455,0,0,'scenario','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(456,0,0,'scenario','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(457,0,0,'scenario','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(458,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(459,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(460,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(461,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(462,0,15,'link','on','boiler','Котел',1,'{"page":"boiler","params":{"boiler_id":1}}','','','',0,0,0,1);
INSERT INTO view_items VALUES(463,0,13,'group','off','curtains','Окно 1',1,'{"is_setted_as_group": false}','','','',0,0,0,1);
INSERT INTO view_items VALUES(464,0,13,'group','off','curtains','Окно 2',5,'{"is_setted_as_group": false}','','','',0,0,0,1);
INSERT INTO view_items VALUES(465,312,13,'curtain','off','curtains','Левая штора',2,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(466,312,13,'curtain','off','curtains','Правая штора',6,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(467,313,13,'curtain','off','curtains','Левая штора',3,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(468,313,13,'curtain','off','curtains','Правая штора',7,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(469,0,13,'switch','on','gal2','Экран проектора',8,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(470,0,13,'switch','on','tv','Проектор',9,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(471,0,13,'light','off','desk-lamp','Светильник',11,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(472,0,13,'light','off','chandelier1','Люстра',10,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(473,0,14,'light','on','chandelier2','Свет',1,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(474,0,14,'switch','on','refregerator','Холодильник',2,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(475,0,14,'button','off','chaynik','Чайник',4,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(476,0,14,'switch','off','kuhnya','Плита',3,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(477,0,15,'light','off','lamp','Свет',2,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(478,0,15,'button','on','gal-up','Открыть',3,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(479,0,15,'button','on','gal-down','Закрыть',5,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(480,0,15,'switch','on','nasos-skv','Насос',4,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(481,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(482,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(483,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
INSERT INTO view_items VALUES(484,0,0,'sensor','','','',0,'{}','','','',0,0,0,1);
CREATE TABLE boilers (
    id                   INTEGER not null primary key autoincrement,
    heating_status       TEXT    not null default '',
    water_status         TEXT    not null default '',
    heating_current_temp REAL    not null default 0,
    heating_optimal_temp REAL    not null default 0,
    water_current_temp   REAL    not null default 0,
    heating_mode         TEXT    not null default '',
    indoor_temp          REAL    not null default 0,
    outdoor_temp         REAL    not null default 0,
    min_threshold        REAL    not null default 0,
    max_threshold        REAL    not null default 0,
    icon                 TEXT    not null default '',
    title                TEXT    not null default '',
    color                TEXT    not null default '',
    auth                 TEXT    not null default ''
);
INSERT INTO boilers VALUES(1,'off','on',75.0,65.0,30.0,'auto',24.0,-14.0,30.0,90.0,'4','','','');
CREATE TABLE boiler_presets (
    id           INTEGER not null primary key autoincrement,
    boiler_id    INTEGER not null references boilers(id) on update cascade on delete cascade,
    temp_out     REAL    not null default 0,
    temp_coolant REAL    not null default 0
);
INSERT INTO boiler_presets VALUES(1,1,-20.0,60.0);
INSERT INTO boiler_presets VALUES(2,1,-10.0,50.0);
INSERT INTO boiler_presets VALUES(3,1,0.0,45.0);
INSERT INTO boiler_presets VALUES(4,1,10.0,30.0);
CREATE TABLE boiler_properties (
    id         INTEGER not null primary key autoincrement,
    boiler_id  INTEGER not null references boilers(id) on update cascade on delete cascade,
    title      TEXT    not null default '',
    image_name TEXT    not null default '',
    value      TEXT    not null default '',
    status     TEXT    not null default ''
);
INSERT INTO boiler_properties VALUES(1,1,'Насос','pump','Вкл.','on');
INSERT INTO boiler_properties VALUES(2,1,'Горелка','fire','Выкл.','off');
INSERT INTO boiler_properties VALUES(3,1,'Горелка ГВС','fire_2','Выкл.','off');
INSERT INTO boiler_properties VALUES(4,1,'Модуляция Горелки','fire_3','23%','on');
INSERT INTO boiler_properties VALUES(5,1,'Давление','pressure','45 бар','on');
CREATE TABLE dimmers (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    name         TEXT    not null default '',
    value        INTEGER not null default 0,
    enabled      bool    not null default false
);
INSERT INTO dimmers VALUES(1,323,'Димер',50,0);
CREATE TABLE local_users (
    id       INTEGER not null primary key autoincrement ,
    name     TEXT    not null default '',
    password TEXT    not null default ''
);
INSERT INTO local_users VALUES(1,'web','12345');
CREATE TABLE menus (
    id        INTEGER not null primary key autoincrement,
    parent_id INTEGER not null default 0 references menus on update cascade on delete cascade,
    page      TEXT    not null default '',
    title     TEXT    not null default '',
    image     TEXT    not null default '',
    sort      INTEGER not null default 0,
    params    JSON    not null default '{}',
    enabled   bool    not null default true
);
INSERT INTO menus VALUES(7,0,'menu','Инженерное','engineer_menu',1,'{}',1);
INSERT INTO menus VALUES(8,7,'boiler','Котёл','boiler_menu',5,'{"boiler_id": 1}',1);
INSERT INTO menus VALUES(23,0,'counters','Счетчики','counters_menu',1,'{}',1);
CREATE TABLE notifications (
    id      INTEGER not null primary key autoincrement,
    type    TEXT    not null default '',
    date    TEXT    not null default '',
    text    TEXT    not null default '',
    is_read bool    not null default false
);
INSERT INTO notifications VALUES(6,'critical','2024-06-05T14:23:12','Критический показатель давления в теплоносителе WBN 6000-24 C. Котел был автоматически выключен! Критический показатель давления в теплоносителе WBN 6000-24 C. Котел был автоматически выключен!',1);
INSERT INTO notifications VALUES(7,'','2024-06-05T13:56:00','Сработал датчик дыма в техпомещении',1);
INSERT INTO notifications VALUES(8,'','2024-06-04T16:23:59','Сработал датчик дыма в техпомещении',1);
INSERT INTO notifications VALUES(9,'critical','2024-06-03T12:23:30','Критический показатель давления в теплоносителе WBN 6000-24 C. Котел был автоматически выключен! Критический показатель давления в теплоносителе WBN 6000-24 C. Котел был автоматически выключен! Критический показатель давления в теплоносителе WBN 6000-24 C. Котел был автоматически выключен! Критический показатель давления в теплоносителе WBN 6000-24 C. Котел был автоматически выключен!',1);
INSERT INTO notifications VALUES(10,'','2024-08-26T12:00:34','Сработал датчик дыма в техпомещении',1);
INSERT INTO notifications VALUES(11,'','2024-08-26T13:00:42','Сработал датчик дыма в техпомещении',1);
INSERT INTO notifications VALUES(12,'','2024-08-27T07:44:54','Сработал датчик дыма в техпомещении',1);
CREATE TABLE scenarios (
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    type         TEXT    not null default '',
    description  TEXT    not null default '',
    icon         TEXT    not null default '',
    title        TEXT    not null default '',
    sort         INTEGER not null default 0,
    color        TEXT    not null default '',
    auth         TEXT    not null default '',
    enabled      bool    not null default false
);
INSERT INTO scenarios VALUES(1,226,'switch','','door','Я ухожу',1,'orange','',1);
INSERT INTO scenarios VALUES(2,227,'switch','','moon','Ночной режим',2,'dark_purple','',1);
INSERT INTO scenarios VALUES(3,228,'switch','','sun','Доброе утро',3,'yellow','',1);
INSERT INTO scenarios VALUES(4,229,'switch','','binoculars','Камеры',4,'light_blue','',1);
INSERT INTO scenarios VALUES(5,230,'button','','kino','Кино',5,'light_purple','',1);
INSERT INTO scenarios VALUES(6,231,'switch','','book','Чтение',6,'light_green','',1);
CREATE TABLE sensors (
    id            INTEGER not null primary key autoincrement,
    view_item_id  INTEGER not null unique references view_items on update cascade on delete cascade,
    object_id     INTEGER not null unique references om_objects on update cascade on delete cascade,
    type          TEXT    not null default '',
    min_threshold REAL    not null default 0,
    max_threshold REAL    not null default 0,
    adjustment    BOOL    not_null default false
);
CREATE TABLE temp_presets (
    id      INTEGER not null primary key autoincrement,
    zone_id INTEGER not null unique references zones on update cascade on delete cascade,
    normal  REAL    not null default 0,
    night   REAL    not null default 0,
    eco     REAL    not null default 0,
    sort    INTEGER not null default 0
);
CREATE TABLE users (
    id            INTEGER  not null primary key autoincrement,
    login         TEXT     not null default '',
    password      TEXT     not null default '',
    role          TEXT     not null default '',
    send_push     bool     not null default true,
    refresh_token TEXT     not null default '',
    token_expired datetime not null default '',
    device_id     INTEGER  not null default 0,
    device_type   TEXT     not null default '',
    device_token  TEXT     not null default ''
, comment TEXT default '' not null);
INSERT INTO users VALUES(43,'','','',1,'416189902aa5702024a5c7400d7ae75c4ed8cd38be9c6a4e3a231e4ecab2b612','2024-09-18 14:31:40.653786 +0700 +07 m=+2592004.894059501',12345,'ios','6970554df6a6320cdb9e6ca3fea764210022c3b08e6718a13a85f3c08854e3b7','');
INSERT INTO users VALUES(44,'','','',1,'','',88691,'android','12345','');
INSERT INTO users VALUES(45,'web','827ccb0eea8a706c4c34a16891f84e7b','',1,'cf2fc1870eca530ab4b2d07ff52ea43b1f0b17ed3d05186e3dead7723ded02ba','2025-02-12 14:46:23.447303183+03:00',10,'android','c7v8GRpOShKdJ3FOINZDsv:APA91bGJ6Phxsqz9bG1tm1Qlv_Hqen2zmMGgNWS2CK9UMc5R1hNDITeS6JfalIUo2sPXNJoozcE00q83xfq-55YDBTq0sUQz1VRfNHE1_ng0KGmQ6cKtkxPZ4zBlPmkLZAzFn5njzU8-','');
INSERT INTO users VALUES(46,'','','',1,'','',66666,'','','');
INSERT INTO users VALUES(47,'','','',1,'','',77777,'','','');
INSERT INTO users VALUES(48,'','','',1,'','',47050,'','','');
INSERT INTO users VALUES(49,'','','',1,'','',111,'','','');
INSERT INTO users VALUES(50,'','','',1,'','',222,'','','');
INSERT INTO users VALUES(51,'','','',1,'','',333,'','','');
INSERT INTO users VALUES(52,'','','',1,'','',444,'','','');
INSERT INTO users VALUES(53,'','','',1,'','',555,'','','');
INSERT INTO users VALUES(54,'','','',1,'','',65808,'','','');
INSERT INTO users VALUES(55,'','','',1,'','',57484,'','','');
INSERT INTO users VALUES(56,'osh88','827ccb0eea8a706c4c34a16891f84e7b','',1,'bf50645f18d07396281071dfa1680ab0a0dbf769655e27f4a8ce060025bc23fd','2024-11-22T15:15:11.67567+03:00',159206,'android','c7v8GRpOShKdJ3FOINZDsv:APA91bGJ6Phxsqz9bG1tm1Qlv_Hqen2zmMGgNWS2CK9UMc5R1hNDITeS6JfalIUo2sPXNJoozcE00q83xfq-55YDBTq0sUQz1VRfNHE1_ng0KGmQ6cKtkxPZ4zBlPmkLZAzFn5njzU8-','TouchOn App');
INSERT INTO users VALUES(57,'ochen_dobriy','827ccb0eea8a706c4c34a16891f84e7b','',1,'80a76ab0fcad7c0b7bf674bde8482e84d917ee81e787b21b87784fa9674c25c0','2025-02-12 14:37:50.680789715+03:00',860742,'android','dKC7YI46TKmfYHfffyQPQb:APA91bG08NaLdyVhwR4Zpow2yjW8x5poSYw073XdAQFxLaOIXj5Xb5Wd9ZHsP6rR_Fr__EGTI5DZwUB9sYg32ikralh-MFdBOh_gxqybyqvG_xjritKzfzk','');
INSERT INTO users VALUES(58,'test','827ccb0eea8a706c4c34a16891f84e7b','',1,'0f4edd81d924574710c8128644cfc65ae2cb55b733200ab7b037b93f2eb6b890','2024-11-17 10:36:49.891327+03:00',38376,'ios','','');
INSERT INTO users VALUES(59,'osh88','827ccb0eea8a706c4c34a16891f84e7b','',1,'e900f830c54f54547a88600c21aca3a2ee5618e55ace1fa2950167f4c14c7f12','2024-11-22 18:49:28.555979+03:00',136584,'android','','TouchOn debug App');
INSERT INTO users VALUES(60,'vitaliy_lq','827ccb0eea8a706c4c34a16891f84e7b','',1,'','',457377,'android','','');
INSERT INTO users VALUES(61,'vitaliy_lq','827ccb0eea8a706c4c34a16891f84e7b','',1,'','',24374,'ios','','');
INSERT INTO users VALUES(62,'yp','827ccb0eea8a706c4c34a16891f84e7b','',1,'','',64147,'','','');
INSERT INTO users VALUES(63,'vd','827ccb0eea8a706c4c34a16891f84e7b','',1,'','',17033,'','','');
CREATE TABLE IF NOT EXISTS "curtain_params"
(
    id           INTEGER         not null primary key autoincrement,
    view_item_id INTEGER         not null unique references view_items on update cascade on delete cascade,
    type         TEXT default '' not null,
    control_type TEXT default '' not null,
    open_percent REAL default 0
);
INSERT INTO curtain_params VALUES(1,10,'default','rs485',NULL);
INSERT INTO curtain_params VALUES(2,11,'default','port',NULL);
INSERT INTO curtain_params VALUES(3,12,'blinds','phase',NULL);
INSERT INTO curtain_params VALUES(4,314,'default','rs485',30.0);
INSERT INTO curtain_params VALUES(5,315,'default','rs485',20.0);
INSERT INTO curtain_params VALUES(6,316,'default','rs485',100.0);
INSERT INTO curtain_params VALUES(7,317,'default','rs485',100.0);
CREATE TABLE IF NOT EXISTS "conditioner_params"
(
    id                    INTEGER           not null primary key autoincrement,
    view_item_id          INTEGER           not null unique references view_items on update cascade on delete cascade,
    inside_temp           REAL default 0    not null,
    outside_temp          REAL default 0,
    current_temp          REAL default 0    not null,
    optimal_temp          REAL default 0    not null,
    min_threshold         REAL default 0    not null,
    max_threshold         REAL default 0    not null,
    silent_mode           bool default false,
    eco_mode              bool default false,
    turbo_mode            bool default false,
    sleep_mode            bool default false,
    fan_speeds            JSON default '[]' not null,
    fan_speed             TEXT default ''   not null,
    vertical_directions   JSON default '[]',
    vertical_direction    TEXT default '',
    horizontal_directions JSON default '[]',
    horizontal_direction  TEXT default '',
    operating_modes       JSON default '[]',
    operating_mode        TEXT default '',
    ionisation            bool default false,
    self_cleaning         bool default false,
    anti_mold             bool default false,
    sound                 bool default false,
    on_duty_heating       bool default false,
    soft_top              bool default false
);
INSERT INTO conditioner_params VALUES(1,16,24.0,10.0,23.0,27.0,17.0,30.0,NULL,1,0,0,'["auto", "first", "second", "third", "fourth", "fifth"]','second','["auto", "swing", "first_position", "second_position", "third_position", "fourth_position", "fifth_position", "sixth_position", "seventh_position"]','swing','["auto", "swing", "first_position", "second_position", "third_position", "fourth_position", "fifth_position", "soft_top"]','third_position','["auto", "cooling", "heating", "dehumidification", "ventilation"]','heating',NULL,1,1,1,1,1);
INSERT INTO conditioner_params VALUES(2,31,0.0,-2.0,18.0,23.0,17.0,30.0,0,0,1,0,'["auto", "first", "second", "third", "fourth", "fifth"]','fifth','["auto", "swing", "first_position", "second_position", "third_position", "fourth_position", "fifth_position", "sixth_position", "seventh_position"]','fifth_position','[]','','["auto", "cooling", "heating", "dehumidification", "ventilation"]','heating',0,0,1,1,0,0);
INSERT INTO conditioner_params VALUES(3,79,1.0,1.0,1.0,1.0,17.0,30.0,0,0,0,0,'[]','фывфыв','[]','фывфыв','[]','','[]','фывфыв',0,0,0,0,0,0);
CREATE TABLE IF NOT EXISTS "counter_daily_history"
(
    id         INTEGER         not null primary key autoincrement,
    counter_id INTEGER         not null references counters on update cascade on delete cascade,
    datetime   TEXT default '' not null,
    value      REAL default 0
);
CREATE TABLE IF NOT EXISTS "counter_monthly_history"
(
    id         INTEGER         not null primary key autoincrement,
    counter_id INTEGER         not null references counters on update cascade on delete cascade,
    datetime   TEXT default '' not null,
    value      REAL default 0,
    unique (counter_id, datetime)
);
CREATE TABLE IF NOT EXISTS "device_daily_history"
(
    id           INTEGER         not null primary key autoincrement,
    view_item_id INTEGER         not null references view_items on update cascade on delete cascade,
    datetime     TEXT default '' not null,
    value        REAL default 0,
    unique (view_item_id, datetime)
);
CREATE TABLE IF NOT EXISTS "device_hourly_history"
(
    id           INTEGER         not null primary key autoincrement,
    view_item_id INTEGER         not null references view_items on update cascade on delete cascade,
    datetime     TEXT default '' not null,
    value        REAL default 0,
    unique (view_item_id, datetime)
);
CREATE TABLE IF NOT EXISTS "counters"
(
    id             INTEGER              not null primary key autoincrement,
    name           TEXT    default ''   not null,
    type           TEXT    default ''   not null,
    unit           TEXT    default '',
    today_value    REAL    default 0    not null,
    week_value     REAL    default 0    not null,
    month_value    REAL    default 0    not null,
    year_value     REAL    default 0    not null,
    price_for_unit REAL    default 0,
    impulse        REAL    default 0    not null,
    sort           INTEGER default 0    not null,
    enabled        bool    default true not null
);
INSERT INTO counters VALUES(1,'Газ','gas','м³',5.299999999999999823,60.60000000000000142,240.3000000000000113,2880.09999999999991,24.35999999999999944,834.2000000000000454,4,1);
INSERT INTO counters VALUES(2,'Холодная вода','cold_water','м³',10.5,92.4000000000000056,368.1000000000000227,4416.0,50.92999999999999972,9431.29999999999927,1,1);
INSERT INTO counters VALUES(3,'Электрическая энергия','electricity','кВт/ч',23.30000000000000071,178.5999999999999944,712.2000000000000454,8544.600000000000363,5.660000000000000142,608.7000000000000454,3,1);
INSERT INTO counters VALUES(4,'Горячая вода','hot_water','м³',9.199999999999999289,70.79999999999999716,280.8000000000000113,3360.0,243.1599999999999966,8691.100000000000363,2,1);
CREATE TABLE IF NOT EXISTS "light_params"
(
    id           INTEGER not null primary key autoincrement,
    view_item_id INTEGER not null unique references view_items on update cascade on delete cascade,
    hue          INTEGER default 0,
    saturation   REAL    default 0,
    brightness   REAL    default 0,
    cct          INTEGER default 0
);
INSERT INTO light_params VALUES(1,2,0,0.0,0.5,3500);
INSERT INTO light_params VALUES(2,3,45,0.5,1.0,0);
INSERT INTO light_params VALUES(3,4,300,1.0,0.5,0);
INSERT INTO light_params VALUES(4,5,360,0.6999999880790710449,0.6000000238418579102,0);
INSERT INTO light_params VALUES(5,6,360,0.3000000119209289551,0.6000000238418579102,7000);
INSERT INTO light_params VALUES(6,14,0,0.0,0.5,0);
INSERT INTO light_params VALUES(7,15,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(8,20,0,0.0,0.0,4300);
INSERT INTO light_params VALUES(9,22,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(10,23,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(11,24,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(12,27,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(13,28,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(14,29,0,0.0,0.0,4000);
INSERT INTO light_params VALUES(15,320,0,0.0,1.0,0);
INSERT INTO light_params VALUES(16,321,0,0.0,1.0,0);
INSERT INTO light_params VALUES(17,322,0,0.0,1.0,0);
INSERT INTO light_params VALUES(18,326,0,0.0,1.0,0);
CREATE TABLE IF NOT EXISTS "events"
(
    id          INTEGER            not null primary key autoincrement,
    target_type TEXT    default '' not null,
    target_id   INTEGER default 0  not null,
    event       TEXT    default '' not null,
    value       TEXT    default '' not null,
    item_id     integer constraint events_view_items_id_fk references view_items (id) on update cascade on delete cascade
);
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('zones',16);
INSERT INTO sqlite_sequence VALUES('view_items',484);
INSERT INTO sqlite_sequence VALUES('boilers',1);
INSERT INTO sqlite_sequence VALUES('boiler_presets',4);
INSERT INTO sqlite_sequence VALUES('boiler_properties',5);
INSERT INTO sqlite_sequence VALUES('dimmers',1);
INSERT INTO sqlite_sequence VALUES('local_users',1);
INSERT INTO sqlite_sequence VALUES('menus',23);
INSERT INTO sqlite_sequence VALUES('notifications',12);
INSERT INTO sqlite_sequence VALUES('scenarios',6);
INSERT INTO sqlite_sequence VALUES('sensors',11);
INSERT INTO sqlite_sequence VALUES('users',63);
INSERT INTO sqlite_sequence VALUES('curtain_params',7);
INSERT INTO sqlite_sequence VALUES('conditioner_params',3);
INSERT INTO sqlite_sequence VALUES('counter_daily_history',3052);
INSERT INTO sqlite_sequence VALUES('counter_monthly_history',104);
INSERT INTO sqlite_sequence VALUES('device_daily_history',1835);
INSERT INTO sqlite_sequence VALUES('device_hourly_history',41079);
INSERT INTO sqlite_sequence VALUES('counters',4);
INSERT INTO sqlite_sequence VALUES('light_params',18);
INSERT INTO sqlite_sequence VALUES('events',73);
INSERT INTO sqlite_sequence VALUES('om_props',27121);
CREATE UNIQUE INDEX object_id_code ON "om_props"(object_id, code);
CREATE UNIQUE INDEX code ON "om_scripts"(code);
CREATE INDEX parent_id on "om_objects" (parent_id);
CREATE INDEX tags on "om_objects" (tags);
CREATE INDEX zone_id on "om_objects" (zone_id);
CREATE INDEX task_id ON ar_cron_actions(task_id);
CREATE UNIQUE INDEX tt_ti_en ON ar_events(target_type, target_id, event_name);
CREATE INDEX event_id ON ar_event_actions(event_id);
CREATE UNIQUE INDEX counter_id_datetime on counter_daily_history (counter_id, datetime);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table 'ar_cron_actions';
drop table 'ar_cron_tasks';
drop table 'ar_event_actions';
drop table 'ar_events';
drop table 'boiler_presets';
drop table 'boiler_properties';
drop table 'boilers';
drop table 'conditioner_params';
drop table 'counter_daily_history';
drop table 'counter_monthly_history';
drop table 'counters';
drop table 'curtain_params';
drop table 'device_daily_history';
drop table 'device_hourly_history';
drop table 'dimmers';
drop table 'events';
drop table 'light_params';
drop table 'local_users';
drop table 'menus';
drop table 'notifications';
drop table 'om_objects';
drop table 'om_props';
drop table 'om_scripts';
drop table 'scenarios';
drop table 'sensors';
drop table 'temp_presets';
drop table 'users';
drop table 'view_items';
drop table 'zones';
-- +goose StatementEnd