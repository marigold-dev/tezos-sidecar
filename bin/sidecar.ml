open Lwt.Syntax
open Piaf

let request_uri uri =
  let open Lwt_result.Syntax in
  let* response = Piaf.Client.Oneshot.get (Uri.of_string uri) in
  if Piaf.Status.is_successful response.status
  then Piaf.Body.to_string response.body
  else Lwt.return (Error (`Msg (Piaf.Status.to_string response.status)))
;;

let create_handler tezos_uri socket =
  let handler ({ request; _ } : Unix.sockaddr Server.ctx) =
    match request.target, request.meth with
    | "/healthz", `GET ->
      let* is_bootstrapped =
        request_uri (tezos_uri ^ "chains/main/is_bootstrapped")
      in
      let response =
        match is_bootstrapped with
        | Ok json_str ->
          let open Yojson.Safe.Util in
          let json = Yojson.Safe.from_string json_str in
          let bootstrapped = member "bootstrapped" json |> to_bool in
          let sync_state = member "sync_state" json |> to_string in
          (match bootstrapped, sync_state with
          | true, "synced" ->
            Response.of_string ~body:"Tezos is bootstrapped and synced" `OK
          | false, _ ->
            Response.of_string
              ~body:"Tezos is not bootstrapped"
              `Internal_server_error
          | true, _ ->
            Response.of_string
              ~body:"Tezos is not synced"
              `Internal_server_error)
        | _ ->
          Response.of_string
            ~body:"Can't request tezos node"
            `Internal_server_error
      in
      Lwt.return response
    | _ ->
      Lwt.return (Response.of_string ~body:"Not found this route" `Not_found)
  in
  handler socket
;;

let main host port tezos_uri =
  let listen_address = Unix.(ADDR_INET (inet_addr_of_string host, port)) in
  Lwt.async (fun () ->
      let+ _server =
        Lwt_io.establish_server_with_client_socket
          listen_address
          (Server.create ?config:None (create_handler tezos_uri))
      in
      Printf.printf "Listening on port %i.\n%!" port;
      let forever, _ = Lwt.wait () in
      Lwt_main.run forever)
;;

let () =
  let port =
    match Sys.getenv_opt "PORT" with
    | Some port -> int_of_string port
    | None -> 8080
  in
  let host =
    match Sys.getenv_opt "HOST" with
    | Some host -> host
    | None -> "127.0.0.1"
  in
  let tezos_uri =
    match Sys.getenv_opt "TEZOS_URI" with
    | Some host -> host
    | None -> "https://mainnet.tezos.marigold.dev"
  in
  main host port tezos_uri
;;
