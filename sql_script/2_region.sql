/* Dumping data for table `Currency` */

LOCK TABLES `Currency` WRITE;

INSERT INTO `Currency` (code, name, format) VALUES
('ALL','ALL',0),
('AED','阿拉伯幣',0),
('ARS','阿根廷披索',0),
('ATS','先令',0),
('AUD','澳幣',0),
('BEF','比利時法郎',0),
('BRL','巴西幣',0),
('CAD','加拿大幣',0),
('CHF','瑞士法郎',0),
('CLP','智利披索',0),
('CNY','人民幣',0),
('CZK','捷克克郎',0),
('DEM','德國馬克',0),
('DKK','丹麥克郎',0),
('ESP','西班牙比塞塔',0),
('EUR','歐元',0),
('FIM','芬蘭馬克',0),
('FRF','法國法郎',0),
('GBP','英鎊',0),
('HKD','港幣',0),
('HUF','富林特',0),
('IDR','印尼盾',0),
('IEP','愛爾蘭磅',0),
('INR','印度盧比',0),
('ITL','義大利里拉',0),
('JPY','日圓',0),
('KRW','韓圜',0),
('KWD','科威特幣',0),
('MMK','缅甸元',2),
('MOP','澳門幣',0),
('MXN','墨西哥披索',0),
('MYR','馬來西亞幣',0),
('NLG','荷蘭盾',0),
('NOK','挪威克郎',0),
('NZD','紐西蘭幣',0),
('PHP','菲律賓披索',0),
('PLN','茲羅提',0),
('RUB','俄羅斯盧比',0),
('SAR','沙烏地里亞爾',0),
('SEK','瑞典克朗',0),
('SGD','新加坡幣',0),
('THB','泰銖',0),
('TRY','新土耳其里拉',0),
('TWD','新臺幣',0),
('USD','美元',0),
('VEB','委內瑞拉幣',0),
('VND','越南幣',0),
('ZAR','南非幣',0);

UNLOCK TABLES;

/* Dumping data for table `Region` */

LOCK TABLES `Region` WRITE;

INSERT INTO `Region` VALUES
('ALL','ALL','ALL',0),
('MENA', '中東和北非', NULL, 1),
('MM', '緬甸', 'MMK', 1),
('MY','馬來西亞','MYR',1),
('ID','印尼','IDR',1),
('IN','印度','INR',1),
('JP','日本','JPY',1),
('HK','香港','HKD',1),
('PH','菲律賓','PHP',1),
('SG','新加坡','SGD',1),
('TH','泰國','THB',1),
('TW','中華民國','TWD',1),
('US','美國','USD',1),
('VN','越南','VND',1);

UNLOCK TABLES;
