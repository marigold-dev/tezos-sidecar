FROM ocaml/opam:alpine AS build-opam
RUN sudo apk update && sudo apk upgrade
RUN sudo apk add gmp-dev libev-dev openssl-dev
COPY tezos-sidecar.opam .
RUN opam install . --deps-only --locked 
COPY . .
RUN eval $(opam env) && sudo dune build bin/sidecar.exe
RUN sudo cp ./_build/default/bin/sidecar.exe /usr/bin/sidecar.exe
# --- #
FROM alpine AS ocaml-app
COPY --from=build-opam /usr/bin/sidecar.exe /home/app/sidecar.exe
RUN apk update && apk upgrade
RUN apk add gmp-dev libev-dev openssl-dev
RUN adduser -D app
RUN chown app:app /home/app/sidecar.exe
WORKDIR /home/app
USER app
EXPOSE 8080
ENTRYPOINT ["/home/app/sidecar.exe"]

