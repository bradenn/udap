# UDAP v2.15.2
[![Go](https://github.com/bradenn/udap/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/bradenn/udap/actions/workflows/go.yml)
[![Go CodeQL](https://github.com/bradenn/udap/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/bradenn/udap/actions/workflows/codeql-analysis.yml)
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

## Glossary

| Command  | Description |
|----------| --- |
| U.D.A.P. | Universal Device Aggregation Platform (encompasses all below terms) |
| Core     | The physical running UDAP program instance |
| Nexus    | The front-end interface for all of UDAP |
| Terminal | An authoritative nexus instance (Used for configuration and management) |
| Pad      | A general use nexus instance, can be used by anyone without authentication if configured by terminal. |

## Front-end elements

#### Element

An element is a super class of a Plot, Pane, or other ui element that has a blurred background.

### Plots

A Plot is a grid-like element that contains a fixed number of cells defined by a number of rows and columns.
Plots can be configured to have a title and alt button. Plots are usually used to hold buttons or other contextual
elements.

##### Plot Selection (from Settings->Preferences page)

![Plot Buttons](./docs/images/plot_buttons.png)

##### Plot Module (from Settings->Modules page)

Plots can contain custom dom to serve whatever purpose is needed:
![Plot Buttons](./docs/images/plot_module.png)

##### Plot Buttons (from Global context Menu)

Plots are best used for providing many buttons for easy selection.

![Plot Buttons](./docs/images/plot_multi.png)

#### Copyright &copy; 2019-2022 Braden Nicholson
