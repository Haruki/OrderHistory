delete from t_purchase where id in (select id from t_purchase t1
where exists (select * from t_purchase t2
  where t2.id > t1.id
    and t1.purchase_date = t2.purchase_date 
    and t1.item_name = t2.item_name 
    and t1.vendor_platform = t2.vendor_platform 
    and t1.price = t2.price))