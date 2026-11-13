# TSETMC Cookbook (E1-1)

> وضعیت: پیش‌نویس. تاریخ بررسی: 1405/08/12 (2026-11-03)، بعد از ساعت بازار.
> روش: صفحه‌های tsetmc.com در مرورگر ← F12 / Network ← بازتولید با curl ← ذخیره‌ی پاسخ خام در `docs/samples/`.
> هر ادعای این فایل باید به یک نمونه‌ی خام ارجاع بدهد. موارد «تأییدنشده» فرضیه‌اند.
> برای مرتب‌کردن متن از AI کمک گرفته شده (در README اعلام شود)؛ همه‌ی یافته‌ها از درخواست‌های واقعی است.

Base URL همه‌ی endpointها: `https://cdn.tsetmc.com/api`
همه: `GET`، پاسخ JSON، بدون لاگین.

---

## 1. جستجوی ابزار
- **هدف:** پیدا کردن ابزارها (سهم، اختیار) با بخشی از نماد؛ گرفتن `insCode`.
- **آدرس:** `GET /Instrument/GetInstrumentSearch/{text}` — `text` باید URL-encode شود.
- **curl:**
  ```
  curl -s "https://cdn.tsetmc.com/api/Instrument/GetInstrumentSearch/%D8%B6%D8%AE" -o docs/samples/search_zakh.json
  ```
- **نمونه‌ی خام:** `docs/samples/search_zakh.json` (متن جستجو: «ضخ»)
- **ساختار:** `{"instrumentSearch": [ ... ]}`

| فیلد | معنی | مثال |
|---|---|---|
| insCode | شناسه‌ی یکتای ابزار | 31459957843277944 |
| lVal18AFC | نماد فارسی | ضخود1161 |
| lVal30 | نام کامل؛ برای اختیار: `اختيارخ <پایه>-<strike>-<سررسید شمسی>` | اختيارخ خودرو-450-1405/11/07 |
| flow / flowTitle | بازار؛ 3 = «بازار مشتقه» | 3 |
| cgrValCot | گروه؛ در همه‌ی اختیارهای نمونه 31 | 31 |
| lastDate | تأییدنشده: 1 = فعال، 0 = منقضی | 1 / 0 |
| insCode2..4 | در نمونه همه برابر insCode؛ معنی نامعلوم | — |

- **رفتار:**
  - قراردادهای **منقضی** هم برمی‌گردند (`ضخود3069`، سررسید 1402/03/03، `lastDate: 0`).
  - در پاسخ جستجو `instrumentID` و `cIsin` برابر `null` هستند؛ ISIN را باید از GetInstrumentInfo گرفت.
- **سؤال باز:** آیا تعداد نتایج محدود است؟ (برای «ضخ» حدود 40 نتیجه آمد؛ آیا همه‌ی منقضی‌های خودرو را برمی‌گرداند؟)

---

## 2. اطلاعات ابزار
- **هدف:** مشخصات کلی ابزار برای رجیستری.
- **آدرس:** `GET /Instrument/GetInstrumentInfo/{insCode}`
- **curl:**
  ```
  curl -s "https://cdn.tsetmc.com/api/Instrument/GetInstrumentInfo/31459957843277944" -o docs/samples/instrument_info_31459957843277944.json
  ```
- **نمونه‌ی خام:** `docs/samples/instrument_info_31459957843277944.json`

| فیلد | معنی | مثال |
|---|---|---|
| insCode | شناسه‌ی یکتا | 31459957843277944 |
| lVal18AFC | نماد فارسی | ضخود1161 |
| lVal30 | نام کامل | اختيارخ خودرو-450-1405/11/07 |
| lVal18 | نماد لاتین؛ برای اختیار شامل سررسید | IKCO-O-14051107 |
| instrumentID | کد ISIN ابزار | IRO9IKCO6P01 |
| cIsin | کد ISIN شرکت/گروه (≠ instrumentID) | IRO9IKCO6P02 |
| contractSize | اندازه‌ی قرارداد | 1000 |
| flow / flowTitle | بازار | 3 / بازار مشتقه |
| sector.lSecVal | صنعت | خودرو و ساخت قطعات |
| staticThreshold.psGelStaMax/Min | سقف/کف قیمت مجاز | 100000 / 1 |
| minWeek/maxWeek/minYear/maxYear | بازه‌ی قیمت | 218 / 245 |

---

## 3. مشخصات قرارداد اختیار
- **هدف:** strike، سررسید، اندازه‌ی قرارداد و **دارایی پایه**.
- **آدرس:** `GET /Instrument/GetInstrumentOptionByInstrumentID/{instrumentID}` — ورودی ISIN است، نه insCode.
- **curl:**
  ```
  curl -s "https://cdn.tsetmc.com/api/Instrument/GetInstrumentOptionByInstrumentID/IRO9IKCO6P01" -o docs/samples/instrument_option_IRO9IKCO6P01.json
  ```
- **نمونه‌ی خام:** `docs/samples/instrument_option_IRO9IKCO6P01.json`

| فیلد | معنی | مثال |
|---|---|---|
| insCode | شناسه‌ی قرارداد | 31459957843277944 |
| instrumentID | ISIN | IRO9IKCO6P01 |
| strikePrice | قیمت اعمال | 450.00 |
| contractSize | اندازه‌ی قرارداد | 1000 |
| uaInsCode | insCode دارایی پایه | 65883838195688438 |
| beginDate / endDate | شروع / سررسید (میلادی YYYYMMDD) | 20261031 / 20270127 |
| buyOP / sellOP | تأییدنشده: موقعیت باز (صفحه «موقعیت‌های باز 12,375» نشان داد) | 12375 |
| aFactor / bFactor / cFactor | تأییدنشده: ضرایب وجه تضمین | 0.2 / 0.1 / 10000 |

- **سؤال باز:** برای قرارداد منقضی هم جواب می‌دهد؟ برای put (ط...) هم همین ساختار است؟

---

## 4. خلاصه‌ی معاملات امروز
- **هدف:** قیمت پایانی و حجم/تعداد/ارزش امروز.
- **آدرس:** `GET /ClosingPrice/GetClosingPriceInfo/{insCode}`
- **نمونه‌ی خام:** `docs/samples/closing_price_info_31459957843277944.json`

| فیلد | معنی | مثال |
|---|---|---|
| dEven / hEven | تاریخ میلادی / ساعت (HHMMSS) | 20261103 / 114310 |
| pClosing | قیمت پایانی | 244 |
| pDrCotVal | آخرین قیمت معامله | 244 |
| priceFirst / priceMin / priceMax | اولین / کمترین / بیشترین | 245 / 244 / 245 |
| priceYesterday | قیمت دیروز | 222 |
| zTotTran | تعداد معاملات | 7 |
| qTotTran5J | حجم معاملات | 6100 |
| qTotCap | ارزش معاملات | 1488500000 |
| instrumentState.cEtavalTitle | وضعیت نماد | مجاز |

---

## 5. تاریخچه‌ی روزانه
- **هدف:** یک ردیف برای هر روز.
- **آدرس:** `GET /ClosingPrice/GetClosingPriceDailyList/{insCode}/{n}` — معنی `n` تأییدنشده؛ با `0` کل تاریخچه‌ی این قرارداد آمد.
- **curl:**
  ```
  curl -s "https://cdn.tsetmc.com/api/ClosingPrice/GetClosingPriceDailyList/56490941862491615/0" -o docs/samples/daily_list_56490941862491615_0.json
  ```
- **نمونه‌ی خام:** `docs/samples/daily_list_56490941862491615_0.json` (ضخود9060، اختيارخ خودرو-1200-1405/09/04)
- **ساختار:** `{"closingPriceDaily": [ ... ]}`، **جدیدترین روز اول**. فیلدها مثل بخش 4.
- **رفتار:**
  - جدیدترین ردیف 20261103، قدیمی‌ترین 20260826.
  - پنج‌شنبه/جمعه (مثلاً 20261029–30) ردیف ندارند.
  - روزی با `zTotTran: 0` هم ردیف دارد (20260826) → وجود ردیف ≠ روز معامله‌دار.
- **سؤال باز:** برای سهم پایه (مثلاً خودرو) تاریخچه تا 1402 عقب می‌رود؟ `n` غیر 0 چه می‌کند؟

---

## 6. دیده‌شده، فعلاً استفاده نمی‌شود
- `GET /ClientType/GetClientType/{insCode}/1/0` — خرید/فروش حقیقی و حقوقی امروز. نمونه: `docs/samples/client_type_31459957843277944.json`
- `GET /MarketData/GetMarketOverview/1` — شاخص و وضعیت بازار (`marketStateTitle`: بسته). نمونه: `docs/samples/market_overview_1.json`
- در Console دیده شد (آدرس کامل ثبت نشده): `BestLimits`، `GetRelatedCompany`.

---

## 7. محدودیت‌ها و بلاک
- 1405/08/12 حدود 20:08، بعد از حدود 20 درخواست در 40 دقیقه (آخرین‌ها با فاصله‌ی 3 ثانیه)، چند endpoint به‌جای JSON یک صفحه‌ی HTML برگرداندند: «دسترسی شما بدلیل مسائل امنیتی مسدود شده است». نمونه: `docs/samples/blocked/`.
- حدود 30 دقیقه بعد دوباره جواب داد → بلاک موقت بود.
- **نتیجه برای خزنده:**
  - پاسخ HTML یا متن «مسدود» = بلاک → توقف و اطلاع به منتور.
  - نرخ باید خیلی کمتر از 1 درخواست در 3 ثانیه باشد (نرخ دقیق: از منتور پرسیده شود).
  - بدنه‌ی HTTP 200 هم ممکن است بلاک باشد → Content-Type و محتوا چک شود.

---

## 8. جواب سؤال‌های PDF (با شاهد)
1. **شناسه‌ی یکتا:** `insCode`. نماد کافی نیست (تأیید کامل: آیا نماد بعد از سررسید دوباره استفاده می‌شود؟ — باز). ISIN هم کافی نیست: `instrumentID` ≠ `cIsin` و در جستجو `null` است. شاهد: بخش 1 و 2.
2. **call/put، strike، سررسید، پایه:** از بخش 3: `strikePrice`، `endDate`، `uaInsCode`. call با «ض» و «اختيارخ»، put با «ط» (put هنوز نمونه ندارد — باز).
3. **فعال و منقضی:** جستجو هر دو را می‌دهد؛ `lastDate` احتمالاً وضعیت است (تأییدنشده). منبع فهرست کامل همه‌ی اختیارها — **هنوز پیدا نشده**.
4. **فیلدهای روزانه و تیک:** روزانه در بخش 5. تیک — **هنوز پیدا نشده**. عمق تاریخچه — باز.
5. **تطبیق تیک با روزانه:** باز (به endpoint تیک نیاز است). فیلدهای مقایسه: `zTotTran`، `qTotTran5J`، `qTotCap`، `priceFirst/Min/Max`، `pDrCotVal`.
6. **روز غیرمعاملاتی:** نبود ردیف در تاریخچه، یا ردیف با `zTotTran = 0`. شاهد: بخش 5. تعطیلی رسمی وسط هفته — باز.

---

## 9. کارهای مانده
- [ ] endpoint تیک (ریز معاملات یک روز)
- [ ] منبع فهرست همه‌ی اختیارها و دارایی‌های پایه
- [ ] نمونه‌ی put و نمونه‌ی منقضی برای بخش 3
- [ ] تاریخچه‌ی روزانه‌ی یک سهم پایه
- [ ] معنی `lastDate`، `n` در DailyList، `buyOP/sellOP`
