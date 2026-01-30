# Architecture

`toolruntime` is the execution boundary underneath `toolcode`.
It chooses a backend based on the requested security profile.

## Runtime selection


![Diagram](assets/diagrams/runtime-selection.svg)

## WASM backend

The WASM backend runs in-process using a pluggable runtime (wazero/wasmtime/wasmer).
It relies on the `backend/wasm` interfaces for module loading, health checks,
and streaming output, keeping runtime bindings out of core `toolruntime`.


## Execution sequence


![Diagram](assets/diagrams/runtime-selection.svg)


## Profiles

- `dev`: convenience, unsafe host allowed
- `standard`: safer defaults, may deny unsafe backend
- `hardened`: expected to require strong isolation
