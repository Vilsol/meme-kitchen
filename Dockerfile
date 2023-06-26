FROM nixos/nix as backend

RUN nix-channel --update

WORKDIR /app

COPY . .

RUN nix-shell --command "CGO_ENABLED=1 go build -v -a -installsuffix cgo -o /go/bin/api ."


FROM node:18-alpine as frontend

WORKDIR /app

RUN npm i -g pnpm

COPY frontend/package.json package.json
COPY frontend/pnpm-lock.yaml pnpm-lock.yaml

RUN pnpm i

COPY frontend/ .

ENV NODE_ENV=production

RUN pnpm run build


FROM nixos/nix

RUN nix-channel --update

COPY --from=backend /go/bin/api /api
COPY --from=frontend /app/build/ static/

COPY fonts/ fonts/
COPY nsfw_model/ nsfw_model/

COPY shell.nix /shell.nix

RUN nix-shell --command "ldd /api"

CMD ["nix-shell", "--command", "/api"]