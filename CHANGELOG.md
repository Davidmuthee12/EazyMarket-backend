# Changelog

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
