# Rand

Rand ("edge" in dutch) is a simple Golang frontend deployment tool. It is being
used to deploy the [Pyrite](https://github.com/garage44/pyrite) frontend with.
All it does is unpacking an uploaded zip behind a basic auth endpoint to
a specified location.

## Usage

* Build rand

  ```bash
  git clone git@github.com:garage44/rand.git
  cd rand
  go build -o rand .
  ```

* Define a configuration

  ```bash
  vim /home/galene/.randrc
  ```

  ```bash
  RAND_PATH=/srv/http/pyrite
  RAND_LISTEN=127.0.0.1:8080
  RAND_USER=pyrite
  RAND_Pw=somesecret
  ```

* Define a service

  ```bash
  [Unit]
  Description=Rand
  After=network.target

  [Service]
  Type=simple
  WorkingDirectory=/home/galene/rand
  User=galene
  Group=galene
  ExecStart=/home/galene/rand/rand
  LimitNOFILE=65536

  [Install]
  WantedBy=multi-user.target
  ```

### Npm

```bash
  "scripts": {
    "deploy": "cd dist;zip -r pyrite.zip .;curl -X POST -F 'distFile=@pyrite.zip' $RAND_ENDPOINT -H \"Authorization: Basic $(echo \"$RAND_USER:$RAND_PW\" | base64)\"",
    ...
  }
```

### Github actions

```bash
- name: Deploy artifacts
  env:
    RAND_ENDPOINT: ${{secrets.RAND_ENDPOINT}}
    RAND_PW: ${{secrets.RAND_PW}}
    RAND_USER: ${{secrets.RAND_USER}}
  run: |
    npm run build
    npm run deploy
```
