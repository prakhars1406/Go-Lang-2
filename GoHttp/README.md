

## Account service
Contains account settings for organisations and members

```
route   /api/account-service/
```

Method	| Path	| Description	| User authenticated	| Available from UI | Done
------------- | ------------------------------- | ------------- |:-------------:|:----------------:|---------:|
POST	| /tempaccounts/	| Create new Temp Account	|   | × | []
GET	| /tempaccounts/	| Get all temporary accounts (based on access)	| x  |  | []
GET	| /tempaccounts/{tempaccount}	| Get temporary account detail	| x  |  | []
GET	| /accounts/	| Get list of all authorised accounts	|  | | []
POST	| /accounts/	| Register new account	| × | × | [x]
GET	| /accounts/{account}	| Get account detail for provided account	|   | 	× | []
PUT	| /accounts/{account}	| Save account detail	| × | × | []
GET	| /accounts/{account}/members	| Get all members for the account	| x  |  | []
POST	| /accounts/{account}/members	| Enroll new Member| x  |  | [] 
GET	| /accounts/{account}/members/{member}	| Get member detail | x  |  | []
PUT	| /accounts/{account}/members/{member}	| Update member detail | x  |  | []
GET	| /features	| Get all features master	| x  |  | []
GET	| /accounts/{account}/features	| Get all features for the account	| x  |  | []
POST	| /accounts/{account}/features	| Add new feature| x  |  | [] 
GET	| /accounts/{account}/features/{feature}	| Get Feature detail | x  |  | []
PUT	| /accounts/{account}/features/{feature}	| Update feature detail | x  |  | []
POST	| /membership/{member}/enroll/default	| Enroll the member as default to emtrust account| x  |  | [] 