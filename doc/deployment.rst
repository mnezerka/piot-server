Deployment of PIOT
==================

MySql Persistent Storage
------------------------

Schema::

    SET NAMES utf8;
    SET time_zone = '+00:00';
    SET foreign_key_checks = 0;
    SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

    SET NAMES utf8mb4;

    DROP TABLE IF EXISTS `piot_sensors`;
    CREATE TABLE `piot_sensors` (
      `id` varchar(100) NOT NULL,
      `org` varchar(150)  NOT NULL,
      `class` varchar(20) NOT NULL,
      `value` float NOT NULL,
      `time` int NOT NULL,
      PRIMARY KEY (`id`,`time`)
    ) ENGINE=InnoDB;


    DROP TABLE IF EXISTS `piot_switches`;
    CREATE TABLE `piot_switches` (
      `id` varchar(100) NOT NULL,
      `org` varchar(150) NOT NULL,
      `value` int NOT NULL,
      `time` int NOT NULL,
      PRIMARY KEY (`id`,`time`)
    ) ENGINE=InnoDB;
