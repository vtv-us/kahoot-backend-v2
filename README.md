# kahoot-backend

## Requirements

- Install [Docker](https://docs.docker.com/engine/install/)
  - Remember to [add Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user) when using Linux
- Install [Docker compose](https://docs.docker.com/compose/install/linux/) when using Linux

## Run on locally

Clone the project

```bash
git clone https://github.com/vtv-us/kahoot-backend.git
cd kahoot-backend
```

Run service

```bash
docker compose up -d
```
Run after pull for remove old images

```bash
docker compose up -d --build
docker images -a | grep none | awk '{ print $3; }' | xargs docker rmi
```
