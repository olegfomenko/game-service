---
    title: Game service for Blockchain Hackathon 2021
    author: olegfomenko 
    date: 30 March 2021 
---


# Game service

game-service implements main logic of creating game, finishing game and user donating.

Main endpoint logic description:
1. Crate game -
       
   will crate a 'GAM' asset that equal's to 1 donated (from an organizer or user) 'USD'. 'GAM' asset will store 
   team player's list, game name, date and stream url in asset detail field. 'GAM' asset is not transferable.
   
   After creating asset the count, equals to donated from organizer USD will be issued, 
   and organizers USD will be moved to admin account. This donation will be a prize for winner team.
   
2. Pay game -

    users can donate to prize bank. That donation causes new 'GAM's issuance and 
    transferring USD from user to admin account.

3. Select winner - 

    selecting winner team and transferring parts of prize to player accounts. Also will update GAM asset details.
   
   
## Configuration

```yaml
log:
  disable_sentry: true
  level: debug

# configurate service port
listener:
  addr: :80

cop:
  endpoint: http://cop
  upstream: http://game-service
  service_name: "game-service"
  service_port: "80"
  service_prefix: "/integrations/game-service/"

horizon:
  endpoint: "http://traefik"
  # master account seed
  signer: SAMJKTZVW5UOHCDK5INYJNORF2HRKYI72M5XSZCBYAHQHR34FFR4Z6G4
  # master account id
  source: GBA4EX43M25UPV4WIE6RRMQOFTWXZZRIPFAI5VPY6Z2ZVVXVWZ6NEOOB
```

docker-composer for running (edited from TokenD dev edition)
```yaml
version: '3.3'

services:
  traefik:
    image: traefik:v2.0
    ports:
      - "80:80"
      - "8081:8080"
    volumes:
      - ./configs/traefik.yaml:/traefik.yaml
  cop:
    image: tokend/traefik-cop:1.0.0
    restart: unless-stopped
    environment:
      - KV_VIPER_FILE=/config.yaml
    volumes:
      - ./configs/cop.yaml:/config.yaml
    entrypoint: sh -c "traefik-cop run"
  upstream:
    image: nginx
    restart: unless-stopped
    volumes:
    - ./configs/nginx.conf:/etc/nginx/nginx.conf
    ports:
    - "8000:80"
  adks:
    image: tokend/adks:1.0.2
    restart: unless-stopped
    depends_on:
      - horizon
      - adks_db
    ports:
      - 8006:80
    volumes:
      - ./configs/adks.toml:/config.toml
  adks_db:
    image: tokend/postgres-ubuntu:9.6
    restart: unless-stopped
    environment:
      - POSTGRES_USER=adks
      - POSTGRES_PASSWORD=adks
      - POSTGRES_DB=adks
      - PGDATA=/pgdata
    volumes:
      - adks-data:/pgdata
  redis:
    image: redis:5.0-alpine
    restart: unless-stopped
    volumes:
      - redis-data:/data
    command:
      - redis-server
      - --appendonly
      - "yes"
  core:
    image: tokend/core:3.7.0-x.10
    depends_on:
      - traefik
    restart: unless-stopped
    environment:
      - POSTGRES_USER=core
      - POSTGRES_PASSWORD=core
      - POSTGRES_DB=core
      - PGDATA=/data/pgdata
      - ENSUREDB=1
      - CONFIG=/core-config.ini
    volumes:
      - ./configs/core.ini:/core-config.ini
      - core-data:/data
    labels:
      - "autoheal=true"
  horizon:
    image: tokend/horizon:3.10.0-x.10
    depends_on:
      - core
    restart: unless-stopped
    environment:
      - POSTGRES_USER=horizon
      - POSTGRES_PASSWORD=horizon
      - POSTGRES_DB=horizon
      - PGDATA=/data/pgdata
    volumes:
      - ./configs/horizon.yaml:/config.yaml
      - horizon-data:/data
  api:
    image: tokend/identity:4.7.0-rc.0
    restart: unless-stopped
    depends_on:
      - horizon
      - api_db
    environment:
      - KV_VIPER_FILE=/config.yaml
    volumes:
      - ./configs/api.yml:/config.yaml
    command: run
  api_db:
    image: tokend/postgres-ubuntu:9.6
    restart: unless-stopped
    environment:
      - POSTGRES_USER=api
      - POSTGRES_PASSWORD=api
      - POSTGRES_DB=api
      - PGDATA=/pgdata
    volumes:
      - api-data:/pgdata
  initscripts:
    image: tokend/terraform-provider-tokend:1.3.4
    depends_on:
      - horizon
      - storage
    restart: on-failure
    volumes:
      - ./terraform:/opt/config
    entrypoint: ""
    command: /opt/config/apply.sh
  admin_client:
    image: tokend/admin-client:1.14.0-rc.0
    restart: unless-stopped
    volumes:
      - ./configs/client.js:/usr/share/nginx/html/static/env.js
    ports:
      - 8070:80
  storage:
    image: minio/minio:RELEASE.2019-01-31T00-31-19Z
    restart: unless-stopped
    entrypoint: "sh"
    command: -c "mkdir -p /data/tfstate && minio server /data"
    environment:
      - MINIO_ACCESS_KEY=miniominio
      - MINIO_SECRET_KEY=sekritsekrit
      - MINIO_BROWSER=off
    volumes:
      - storage-data:/data

  game-service:
    image: olegfomenko2002/bua2021:game-service
    restart: unless-stopped
    environment:
      - KV_VIPER_FILE=/config.yaml
    volumes:
      - ./configs/game-service.yaml:/config.yaml
    command: ["run", "service"]

  web_client:
    image: olegfomenko2002/bua2021:game-web-client
    restart: unless-stopped
    volumes:
      - ./configs/client.js:/usr/share/nginx/html/static/env.js
    ports:
      - 8060:80
  
volumes:
  adks-data:
  api-data:
  core-data:
  horizon-data:
  storage-data:
  redis-data:
```
Also you need to add some rules to corporate user account by adding to docker-compose 
```yaml
   init-rule:
       image: olegfomenko2002/bua2021:init-rule 
```

and launch this only once.


  