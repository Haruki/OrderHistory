create table if not exists t_purchase (id integer primary key, purchase_date text, item_name text, vendor_platform text, div text, price integer, currency text, img_url text, img_file text,img_hash text, UNIQUE(item_name, purchase_date, vendor_platform, price, img_hash) ON CONFLICT IGNORE);

select * from t_purchase tp 

update t_purchase set item_name = 'Original Samsung S3 Akku Galaxy EB-L1G6LLU GT-I9300 Neo Ersatz Batterie Accu NEU' where id = 76

drop table t_purchase 

delete  from t_purchase

select * from t_purchase where id >= 82 -25 order by id

--t_ignoreHash tih 

create table if not exists t_ignoreHash (id integer primary key, vendor_platform text, img_hash text)

select * from t_ignoreHash 

insert into t_ignoreHash (vendor_platform, img_hash) values ('ebay', 'a567462f4edd496bdf5cd00da5bbde64131c283e3cf396bfd58c0fac26b13d9a')