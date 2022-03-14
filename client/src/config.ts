// Copyright (c) 2022 Braden Nicholson

import jsonFile from "./config.json"
import type {Config} from "@/types"

export const config = <Config>jsonFile
