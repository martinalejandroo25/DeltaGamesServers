#!/usr/bin/env bash
set -e

echo "=== DeltaGamesServers Container Entrypoint ==="
echo "Game AppID: ${STEAM_APP_ID:-Not specified}"
echo "Auto Update: ${AUTO_UPDATE:-true}"

STEAMCMD="/home/steam/steamcmd/steamcmd.sh"
GAME_DIR="/home/steam/game"

if [ -n "${STEAM_APP_ID}" ]; then
    if [ "${AUTO_UPDATE}" = "true" ] || [ ! -f "${GAME_DIR}/${START_BINARY}" ]; then
        echo "[SteamCMD] Instalando/Actualizando aplicación SteamAppID: ${STEAM_APP_ID}..."
        ${STEAMCMD} +force_install_dir ${GAME_DIR} \
                   +login ${STEAM_USER:-anonymous} ${STEAM_PASS:-} \
                   +app_update ${STEAM_APP_ID} ${STEAM_BETA_FLAG:-} validate \
                   +quit
    fi
fi

cd ${GAME_DIR}

if [ -n "${START_COMMAND}" ]; then
    echo "[DeltaGames] Ejecutando comando de inicio..."
    exec bash -c "${START_COMMAND}"
elif [ -n "${START_BINARY}" ]; then
    echo "[DeltaGames] Iniciando binario ${START_BINARY}..."
    exec ./${START_BINARY} ${START_PARAMS}
else
    echo "[DeltaGames ERROR] No se especificó START_COMMAND ni START_BINARY."
    exit 1
fi
