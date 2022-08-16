/* Insert first admin user */

INSERT INTO User VALUES (1, 'system', 'system', 'system@17.media', 1, '2022-03-25 15:00:00.000', 1, '2022-03-25 15:00:00.000', 1);

INSERT INTO UserAuth VALUES (1, 1, 'ALL', 99, 3, 0, '2022-03-25 15:00:00.000', 1);

/* Execute after insert first admin user */

ALTER TABLE User ADD FOREIGN KEY (creatorUID) REFERENCES User(id);

ALTER TABLE User ADD FOREIGN KEY (modifierUID) REFERENCES User(id);
