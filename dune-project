(lang dune 3.2)

(name tezos-sidecar)

(version 0.1)

(maintainers "contact@marigold.dev")

(generate_opam_files true)

(package
 (name tezos-sidecar)
 (synopsis "Tezos sidecar")
 (description "Tezos Sidecar")
 (depends
  (dune
   (and
    :build
    (>= 3.1.1)))
  piaf
  yojson
  ppx_let
  ))

(package
 (name dev)
 (synopsis "A package to install dev dependencies")
 (description "THIS PACKAGE IS FOR DEVELOPMENT PURPOSE")
 (allow_empty)
 (depends
  (ocamlformat
   (>= 0.20))
  (ocamlformat-rpc
   (>= 0.19))
  (ocaml-lsp-server
   (>= 1.10.3))))
