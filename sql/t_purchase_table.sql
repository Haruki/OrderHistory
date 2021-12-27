create table if not exists t_purchase (id integer primary key, purchase_date text, item_name text, vendor_platform text, div text, price float, img_url text);

select * from t_purchase tp 