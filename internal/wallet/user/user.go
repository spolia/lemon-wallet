package user

import (
	"context"
	"errors"
)

var ErrorUserNotFound = errors.New("user: not found")
var ErrorAlreadyExist = errors.New("user: already exist")

type Repository interface {
	Save(ctx context.Context, firstName, lastName, alias, email string) (int64, error)
	Get(ctx context.Context, id int64) (User, error)
	Delete(ctx context.Context, id int64) error
}

type User struct {
	ID              int64              `json:"id"`
	FirstName       string             `json:"firstname" binding:"required"`
	LastName        string             `json:"lastname" binding:"required"`
	Alias           string             `json:"alias" binding:"required"`
	Email           string             `json:"email" binding:"required"`
	WalletStatement map[string]float64 `json:"walletstatement"`
}

/*

CRUD
CREATE TABLE `wallet`.`users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `last_name` VARCHAR(45) NOT NULL,
  `alias` VARCHAR(45) NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `alias_UNIQUE` (`alias` ASC),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC));

CREATE TABLE `wallet`.`movements_btc` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `mov_type` ENUM("deposit", "extract") NOT NULL,
  `currency_name` VARCHAR(20) NOT NULL DEFAULT 'USDT',
  `date_created` DATETIME NOT NULL DEFAULT current_timestamp,
  `tx_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `total_amount` DECIMAL(18,8) ZEROFILL NOT NULL,
  `user_id` BIGINT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `user_id_idx` (`user_id` ASC),
  CONSTRAINT `fku_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `wallet`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);



DELIMITER $$
USE `wallet`$$
CREATE DEFINER=`root`@`%` TRIGGER `wallet`.`movements_usdt_BEFORE_INSERT` BEFORE INSERT ON `movements_usdt` FOR EACH ROW
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
*/
