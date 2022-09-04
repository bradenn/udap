// Copyright (c) 2022 Braden Nicholson


import axios from "axios";

const token = localStorage.getItem("token")

let headers = {
    headers: {
        Authorization: "Bearer " + token
    }
}

export default {
    async post(url: string, data?: {} | undefined): Promise<void> {
        return await axios.post(url, data, headers)
    }
}

