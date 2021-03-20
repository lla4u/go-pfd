#!/bin/bash
set -e

influxd restore -portable /backup
