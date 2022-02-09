CREATE SCHEMA `wallet` ;
CREATE TABLE `wallet`.`users` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `first_name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `alias` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `alias_UNIQUE` (`alias` ASC),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC));

CREATE TABLE `wallet`.`movements_btc` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "extract","init") NOT NULL,
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'BTC',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `user_id` BIGINT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id_idx` (`user_id` ASC),
  CONSTRAINT `fk_btc_user_id`
      FOREIGN KEY (`user_id`)
          REFERENCES `wallet`.`users` (`id`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);

CREATE TABLE `wallet`.`movements_usdt` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "extract","init") NOT NULL,
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'USDT',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
  `user_id` BIGINT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id_idx` (`user_id` ASC),
  CONSTRAINT `fk_usdt_user_id`
      FOREIGN KEY (`user_id`)
          REFERENCES `wallet`.`users` (`id`)
          ON DELETE CASCADE
          ON UPDATE CASCADE);

CREATE TABLE `wallet`.`movements_ars` (
   `id` BIGINT NOT NULL AUTO_INCREMENT,
   `mov_type` ENUM("deposit", "extract") NOT NULL,
   `currency_name` VARCHAR(20) NOT NULL DEFAULT 'ARS',
   `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
   `tx_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
   `total_amount` DECIMAL(18,2) ZEROFILL NOT NULL,
   `user_id` BIGINT NOT NULL,
   PRIMARY KEY (`id`),
   INDEX `user_id_idx` (`user_id` ASC),
   CONSTRAINT `fk_ars_user_id`
       FOREIGN KEY (`user_id`)
           REFERENCES `wallet`.`users` (`id`)
           ON DELETE CASCADE
           ON UPDATE CASCADE);


/*Triggers*/


DROP TRIGGER IF EXISTS `wallet`.`movements_usdt_BEFORE_INSERT`;
DELIMITER $$
USE `wallet`$$
CREATE DEFINER = CURRENT_USER TRIGGER `wallet`.`movements_usdt_BEFORE_INSERT` BEFORE INSERT ON `movements_usdt` FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_usdt WHERE id =(SELECT max(id) from movements_usdt WHERE user_id = NEW.user_id));
END IF;
IF NEW.mov_type='extract' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_usdt WHERE id =(SELECT max(id) from movements_usdt WHERE user_id = NEW.user_id))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;



DROP TRIGGER IF EXISTS `wallet`.`movements_btc_BEFORE_INSERT`;

DELIMITER $$
USE `wallet`$$
CREATE DEFINER = CURRENT_USER TRIGGER `wallet`.`movements_btc_BEFORE_INSERT` BEFORE INSERT ON `movements_btc` FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_btc WHERE id =(SELECT max(id) from movements_btc WHERE user_id = NEW.user_id));
END IF;
IF NEW.mov_type='extract' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_btc WHERE id =(SELECT max(id) from movements_btc WHERE user_id = NEW.user_id))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;



DROP TRIGGER IF EXISTS `wallet`.`movements_ars_BEFORE_INSERT`;

DELIMITER $$
USE `wallet`$$
CREATE DEFINER = CURRENT_USER TRIGGER `wallet`.`movements_ars_BEFORE_INSERT` BEFORE INSERT ON `movements_ars` FOR EACH ROW
BEGIN
IF NEW.mov_type='deposit' THEN
		SET NEW.total_amount = NEW.tx_amount +
		(SELECT total_amount FROM movements_ars WHERE id =(SELECT max(id) from movements_ars WHERE user_id = NEW.user_id));
END IF;
IF NEW.mov_type='extract' THEN
		SET NEW.total_amount = (SELECT total_amount FROM movements_ars WHERE id =(SELECT max(id) from movements_ars WHERE user_id = NEW.user_id))
        - NEW.tx_amount;
END IF;
END$$
DELIMITER ;