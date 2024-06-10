## [0.7.2](https://github.com/PUFA-Computing/backend/compare/v0.7.1...v0.7.2) (2024-06-10)


### Bug Fixes

* documentation of event service ([9f7cf72](https://github.com/PUFA-Computing/backend/commit/9f7cf725cb544431d1ba7568e44b57ff5516965b))

## [0.7.1](https://github.com/PUFA-Computing/backend/compare/v0.7.0...v0.7.1) (2024-06-04)


### Bug Fixes

* bugfix not return the profile picture ([6f0cb68](https://github.com/PUFA-Computing/backend/commit/6f0cb687467a383319da5b8cd5c3608ed76b4f66))

# [0.7.0](https://github.com/PUFA-Computing/backend/compare/v0.6.2...v0.7.0) (2024-05-25)


### Bug Fixes

* Add return statement after email verification ([525049d](https://github.com/PUFA-Computing/backend/commit/525049d20cf76e86bd910c84934d7f78bf5ef06f))
* after authenticate for email verification ([37cb828](https://github.com/PUFA-Computing/backend/commit/37cb82836b8e8b6a0d999646eb1c1312b60241fb))
* bug on login not return error if email not valid ([a390119](https://github.com/PUFA-Computing/backend/commit/a39011915c3b8d8e09010d5bd3773f71bb12fb43))
* **ci/cd:** fix env for staging ([34202cc](https://github.com/PUFA-Computing/backend/commit/34202cc8664abcd6b098cd9707b888effb1f8337))
* **cors:** cors for production added ([3f1fa2b](https://github.com/PUFA-Computing/backend/commit/3f1fa2b518c9adf949546cc9bf24cffd3745ab42))
* **cors:** fix cors to staging ([8ece9a6](https://github.com/PUFA-Computing/backend/commit/8ece9a6dac1e79154e66d7ef9e1d782892fccc2b))
* mail template url ([0a63c4a](https://github.com/PUFA-Computing/backend/commit/0a63c4a19b446cba0ef300b0f2d360ab828318b4))
* **news:** fixing get news unhandle 7 by 8 ([90c0869](https://github.com/PUFA-Computing/backend/commit/90c0869725d3bd02d203b99662ff490d0d4d702a))
* **port:** fix api port ([fc64f1c](https://github.com/PUFA-Computing/backend/commit/fc64f1ce97d244b80da86522af452cb1524ee9ef))
* **port:** workflows fix for api port ([43ca8ac](https://github.com/PUFA-Computing/backend/commit/43ca8ac9a68dc53d9291511fab96fee02d239bbf))
* **register event:** adding *int to achieve nullable data ([be2e0b8](https://github.com/PUFA-Computing/backend/commit/be2e0b864e08312efe5a8271b77c8c5894a9c663))
* Return error more clear ([bdde16e](https://github.com/PUFA-Computing/backend/commit/bdde16ee3828ce004f9f7cac2a34c2f029879ec7))
* Return error more clear ([92c6974](https://github.com/PUFA-Computing/backend/commit/92c697416e31261e2dfa478727168fa80e634640))
* **s3:** Modify DeleteEvent and GetFileR2 functions ([9c85e27](https://github.com/PUFA-Computing/backend/commit/9c85e277cf65fc2edc32a7fc0a8be2bec257497d))
* **staging:** staging not the latest build ([3eb7893](https://github.com/PUFA-Computing/backend/commit/3eb789389031ace065d6f87f238546149b77f528))


### Features

* **CI/CD:** production and staging workflows ([dbf5a5e](https://github.com/PUFA-Computing/backend/commit/dbf5a5edb473fbd044d7d4f88154a5486c599d17))
* **event status:** make a new goroutine for event status updater every 5 minut ([dfd520c](https://github.com/PUFA-Computing/backend/commit/dfd520c04a6a18ae283b58b7d6f5a37ac9feddba))
* **list events:** add total page for pagination ([02eb520](https://github.com/PUFA-Computing/backend/commit/02eb52079ff71580fc2e9bd24dc048839b82ecc9))
* **news:** adding thumbnail and slug on the news ([44103d1](https://github.com/PUFA-Computing/backend/commit/44103d1ab274cd8dc3ad63b6b49d6686331d43f9))
* **register event:** adding limit registration ([4888766](https://github.com/PUFA-Computing/backend/commit/488876678f3b83a97d089ae02d003ebe5b0428f0))
* **register:** validation, and check exists ([9bdbb69](https://github.com/PUFA-Computing/backend/commit/9bdbb69ebb9ca802f6997555258fdbf33e379291))
* **s3:** adding s3 bucket storage ([39fbbb2](https://github.com/PUFA-Computing/backend/commit/39fbbb21a8053f49c8ff4b76acbe14ce2f60114a))
* **s3:** integrate with handler ([97ddc29](https://github.com/PUFA-Computing/backend/commit/97ddc294ccf0f76d59c2c0b4ca495b9e082b2495))
* **uploads:** implement storage Distribution in Event Handlers ([33d8118](https://github.com/PUFA-Computing/backend/commit/33d81182d1f976fe8ada1cfa25a9b2e0a6df1a8b))
* Version fetch from github and save to db ([5f0fdc4](https://github.com/PUFA-Computing/backend/commit/5f0fdc4e4daea5e5bb41d2c65bd01d742619e667))
