create table if not exists t_purchase (id integer primary key, purchase_date text, item_name text, vendor_platform text, div text, price integer, currency text, img_url text, img_file text);

select * from t_purchase tp 

drop table t_purchase 