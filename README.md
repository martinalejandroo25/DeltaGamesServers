# DeltaGamesServers

Plataforma self-hosted para el despliegue, gestion y monitorizacion de servidores de videojuegos dedicados en contenedores Podman en modo rootless.

## Caracteristicas

- Gestion de contenedores sin privilegios de superusuario mediante Podman socket.
- Plantillas OCI listas para servidores de juegos (CS2, TF2, Garry's Mod, etc.).
- Persistencia e integracion con Systemd y Quadlet.
- Panel web con consola interactiva y metricas de consumo de recursos.
- Backend en Go con persistencia en SQLite (WAL mode).

## Estructura

- `cmd/server/`: Punto de entrada del backend en Go.
- `internal/`: Logica de conexion con Podman, base de datos y API.
- `containers/templates/`: Definiciones JSON para cada servidor de juego.
- `web/`: Interfaz web de usuario.

## Requisitos y Despliegue

1. Habilitar socket de usuario de Podman:
   ```bash
   systemctl --user enable --now podman.socket
   ```
2. Configurar variables de entorno:
   ```bash
   cp .env.example .env
   ```
3. Compilar y ejecutar backend:
   ```bash
   go run cmd/server/main.go
   ```
