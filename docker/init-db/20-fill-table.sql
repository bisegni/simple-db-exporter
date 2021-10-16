-- Set params
set session my.number_of_test_data = '300';

-- Filling of test_data
INSERT INTO test_data
select id, concat('data ', id) 
FROM GENERATE_SERIES(1, 300) as id;