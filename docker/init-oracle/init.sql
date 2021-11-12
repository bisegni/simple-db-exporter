CREATE TABLESPACE CUSTOM_TABLE_TABLESPACE
  DATAFILE 'CUSTOM_TABLE_TABLESPACE.dat'
  SIZE 10M AUTOEXTEND on;

CREATE TEMPORARY TABLESPACE CUSTOM_TABLESPACE_TEMP
  TEMPFILE 'CUSTOM_TABLESPACE_TEMP.dat'
  SIZE 5M AUTOEXTEND on;

CREATE USER user_test
  IDENTIFIED BY user_test_password
  DEFAULT TABLESPACE CUSTOM_TABLE_TABLESPACE
  TEMPORARY TABLESPACE CUSTOM_TABLESPACE_TEMP
  QUOTA 20M on CUSTOM_TABLE_TABLESPACE;

GRANT create session TO user_test;
GRANT create table TO user_test;
GRANT create view TO user_test;
GRANT create any trigger TO user_test;
GRANT create any procedure TO user_test;
GRANT create sequence TO user_test;
GRANT create synonym TO user_test;

CREATE TABLE user_test.table_test (
    id NUMBER,
    first_name VARCHAR2(50) NOT NULL,
    last_name VARCHAR2(50) NOT NULL,
    last_date DATE NOT NULL,
    PRIMARY KEY(id)
);

insert into user_test.table_test values (1, 'nome_1', 'cognome_1', SYSDATE);
insert into user_test.table_test values (2, 'nome_2', 'cognome_2', SYSDATE);
insert into user_test.table_test values (3, 'nome_3', 'cognome_3', SYSDATE);