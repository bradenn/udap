// Copyright (c) 2022 Braden Nicholson

import jsonFile from "./config.json"
import type {ConfigFile} from "./types"

export const config: ConfigFile = jsonFile as ConfigFile
