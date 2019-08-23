-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Aug 23, 2019 at 10:08 AM
-- Server version: 10.3.16-MariaDB
-- PHP Version: 7.3.7

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `be_test`
--

-- --------------------------------------------------------

--
-- Table structure for table `exchange`
--

CREATE TABLE `exchange` (
  `id` int(11) NOT NULL,
  `from` varchar(3) COLLATE utf8mb4_unicode_ci NOT NULL,
  `to` varchar(3) COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `exchange`
--

INSERT INTO `exchange` (`id`, `from`, `to`) VALUES
(1, 'USD', 'IDR'),
(2, 'IDR', 'USD');

-- --------------------------------------------------------

--
-- Table structure for table `rates`
--

CREATE TABLE `rates` (
  `id` int(11) NOT NULL,
  `exchange_id` int(11) NOT NULL,
  `rate` double NOT NULL,
  `date` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `rates`
--

INSERT INTO `rates` (`id`, `exchange_id`, `rate`, `date`) VALUES
(1, 1, 14000, '2019-08-01'),
(2, 1, 14020, '2019-08-02'),
(3, 1, 14045, '2019-08-03'),
(4, 1, 14105, '2019-08-04'),
(5, 1, 14130, '2019-08-05'),
(6, 1, 14160, '2019-08-06'),
(7, 1, 14504, '2019-08-07'),
(8, 2, 0.00007, '2019-08-01'),
(9, 2, 0.000074, '2019-08-02'),
(11, 2, 0.000076, '2019-08-03'),
(12, 2, 0.000078, '2019-08-04'),
(13, 2, 0.00008, '2019-08-05'),
(14, 2, 0.0001, '2019-08-06'),
(15, 2, 0.000501, '2019-08-07');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `exchange`
--
ALTER TABLE `exchange`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `rates`
--
ALTER TABLE `rates`
  ADD PRIMARY KEY (`id`),
  ADD KEY `exchange_id` (`exchange_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `exchange`
--
ALTER TABLE `exchange`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `rates`
--
ALTER TABLE `rates`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=19;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `exchange`
--
ALTER TABLE `exchange`
  ADD CONSTRAINT `Exchange Rates` FOREIGN KEY (`id`) REFERENCES `rates` (`exchange_id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `rates`
--
ALTER TABLE `rates`
  ADD CONSTRAINT `rates_ibfk_1` FOREIGN KEY (`exchange_id`) REFERENCES `exchange` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
