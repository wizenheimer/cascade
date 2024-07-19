# Dockerfile Readme

This readme provides an overview of the Dockerfile for directory.

## Purpose
The Dockerfile is used to define the instructions for building a Docker image. It specifies the base image, installs dependencies, sets environment variables, and configures the container.

We utilize a multi-stage build for both development and production environments. The development stage includes auto-reload functionality backed by air, allowing for seamless code changes during development. Additionally, we have a Docker Compose configuration that includes a PostgreSQL container for database management.



## Usage
To build the Docker image for development, navigate to the directory containing the Dockerfile and run the following command:

```bash
docker build --target development -t cascade:dev .
```

To build the Docker image for production, use the following command:

```bash
docker build --target production -t cascade:prod .
```

To use the docker compose

```bash
docker compose up -d --build
```

For compose teardown

```bash
docker compose down -v
```

## Customization
Feel free to modify the Dockerfile to suit your specific needs. You can add additional dependencies, configure ports, or include any other necessary instructions.

## Maintenance
Please ensure that the Dockerfile is kept up to date with any changes or updates to your application. Regularly review and test the Dockerfile to ensure it builds the desired image correctly.
