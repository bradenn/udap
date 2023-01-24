// Copyright (c) 2022 Braden Nicholson
export interface Wiki {
    pageID: any
    categories: string[]
    sections: Section[]
}

export interface Section {
    title: string
    depth: number
    paragraphs?: Paragraph[]
    infoboxes?: Infobox[]
    references?: DataBox[]
    images?: Image[]
    lists?: List[][]
}

export interface Paragraph {
    sentences: Sentence[]
}

export interface Sentence {
    text: string
    links?: Link[]
    formatting?: Formatting
}

export interface Link {
    text?: string
    type: string
    page?: string
    anchor?: string
    wiki?: string
}

export interface Formatting {
    italic: string[]
    bold?: string[]
}


export interface Infobox {
    name: Name
    fossil_range: DataBox
    status: DataBox
    genus: DataBox
    species: DataBox
    authority: DataBox
    synonyms: DataBox
}

export interface Name {
    text: string
    links: Link[]
}

export interface DataBox {
    text: string
}

export interface Image {
    file: string
    thumb: string
    url: string
    caption: string
    links: any[]
    alt?: string
}

export interface List {
    text: string
    links: Link[]
}


