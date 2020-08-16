## Voucher Service
---

Services related to Vouchers and exposes API's related to vouchers

#### How to run
Install docker and run
```
docker-compose up
```

#### Access URL

After running can access the url at localhost:8080

#### API Exposed
 

```
route   /api/voucher-service/
```

Method	| Path	| Description	| User authenticated	| Available from UI | Done
------------- | ------------------------------- | ------------- |:-------------:|:----------------:|---------:|
GET	| /deals	| Get all voucher templates by params, filters can be ?issuer=xx or ?owner=xxx	|  | | [x]
POST	| /deals	| Record new deals (voucher templates)	| × | × | [x]
GET	| /deals/{deal}	| Get respective voucher template details	|   | 	× | [x]
PUT	| /deals/{deal}	| Save deals (template) detail	| × | × | [x]
GET	| /vouchers	| Get all vouchers (claimed vouchers) , filter with ?did=xxx	| x  |  | [x]
POST	| /vouchers	| Claim new voucher , pass template details as post body| x  |  | [x] 
GET	| /vouchers/{voucher}	| get voucher details | x  |  | [x]
PUT	| /vouchers/{voucher}	| update voucher details | x  |  | [x]
POST	| /vouchers/{voucher}/transfer	| Transfer vouchers | x  |  | [x]
GET	| /vouchers/{voucher}/history	| Get history of voucher	| x  |  | [x]

