# udap v2.9.5 beta
[![Go](https://github.com/bradenn/udap/actions/workflows/go.yml/badge.svg)](https://github.com/bradenn/udap/actions/workflows/go.yml)
[![Typescript](https://github.com/bradenn/udap/actions/workflows/ts.yml/badge.svg)](https://github.com/bradenn/udap/actions/workflows/ts.yml)
## Universal Device Aggregation Platform

Udap aims to efficiently link and aggregate many unlike interfaces into a heuristic model that can be manipulated and
supplemented from add-on plugins called modules. These modules are written in go and compile during Udap's runtime.
Modules can be configured to control computer settings, lights, music, air conditioners, phone systems, access points,
media, or even spaceships.

## Entities & Attributes

Any state within udap is stored within an attribute. These attributes belong to a parent entity.

An example:

You have a smart light bulb. It has two settings through it's api (which you've linked via a module), brightness and
color hue. Each of these settings becomes its own attribute. The smart bulb is represented as an entity whose id is
linked to the aforementioned attributes.

Another examples:

You want to connect spotify to udap. Spotify's API has a lot of options, but we can just focus three attributes:
playing, currentSong First you create an entity to represent the api, then you create and provide functions for
resolving each attribute.

## Modules & Endpoint

Both modules and Endpoints are permitted to make modification to entities and attributes. Any module and any endpoint
can concurrently modify multiple entities and attributes at a time. The command buffer can be modified to accept 4096
concurrent commands, but for larger loads, instancing is recommended.

#### Copyright &copy; 2019-2022 Braden Nicholson
