Weather Monster Project
=======================

The project is going to be used by different institutions to monitor weather temperature and get forecasts.

1. Setting up the project
- clone the project `git clone https://github.com/serega-cpp/tt-wm.git`

2. Running the tests
- prepare config file for test database (using `config.yaml` as an example)
- run `go test -cfg config_test.yaml`

3. Running the app
- prepare config file for production database (using `config.yaml` as an example)
- build the binary `go build -o finleap`
- run `./finleap -cfg config.yaml` 
