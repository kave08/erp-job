


جهت رویت راهنمای API آدرس لوکال http://localhost:80/api یا آدرس دریافتی از شرکت مربوطه را در Browser خود وارد کرده و راهنمای کلیات API را در آنجا ملاحظه نمایید.

لیست API های مورد نیاز مجموعه تابش نور قم به شرح ذیل می باشد. و همه Apiها, Post  میباشند.

نمونه کدهای   postMan جهت استفاده API


Address   : http://localhost:80/api
UserName  : Api
Pass      : 123


سفارش فروش



MethodName : SaleOrder: سفارش فروش
ردیف	نام فیلد	توضیحات	Data Type
1	CustomerId	مشتری	Int
2	VoucherDate	تاریخ ثبت فاکتور	Nvarchar(10)
3	SecondNumber	شماره فرعی	Nvarchar(50)
4	VoucherDesc	شرح سند	Nvarchar(1024)
5	StockID	آی دی انبار	INT
6	SaleTypeId	نوع فروش	INT
7	DeliveryCenterID	محل تحویل	INT
8	SaleCenterID	مرکز فروش	INT
9	PaymentWayID	نحوه پرداخت	INT
10	SellerVisitorID	فروشنده	INT
11	ServiceGoodsID	آی دی کالا	INT
12	Quantity	تعداد	Decimal(10,24)
13	Fee	مبلغ فی	Decimal(23,8)
14	DetailDesc	شرح	Nvarchar(512)





MethodName : SaleCustomer : لیست مشتریان
ردیف	نام فیلد	توضیحات	Data Type
1	CustomerID	آی دی مشتری	INT
2	CustomerCode	کد مشتری	NVarchar(50)



MethodName : SaleTypeSelect : نوع فروش
ردیف	نام فیلد	توضیحات	Data Type
1	BuySaleTypeID	آی دی نوع فروش	INT
2	BuySaleTypeCode	کد نوع فروش	NVarchar(50)
3	BuySaleTypeDesc	شرح نوع فروش	NVarchar(200)

MethodName : SaleCenter4SaleSelectانبار فروش :   لیست
ردیف	نام فیلد	توضیحات	Data Type
1	StockID	آی دی انبار	INT
2	StockCode	کد انبار	NVarchar(100)
3	StockDesc	شرح انبار	NVarchar(100)

MethodName : SalePaymentSelect :نحوه پرداخت
ردیف	نام فیلد	توضیحات	Data Type
1	PaymentWayID	آی دی پرداخت	INT
2	PaymentwayDesc	شرح پرداخت	NVarchar(50)

MethodName : SaleCenterSelect: مرکز فروش
ردیف	نام فیلد	توضیحات	Data Type
1	CentersID	آی دی مرکز فروش	INT
2	CentersCode	کد مرکز فروش	NVarchar(20)
    CenterDesc	شرح مرکز فروش	NVarchar(50)

MethodName : DeliverCenter_SaleSelect: محل تحویل
ردیف	نام فیلد	توضیحات	Data Type
1	CentersID	آی دی محل تحویل	INT
2	CentersCode	شرح محل تحویل	NVarchar(50)

MethodName : SalerSelect: فروشنده
ردیف	نام فیلد	توضیحات	Data Type
1	SaleVisitorID	آی دی فروشنده	INT
2	SaleVisitorDesc	کد فروشنده	NVarchar(50)

MethodName : SaleSellerVisitor :بازاریاب
ردیف	نام فیلد	توضیحات	Data Type
1	CentersID	آی دی محل تحویل	INT
2	CentersCode	شرح محل تحویل	NVarchar(50)

MethodName : Goods : کالا
ردیف	نام فیلد	توضیحات	Data Type
1	ServiceGoodsID	آیدی کالا	INT
2	ServiceGoodsCode	کد کالا	NVarchar(50)
3	ServiceGoodsDesc	نام کالا	NVarchar(255)
4	GroupId	گروه کالا	INT
5	TypeID	نوع کالا	INT
6	SecUnitType	نمایش مقدار واحد فرعی	INT
7	Level1	کد طبقه اول	INT
8	Level2	کد طبقه دوم	INT
9	Level3	کد طبقه سوم	INT



{
"id": "AddSubElementSelect",//عوامل افزاینده کاهنده
     "Params":[{"Name":"IsStruct" ,"Value":1}
           //{"Name":"TafsiliCode", "Type" : "__Query__", "FilterType": 0 ,"Value":100}
          ]
}

----------------------------- پیش فاکتور فروش--------------------------------
{
"id":"SaleProforma", // درج  فاکتور فروش
"Params":[{"Name":"isrow" ,"Value":"1"},
          {"Name":"CustomerId" ,"Value":10000001}, // مشتری
          {"Name":"VoucherDate" ,"Value":"1401/06/19"}, // تاریخ سند
          {"Name":"SecondNumber" ,"Value":3211}, // شماره فرعی
          {"Name":"StockId" ,"Value":10000001}, // آی دی انبار
          {"Name":"VoucherDesc" ,"Value":"شماره اشتراک"}, // شرح سند
          {"Name":"SaleTypeId" ,"Value":10000001}, // نوع فروش
          {"Name":"DeliveryCenterID" ,"Value":10000002},  // محل تحویل
          {"Name":"SaleCenterID" ,"Value":10000001}, // مرکز فروش
          {"Name":"PaymentWayID" ,"Value":10000003}, // نحوه پرداخت
          {"Name":"SellerID" ,"Value":10000001}, // آی دی فروشنده
          {"Name":"[Inserted]" ,"Array_Value":[1,10000384,2,5000,"تست"] },
          {"Name":"[Inserted]" ,"Array_Value":[2,10000307,2,5000,"تست"] },
          {"Name":"[el_Inserted]" ,"Array_Value":[10000001,1000,0]}
          ]
}



MethodName : SaleProforma: پیش فاکتور فروش
ردیف	نام فیلد	توضیحات	Data Type
1	CustomerId	مشتری	Int
2	VoucherDate	تاریخ ثبت فاکتور	Nvarchar(10)
3	SecondNumber	شماره فرعی	Nvarchar(50)
4	VoucherDesc	شرح سند	Nvarchar(1024)
5	StockID	آی دی انبار	INT
6	SaleTypeId	نوع فروش	INT
7	DeliveryCenterID	محل تحویل	INT
8	SaleCenterID	مرکز فروش	INT
9	PaymentWayID	نحوه پرداخت	INT
10	SellerVisitorID	فروشنده	INT
11	ServiceGoodsID	آی دی کالا	INT
12	Quantity	تعداد	Decimal(10,24)
13	Fee	مبلغ فی	Decimal(23,8)
14	DetailDesc	شرح	Nvarchar(512)


MethodName : SaleFactor: فاکتور فروش
ردیف	نام فیلد	توضیحات	Data Type
1	CustomerId	مشتری	Int
2	VoucherDate	تاریخ ثبت فاکتور	Nvarchar(10)
3	SecondNumber	شماره فرعی	Nvarchar(50)
4	VoucherDesc	شرح سند	Nvarchar(1024)
5	StockID	آی دی انبار	INT
6	SaleTypeId	نوع فروش	INT
7	DeliveryCenterID	محل تحویل	INT
8	SaleCenterID	مرکز فروش	INT
9	PaymentWayID	نحوه پرداخت	INT
10	SellerVisitorID	فروشنده	INT
11	ServiceGoodsID	آی دی کالا	INT
12	Quantity	تعداد	Decimal(10,24)
13	Fee	مبلغ فی	Decimal(23,8)
14	DetailDesc	شرح	Nvarchar(512)


-------------------------------------فاکتور فروش----------------------------------
{
"id":"SaleFactor", // درج  فاکتور فروش
"Params":[{"Name":"isrow" ,"Value":"1"},
          {"Name":"CustomerId" ,"Value":10000001}, // مشتری
          {"Name":"VoucherDate" ,"Value":"1401/06/19"}, // تاریخ سند
          {"Name":"SecondNumber" ,"Value":3211}, // شماره فرعی
          {"Name":"StockId" ,"Value":10000001}, // آی دی انبار
          {"Name":"VoucherDesc" ,"Value":"شماره اشتراک"}, // شرح سند
          {"Name":"SaleTypeId" ,"Value":10000001}, // نوع فروش
          //{"Name": "SaleCenterID","value":"10000004"},//مشتری مرتبط
          {"Name":"DeliveryCenterID" ,"Value":10000002},  // محل تحویل
          {"Name":"SaleCenterID" ,"Value":10000001}, // مرکز فروش
          {"Name":"PaymentWayID" ,"Value":10000003}, // نحوه پرداخت
          {"Name":"SellerID" ,"Value":10000001}, // آی دی فروشنده
          {"Name":"[Inserted]" ,"Array_Value":[1,10000384,2,5000,"تست"] },
          {"Name":"[Inserted]" ,"Array_Value":[2,10000307,2,5000,"تست"] },
          {"Name":"[el_Inserted]" ,"Array_Value":[10000001,1000,0]}
          ]
}


MethodName : SaleFactor: فاکتور فروش
ردیف	نام فیلد	توضیحات	Data Type
1	CustomerId	مشتری	Int
2	VoucherDate	تاریخ ثبت فاکتور	Nvarchar(10)
3	SecondNumber	شماره فرعی	Nvarchar(50)
4	VoucherDesc	شرح سند	Nvarchar(1024)
5	StockID	آی دی انبار	INT
6	SaleTypeId	نوع فروش	INT
7	DeliveryCenterID	محل تحویل	INT
8	SaleCenterID	مرکز فروش	INT
9	PaymentWayID	نحوه پرداخت	INT
10	SellerVisitorID	فروشنده	INT
11	ServiceGoodsID	آی دی کالا	INT
12	Quantity	تعداد	Decimal(10,24)
13	Fee	مبلغ فی	Decimal(23,8)
14	DetailDesc	شرح	Nvarchar(512)


-------------------------------------فاکتور فروش----------------------------------
{
"id":"SaleFactor", // درج  فاکتور فروش
"Params":[{"Name":"isrow" ,"Value":"1"},
          {"Name":"CustomerId" ,"Value":10000001}, // مشتری
          {"Name":"VoucherDate" ,"Value":"1401/06/19"}, // تاریخ سند
          {"Name":"SecondNumber" ,"Value":3211}, // شماره فرعی
          {"Name":"StockId" ,"Value":10000001}, // آی دی انبار
          {"Name":"VoucherDesc" ,"Value":"شماره اشتراک"}, // شرح سند
          {"Name":"SaleTypeId" ,"Value":10000001}, // نوع فروش
          //{"Name": "SaleCenterID","value":"10000004"},//مشتری مرتبط
          {"Name":"DeliveryCenterID" ,"Value":10000002},  // محل تحویل
          {"Name":"SaleCenterID" ,"Value":10000001}, // مرکز فروش
          {"Name":"PaymentWayID" ,"Value":10000003}, // نحوه پرداخت
          {"Name":"SellerID" ,"Value":10000001}, // آی دی فروشنده
          {"Name":"[Inserted]" ,"Array_Value":[1,10000384,2,5000,"تست"] },
          {"Name":"[Inserted]" ,"Array_Value":[2,10000307,2,5000,"تست"] },
          {"Name":"[el_Inserted]" ,"Array_Value":[10000001,1000,0]}
          ]
}
