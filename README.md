# DeltaGamesServers 🎮🐳

Plataforma self-hosted de alto rendimiento para el despliegue, gestión y monitorización automatizada de servidores de videojuegos dedicados en contenedores **Podman** (Rootless).

## 🚀 Características Principales

- **Podman Rootless Native**: Aislamiento total de seguridad en el Home Server sin requerir privilegios de superusuario (`root`) ni el demonio Docker.
- **Soporte Universal SteamCMD**: Imágenes base OCI Debian 12 optimizadas para Source Engine (CS2, Garry's Mod, L4D2, TF2), Unity y Unreal Engine.
- **Systemd & Quadlet Integration**: Generación transparente de archivos `.container` de Quadlet para persistencia e inicio automático con la sesión del usuario host.
- **Control Panel React + TS**: Interfaz moderna con emulación de consola ANSI `xterm.js`, métricas de CPU/RAM/Red en vivo y editor de configuraciones.
- **Backend Ligero en Go**: Consumo de memoria mínimo (< 30 MB) con SQLite en modo WAL y comunicación vía Unix Socket (`/run/user/$UID/podman/podman.sock`).

## 📁 Estructura del Proyecto

Ver detalle completo en [implementation_plan.md](file:///home/martin/.gemini/antigravity/brain/973a2c96-a200-4eec-854b-525cb9eabc8b/implementation_plan.md).

```text
DeltaGamesServers/
├── cmd/server/main.go          # Servidor Backend principal Go
├── containers/
│   ├── base/
│   │   ├── Containerfile.steamcmd
│   │   └── entrypoint.sh
│   └── templates/              # Plantillas JSON por juego
├── internal/
│   ├── api/                    # Gin REST Handlers & WebSockets
│   ├── database/               # SQLite WAL Manager
│   ├── gameserver/             # Lifecycle Manager & SteamCMD
│   └── podman/                 # Podman Unix Socket REST Client
└── web/                        # React + TypeScript Frontend Dashboard
```

## 🛠️ Requisitos Previos

- **OS**: Linux (Fedora, Ubuntu, Debian, Arch, RHEL)
- **Podman**: v4.0+ con el socket de usuario activado:
  ```bash
  systemctl --user enable --now podman.socket
  ```
- **Go**: 1.22+ (para compilar el backend)
- **Node.js**: 20+ (para compilar el frontend React)

---
*Desarrollado con arquitectura de software orientada a Home Servers estables, seguros y eficientes.*
