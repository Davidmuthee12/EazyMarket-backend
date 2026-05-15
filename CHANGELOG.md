# Changelog

## [1.8.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.7.0...v1.8.0) (2026-05-15)


### Features

* **api:** add rate limiter middleware ([003a247](https://github.com/Davidmuthee12/EazyMarket-backend/commit/003a2472c48f73598838b0f603d7fab6f2d9f88b))
* **api:** allow storefront subdomain origins ([6070ba2](https://github.com/Davidmuthee12/EazyMarket-backend/commit/6070ba255b14b6ff3889ddecb36bd94191699ff4))
* **api:** configure rate limiter ([3cb8601](https://github.com/Davidmuthee12/EazyMarket-backend/commit/3cb860101920dd24add4c53a5d0f1581ce83db91))
* **api:** mount rate limiter middleware ([00e0e9f](https://github.com/Davidmuthee12/EazyMarket-backend/commit/00e0e9f92d150f7fb3d0e2312fb7e9c3bb775cae))
* **api:** validate vendor subdomains ([dc97cd5](https://github.com/Davidmuthee12/EazyMarket-backend/commit/dc97cd510652efa2c090d6f383d43f305e081b1c))
* **db:** harden vendor applications ([3613390](https://github.com/Davidmuthee12/EazyMarket-backend/commit/36133903c6c6ed53b6e3f3e61c44dbec0bb529f8))
* **ratelimiter:** add fixed window limiter ([e9cefd6](https://github.com/Davidmuthee12/EazyMarket-backend/commit/e9cefd65bf26d6cb8d3d517ada85a6b72b38ce41))
* **ratelimiter:** define limiter contract ([44ea35c](https://github.com/Davidmuthee12/EazyMarket-backend/commit/44ea35cfb9a402ee763e6065cab821e8f7f8b6f2))
* **store:** add pagination query parser ([1d2e867](https://github.com/Davidmuthee12/EazyMarket-backend/commit/1d2e867ce5db4f33cb4714238e5ed187a1871b9c))
* **users:** activate storefronts on vendor approval ([9ad55f3](https://github.com/Davidmuthee12/EazyMarket-backend/commit/9ad55f3af9375bb26d5085a707b7cc28f948c683))
* **users:** collect storefront details for vendor applications ([ac2d01b](https://github.com/Davidmuthee12/EazyMarket-backend/commit/ac2d01b7439e63d70b9754e75f0e96da0e1c0e4a))


### Bug Fixes

* **vendor:** map duplicate profile errors ([f8eb6ae](https://github.com/Davidmuthee12/EazyMarket-backend/commit/f8eb6aefe1772a228ef440e594f4a5228b08b271))
* **vendor:** normalize profile subdomains ([0113d18](https://github.com/Davidmuthee12/EazyMarket-backend/commit/0113d1839fa149753c9196e89fd77c13122e1d21))

## [1.7.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.6.0...v1.7.0) (2026-05-13)


### Features

* **api:** adds new api routes ([3d2025d](https://github.com/Davidmuthee12/EazyMarket-backend/commit/3d2025d221d86ee5b1b0b58471e708e1cb84567b))
* **api:** adds new api routes for wishlists ([4947a71](https://github.com/Davidmuthee12/EazyMarket-backend/commit/4947a71a5371ff58df6aa40c0f517ce5ac32e5c2))
* **api:** adds new status update handler ([18537f0](https://github.com/Davidmuthee12/EazyMarket-backend/commit/18537f0fd3d11a85f80109984813c828ed99595c))
* **api:** adds wishlists handlers ([f1bd091](https://github.com/Davidmuthee12/EazyMarket-backend/commit/f1bd091430d976fc34e49ff6a3b9b25f724a76f2))
* **api:** mount storefront-scoped routes ([c17662b](https://github.com/Davidmuthee12/EazyMarket-backend/commit/c17662bcc5491e9e9c8aaf4bb544e6f6096259f4))
* **api:** resolve storefront vendor context ([0852211](https://github.com/Davidmuthee12/EazyMarket-backend/commit/0852211d63d284c415b0366fe8b96e079428cbc9))
* **cart:** scope cart handlers to storefront ([b5ac5ae](https://github.com/Davidmuthee12/EazyMarket-backend/commit/b5ac5ae27ceea81738b948a906a5aff8886f8ce3))
* **cart:** scope cart storage by vendor ([226fa79](https://github.com/Davidmuthee12/EazyMarket-backend/commit/226fa791bb0ef49ede8e00f8933b13804e380110))
* **db:** scope storefront data by vendor ([a926b99](https://github.com/Davidmuthee12/EazyMarket-backend/commit/a926b999bf637bad379112329e8c991d765b5bd7))
* **middleware:** adds new middleware to prevent suspended vendor accounts ([5fa7be4](https://github.com/Davidmuthee12/EazyMarket-backend/commit/5fa7be4181c34345d615325c4dadc9d9538c21c3))
* **migration:** adds new migration for wishlist table ([35012eb](https://github.com/Davidmuthee12/EazyMarket-backend/commit/35012eb8ff8c029f05348bcb6ff0daa965053687))
* **migration:** adds new migration to alter users table with status column ([2649f61](https://github.com/Davidmuthee12/EazyMarket-backend/commit/2649f612b74edf2cf59bd7fa2a6a4517bd4cf1eb))
* **orders:** filter customer orders by vendor ([081639c](https://github.com/Davidmuthee12/EazyMarket-backend/commit/081639cd62b45f501360b589d0c3b66398246f54))
* **orders:** scope customer orders to storefront ([d106462](https://github.com/Davidmuthee12/EazyMarket-backend/commit/d1064626d6f3bdea1802fbe54955c9c9e7cee4b1))
* **products:** add storefront product queries ([6d756e4](https://github.com/Davidmuthee12/EazyMarket-backend/commit/6d756e4cfa91fb070affc4c7ed3c662a964d5e2d))
* **products:** allow vendor product statuses ([f739b74](https://github.com/Davidmuthee12/EazyMarket-backend/commit/f739b74ed5b913fd78d15d9f07e1ef047166e357))
* **redis-cache:** ensures suspended accounts are unathorized immediately in cache ([64c5b05](https://github.com/Davidmuthee12/EazyMarket-backend/commit/64c5b05562978460dd1af06a3e109d87685eeef5))
* **store:** adds new status update function for vendors ([9205924](https://github.com/Davidmuthee12/EazyMarket-backend/commit/9205924c326bd492302867d0b85f255e673de350))
* **store:** adds new store function for status update ([b7eb259](https://github.com/Davidmuthee12/EazyMarket-backend/commit/b7eb2595f2752f3b2a7697d1268ad707e5ea9807))
* **store:** adds wishlist db level queries ([0570380](https://github.com/Davidmuthee12/EazyMarket-backend/commit/0570380bf82d2c4e3c6c1a2e644084a62e461e3f))
* **storefront:** add public storefront handlers ([f35b029](https://github.com/Davidmuthee12/EazyMarket-backend/commit/f35b02972408c4b8ddd531a8dba0ae95317f753e))
* **store:** wires status update handler to store ([883900a](https://github.com/Davidmuthee12/EazyMarket-backend/commit/883900a901dc110fe071d4f4983e6a8659b8e215))
* **store:** wires wishlist api level to store level ([0b271fe](https://github.com/Davidmuthee12/EazyMarket-backend/commit/0b271fe6094a87f62e89b06595cb7c34c3f6dd90))
* **vendor:** lookup storefronts by subdomain ([c455ec5](https://github.com/Davidmuthee12/EazyMarket-backend/commit/c455ec5e893c01e04b7254849fd16d149e2838ae))
* **wishlist:** filter wishlist by vendor ([701f68c](https://github.com/Davidmuthee12/EazyMarket-backend/commit/701f68c134a93133c1958c300fca351cdbed55ba))
* **wishlist:** scope wishlist handlers to storefront ([675bbfa](https://github.com/Davidmuthee12/EazyMarket-backend/commit/675bbfa183bf45977862cae7fd52e906967b17b1))


### Bug Fixes

* **lint:** fixes lint errors ([f88ab43](https://github.com/Davidmuthee12/EazyMarket-backend/commit/f88ab4341326f480a43aceb996e563fa4b19131a))
* **vendor:** handle duplicate subdomains ([6c501e7](https://github.com/Davidmuthee12/EazyMarket-backend/commit/6c501e768b14e28b579d5f3d1b8b52272d2f8b5f))

## [1.6.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.5.0...v1.6.0) (2026-05-10)


### Features

* **api:** add cart endpoints ([c6b90ff](https://github.com/Davidmuthee12/EazyMarket-backend/commit/c6b90fff15e08b32fa0e5dc9295b214212ed1234))
* **api:** add order endpoints ([bae8f99](https://github.com/Davidmuthee12/EazyMarket-backend/commit/bae8f996edbc1836697c2e911547cc8f18a8fde0))
* **api:** adds new delete product handler ([e25f2e5](https://github.com/Davidmuthee12/EazyMarket-backend/commit/e25f2e5a7c4f14c48fc80cd714c20783c6bdc91e))
* **api:** adds new vendor delete product endpoint ([1387142](https://github.com/Davidmuthee12/EazyMarket-backend/commit/1387142ae6389bce6cc268da0a485ba0becfd247))
* **cart:** add cart store operations ([c4bc151](https://github.com/Davidmuthee12/EazyMarket-backend/commit/c4bc1515f8bb16ae3aa29ff6dbeb0bccb73ced21))
* **db:** add cart persistence migrations ([44c3888](https://github.com/Davidmuthee12/EazyMarket-backend/commit/44c38883d5e91a0800f11165360c0beb9afeff30))
* **db:** add order persistence migrations ([47cee4a](https://github.com/Davidmuthee12/EazyMarket-backend/commit/47cee4a0c896f5efb3433369663859550dde5008))
* **DB:** adds db query to delete vendor products ([569d36a](https://github.com/Davidmuthee12/EazyMarket-backend/commit/569d36a0f3b1a33fc681517dbb1c26574b08dc77))
* **orders:** add order store operations ([45e7c76](https://github.com/Davidmuthee12/EazyMarket-backend/commit/45e7c76f651dcd68093f23a401309cc85dd65daa))
* **store:** registers the delete handler in store ([952490e](https://github.com/Davidmuthee12/EazyMarket-backend/commit/952490e12e04070b4457959eacdab2677ebc71d4))


### Bug Fixes

* **api:** handle nil internal server errors ([725a81d](https://github.com/Davidmuthee12/EazyMarket-backend/commit/725a81d0ae17176ffc853650eb3891d2cbf91f0b))

## [1.5.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.4.0...v1.5.0) (2026-05-09)


### Features

* **api:** add vendor products route ([11c9edc](https://github.com/Davidmuthee12/EazyMarket-backend/commit/11c9edca654a5351f8572490d4f7313cce74a082))
* **api:** expose product detail route ([cf009a4](https://github.com/Davidmuthee12/EazyMarket-backend/commit/cf009a4115e1ac548789eab09af772a954e9e38e))
* **api:** expose product update route ([e412874](https://github.com/Davidmuthee12/EazyMarket-backend/commit/e412874022205e1e545ff7a105f43bb4380ede44))
* **api:** expose vendor product listing route ([bff40e0](https://github.com/Davidmuthee12/EazyMarket-backend/commit/bff40e03d13e34a1036433bdcf0bb1aa436b9ed9))
* **db:** add product tables migration ([385950f](https://github.com/Davidmuthee12/EazyMarket-backend/commit/385950f3fb2be13a5561c471f19165e75a03c9d9))
* **products:** add create product handler ([40f1459](https://github.com/Davidmuthee12/EazyMarket-backend/commit/40f1459ecdad5946727a1f08a7fda08ed15fea74))
* **products:** add product store creation ([027363f](https://github.com/Davidmuthee12/EazyMarket-backend/commit/027363f448e4ec3df05b30ece7b050a7060ed416))
* **products:** fetch product details by id ([2ff36ef](https://github.com/Davidmuthee12/EazyMarket-backend/commit/2ff36ef4ff004c6ffbbb4c515fb7150f564413f1))
* **products:** handle vendor product listing ([52d38bf](https://github.com/Davidmuthee12/EazyMarket-backend/commit/52d38bfd557e46b74b80d216068649fbdc1d0a34))
* **products:** handle vendor product updates ([2ccdc12](https://github.com/Davidmuthee12/EazyMarket-backend/commit/2ccdc12874591f7f0c08a68e4dca6badc5bc6528))
* **products:** update vendor-owned products ([5e31c98](https://github.com/Davidmuthee12/EazyMarket-backend/commit/5e31c980bd7c38c638bc82c99d89155f36c3ae79))
* **store:** expose product detail lookup ([d7dec55](https://github.com/Davidmuthee12/EazyMarket-backend/commit/d7dec55b9d4792b016e99d7a454236ab045c3ad0))
* **store:** expose product listing contract ([2c1175a](https://github.com/Davidmuthee12/EazyMarket-backend/commit/2c1175a80d58474b81d4d83e2021f2215d0f96a4))
* **store:** expose product update contract ([3b143b0](https://github.com/Davidmuthee12/EazyMarket-backend/commit/3b143b06790a24871a793dd90e8fe9190bb52f15))
* **store:** register product storage ([eafb608](https://github.com/Davidmuthee12/EazyMarket-backend/commit/eafb608b08dcbd26c912bf1396afccab72554793))


### Bug Fixes

* **products:** query vendor products correctly ([58a9811](https://github.com/Davidmuthee12/EazyMarket-backend/commit/58a9811aea11c1b5c97aa0938f86924712c693ef))

## [1.4.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.3.0...v1.4.0) (2026-05-08)


### Features

* **api:** add category HTTP handlers in category.go ([b4a93a6](https://github.com/Davidmuthee12/EazyMarket-backend/commit/b4a93a6a8b07282456e6d048084f7c5f2d44a4f4))
* **api:** add getCategoriesHandler, deleteCategoryHandler, and updateCategoryHandler ([e6838f8](https://github.com/Davidmuthee12/EazyMarket-backend/commit/e6838f8c9e08f1aa9ed2507fd2995a22856bcde8))
* **api:** expand vendor handler responses in vendor.go ([69e1f90](https://github.com/Davidmuthee12/EazyMarket-backend/commit/69e1f90c9b17ca3496bb9ebc87a1f20e9ef95469))
* **router:** register GET, DELETE, and PUT routes for /admin/categories ([4794da4](https://github.com/Davidmuthee12/EazyMarket-backend/commit/4794da4dd58d978c75d630da3f8a1e99df5355c8))
* **store:** add GetCategories, DeleteCategory, and UpdateCategory store methods ([32038ce](https://github.com/Davidmuthee12/EazyMarket-backend/commit/32038ce0397ed8dde52762c029ad5ac9dffc97ba))
* **store:** add vendor profile retrieval in vendor.go ([a940f1a](https://github.com/Davidmuthee12/EazyMarket-backend/commit/a940f1a58ce2a8b12c89bf263148c253048d81d4))
* **store:** extend Category interface with GetCategories, DeleteCategory, and UpdateCategory ([1a42e8c](https://github.com/Davidmuthee12/EazyMarket-backend/commit/1a42e8c908b3feb1375ade0046d15570cf876dfd))
* **store:** implement category store operations in category.go ([7b3fe07](https://github.com/Davidmuthee12/EazyMarket-backend/commit/7b3fe074a2d8ad1ea27214fef79130a693f42af6))

## [1.3.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.2.0...v1.3.0) (2026-05-07)


### Features

* **api:** add vendorProfileHandler to create vendor profile ([301b6d1](https://github.com/Davidmuthee12/EazyMarket-backend/commit/301b6d1b9eebd82d855366ea80a71c9a979d25df))
* **api:** register POST /vendor/profile route ([5ec639e](https://github.com/Davidmuthee12/EazyMarket-backend/commit/5ec639e8a08561e5f7a06013ef1cc20702d2fb21))
* **migrations:** add migration 007 to create vendor_profiles table ([8419277](https://github.com/Davidmuthee12/EazyMarket-backend/commit/841927750f1eae20076ea6e195deb9b120c0979a))
* **store:** add Vendor model and VenderStore with CreateVendorProfile ([2184d95](https://github.com/Davidmuthee12/EazyMarket-backend/commit/2184d95b8f4e246b1e1ebaac6d10727860683c21))


### Bug Fixes

* **store:** wire VenderStore into NewStorage ([67677a6](https://github.com/Davidmuthee12/EazyMarket-backend/commit/67677a6e80cfd27861dc80462238e4dc9ad41f3d))

## [1.2.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.1.0...v1.2.0) (2026-05-06)


### Features

* add reject vendor-request route in admin router ([546102b](https://github.com/Davidmuthee12/EazyMarket-backend/commit/546102bfa755d622636aa093a5cd4526cea9a4d1))
* add RejectRequest store method for role upgrade requests ([e7625db](https://github.com/Davidmuthee12/EazyMarket-backend/commit/e7625db64481cf3013f3c3904a4936d948a05371))
* Add reviewer tracking to role upgrade requests with updateRequestTable function ([5bab53c](https://github.com/Davidmuthee12/EazyMarket-backend/commit/5bab53c63964f69193006e26c33a5f609b59c4d2))
* Extract and pass authenticated admin reviewer to approveVendorHandler for audit tracking ([19c0042](https://github.com/Davidmuthee12/EazyMarket-backend/commit/19c0042f8ec4d00f10ea2c156272910b5f183dd1))
* implement rejectVendorHandler with reviewer context ([e1557b9](https://github.com/Davidmuthee12/EazyMarket-backend/commit/e1557b9f3c766c2b9ce1d31544a6199e61d98af5))

## [1.1.0](https://github.com/Davidmuthee12/EazyMarket-backend/compare/v1.0.0...v1.1.0) (2026-05-04)


### Features

* **api:** add user HTTP handlers for get, activate, and list all users ([14fcae6](https://github.com/Davidmuthee12/EazyMarket-backend/commit/14fcae644939a185869f09b57233b5320752c52b))
* **api:** wire Redis cache storage and register user routes with auth middleware ([de94074](https://github.com/Davidmuthee12/EazyMarket-backend/commit/de9407460c576a369e6864b4a418bd9692eba257))
* **cache:** add Redis-backed cache layer for user storage with get/set operations ([653d2c5](https://github.com/Davidmuthee12/EazyMarket-backend/commit/653d2c5ad6c25f3fce9e9de80c05bffdf3289d71))
* **middleware:** implement JWT auth middleware with cache-aware user resolution ([de971d5](https://github.com/Davidmuthee12/EazyMarket-backend/commit/de971d55adbc16051f8123496bae02815c961401))
* **store:** add GetAllUsers query joining users and roles tables ([9dd33a9](https://github.com/Davidmuthee12/EazyMarket-backend/commit/9dd33a9e0ca08dcf5d36e1847570d1b8f393eb16))
* **store:** add GetAllUsers to Users store interface ([48e0a84](https://github.com/Davidmuthee12/EazyMarket-backend/commit/48e0a841d77d0924721cd98edddbc2e17ae3bdfd))


### Bug Fixes

* **api:** wire Redis config and conditional cache initialization ([3395513](https://github.com/Davidmuthee12/EazyMarket-backend/commit/3395513593d41716aeff38e8c051cb92feb0a443))

## 1.0.0 (2026-05-03)


### Features

* bootstrap Go API with auth, storage, mailer, and migrations ([527e3e1](https://github.com/Davidmuthee12/EazyMarket-backend/commit/527e3e1a531da60159b8e52a252c57b910cefbbe))


### Bug Fixes

* remove unused mailTrap field and mailTrapConfig type ([efc7c6e](https://github.com/Davidmuthee12/EazyMarket-backend/commit/efc7c6e25efda2985700230fbf618104c68d6e56))
* suppress U1000 on reserved error response helpers ([ca2b478](https://github.com/Davidmuthee12/EazyMarket-backend/commit/ca2b4784feb2aed110462d49d5f33cd61db40ce6))
* suppress U1000 on RoleStore.db pending role query implementation ([dc8097d](https://github.com/Davidmuthee12/EazyMarket-backend/commit/dc8097d2987aac09077db7170e1d5abf169ee107))
* use version const in startup log to satisfy staticcheck ([414a8bb](https://github.com/Davidmuthee12/EazyMarket-backend/commit/414a8bbdfab99f395ce70e6fee9a7a15b5cef653))
