# Changelog

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
