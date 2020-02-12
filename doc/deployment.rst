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

    DROP TABLE IF EXISTS `sensors`;
    CREATE TABLE `sensors` (
      `id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
      `org` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
      `class` varchar(20) NOT NULL,
      `value` float NOT NULL,
      `time` int NOT NULL,
      PRIMARY KEY (`id`,`time`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


    DROP TABLE IF EXISTS `switches`;
    CREATE TABLE `switches` (
      `id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
      `org` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
      `value` int NOT NULL,
      `time` int NOT NULL,
      PRIMARY KEY (`id`,`time`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
