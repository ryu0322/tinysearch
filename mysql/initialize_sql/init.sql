create schema if not exists tinysearch default character set utf8mb4;

use tinysearch;

create table if not exists documents(
    PRIMARY KEY (document_id),
    document_id INT UNSIGNED AUTO_INCREMENT NOT NULL,
    document_title TEXT NOT NULL
) ENGINED=InnoDB DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin;