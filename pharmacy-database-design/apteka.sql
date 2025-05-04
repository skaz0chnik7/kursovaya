-- phpMyAdmin SQL Dump
-- version 3.5.1
-- http://www.phpmyadmin.net
--
-- Хост: 127.0.0.1
-- Время создания: Дек 02 2023 г., 19:37
-- Версия сервера: 5.5.25
-- Версия PHP: 5.3.13

SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- База данных: `apteka`
--

-- --------------------------------------------------------

--
-- Структура таблицы `заказ`
--

CREATE TABLE IF NOT EXISTS `заказ` (
  `Код заказа` int(11) NOT NULL,
  `Код сотрудника` int(11) NOT NULL,
  `Код клиента` int(11) NOT NULL,
  `Дата размещения` date NOT NULL,
  PRIMARY KEY (`Код заказа`),
  KEY `Код сотрудника` (`Код сотрудника`,`Код клиента`),
  KEY `Код клиента` (`Код клиента`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `клиент`
--

CREATE TABLE IF NOT EXISTS `клиент` (
  `Код клиента` int(11) NOT NULL,
  `Название фирмы/представителя` varchar(100) NOT NULL,
  `Адрес` varchar(100) NOT NULL,
  `Телефон` varchar(20) NOT NULL,
  PRIMARY KEY (`Код клиента`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `лекарство`
--

CREATE TABLE IF NOT EXISTS `лекарство` (
  `Код лекарства` int(11) NOT NULL,
  `Название` varchar(100) NOT NULL,
  `Код производителя` int(11) NOT NULL,
  `Страна производитель` varchar(50) NOT NULL,
  `Форма выпуска` varchar(50) NOT NULL,
  `Срок годности` date NOT NULL,
  `Действующее вещество` varchar(100) NOT NULL,
  `Дозировка` float NOT NULL,
  `Стоимость` float NOT NULL,
  `Код назначения` int(11) NOT NULL,
  PRIMARY KEY (`Код лекарства`),
  KEY `Код назначения` (`Код назначения`),
  KEY `Код производителя` (`Код производителя`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `назначение лекарства`
--

CREATE TABLE IF NOT EXISTS `назначение лекарства` (
  `Код назначения` int(11) NOT NULL,
  `Группа` varchar(50) NOT NULL,
  `Описание` varchar(200) NOT NULL,
  PRIMARY KEY (`Код назначения`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `поставка`
--

CREATE TABLE IF NOT EXISTS `поставка` (
  `Код поставки` int(11) NOT NULL,
  `Код поставщика` int(11) NOT NULL,
  `Дата поставки` date NOT NULL,
  PRIMARY KEY (`Код поставки`),
  KEY `Код поставщика` (`Код поставщика`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `поставщик`
--

CREATE TABLE IF NOT EXISTS `поставщик` (
  `Код поставщика` int(11) NOT NULL,
  `Название фирмы` varchar(100) NOT NULL,
  `Адрес` varchar(100) NOT NULL,
  `Телефон` varchar(20) NOT NULL,
  `ИНН` int(12) NOT NULL,
  PRIMARY KEY (`Код поставщика`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `продажа заказа`
--

CREATE TABLE IF NOT EXISTS `продажа заказа` (
  `Код товарной накладной` int(11) NOT NULL,
  `Код лекарства` int(11) NOT NULL,
  `Количество` int(11) NOT NULL,
  `Стоимость` float NOT NULL,
  KEY `Код товарной накладной` (`Код товарной накладной`),
  KEY `Код лекарства` (`Код лекарства`),
  KEY `Код лекарства_2` (`Код лекарства`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `продажа поставки`
--

CREATE TABLE IF NOT EXISTS `продажа поставки` (
  `Код товарной накладной` int(11) NOT NULL,
  `Код лекарства` int(11) NOT NULL,
  `Количество` int(11) NOT NULL,
  `Стоимость` float NOT NULL,
  UNIQUE KEY `Код товарной накладной` (`Код товарной накладной`),
  KEY `Код лекарства` (`Код лекарства`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `сотрудник`
--

CREATE TABLE IF NOT EXISTS `сотрудник` (
  `Код сотрудника` int(11) NOT NULL,
  `Фамилия` varchar(50) NOT NULL,
  `Имя` varchar(50) NOT NULL,
  `Отчество` varchar(50) NOT NULL,
  `Должность` varchar(50) NOT NULL,
  `Дата найма` date NOT NULL,
  `Телефон` varchar(20) NOT NULL,
  `Образование` varchar(100) NOT NULL,
  `Оклад` float NOT NULL,
  PRIMARY KEY (`Код сотрудника`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `товарная накладная заказа`
--

CREATE TABLE IF NOT EXISTS `товарная накладная заказа` (
  `Код товарной накладной` int(11) NOT NULL,
  `Код заказа` int(11) DEFAULT NULL,
  PRIMARY KEY (`Код товарной накладной`),
  KEY `Код заказа` (`Код заказа`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Структура таблицы `товарная накладная поставки`
--

CREATE TABLE IF NOT EXISTS `товарная накладная поставки` (
  `Код товарной накладной` int(11) NOT NULL,
  `Код поставки` int(11) NOT NULL,
  PRIMARY KEY (`Код товарной накладной`),
  KEY `Код поставки` (`Код поставки`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Ограничения внешнего ключа сохраненных таблиц
--

--
-- Ограничения внешнего ключа таблицы `заказ`
--
ALTER TABLE `заказ`
  ADD CONSTRAINT `@n0@g0@q0@g0@n0_ibfk_1` FOREIGN KEY (`Код клиента`) REFERENCES `клиент` (`Код клиента`),
  ADD CONSTRAINT `@n0@g0@q0@g0@n0_ibfk_2` FOREIGN KEY (`Код сотрудника`) REFERENCES `сотрудник` (`Код сотрудника`);

--
-- Ограничения внешнего ключа таблицы `лекарство`
--
ALTER TABLE `лекарство`
  ADD CONSTRAINT `@r0@l0@q0@g0@w0@x0@y0@i0@u0_ibfk_1` FOREIGN KEY (`Код назначения`) REFERENCES `назначение лекарства` (`Код назначения`);

--
-- Ограничения внешнего ключа таблицы `поставка`
--
ALTER TABLE `поставка`
  ADD CONSTRAINT `@v0@u0@x0@y0@g0@i0@q0@g0_ibfk_1` FOREIGN KEY (`Код поставщика`) REFERENCES `поставщик` (`Код поставщика`);

--
-- Ограничения внешнего ключа таблицы `продажа заказа`
--
ALTER TABLE `продажа заказа`
  ADD CONSTRAINT `@v0@w0@u0@k0@g0@m0@g0@0020@n0@g0@q0@g0@n0@g0_ibfk_3` FOREIGN KEY (`Код товарной накладной`) REFERENCES `товарная накладная заказа` (`Код товарной накладной`),
  ADD CONSTRAINT `@v0@w0@u0@k0@g0@m0@g0@0020@n0@g0@q0@g0@n0@g0_ibfk_2` FOREIGN KEY (`Код лекарства`) REFERENCES `лекарство` (`Код лекарства`);

--
-- Ограничения внешнего ключа таблицы `продажа поставки`
--
ALTER TABLE `продажа поставки`
  ADD CONSTRAINT `@v0@w0@u0@k0@g0@m0@g0@0020@v0@u0@x0@y0@g0@i0@q0@o0_ibfk_2` FOREIGN KEY (`Код лекарства`) REFERENCES `лекарство` (`Код лекарства`),
  ADD CONSTRAINT `@v0@w0@u0@k0@g0@m0@g0@0020@v0@u0@x0@y0@g0@i0@q0@o0_ibfk_1` FOREIGN KEY (`Код товарной накладной`) REFERENCES `товарная накладная поставки` (`Код товарной накладной`);

--
-- Ограничения внешнего ключа таблицы `товарная накладная заказа`
--
ALTER TABLE `товарная накладная заказа`
  ADD CONSTRAINT `@y0@u0@i0@g0@w0@t0@g0@r1@0020@t0@g0@q0@r0@g0@k0@t0@g0@r1@0020@n0@g0@q0@g0@n0@g0_ibfk_2` FOREIGN KEY (`Код заказа`) REFERENCES `заказ` (`Код заказа`);

--
-- Ограничения внешнего ключа таблицы `товарная накладная поставки`
--
ALTER TABLE `товарная накладная поставки`
  ADD CONSTRAINT `@y0@u0@i0@g0@w0@t0@g0@r1@0020@t0@g0@q0@r0@g0@k0@t0@g0@r1@0020@v0@u0@x0@y0@g0@i0@q0@o0_ibfk_1` FOREIGN KEY (`Код поставки`) REFERENCES `поставка` (`Код поставки`);

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
