#!/usr/bin/env bash
export $(grep -v '^#' .env.local)
docker compose -f docker-compose.yml up