-- MySQL dump 10.13  Distrib 8.0.28, for Win64 (x86_64)
--
-- Host: localhost    Database: portto_shop
-- ------------------------------------------------------
-- Server version	8.0.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cart`
--

DROP TABLE IF EXISTS `cart`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cart` (
  `User_id` int NOT NULL,
  `Merchandise_id` int NOT NULL,
  `Number` int NOT NULL,
  PRIMARY KEY (`User_id`,`Merchandise_id`),
  KEY `cart_fk_merchandise` (`Merchandise_id`),
  CONSTRAINT `cart_fk_merchandise` FOREIGN KEY (`Merchandise_id`) REFERENCES `merchandise` (`Merchandise_id`) ON DELETE CASCADE,
  CONSTRAINT `cart_fk_user` FOREIGN KEY (`User_id`) REFERENCES `user` (`User_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `merchandise`
--

DROP TABLE IF EXISTS `merchandise`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `merchandise` (
  `Merchandise_id` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `Price` int NOT NULL,
  `Number_in_stock` int NOT NULL,
  `Number_sold` int NOT NULL,
  `User_id` int NOT NULL,
  `Description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci,
  PRIMARY KEY (`Merchandise_id`),
  UNIQUE KEY `Merchandise_id_UNIQUE` (`Merchandise_id`),
  KEY `merchandise_fk_user` (`User_id`),
  CONSTRAINT `merchandise_fk_user` FOREIGN KEY (`User_id`) REFERENCES `user` (`User_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order`
--

DROP TABLE IF EXISTS `order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order` (
  `Order_id` int NOT NULL AUTO_INCREMENT,
  `User_id` int NOT NULL,
  `Order_time` datetime NOT NULL,
  PRIMARY KEY (`Order_id`),
  KEY `order_fk_user` (`User_id`),
  CONSTRAINT `order_fk_user` FOREIGN KEY (`User_id`) REFERENCES `user` (`User_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orderitem`
--

DROP TABLE IF EXISTS `orderitem`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orderitem` (
  `Order_id` int NOT NULL,
  `Merchandise_id` int NOT NULL,
  `Trade_price` int NOT NULL,
  `Number` int NOT NULL,
  PRIMARY KEY (`Order_id`,`Merchandise_id`),
  KEY `orderitem_fk_merchandise` (`Merchandise_id`),
  CONSTRAINT `orderitem_fk_merchandise` FOREIGN KEY (`Merchandise_id`) REFERENCES `merchandise` (`Merchandise_id`) ON DELETE CASCADE,
  CONSTRAINT `orderitem_fk_order` FOREIGN KEY (`Order_id`) REFERENCES `order` (`Order_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `User_id` int NOT NULL AUTO_INCREMENT,
  `Account` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `Password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `Name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `Register_time` datetime DEFAULT NULL,
  `Bill_info` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci,
  `Email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci DEFAULT NULL,
  PRIMARY KEY (`User_id`),
  UNIQUE KEY `Account_UNIQUE` (`Account`),
  UNIQUE KEY `mail_UNIQUE` (`Email`)
) ENGINE=InnoDB AUTO_INCREMENT=312 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-02-18 18:17:31
