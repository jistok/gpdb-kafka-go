-- DDL for the Chicago crimes data set
-- Ref. https://data.cityofchicago.org/Public-Safety/Crimes-2001-to-present/ijzp-q8t2/data
DROP TABLE IF EXISTS crimes;
CREATE TABLE crimes
(
  id INT
  , case_number VARCHAR (20)
  , crime_date TIMESTAMP
  , block VARCHAR(50)
  , IUCR VARCHAR(10)
  , primary_type VARCHAR(50)
  , description VARCHAR(75)
  , location_desc VARCHAR (75)
  , arrest BOOLEAN
  , domestic BOOLEAN
  , beat VARCHAR(7)
  , district VARCHAR(7)
  , ward SMALLINT
  , community_area VARCHAR(10)
  , fbi_code VARCHAR(5)
  , x_coord FLOAT
  , y_coord FLOAT
  , crime_year SMALLINT
  , record_update_date TIMESTAMP
  , latitude FLOAT
  , longitude FLOAT
  , location VARCHAR (60)
)
distributed by (id);

-- External table for loading a small subset of the Chicago crimes data
DROP EXTERNAL TABLE IF EXISTS crimes_kafka;
CREATE EXTERNAL WEB TABLE crimes_kafka
(LIKE crimes)
EXECUTE '$HOME/kafka_consumer -zookeeper my_zk_host:2181 -topic chicago_crimes 2>>$HOME/`printf "kafka_consumer_%02d.log" $GP_SEGMENT_ID`'
ON ALL FORMAT 'CSV' (DELIMITER ',' NULL '')
LOG ERRORS INTO err_crimes SEGMENT REJECT LIMIT 1 PERCENT;
