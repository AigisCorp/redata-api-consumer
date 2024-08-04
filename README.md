<h1 align="center">REData API Consumer</h1>

- [REData API Consumer](#redata-api-consumer)
  - [Introduction](#introduction)
  - [Documentation](#documentation)
  - [Contributing](#contributing)
  - [License](#license)

## Introduction

REData API Consumer is a tool that process the daily information provided by [red eléctrica](https://www.ree.es/es/apidatos).

Its purpose is to provide an API service for being able to take better decissions on residential power consumption, even for self-consumption.

## Documentation

You can run 'redata-api-consumer' with a docker container:

```shell
docker build . --tag redata-api-consumer

docker run --rm -p 8080:8080 -e location=Europe/Madrid redata-api-consumer
```

Check the [Docs](docs/README.md) to explore the endpoints.

## Contributing

We welcome contributions from the community. If you have ideas for improvements, feature requests, or bug reports, please open an issue in this repository. If you'd like to contribute code, please fork the repository, make your changes, and submit a pull request.

Please review our [Contribution Guidelines](CONTRIBUTING) for more information on how to get involved.

## License

This project is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for more information.

